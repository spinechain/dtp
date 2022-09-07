package tasknet

import (
	"fmt"
	"spinedtp/util"
	"time"

	"github.com/lithammer/shortuuid/v3"
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
		taskForExecutionAvailable <- 1
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
			task.GlobalStatus = StatusBidsSelected
			task.LocalStatus = StatusWaitingForExecution

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

		switch task.LocalStatus {
		// A task comes in that we need to bid for. In this iteration we bid for all tasks, but later
		// we will discriminate a bit
		case StatusNewFromLocal:
			// Send all tasks that have not been propagated yet to peers.
			SendNewTaskToPeers(tasks)

		case StatusNewFromNetwork:
			util.PrintYellow("Found a new unprocessed task: " + task.Command)

			BidForTask(task)

		}

	}

}

// Call when a new task bid arrives. We add it to the database. When our bid
// period expires is when we check for all bids and select the best
func NewTaskBidArrived(tb *TaskBid) {

	if tb.TaskOwnerID == NetworkSettings.MyPeerID {
		// This is a bid for a task of mine

		OpenTaskPool.AddBid(tb)

	} else {
		// This is a bid for another peer that is not me. We route
		// it to the best connection we have

		RoutePacketOn()
	}
}

func BidForTask(task *Task) {

	OpenTaskPool.UpdateTaskStatus(task, task.GlobalStatus, StatusSentBid)

	for _, peer := range Peers {
		peer.BidForTask(task)
		break
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
			task.LocalStatus = StatusWaitingForBidsForMe

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

		var acceptedTasks []*Task
		// retrieve the tasks from the taskpool

		for _, task := range acceptedTasks {

			var tt TaskSubmission
			tt.ID = shortuuid.New()
			tt.Created = time.Now()

			if task.LocalStatus == StatusApprovedForMe {
				util.PrintYellow("Executing Task: " + task.Command)
				tt.Submission = []byte("This would be my submission")

				if NetworkSettings.TaskReadyForProcessing != nil {
					NetworkSettings.TaskReadyForProcessing(task.Command)
				} else {
					fmt.Println("No callback available for task processing")
				}
			}

		}

		if shutDownTaskThread {
			util.PrintYellow("Task execution Thread Shutdown")
			return
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
