package tasknet

import (
	"bufio"
	"errors"
	"fmt"
	"html"
	"io"
	"spinedtp/util"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid/v3"
)

// The header of a spine packet
type SpineHeader struct {
	Version           int
	Id                string
	TimeStamp         string
	Bodylength        int
	PastRoutelength   int
	FutureRoutelength int
}

// The data body of a spine packet
type SpineBody struct {
	Id    string
	Type  string
	Items map[string]string
}

// The route this spine packet followed to reachus
type SpineRoute struct {
	Nodes []*Peer // the peer info here will be very limited
}

// A spine protocol packet
type SpinePacket struct {
	Header    SpineHeader
	Body      SpineBody
	PastRoute SpineRoute // This is the route over which we received this packet

	// This is the route over which we want to send this packet
	// It must not be obeyed, as some of the peers may have disappeared
	// in the meantime
	FutureRoute SpineRoute
}

type TaskBid struct {
	ID            string
	Created       time.Time
	Fee           float64
	BidValue      float64
	Reputation    float64 // will be verified by the task owner
	TaskID        string
	TaskOwnerID   string
	BidderID      string
	EstimatedTime int
	Geo           string
	ArrivalRoute  []*Peer
	Selected      int

	// this must never be propagated or read from the network
	// Even though it can be infered from bidderID, bidderID is propagated
	// from the network and could be faked. This one is saved in local db
	// but never retrieved from the db. !17 = not my bid. 17 = My Bid
	MyBid int
}

type TaskCompleted struct {
	ID           string
	Created      time.Time
	Fee          float64
	Reward       float64
	TaskOwnerID  string
	BidderID     string
	TaskID       string
	ArrivalRoute []*Peer
}

type TaskSubmissionMedia struct {
	mimeType string
	data     []byte
}

type TaskSubmission struct {
	ID           string
	Created      time.Time
	Fee          float64
	TaskID       string
	TaskOwnerID  string
	BidderID     string
	WalletID     string
	Submissions  []TaskSubmissionMedia
	Geo          string
	ArrivalRoute []*Peer
}

type TaskBidApproval struct {
	ID           string
	Created      time.Time
	Fee          float64
	TaskID       string
	TaskOwnerID  string
	BidderID     string
	Value        float64
	Reward       float64
	Geo          string
	ArrivalRoute []*Peer
}

// Creates a very basic spine packet
func CreateSpinePacket(pastroute []*Peer, futureroute []*Peer) *SpinePacket {

	t := time.Now()
	timestamp := t.Format(time.RFC3339)

	var p SpinePacket
	p.Header.Id = shortuuid.New()
	p.Header.TimeStamp = timestamp
	p.Body.Id = p.Header.Id
	p.Body.Items = make(map[string]string)

	p.PastRoute.Nodes = pastroute
	p.FutureRoute.Nodes = futureroute

	return &p
}

// Like the above, but this packet is being forwarded. We need to add
// ourselves as part of the route
func CreateSpinePacketForForwarding(packet *SpinePacket) *SpinePacket {
	return nil
}

// Takes the data portion and wraps it in a spine packet
func CreateSpinePacketWithData(data map[string]string, pastroute []*Peer, futureroute []*Peer) *SpinePacket {

	packet := CreateSpinePacket(pastroute, futureroute)
	packet.Body.Items = data

	return packet
}

func GetMePeer() *Peer {
	var mepeer Peer
	mepeer.ID = NetworkSettings.MyPeerID
	mepeer.Address = NetworkSettings.ServerHost
	mepeer.ConnectPort = int(NetworkSettings.ServerPort)
	return &mepeer
}

// A task propagation packet is used to send a task through the entire network
func ConstructTaskPropagationPacket(tt *Task) (*SpinePacket, error) {

	pastRoute := []*Peer{GetMePeer()}
	packet := CreateSpinePacket(pastRoute, nil)

	packet.Body.Type = "task"
	packet.Body.Items["task.Command"] = tt.Command
	packet.Body.Items["task.ID"] = tt.ID
	packet.Body.Items["task.Created"] = tt.Created.Format(time.RFC3339)
	packet.Body.Items["task.Fee"] = strconv.FormatFloat(tt.Fee, 'E', -1, 64)
	packet.Body.Items["task.Reward"] = strconv.FormatFloat(tt.Reward, 'E', -1, 64)
	packet.Body.Items["task.TaskOwnerID"] = tt.TaskOwnerID
	packet.Body.Items["task.Status"] = strconv.Itoa(int(tt.GlobalStatus))
	packet.Body.Items["task.Index"] = strconv.FormatUint(tt.Index, 10)
	packet.Body.Items["task.GeoRequirement"] = "DE, FR, US"
	packet.Body.Items["task.Hash"] = tt.TaskHash

	return packet, nil
}

