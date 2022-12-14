package tasknet

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"spinedtp/util"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid/v3"
)

type Peer struct {
	conn                net.Conn
	reader              *bufio.Reader
	ID                  string
	Address             string
	ConnectPort         int // to connect to peer
	ActivePort          int // to maintain connections with this peer
	Connected           bool
	ConnectSuccessCount int
	ConnectFailCount    int
	LastConnected       int
	Deadness            uint
	IsMetaTracker       bool     // indicates if I should publish my peer list here
	PeerList            []string // a list of all the peers this peer is connected to. Not guaranteed to be there
	SiblingCount        uint64   // number of people it is connected to
	OutConnection       bool     // true if we built connection, false if it connected to us
	InConnection        bool
	FirstCommand        string // if we connect to this peer for a specific reason, we can specify it here. E.g "peers"

}

// All peers I know, connected or not
var Peers []*Peer
var defaultPeersLoaded bool = false

// Maximum number of peers we will remember at all (not neccessarily)
// connected right now
// var maxPeersMemory uint

func CreatePeer(c net.Conn) *Peer {

	saddr := c.RemoteAddr().String()
	saddrl := strings.Split(saddr, ":")

	var p Peer
	p.Address = saddrl[0]
	p.ActivePort, _ = strconv.Atoi(saddrl[1])

	p.conn = c

	return &p
}

func CreatePeerFromIPAndPort(ip string, port int) *Peer {
	var peer Peer
	peer.Address = ip
	peer.ConnectPort = port

	return &peer
}

func LoadPeerTable() error {

	filePath := filepath.Join(NetworkSettings.DbFolder, "peers.db")
	var create bool
	if !util.FileExists(filePath) {
		create = true
	}

	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}

	defer db.Close()

	if create {
		sqlStmt := `drop table peers;`
		_, err := db.Exec(sqlStmt)
		if err != nil {
			// util.PrintRed("Could not drop peer table")
		}

		sqlStmt = "create table peers (pid string not null unique primary key,address text, port int, connect_success int, connect_fail int, last_connected int);"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			util.PrintRed("Could not create peer table")
			return err
		}
	}

	rows, err := db.Query("SELECT * FROM peers")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {

		var peer Peer

		err = rows.Scan(&peer.ID, &peer.Address,
			&peer.ConnectPort, &peer.ConnectSuccessCount,
			&peer.ConnectFailCount,
			&peer.LastConnected)
		if err == nil {
			AddToPeerTable(&peer)
		}
	}

	return nil
}

func LoadDefaultPeerTable(default_peers string) {

	if defaultPeersLoaded {
		return
	}

	defaultPeersLoaded = true
	fmt.Println("Loading default peers: " + default_peers)

	// Check if the peers.txt in the data folder file exists
	// If it does, we will use that instead of the default peers

	peers_file := filepath.Join(NetworkSettings.DbFolder, "peers.txt")
	if util.FileExists(peers_file) {
		// read the file content into the default_peers string
		fileData, _ := util.ReadFile(peers_file)

		if len(fileData) > 0 {
			default_peers = fileData
		}
	}

	// Load default peers from file
	single_peers := strings.Split(default_peers, "\n")
	for _, single_peer := range single_peers {
		// Split the peer into address and port
		single_peer = strings.TrimSpace(single_peer)
		if single_peer == "" {
			continue
		}

		single_peer_l := strings.Split(single_peer, ":")
		if len(single_peer_l) != 2 {
			continue
		}

		port, err := strconv.Atoi(single_peer_l[1])
		if err != nil {
			continue
		}

		peer := CreatePeerFromIPAndPort(single_peer_l[0], port)
		AddToPeerTable(peer)
	}

}

