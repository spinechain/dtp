package tasknet

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	util "spinedtp/util"

	"github.com/gotk3/gotk3/glib"
)

// Types of Callback functions
type StatusUpdateFn func(string, int)

type NetSettings struct {
	ServerPort             uint
	ServerHost             string
	MyPeerID               string
	MaximumPeers           uint
	BidTimeoutSeconds      time.Duration
	AcceptedBidsPerTask    int
	OnStatusUpdate         StatusUpdateFn
	TaskReadyForProcessing func(string)
	DataFolder             string
}

type TaskReceiveFn func(string)
type TaskApprovedFn func(string)

type NetCallbacks struct {
	OnTaskReceived TaskReceiveFn
	OnTaskApproved TaskApprovedFn
}

// For storing everything we need to participate in the Spine network
var NetworkSettings NetSettings
var NetworkCallbacks NetCallbacks

var listeningForPeers bool = false
var requestDisconnect = false

// This will connect this node into the
func Connect() {

	if len(NetworkSettings.MyPeerID) < 3 {
		panic("Network settings have not been set!")
	}

	// Create new channel to wait for tasks
	taskForProcessingAvailable = make(chan int)
	taskForExecutionAvailable = make(chan int)

	// listen for anyone connecting to us
	go listenForPeers()

	// process any tasks in the queue
	go ProcessTasks()

	// connect to any known peers
	go connectToPeers()
}

func Disconnect() {
	// SleepTasks()

	fmt.Println("Shutting down TaskPool...")
	requestDisconnect = true

	SavePeerTable()

	ShutDownTaskRunner()
}

func ExecNetworkCommand(cmd string) {
	SendTaskToNetwork(cmd)

	// tasks := GetMyTaskList()
	// OnNewTaskReceivedMine(tasks)
}

func RoutePacketOn() {

}

func connectToPeers() {

	time.Sleep(5 * time.Second)

	StatusBarUpdate("📺 Connecting to peers...", 1)

	// Get all known peers from the DB
	LoadPeerTable()
	SavePeerTable()

	StatusBarUpdate(fmt.Sprint(len(Peers))+" local peer(s) found", 1)

	// attempt to build connections to each peer
	for _, peer := range Peers {

		if peer.Connected {
			continue
		}

		// connect to a single peer
		c, err := net.Dial("tcp", peer.GetFullAddress())
		if err != nil {
			fmt.Println(err)
			return
		}

		peer.conn = c
		go handlePeerConnection(peer, true)

	}

}

func StatusBarUpdate(str string, section int) {
	fmt.Println(str)

	glib.TimeoutAdd(10, func() bool {

		if NetworkSettings.OnStatusUpdate != nil {
			NetworkSettings.OnStatusUpdate(str, section)
		}

		return false
	})
}

// This will create a new task in our local task pool.
// Immediately after, it will propagate open tasks in
// our task pool into the next nodes. It will favour
// the latest task we just droppped of course, but it will
// also handle other open tasks it may have received
// from other nodes.
func SendTaskToNetwork(text string) {

	task := CreateNewTask(text)
	OpenTaskPool.AddTask(task)
	CheckForNewTasks()
}

// This function starts this client listening on this port for other clients.
func listenForPeers() {

	if listeningForPeers {
		StatusBarUpdate("📡 Listening for peers on "+fmt.Sprint(NetworkSettings.ServerPort), 0)
		return
	}

	StatusBarUpdate("📡 Listening for peers on "+fmt.Sprint(NetworkSettings.ServerPort), 0)

	NetworkSettings.ServerHost = "" // TODO
	l, err := net.Listen("tcp4", NetworkSettings.ServerHost+":"+fmt.Sprint(NetworkSettings.ServerPort))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	listeningForPeers = true
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		peer := CreatePeer(c)

		StatusBarUpdate("🔄 New connection from "+peer.Address, 0)
		go handlePeerConnection(peer, false)

		if requestDisconnect {
			break
		}
	}

	listeningForPeers = false
}

/*
func ConnectAndRequestPeers(tgt_ip string, tgt_port int) {

	print("Requesting peers from " + tgt_ip)

	peer := AddPeerWithIPAndPort(tgt_ip, tgt_port)
	// add this as a new peer
	// connect to it
	// request new peers

	peer.FirstCommand = "peers" // this will make it request peers once it connects

	connectToPeers()

}
*/

func MakeRoute(peer *Peer) []*Peer {

	route := make([]*Peer, 0)

	route = append(route, peer)
	return route
}

// Sends a spine packet to every single connected peer
func SendPacketToAllPeers(packet *SpinePacket) {

	for _, peer := range Peers {
		peer.conn.Write([]byte(packet.ToString()))
	}
}

func handlePeerConnection(peer *Peer, weConnected bool) {

	peer.reader = bufio.NewReader(peer.conn)

	if weConnected {
		util.PrintYellow("🔄 Connected to peer. Sending connect string...")

		// When we connect, we should send a message to the other side
		// identifying ourselves
		peer.SendConnectString()

		StatusBarUpdate("🔄 Waiting for connect string response...", 1)
		val := peer.ValidateConnectString(peer.reader)
		if val != nil {
			peer.SetBad()
			fmt.Println("Invalid connection string")
			return
		}

		peer.WeConnected = true
		peer.Connected = true

		UpdatePeerCount()

		if peer.FirstCommand == "peers" {
			util.PrintPurple("Requesting list of peers from this peer")
			peer.RequestPeerList()
		}

	} else {
		// If not, then we expect a message from the other side identifying
		// themselves

		StatusBarUpdate("🔄 Validating connection string...", 0)
		val := peer.ValidateConnectString(peer.reader)
		if val != nil {
			peer.SetBad()
			fmt.Println("Invalid connection string")
			return
		}

		// This is a valid peer. We connect it on our side
		AddToPeerTable(peer)
		peer.WeConnected = false
		peer.Connected = true

		UpdatePeerCount()

		// We now respond with the same connection protocol
		peer.SendConnectString()

	}

	for {

		// Receive a message from this peer
		packet, err := peer.Read(peer.reader)
		if err == io.EOF {

			util.PrintYellow("Disconnected from peer")
			// disconnected
			peer.Connected = false

			UpdatePeerCount()
			break
		}

		if err != nil {
			continue
		}

		// Check if there was actually any found command, if not skip
		if len(packet.Body.Items) <= 0 {
			continue
		}

		util.PrintYellow("New Packet Arrived")

		// Process packet
		ReceivePacket(packet, peer)

	}

}

func UpdatePeerCount() {
	WeConnectedCount, ConnectedToUsCount := CountPeers()

	if ConnectedToUsCount == 0 {
		StatusBarUpdate("Server: 🌐 No connections", 0)
	} else {
		StatusBarUpdate(fmt.Sprintf("Server: 🌐 Connected - %d peers", ConnectedToUsCount), 0)
	}

	if WeConnectedCount == 0 {
		StatusBarUpdate("Client: 🌐 No Connections", 1)
	} else {
		StatusBarUpdate(fmt.Sprintf("Client: 🌐 Connected - %d peers", WeConnectedCount), 1)
	}

}