func ConstructPeerListRequestPacket() (*SpinePacket, error) {
	pastRoute := []*Peer{GetMePeer()}
	packet := CreateSpinePacket(pastRoute, nil)

	packet.Body.Type = "peer-list-request"

	return packet, nil
}

func ConstructPeerListPacket(peerList []string) (*SpinePacket, error) {

	pastRoute := []*Peer{GetMePeer()}
	packet := CreateSpinePacket(pastRoute, nil)

	slist := ""
	for _, s := range peerList {
		slist = slist + s + ";"
	}
	packet.Body.Type = "peer-list"
	packet.Body.Items["peer-list"] = slist

	return packet, nil
}

// When a peer receives a task that it can do, it sends a task bid to the originator
// of the task, informing it that it can do this task. This will either be ignored or
// it will get a task acceptance.
func ConstructTaskBidPacket(taskbid *TaskBid, returnRoute []*Peer) (*SpinePacket, error) {

	pastRoute := []*Peer{GetMePeer()}
	packet := CreateSpinePacket(pastRoute, returnRoute)

	packet.Body.Type = "task-bid"
	packet.Body.Items["task-bid.ID"] = taskbid.ID
	packet.Body.Items["task-bid.Created"] = taskbid.Created.Format(time.RFC3339)
	packet.Body.Items["task-bid.Fee"] = fmt.Sprintf("%f", taskbid.Fee)
	packet.Body.Items["task-bid.BidValue"] = fmt.Sprintf("%f", taskbid.BidValue)
	packet.Body.Items["task-bid.TaskOwnerID"] = taskbid.TaskOwnerID
	packet.Body.Items["task-bid.TaskID"] = taskbid.TaskID
	packet.Body.Items["task-bid.BidderID"] = NetworkSettings.MyPeerID
	packet.Body.Items["task-bid.Geo"] = taskbid.Geo
	packet.Body.Items["task-bid.Hash"] = "NOHASHYET"

	return packet, nil
}

func ConstructTaskBidApprovalPacket(taskBidApproval *TaskBidApproval, returnRoute []*Peer) (*SpinePacket, error) {

	pastRoute := []*Peer{GetMePeer()}
	packet := CreateSpinePacket(pastRoute, returnRoute)

	packet.Body.Type = "task-bid-approval"
	packet.Body.Items["task-bid-approval.ID"] = taskBidApproval.ID
	packet.Body.Items["task-bid-approval.Created"] = taskBidApproval.Created.Format(time.RFC3339)
	packet.Body.Items["task-bid-approval.Fee"] = fmt.Sprintf("%f", taskBidApproval.Fee)
	packet.Body.Items["task-bid-approval.Value"] = "0"
	packet.Body.Items["task-bid-approval.TaskOwnerID"] = taskBidApproval.TaskOwnerID
	packet.Body.Items["task-bid-approval.BidderID"] = taskBidApproval.BidderID
	packet.Body.Items["task-bid-approval.Reward"] = fmt.Sprintf("%f", taskBidApproval.Reward)
	packet.Body.Items["task-bid-approval.TaskID"] = taskBidApproval.TaskID
	packet.Body.Items["task-bid-approval.Hash"] = "NOHASHYET"

	return packet, nil
}

func ConstructTaskSubmissionPacket(taskSubmit *TaskSubmission, returnRoute []*Peer) (*SpinePacket, error) {

	pastRoute := []*Peer{GetMePeer()}
	packet := CreateSpinePacket(pastRoute, returnRoute)

	packet.Body.Type = "task-submission"
	packet.Body.Items["task-submission.ID"] = taskSubmit.ID
	packet.Body.Items["task-submission.Created"] = taskSubmit.Created.Format(time.RFC3339)
	packet.Body.Items["task-submission.Fee"] = fmt.Sprintf("%f", taskSubmit.Fee)
	packet.Body.Items["task-submission.TaskOwnerID"] = taskSubmit.TaskOwnerID
	packet.Body.Items["task-submission.BidderID"] = taskSubmit.BidderID
	packet.Body.Items["task-submission.TaskID"] = taskSubmit.TaskID
	packet.Body.Items["task-submission.Geo"] = taskSubmit.Geo
	packet.Body.Items["task-submission.Hash"] = "NOHASHYET"

	packet.Body.Items["task-submission.ResultCount"] = strconv.Itoa(len(taskSubmit.Submissions))

	// Loop through all submissions
	for i, sub := range taskSubmit.Submissions {

		mimeType := sub.mimeType
		packet.Body.Items["task-submission.Submission-"+strconv.Itoa(i)] = string(sub.data)
		packet.Body.Items["task-submission.SubmissionType-"+strconv.Itoa(i)] = mimeType
	}

	return packet, nil
}

