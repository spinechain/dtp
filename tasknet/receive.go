package tasknet

import (
	"fmt"
	util "spinedtp/util"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/glib"
)

// This file has all the packages that arrive in the network

func ReceivePacket(packet *SpinePacket, peer *Peer) {

	switch packetType := packet.Body.Type; packetType {
	case "peer-list-request":
		util.PrintPurple("Request for peers Received")
		ReceiveRequestForPeerList(packet, peer)
	case "peer-list":
		util.PrintPurple("Peer List Received")
		ReceivePeerList(packet, peer)
	case "task":
		ReceiveTask(packet)
	case "task-bid":
		ReceiveTaskBid(packet)
	case "task-approval":
		ReceiveTaskApproval(packet)
	case "task-submission":
		ReceiveTaskSubmission(packet)
	case "task-accept":
		ReceiveTaskAccept(packet)
	default:
		fmt.Printf("Unknown packet type received: %s.\n", packetType)
	}
}

func ReceivePeerList(packet *SpinePacket, peer *Peer) {

	//s_peer_list := packet.Body.Items["peer-list"]
	//peer_list := strings.Split(s_peer_list, ";")

	//for _, ip_port := range peer_list {
	// AddPeerWithIPColonPort(ip_port)
	//}
	SavePeerTable()
}

func ReceiveRequestForPeerList(packet *SpinePacket, peer *Peer) {

	// Find the response peer
	// PRepare the response protocol
	// Get all our peers and put inside
	// send it back to the person
	peer_list := []string{}
	for _, epeer := range Peers {
		peer_list = append(peer_list, epeer.GetFullAddress())
	}

	peer.SendPeerList(peer_list)
}

func ReceiveTask(packet *SpinePacket) {
	util.PrintYellow("Received new task: " + packet.Body.Items["task.Command"])

	var task Task
	task.Command = packet.Body.Items["task.Command"]

	t1, err := time.Parse(time.RFC3339, packet.Body.Items["task.Created"])
	if err != nil {
		util.PrintYellow("Invalid task time received: " + task.Command)
		return
	}

	ffee, err := strconv.ParseFloat(packet.Body.Items["task.Fee"], 64)
	if err != nil {
		util.PrintYellow("Invalid task fee received: " + task.Command)
		return
	}

	freward, err := strconv.ParseFloat(packet.Body.Items["task.Reward"], 64)
	if err != nil {
		util.PrintYellow("Invalid task reward received: " + task.Command)
		return
	}

	status, _ := strconv.Atoi(packet.Body.Items["task.Status"])

	task.ID = packet.Body.Items["task.ID"]
	task.Created = t1
	task.Fee = ffee
	task.Reward = freward
	task.LocalWorkerStatus = StatusNewFromNetwork
	task.LocalWorkProviderStatus = StatusNewFromNetwork
	task.GlobalStatus = GlobalTaskStatus(status)
	task.TaskOwnerID = packet.Body.Items["task.TaskOwnerID"]
	task.FullyPropagated = false
	task.Index = OpenTaskPool.highestIndex + 1
	task.ArrivalRoute = packet.PastRoute.Nodes

	OpenTaskPool.AddToTaskPool(&task)

	// This changes the thread and informs the UI about this new task
	glib.TimeoutAdd(10, func() bool {
		if NetworkCallbacks.OnTaskReceived != nil {
			NetworkCallbacks.OnTaskReceived(packet.Body.Items["task.Command"])
		}
		return false
	})

	// This will tell the network to start processing the task
	taskForProcessingAvailable <- 1
}

// Call when a new task bid arrives. We add it to the database. When our bid
// period expires is when we check for all bids and select the best
func NewTaskBidArrived(tb *TaskBid) {

	util.PrintPurple("New Task Bid Arrived for Task " + tb.TaskID + " from " + PeerIDToDescription(tb.BidderID))

	if tb.TaskOwnerID == NetworkSettings.MyPeerID {
		// This is a bid for a task of mine

		AddBid(taskDb, tb, false)

	} else {
		// This is a bid for another peer that is not me. We route
		// it to the best connection we have
		util.PrintPurple("Task bid for another client: " + tb.BidderID)
		RouteTaskBidOn(tb)
	}
}

