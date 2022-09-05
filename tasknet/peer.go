package tasknet

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"spinedtp/util"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid/v3"
)

type Peer struct {
	conn          net.Conn
	reader        *bufio.Reader
	ID            string
	Address       string
	Port          int
	Connected     bool
	Deadness      uint
	IsMetaTracker bool     // indicates if I should publish my peer list here
	PeerList      []string // a list of all the peers this peer is connected to. Not guaranteed to be there
	SiblingCount  uint64   // number of people it is connected to
	WeConnected   bool     // true if we built connection, false if it connected to us
	FirstCommand  string   // if we connect to this peer for a specific reason, we can specify it here. E.g "peers"

}

// All peers I know, connected or not
var Peers []*Peer

// Maximum number of peers we will remember at all (not neccessarily)
// connected right now
// var maxPeersMemory uint

func CreatePeer(c net.Conn) *Peer {

	saddr := c.RemoteAddr().String()
	saddrl := strings.Split(saddr, ":")

	var p Peer
	p.Address = saddrl[0]
	p.Port, _ = strconv.Atoi(saddrl[1])

	p.conn = c

	return &p
}

func LoadPeerTable() {

	file, err := os.Open(filepath.Join(NetworkSettings.DataFolder, "peers.txt"))
	if err != nil {
		fmt.Println("No Peers.txt file")
		return
	}
	defer func() {
		file.Close()
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() { // internally, it advances token based on sperator

		AddPeerWithIPColonPort(scanner.Text())
	}

}

func AddPeerWithIPColonPort(ip_port string) *Peer {
	saddr := strings.Split(ip_port, ":")
	if len(saddr) != 2 {
		return nil
	}

	p, _ := strconv.Atoi(saddr[1])
	return AddPeerWithIPAndPort(saddr[0], p)
}

func AddPeerWithIPAndPort(ip string, port int) *Peer {
	var peer Peer

	peer.Address = ip
	peer.Port = port

	AddToPeerTable(&peer)

	return &peer
}

func SavePeerTable() {

	file, err := os.Create(filepath.Join(NetworkSettings.DataFolder, "peers.txt"))
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	addr := map[string]bool{}

	for _, peer := range Peers {

		// avoid repeating addresses. addreses can be repeated in online peer table because of
		// IDs, but no need in our file
		_, exists := addr[peer.GetFullAddress()]
		if !exists {
			addr[peer.GetFullAddress()] = true
			_, _ = datawriter.WriteString(peer.GetFullAddress() + "\n")
		}
	}

	datawriter.Flush()
	file.Close()

}

func AddToPeerTable(peer *Peer) {

	// It is possible that the same peer reconnected, but we have the same
	// peer ID in our table already with an old IP. Let's first check that
	for i, epeer := range Peers {
		if epeer.ID == peer.ID {

			// update the peer if the address has changed
			// TODO: This is a potential bug - anyone can claim the ID of another peer
			if epeer.Address != peer.Address || epeer.Port != peer.Port {
				Peers[i] = peer
			}
			return
		}
	}

	Peers = append(Peers, peer)
}

func FindPeer(peerID string) *Peer {

	for _, peer := range Peers {
		if peer.ID == peerID {
			return peer
		}
	}

	return nil
}

func (peer *Peer) ValidateConnectString(reader *bufio.Reader) error {
	packet, err := peer.Read(reader)
	if err != nil {
		return err
	}

	if len(packet.Body.Items) <= 0 {
		return errors.New("no items in the body")
	}

	peerID, ok := packet.Body.Items["PEER-ID"]
	if !ok {
		return errors.New("no peer-ID provided")
	}

	siblingCount, ok := packet.Body.Items["SIBLING-COUNT"]
	if !ok {
		return errors.New("no sibling count provided")
	}

	u, _ := strconv.ParseUint(siblingCount, 10, 64)

	peer.SiblingCount = u
	peer.ID = peerID

	return nil
}

func CountPeers() (int, int) {
	WeConnectedCount := 0
	TheyConnectedCount := 0
	for _, peer := range Peers {
		if peer.WeConnected && peer.Connected {
			WeConnectedCount = WeConnectedCount + 1
		}

		if !peer.WeConnected && peer.Connected {
			TheyConnectedCount = TheyConnectedCount + 1
		}
	}

	return WeConnectedCount, TheyConnectedCount
}

func (peer *Peer) RequestPeerList() error {
	packet, err := ConstructPeerListRequestPacket()
	if err != nil {
		return err
	}

	count, err := peer.conn.Write([]byte(packet.ToString()))
	if count == 0 || err != nil {
		util.PrintRed("Error when writing to socket: " + err.Error())
	}

	return nil
}

func (peer *Peer) SendConnectString() {

	data := map[string]string{
		"PEER-ID":       NetworkSettings.MyPeerID,
		"SIBLING-COUNT": fmt.Sprint(len(peer.PeerList)),
		"SANITY-CHECK":  "(>人<)&¯\\_(ツ)_/¯", // This is to ensure the parser is correctly escaping these characters
	}

	peer.Send(data, nil, MakeRoute(peer))
}

func (peer *Peer) Send(data map[string]string, PastRoute []*Peer, FutureRoute []*Peer) {
	packet := CreateSpinePacketWithData(data, PastRoute, FutureRoute)
	count, err := peer.conn.Write([]byte(packet.ToString()))
	if count == 0 || err != nil {
		util.PrintRed("Error when writing to socket" + err.Error())
	}
}

// Read a packet from this peer.
func (peer *Peer) Read(reader *bufio.Reader) (*SpinePacket, error) {

	var p = CreateSpinePacket(nil, nil)
	err := p.ParsePacket(reader)

	return p, err
}

func (p *Peer) GetFullAddress() string {
	return p.Address + ":" + fmt.Sprint(p.Port)
}

func (peer *Peer) SetBad() {

}

func (peer *Peer) SendPeerList(peerList []string) error {

	packet, err := ConstructPeerListPacket(peerList)
	if err != nil {
		return err
	}

	count, err := peer.conn.Write([]byte(packet.ToString()))
	if count == 0 || err != nil {
		util.PrintRed("Error when writing to socket: " + err.Error())
	}

	return err
}

func (peer *Peer) BidForTask(task *Task) error {

	var t TaskBid
	t.BidValue = task.Reward - 0.0001
	t.BidderID = NetworkSettings.MyPeerID
	t.Fee = 0
	t.Geo = "US"
	t.ID = shortuuid.New()
	t.TaskOwnerID = task.TaskOwnerID
	t.TaskID = task.ID
	t.Created = time.Now()

	packet, err := ConstructTaskBidPacket(&t, task.GetReturnRoute())
	if err != nil {
		return err
	}

	count, err := peer.conn.Write([]byte(packet.ToString()))
	if count == 0 || err != nil {
		util.PrintRed("Error when writing to socket: " + err.Error())
	}

	return err
}

func (peer *Peer) AcceptBid(task *Task, taskbid TaskBid) error {

	var t TaskAccept
	t.BidderID = NetworkSettings.MyPeerID
	t.Created = time.Now()
	t.Fee = 0
	t.ID = shortuuid.New()
	t.Reward = taskbid.BidValue
	t.TaskID = task.ID

	packet, err := ConstructTaskAcceptPacket(&t, taskbid.GetReturnRoute())
	if err != nil {
		return err
	}

	count, err := peer.conn.Write([]byte(packet.ToString()))
	if count == 0 || err != nil {
		util.PrintRed("Error when writing to socket" + err.Error())
	}

	return nil
}

func (peer *Peer) SubmitTaskResult(task *Task) error {

	var taskSubmit TaskSubmission
	taskSubmit.BidderID = NetworkSettings.MyPeerID
	taskSubmit.Submission = task.Result
	taskSubmit.Created = time.Now()
	taskSubmit.ID = shortuuid.New()
	taskSubmit.TaskID = task.ID
	taskSubmit.Fee = 0
	taskSubmit.Geo = "US"
	taskSubmit.TaskOwnerID = task.TaskOwnerID
	taskSubmit.ArrivalRoute = task.ArrivalRoute

	packet, err := ConstructTaskSubmissionPacket(&taskSubmit, taskSubmit.GetReturnRoute())
	if err != nil {
		return err
	}

	count, err := peer.conn.Write([]byte(packet.ToString()))
	if count == 0 || err != nil {
		util.PrintRed("Error when writing to socket" + err.Error())
	}

	return nil
}
