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

	if tb.TaskOwnerID == NetworkSettings.MyPeerID {
		// This is a bid for a task of mine

		AddBid(taskDb, tb, false)

	} else {
		// This is a bid for another peer that is not me. We route
		// it to the best connection we have
		util.PrintPurple("Task bid for another client")
		RoutePacketOn()
	}
}

func BidForTask(task *Task) {

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

			SendTaskResult(task, &data)

		}

		if shutDownTaskThread {
			util.PrintYellow("Task execution Thread Shutdown")
			return
		}
	}

}

// This function will need to be improved a lot. This is because the submissions can be quite large.
// Sending the result through all peers and potentially over the entire network is not going to be good
// Solutions:
//  1. The acceptance packet should contain all peers that the task-giver is connected to.
//  2. The acceptance packet naturally contains the route over which it came
//  3. The worker sends the result to the route it arrived from. It also provides the connection list
//     The next peer checks if it can reach the task giver. If not, it gives up and informs worker.
//  4. In this situation, the worker sends to all the connected peers. If no peer has any of the connected
//     ones connected to it, then they give up. This way we need maximum of two hops between worker and
//     the task giver. This is better for network fairness, so everyone has a chance to get jobs.
//  5. If the above method still congests the network, we will use 'judges' who will provide IP routes. They
//     can also do the transfer for a fee. Generally, there should be a fee for those transferring results.
func SendTaskResult(task *Task, submissionData *[]byte) {

	// Get the target client ID
	targetID := task.TaskOwnerID

	// If we are connected directly to the target client, then send it directly
	// without sending to any other clients
	foundTarget := false
	for _, peer := range Peers {
		if peer.ID == targetID {
			task.Result = *submissionData
			peer.SubmitTaskResult(task)
			foundTarget = true
			break
		}
	}

	// Otherwise send to all connected clients
	if !foundTarget {
		for _, peer := range Peers {
			task.Result = *submissionData
			peer.SubmitTaskResult(task)
		}
	}

}

// The input in here does not know strictly about who triggered it. So it just needs
// to search for what was requested so it can respond
func EngineHasCompletedTask(taskType string, taskCommand string, taskData string) {

	/*
		// var t network.TaskSubmission
		// network.SubmitTaskResult(&t)
		if taskType == "download" {
			for _, task := range taskPool.networkTasks {
				if task.Status == AcceptedForWork {
					if task.Command == taskType+" "+taskCommand {
						// At this point, we have completed a task here internally. We are to
						// propagate it back to the person who requested this task.

						util.PrintPurple("Changing task status to WorkComplete")

						// Set the task result in the task structure here
						task.Result = []byte(taskData)
						task.Status = WorkComplete

						// Trigger the thread to return the result to network
						taskForProcessingAvailable <- 1

					}
				}
			}
		}
	*/
}

func ShutDownTaskRunner() {

	shutDownTaskThread = true
	taskForProcessingAvailable <- 1
}
