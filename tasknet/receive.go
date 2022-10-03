package tasknet

import (
	"bytes"
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
	case "task-bid-approval":
		ReceiveTaskBidApproval(packet)
	case "task-submission":
		ReceiveTaskSubmission(packet)
	case "task-completed":
		util.PrintRed("Task Completed Received - NOT IMPLEMENTED YET")
		// ReceiveTaskCompleted(packet)
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
	// send it back to the person.
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

	// OpenTaskPool.AddToTaskPool(&task)

	// We check if we have this task already
	tasks, err := OpenTaskPool.GetTasks("where tid=?", task.ID)

	if tasks == nil || err != nil || len(tasks) == 0 {
		// we do not have this task in our db. We can add it directly
		OpenTaskPool.AddTask(&task)
	} else {
		// We have this task already
		existingTask := tasks[0]
		fmt.Print("Existing task: " + existingTask.ID + " " + existingTask.Command)

		if existingTask.LocalWorkerStatus == StatusNewFromLocal {
			OpenTaskPool.UpdateTaskStatus(existingTask, task.GlobalStatus, StatusNewFromNetwork, task.LocalWorkProviderStatus)
		}
	}

	OpenTaskPool.IncHighestIndex(task.Index)

	// This changes the thread and informs the UI about this new task
	glib.TimeoutAdd(10, func() bool {
		if NetworkCallbacks.OnTaskStatusUpdate != nil {
			NetworkCallbacks.OnTaskStatusUpdate()
		}
		return false
	})

	// This will tell the network to start processing the task
	taskForProcessingAvailable <- 1
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

	util.PrintPurple("New Task Bid Arrived for Task " + t.TaskID + " from " + PeerIDToDescription(t.BidderID))

	if t.TaskOwnerID == NetworkSettings.MyPeerID {
		// This is a bid for a task of mine

		ProcessBidForMyTask(taskDb, &t)

	} else {
		// This is a bid for another peer that is not me. We route
		// it to the best connection we have
		util.PrintPurple("Task bid for another client: " + t.BidderID)
		RouteTaskBidOn(&t)
	}
}

func ReceiveTaskBidApproval(packet *SpinePacket) {

	created, err := time.Parse(time.RFC3339, packet.Body.Items["task-bid-approval.Created"])

	if err != nil {
		fmt.Println("invalid task bid time received")
		return
	}

	fee, err := strconv.ParseFloat(packet.Body.Items["task-bid-approval.Fee"], 64)
	if err != nil {
		fmt.Println("Invalid task bid received")
		return
	}

	value, err := strconv.ParseFloat(packet.Body.Items["task-bid-approval.Value"], 64)
	if err != nil {
		fmt.Println("invalid bid value received")
		return
	}

	var t TaskBidApproval
	t.ID = packet.Body.Items["task-bid-approval.ID"]
	t.TaskID = packet.Body.Items["task-bid-approval.TaskID"]
	t.Created = created
	t.Fee = fee
	t.Value = value
	t.TaskOwnerID = packet.Body.Items["task-bid-approval.TaskOwnerID"]
	t.BidderID = packet.Body.Items["task-bid-approval.BidderID"]
	t.Geo = packet.Body.Items["task-bid-approval.Geo"]
	t.ArrivalRoute = packet.PastRoute.Nodes

	if NetworkCallbacks.OnTaskStatusUpdate != nil {
		NetworkCallbacks.OnTaskStatusUpdate()
	}

	if t.BidderID == NetworkSettings.MyPeerID {

		util.PrintYellow("Received new task bid approval from: " + PeerIDToDescription(t.TaskOwnerID) + " for task: " + t.TaskID)

		// We need to find the task in our taskpool. If it's not there, we should
		// not do it
		task := OpenTaskPool.GetTask(t.TaskID)
		if task == nil {
			util.PrintRed("Invalid task found!")
		}

		// Let's check if we bid on it
		ourBid, err := GetMyBids("where task_id=?", t.TaskID)
		if err == nil && len(ourBid) > 0 {
			// In this case, we really did bid for this
			// TODO: check that if someone sends us a bid telling us that it is our ID, that we do not
			// accept it
			OpenTaskPool.UpdateTaskStatus(task, task.GlobalStatus, StatusApprovedForMe, task.LocalWorkProviderStatus)
			ourBid[0].MarkAsAccepted()
			taskForExecutionAvailable <- 1

		} else {
			util.PrintRed("Did not find any bid for this task")
		}

	} else {

		util.PrintYellow("Task bid approval for another client received: " + PeerIDToDescription(t.BidderID))
		RouteTaskBidApprovalOn(&t)
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

	if t.TaskOwnerID == NetworkSettings.MyPeerID {

		// Loop over all submissions
		for i := 0; i < len(t.Submissions); i++ {

			// Get a submission
			submission := t.Submissions[i]

			// Convert byte to hex
			peek_val := bytes.NewBuffer(submission.data[:20]).String()

			util.PrintPurple("Received task submission with length: " + fmt.Sprint(len(submission.data)) + ", Peek Data: " + util.Red + peek_val + util.Reset)

			// Get the task
			task := OpenTaskPool.GetTask(t.TaskID)

			NetworkCallbacks.OnTaskResult(task, submission.mimeType, submission.data)

		}

	} else {
		// Forward the task submission to the correct client
		RouteTaskSubmissionOn(&t)
	}

}

func ReceivePeersRequest(packet *SpinePacket) {
	print("Received peers request")
}
