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

var TaskPool *Taskpool
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

	var myTasks []*Task
	var networkTasks []*Task

	// fill the two above from the taskpool db

	SendNewTaskToPeers(myTasks)

	for _, task := range myTasks {
		switch task.Status {
		case BiddingComplete:
			util.PrintPurple("Found a task with bidding period complete")
			SelectWinningBids(task)
			task.Status = BidsSelected
		}
	}

	for _, task := range networkTasks {
		switch task.Status {
		case Received:
			util.PrintYellow("Found a new unprocessed task: " + task.Command)

			task.Status = Bid

			// We now bid for this task. We need to find the route that this
			// packet came in through to respond through the same one
			// TODO: this is likely not needed here
			for _, routePeer := range task.ArrivalRoute {

				peer := FindPeer(routePeer.ID)

				peer.BidForTask(task)
				break
			}

		case WorkComplete:
			util.PrintYellow("Found a completed task. Submitting: " + task.Command)
			for _, routePeer := range task.ArrivalRoute {

				peer := FindPeer(routePeer.ID)

				util.PrintYellow("Submitting task to " + peer.ID)
				peer.SubmitTaskResult(task)

				break
			}
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

			task.Status = WaitingForBids
			// Todo: probably needs to be removed
			// taskPool.tasks = RemoveIndex(taskPool.tasks, i)
			SendPacketToAllPeers(packet)
			task.FullyPropagated = true

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

			if task.Status == AcceptedForWork {
				util.PrintYellow("Executing Task: " + task.Command)
				tt.Submission = []byte("This would be my submission")

				if networkSettings.TaskReadyForProcessing != nil {
					networkSettings.TaskReadyForProcessing(task.Command)
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
