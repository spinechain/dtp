package tasknet

import (
	"bufio"
	"fmt"
	"os"
	"spinedtp/tasktypes"
	"spinedtp/util"
	"strings"
)

// This file contains the thread that watches the taskpool for changes and makes appropriate actions

var taskForProcessingAvailable chan int
var taskForExecutionAvailable chan int
var TaskForSubmissionAvailable chan int
var shutDownTaskThread bool = false
var shutDownSubmissionThread bool = false

var OpenTaskPool *Taskpool
var ProcessingThreadRunning bool = false
var SubmissionThreadRunning bool = false

var TasksToExecute *[]tasktypes.TaskToExecute

// This thread waits for new tasks to come into the network
func ProcessTasks() {

	// Start the thread that will do the actual work (execution of each task)
	go ProcessExecutionTasks()
	go ProcessCompletedTasks()
	go ProcessAvailableTasks()

}

// Triggered when a task has arrived from network or local
func ProcessAvailableTasks() {

	ProcessingThreadRunning = true

	for {

		<-taskForProcessingAvailable

		// Retrieve all open tasks. In future we may want to limit the max tasks retrievable
		// if the taskpool gets too large
		tasks, _ := OpenTaskPool.GetAllTasks()

		// Go through all other tasks and ensure that they are appropriately handled based on their
		// status.
		for _, task := range tasks {

			switch task.GlobalStatus {
			case StatusBiddingComplete:

				if !NetworkSettings.RouteOnly {
					util.PrintPurple("Found a task with bidding period complete")

					SelectWinningBids(task)

					OpenTaskPool.UpdateTaskStatus(task, StatusAcceptedWorkers, task.LocalWorkerStatus, StatusWaitingForExecution)
				}

			}

			switch task.LocalWorkerStatus {
			// A task comes in that we need to bid for. In this iteration we bid for all tasks, but later
			// we will discriminate a bit
			case StatusNewFromNetwork:
				util.PrintYellow("Found a new unprocessed task: " + task.Command)

				if !NetworkSettings.RouteOnly {
					BidForTask(task)
				}

				// We check here if we should route the task to another peer
				RouteTaskOn(task)
			}

			switch task.LocalWorkProviderStatus {
			case StatusNewFromLocal:

				if !NetworkSettings.RouteOnly {
					// Send all tasks that have not been propagated yet to peers.
					SendNewTaskToPeers(tasks)
				}

			}

		}

		if shutDownTaskThread {
			fmt.Println("Task processing Thread Shutdown")
			return
		}
	}

}

// Triggered when we have been accepted to work on a task
func ProcessExecutionTasks() {

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
			bids, err := GetMyBids("where task_id=? and selected=?", task.ID, 1)
			if err != nil {
				util.PrintRed("☢️ Found a task for me, but we never bid on this 😨")
				continue
			}

			if len(bids) == 0 {
				util.PrintRed("🐛 We received a bid approval on a task we did not bid on. How can??? 🙆‍♂️")
				continue
			}

			if len(bids) > 1 {
				// It would only be greater than 1 if there is a bug. Better we know
				util.PrintRed("🐛 It looks like we bid more than once on task. How can??? 🙆‍♂️")
				continue
			}

			// Change status so we know we are executing this task. If there a failure it does not recover
			OpenTaskPool.UpdateTaskStatus(task, task.GlobalStatus, StatusExecuting, task.LocalWorkProviderStatus)

			// Get first word of the command
			firstWord := util.FirstWords(task.Command, 1)
			if strings.ToLower(firstWord) == "draw" {
				tasktypes.AddToTaskExecutionQueue(NetworkSettings.DataFolder, "sd", task.ID, task.Command)
			} else if strings.ToLower(firstWord) == "ping" {
				tasktypes.AddToTaskExecutionQueue(NetworkSettings.DataFolder, "ping", task.ID, task.Command)
			} else {
				// We just do sd for now
				tasktypes.AddToTaskExecutionQueue(NetworkSettings.DataFolder, "sd", task.ID, task.Command)
			}

		}

		if shutDownTaskThread {
			util.PrintYellow("Task execution Thread Shutdown")
			return
		}
	}

}

func ProcessCompletedTasks() {

	for {

		SubmissionThreadRunning = true
		<-TaskForSubmissionAvailable

		// Loop through all completes in TasksToExecute
		for i, task_exec := range *TasksToExecute {
			// Check if the task is complete
			if task_exec.Complete && !task_exec.Sent && task_exec.ResultError == nil {

				task_exec.Sent = true
				(*TasksToExecute)[i] = task_exec

				var resultFile string
				if len(task_exec.ResultFiles) > 0 {
					resultFile = task_exec.ResultFiles[0]
				} else {
					continue
				}

				file, err := os.Open(resultFile)
				if err != nil {
					util.PrintRed("🐛 Could not open file to be uploaded: " + resultFile)
					continue
				}

				defer file.Close()

				fileInfo, _ := file.Stat()
				var size int64 = fileInfo.Size()
				bin := make([]byte, size)

				// read file into bytes
				buffer := bufio.NewReader(file)
				_, err = buffer.Read(bin)
				if err != nil {
					util.PrintRed("🐛 Could not read file to be uploaded: " + resultFile)
					return
				}

				task := OpenTaskPool.GetTask(task_exec.TaskID)

				// Submit the task to the network
				SendTaskSubmission(task, task_exec.MimeType, &bin)
			}
		}

		if shutDownSubmissionThread {
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

func ShutDownTaskRunner() {

	shutDownTaskThread = true
	shutDownSubmissionThread = true
	taskForProcessingAvailable <- 1
	TaskForSubmissionAvailable <- 1
}