func ConstructTaskCompletedPacket(taskCompleted *TaskCompleted, returnRoute []*Peer) (*SpinePacket, error) {

	pastRoute := []*Peer{GetMePeer()}
	packet := CreateSpinePacket(pastRoute, returnRoute)

	packet.Body.Type = "task-completed"
	packet.Body.Items["task-completed.ID"] = taskCompleted.ID
	packet.Body.Items["task-completed.Created"] = taskCompleted.Created.Format(time.RFC3339)
	packet.Body.Items["task-completed.Fee"] = "0"
	packet.Body.Items["task-completed.TaskOwnerID"] = taskCompleted.TaskOwnerID
	packet.Body.Items["task-completed.Value"] = "0"
	packet.Body.Items["task-completed.TaskID"] = taskCompleted.TaskID
	packet.Body.Items["task-completed.Hash"] = "NOHASHYET"

	return packet, nil
}

// Converts a Spine Packet structure into a string that can
// be sent over the network
func (packet *SpinePacket) ToString() string {

	// Construct Header
	var headers string
	headers = "~H:HEADER [\n"
	headers += "PROTOCOL: <SPINE>\n"
	headers += "VERSION: <10>\n"
	headers += "ID: <" + packet.Header.Id + ">\n"
	headers += "TIME: <" + packet.Header.TimeStamp + ">\n"

	// Construct Body
	var body string
	body = "~B:BODY [\n"
	body += "ID: <" + packet.Header.Id + ">\n"
	body += "TYPE: <" + packet.Body.Type + ">\n"

	for k, v := range packet.Body.Items {
		// We html escape the data, and then use angular brackets
		// to wrap it. This will allow us have all kinds of text in
		// the data, including full html
		body += k + ": <" + html.EscapeString(v) + ">\n"
	}
	body += "]\n"

	headers += "BODYLENGTH: <" + fmt.Sprint(len(body)) + ">\n"

	// Construct route
	var pastroute string
	pastroute = "~R:PASTROUTE [\n"
	for _, node := range packet.PastRoute.Nodes {
		pastroute += node.ID + ": <" + node.GetFullAddress() + ">\n"
	}
	pastroute += "]\n"

	headers += "PASTROUTELENGTH: <" + fmt.Sprint(len(pastroute)) + ">\n"

	var futureroute string
	if packet.FutureRoute.Nodes != nil && len(packet.FutureRoute.Nodes) > 0 {
		futureroute = "~R:FUTUREROUTE [\n"
		for _, node := range packet.FutureRoute.Nodes {
			futureroute += node.ID + ": <" + node.GetFullAddress() + ">\n"
		}

		futureroute += "]"

		headers += "FUTUREROUTELENGTH: <" + fmt.Sprint(len(futureroute)) + ">]"

	} else {
		headers += "]"
	}

	return headers + body + pastroute + futureroute
}