func ReceiveTaskBid(packet *SpinePacket) {
	// util.PrintPurple("Received new task bid: " + packet.Body.Items["task-bid.TaskOwnerID"])

	// When we receive a bid, it may be for us, or it may be for another peer. If
	// it's for us, we can process it. Otherwise we look for one of the peers that
	// can route it.
	// If it's for us, we put in our "selection queue" to wait for other bids to come
	// in. When enough bids have come in, we will select. The timeout for this has to be
	// at least 5 minutes or so.

	created, err := time.Parse(time.RFC3339, packet.Body.Items["task-bid.Created"])

	if err != nil {
		util.PrintPurple("invalid task bid time received")
		return
	}

	ffee, err := strconv.ParseFloat(packet.Body.Items["task-bid.Fee"], 64)
	if err != nil {
		util.PrintPurple("Invalid task bid received")
		return
	}

	bidvalue, err := strconv.ParseFloat(packet.Body.Items["task-bid.BidValue"], 64)
	if err != nil {
		util.PrintPurple("invalid bid value received")
		return
	}

	var t TaskBid
	t.ID = packet.Body.Items["task-bid.ID"]
	t.TaskID = packet.Body.Items["task-bid.TaskID"]
	t.Created = created
	t.Fee = ffee
	t.BidValue = bidvalue
	t.BidderID = packet.Body.Items["task-bid.BidderID"]
	t.TaskOwnerID = packet.Body.Items["task-bid.TaskOwnerID"]
	t.Geo = packet.Body.Items["task-bid.Geo"]
	t.ArrivalRoute = packet.PastRoute.Nodes

	NewTaskBidArrived(&t)
}

func ReceiveTaskApproval(packet *SpinePacket) {

	created, err := time.Parse(time.RFC3339, packet.Body.Items["task-bid.Created"])

	if err != nil {
		fmt.Println("invalid task bid time received")
		return
	}

	fee, err := strconv.ParseFloat(packet.Body.Items["task-approval.Fee"], 64)
	if err != nil {
		fmt.Println("Invalid task bid received")
		return
	}

	value, err := strconv.ParseFloat(packet.Body.Items["task-approval.Value"], 64)
	if err != nil {
		fmt.Println("invalid bid value received")
		return
	}

	var t TaskApproval
	t.ID = packet.Body.Items["task-approval.ID"]
	t.TaskID = packet.Body.Items["task-approval.TaskID"]
	t.Created = created
	t.Fee = fee
	t.Value = value
	t.TaskOwnerID = packet.Body.Items["task-approval.TaskOwnerID"]
	t.Geo = packet.Body.Items["task-approval.Geo"]
	t.ArrivalRoute = packet.PastRoute.Nodes

	if NetworkCallbacks.OnTaskApproved != nil {
		NetworkCallbacks.OnTaskApproved("yes")
	}
}

func ReceiveTaskSubmission(packet *SpinePacket) {

	created, err := time.Parse(time.RFC3339, packet.Body.Items["task-submission.Created"])

	if err != nil {
		fmt.Println("invalid task bid time received")
		return
	}

	fee, err := strconv.ParseFloat(packet.Body.Items["task-submission.Fee"], 64)
	if err != nil {
		fmt.Println("Invalid task bid received")
		return
	}

	var t TaskSubmission
	t.ID = packet.Body.Items["task-submission.ID"]
	t.TaskID = packet.Body.Items["task-submission.TaskID"]
	t.Created = created
	t.Fee = fee
	// t.Submission = []byte(packet.Body.Items["task-submission.Submission"])
	t.TaskOwnerID = packet.Body.Items["task-submission.TaskOwnerID"]
	t.Geo = packet.Body.Items["task-submission.Geo"]
	t.ArrivalRoute = packet.PastRoute.Nodes

	SubmissionCount := packet.Body.Items["task-submission.ResultCount"]
	SubmissionCountInt, err := strconv.Atoi(SubmissionCount)
	if err != nil {
		fmt.Println("Invalid task bid received")
		return
	}

	// Loop over all results
	for i := 0; i < SubmissionCountInt; i++ {
		var result TaskSubmissionMedia
		result.data = []byte(packet.Body.Items["task-submission.Submission-"+strconv.Itoa(i)])
		result.mimeType = packet.Body.Items["task-submission.SubmissionType-"+strconv.Itoa(i)]
		t.Submissions = append(t.Submissions, result)
	}

	TaskSubmissionReceived(&t)
}

func ReceiveTaskAccept(packet *SpinePacket) {

	created, err := time.Parse(time.RFC3339, packet.Body.Items["task-accept.Created"])

	if err != nil {
		fmt.Println("invalid task bid time received")
		return
	}

	fee, err := strconv.ParseFloat(packet.Body.Items["task-accept.Fee"], 64)
	if err != nil {
		fmt.Println("Invalid task bid received")
		return
	}

	var t TaskAccept
	t.ID = packet.Body.Items["task-accept.ID"]
	t.TaskID = packet.Body.Items["task-accept.TaskID"]
	t.Created = created
	t.Fee = fee
	t.BidderID = packet.Body.Items["task-accept.BidderID"]
	t.TaskOwnerID = packet.Body.Items["task-accept.TaskOwnerID"]
	t.ArrivalRoute = packet.PastRoute.Nodes

	TaskAcceptanceReceived(&t)
}

func ReceivePeersRequest(packet *SpinePacket) {
	print("Received peers request")
}