func SavePeerTable() error {

	filePath := filepath.Join(NetworkSettings.DbFolder, "peers.db")
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}

	defer db.Close()

	// delete
	stmt, err := db.Prepare("delete from peers")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	for _, peer := range Peers {
		s := fmt.Sprintf("INSERT INTO peers(pid, address, port, connect_success, connect_fail, last_connected) values(?,?,?,?,?,?)")

		// insert
		stmt, err := db.Prepare(s)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(peer.ID, peer.Address, peer.ConnectPort, peer.ConnectSuccessCount,
			peer.ConnectFailCount, peer.LastConnected)
		if err != nil {
		}
	}

	return nil
}

func AddToPeerTable(peer *Peer) *Peer {

	var peerFound bool = false
	// It is possible that the same peer reconnected, but we have the same
	// peer ID in our table already with an old IP. Let's first check that
	for i, epeer := range Peers {
		if epeer.ID != "" && epeer.ID == peer.ID {
			// TODO: We REALLY have to look at the implications of what is being done
			// here. It most likely opens up a vulnerability here.

			// TODO: This is a potential bug - anyone can claim the ID of another peer
			//if epeer.Address != peer.Address || epeer.Port != peer.Port {
			//	Peers[i] = peer
			// }

			peer.ConnectPort = epeer.ConnectPort

			if epeer.InConnection {
				peer.InConnection = true
			}

			if epeer.OutConnection {
				peer.OutConnection = true
			}

			Peers[i] = peer
			return Peers[i]
		}

		if epeer.Address == peer.Address && epeer.ConnectPort == peer.ConnectPort && epeer.ID != "" && epeer.ID == peer.ID {
			peerFound = true
		}

		if epeer.ID == "" && epeer.Address == peer.Address && epeer.ConnectPort == peer.ConnectPort {
			epeer.ID = peer.ID
			return epeer
		}
	}

	if !peerFound {
		Peers = append(Peers, peer)
	}

	return peer
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

func CountPeers() (int, int, string) {
	WeConnectedCount := 0
	TheyConnectedCount := 0
	var peerName string
	for _, peer := range Peers {
		if peer.OutConnection && peer.IsConnected() {
			WeConnectedCount = WeConnectedCount + 1

			peerName = peer.Address
		}

		if peer.InConnection && peer.IsConnected() {
			TheyConnectedCount = TheyConnectedCount + 1
		}
	}

	return WeConnectedCount, TheyConnectedCount, peerName
}

func (peer *Peer) SendPacket(packet *SpinePacket) error {
	if peer.IsConnected() {
		_, err := peer.conn.Write([]byte(packet.ToString()))
		if err != nil {
			util.PrintRed("Error sending packet to peer: " + err.Error())
		}
		return err
	}

	return errors.New("peer is not connected")
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
		"SANITY-CHECK":  "(>???<)&??\\_(???)_/??", // This is to ensure the parser is correctly escaping these characters
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
func (peer *Peer) IsConnected() bool {
	return peer.Connected
}

// Read a packet from this peer.
func (peer *Peer) Read(reader *bufio.Reader) (*SpinePacket, error) {

	var p = CreateSpinePacket(nil, nil)
	err := p.ParsePacket(reader)

	return p, err
}

func (p *Peer) GetFullAddress() string {
	return p.Address + ":" + fmt.Sprint(p.ConnectPort)
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

func (peer *Peer) SubmitTaskResult(task *Task) error {

	if !peer.IsConnected() {
		return errors.New("peer is not connected")
	}

	var taskSubmit TaskSubmission
	taskSubmit.BidderID = NetworkSettings.MyPeerID
	taskSubmit.Created = time.Now()
	taskSubmit.ID = shortuuid.New()
	taskSubmit.TaskID = task.ID
	taskSubmit.Fee = 0
	taskSubmit.Geo = "US"
	taskSubmit.TaskOwnerID = task.TaskOwnerID
	taskSubmit.ArrivalRoute = task.ArrivalRoute

	// Loop over all task results
	for _, taskResult := range task.Results {
		// Copy to taskSubmit.Submissions
		var tm TaskSubmissionMedia
		tm.data = taskResult.Data
		tm.mimeType = taskResult.MimeType
		taskSubmit.Submissions = append(taskSubmit.Submissions, tm)
	}

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