func (packet *SpinePacket) ParsePacket(b *bufio.Reader) error {

	// Let's parse the header. The ] symbol is illegal in the
	// header so can be used as a terminator
	header, err := b.ReadString(']')
	if err != nil {
		return err
	}
	if !strings.HasPrefix(header, "~H:HEADER") {
		util.PrintRed("invalid header")
		return errors.New("invalid header")
	}

	// Find first bracket so we cut away the header leader
	b1 := strings.Index(header, "[")
	data := header[b1:]

	// Remove all newlines, they are not relevant in the header
	data = strings.ReplaceAll(data, "\n", "")

	// split by our comma separator
	dataItems := strings.Split(data, ">")

	// loop over every line
	for _, item := range dataItems {

		// split into two pieces, left and right of colon
		lrItems := strings.Split(item, "<")

		// make sure it's really two
		if len(lrItems) == 2 {
			firstItem := lrItems[0]
			secondItem := lrItems[1]

			// remove all extra characters
			firstItem = strings.Trim(firstItem, "[' \r\n:")
			secondItem = strings.Trim(secondItem, "]', \r\n")

			if firstItem == "ID" {
				packet.Header.Id = secondItem
			}

			if firstItem == "VERSION" {
				val, _ := strconv.Atoi(secondItem)
				if val < 10 {
					return errors.New("not comfortable with this packet - version wrong")
				}
				packet.Header.Version = val
			}

			if firstItem == "TIME" {
				packet.Header.TimeStamp = secondItem
			}

			if firstItem == "BODYLENGTH" {
				le, _ := strconv.Atoi(secondItem)
				packet.Header.Bodylength = le
			}

			if firstItem == "PASTROUTELENGTH" {
				le, _ := strconv.Atoi(secondItem)
				packet.Header.PastRoutelength = le
			}

			if firstItem == "FUTUREROUTELENGTH" {
				le, _ := strconv.Atoi(secondItem)
				packet.Header.FutureRoutelength = le
			}
		}
	}

	if len(packet.Header.Id) < 10 {
		return errors.New("not comfortable with this packet")
	}

	// This piece here is to progressively read the file if it does not all come at once

	// Create buffer large enough
	bodyBytes := make([]byte, packet.Header.Bodylength)

	totalRead := 0
	// try to read full buffer
	count, err := b.Read(bodyBytes)
	totalRead += count

	// Check if we read all. If not, loop through and progressively read more
	for count > 0 && err != io.EOF && totalRead < packet.Header.Bodylength {

		// make buffer big enough for everything
		bodyTmpBytes := make([]byte, packet.Header.Bodylength-totalRead)

		// Read what is left
		count, err = b.Read(bodyTmpBytes)

		// append to the original buffer
		bodyBytes = append(bodyBytes[:totalRead], bodyTmpBytes[:count]...)

		// increment amount read
		totalRead += count
	}

	if err != nil || totalRead != packet.Header.Bodylength {
		return errors.New("not comfortable with this packet. len don't match")
	}

	packet.Body.Items, err = PullDataPartOut("~B:BODY", string(bodyBytes), false)
	if err != nil {
		return err
	}

	// Confirm that header and body ID matched
	bodyID, ok := packet.Body.Items["ID"]
	if !ok {
		return errors.New("invalid body ID")
	}

	if bodyID != packet.Header.Id {
		return errors.New("body and header ID do not match")
	}

	bodyType, ok := packet.Body.Items["TYPE"]
	if !ok {
		return errors.New("invalid body Type")
	}

	packet.Body.Type = bodyType

	// Parse Route
	pastRouteBytes := make([]byte, packet.Header.PastRoutelength)

	count, err = b.Read(pastRouteBytes)
	if err != nil || count != packet.Header.PastRoutelength {
		return errors.New("not comfortable with this route packet. len don't match")
	}

	pastRouteString := string(pastRouteBytes)
	// fmt.Println(routeString)

	pastRouteDataItems, err := PullDataPartOut("~R:PASTROUTE", pastRouteString, true)
	if err != nil {
		return errors.New("could not parse route")
	}

	for pid, addr := range pastRouteDataItems {
		var peer Peer
		peer.ID = pid
		peer.Address = addr
		packet.PastRoute.Nodes = append(packet.PastRoute.Nodes, &peer)
	}

	// Parse Route
	if packet.Header.FutureRoutelength > 0 {
		futureRouteBytes := make([]byte, packet.Header.FutureRoutelength)

		count, err = b.Read(futureRouteBytes)
		if err != nil || count != packet.Header.FutureRoutelength {
			return errors.New("not comfortable with this route packet. len don't match")
		}

		futureRouteString := string(futureRouteBytes)

		futureRouteDataItems, err := PullDataPartOut("~R:FUTUREROUTE", futureRouteString, true)
		if err != nil {
			return errors.New("could not parse route")
		}

		for pid, addr := range futureRouteDataItems {
			var peer Peer
			peer.ID = pid
			peer.Address = addr
			packet.FutureRoute.Nodes = append(packet.FutureRoute.Nodes, &peer)
		}
	}

	return nil

}

func PullDataPartOut(protocolTitle string, dataString string, removeNewLines bool) (map[string]string, error) {

	// Remove the leading newlines
	dataString = strings.TrimLeft(dataString, "\n\r")

	// confirm that protocol is right
	if !strings.HasPrefix(dataString, protocolTitle) {
		return nil, errors.New("invalid body header")
	}

	// find the first [. It follows the protocol
	b1 := strings.Index(dataString, "[")
	if b1 == -1 {
		return nil, errors.New("invalid body")
	}

	// Get the body data, excluding the protocol title
	dataSection := dataString[b1:]

	if removeNewLines {
		// Remove all newlines, they are not relevant in the data
		dataSection = strings.ReplaceAll(dataSection, "\n", "")
	}

	// split by our line separator
	dataItems := strings.Split(dataSection, ">")

	result := make(map[string]string)

	// loop over every line
	for _, item := range dataItems {

		// split into two pieces, left and right of colon
		lrItems := strings.Split(item, "<")

		// make sure it's really two
		if len(lrItems) == 2 {
			firstItem := lrItems[0]
			secondItem := lrItems[1]

			// remove all extra characters
			firstItem = strings.Trim(firstItem, "[' \r\n:")
			secondItem = strings.Trim(secondItem, "]', \r\n")

			result[firstItem] = html.UnescapeString(secondItem)
		}
	}

	return result, nil
}
