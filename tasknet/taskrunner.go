package tasknet

import (
	"fmt"
	"spinedtp/util"
)

// This file contains the thread that watches the taskpool for changes and makes appropriate actions

var taskForProcessingAvailable chan int
var taskForExecutionAvailable chan int
var shutDownTaskThread bool = false

var OpenTaskPool *Taskpool
var ProcessingThreadRunning bool = false

// This thread waits for new tasks to come into the network
func ProcessTasks() {

	// Start the thread that will do the actual work (execution of each task)
	go ProcessAcceptedTasks()

	for {

		ProcessingThreadRunning = true
		<-taskForProcessingAvailable
		ProcessAvailableTasks()

		if shutDownTaskThread {
			fmt.Println("Task processing Thread Shutdown")
			return
		}
	}

}

// This function is called to make the task processor check for new tasks
func CheckForNewTasks() {
	// Tell the thread to check for new tasks

	if ProcessingThreadRunning {
		taskForProcessingAvailable <- 1
		// taskForExecutionAvailable <- 1
	}

}

// This function looks through all tasks, and based on their status, it decides what
// to do with them
func ProcessAvailableTasks() {

	// Retrieve all open tasks. In future we may want to limit the max tasks retrievable
	// if the taskpool gets too large
	tasks, _ := OpenTaskPool.GetAllTasks()

	// Go through all other tasks and ensure that they are appropriately handled based on their
	// status.
	for _, task := range tasks {

		switch task.GlobalStatus {
		case StatusBiddingComplete:
			util.PrintPurple("Found a task with bidding period complete")
			SelectWinningBids(task)
			//task.GlobalStatus = StatusAcceptedWorkers
			// task.LocalStatus = StatusWaitingForExecution
			OpenTaskPool.UpdateTaskStatus(task, StatusAcceptedWorkers, task.LocalWorkerStatus, StatusWaitingForExecution)

			/*
				case WorkComplete:
					util.PrintYellow("Found a completed task. Submitting: " + task.Command)
					for _, routePeer := range task.ArrivalRoute {

						peer := FindPeer(routePeer.ID)

						util.PrintYellow("Submitting task to " + peer.ID)
						peer.SubmitTaskResult(task)

						break
					}
			*/
		}

		switch task.LocalWorkerStatus {
		// A task comes in that we need to bid for. In this iteration we bid for all tasks, but later
		// we will discriminate a bit
		case StatusNewFromNetwork:
			util.PrintYellow("Found a new unprocessed task: " + task.Command)

			BidForTask(task)

		}

		switch task.LocalWorkProviderStatus {
		case StatusNewFromLocal:
			// Send all tasks that have not been propagated yet to peers.
			SendNewTaskToPeers(tasks)
		}

	}

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
		RoutePacketOn()
	}
}

func BidForTask(task *Task) {

	util.PrintBlue("Bidding for Task: " + task.ID + " (" + task.Command + ")")
	OpenTaskPool.UpdateTaskStatus(task, task.GlobalStatus, StatusSentBid, task.LocalWorkProviderStatus)

	task_bid := CreateTaskBid(task)

	AddBid(taskDb, task_bid, true)

	foundPeer := false
	for _, peer := range Peers {
		if peer.ID == task.TaskOwnerID {
			peer.BidForTask(task, task_bid)
			foundPeer = true
			break
		}
	}

	if !foundPeer {
		for _, peer := range Peers {
			peer.BidForTask(task, task_bid)
		}
	}
}

func SendNewTaskToPeers(myTasks []*Task) {

	for _, task := range myTasks {

		if !task.FullyPropagated {
			packet, err := ConstructTaskPropagationPacket(task)
			if err != nil {
				continue
			}

			task.MarkAsPropagated(OpenTaskPool)
			task.GlobalStatus = StatusWaitingForBids

			OpenTaskPool.UpdateTaskStatus(task, task.GlobalStatus, task.LocalWorkerStatus, StatusWaitingForBidsForMe)

			SendPacketToAllPeers(packet)

			// Set a timeout
			go WaitForBidExpiry(task)

			break
		}
	}
}

func ProcessAcceptedTasks() {

	for {

		<-taskForExecutionAvailable

		// retrieve the tasks from the taskpool
		acceptedTasks, err := OpenTaskPool.GetTasks("where local_worker_status=?", StatusApprovedForMe)
		if err != nil {
			continue
		}

		for _, task := range acceptedTasks {

			// We confirm again that we actually bid for this task
			// yes, we checked this before, but we need sanity checks
			bids, err := GetBids("where task_id=? and bidder_id=? and selected=? and my_bid=17", task.ID, NetworkSettings.MyPeerID, 1)
			if err != nil {
				util.PrintRed("☢️ Found a task for me, but we never bid on this 😨")
				continue
			}

			if len(bids) != 1 {
				// It would only be greater than 1 if there is a bug. Better we know
				util.PrintRed("🐛 It looks like we bid more than once on task. How can??? 🙆‍♂️")
				continue
			}

			data := []byte("This would be my submission")

			SendTaskSubmission(task, &data)

		}

		if shutDownTaskThread {
			util.PrintYellow("Task execution Thread Shutdown")
			return
		}
	}

}

func ShutDownTaskRunner() {

	shutDownTaskThread = true
	taskForProcessingAvailable <- 1
}
