package main

import (
	"fmt"
	tasknet "spinedtp/tasknet"
	"spinedtp/ui"

	"github.com/gotk3/gotk3/glib"
)

func Event_SubmitTaskToNetwork(taskStr string) {

	tasksAvailable.AddMyTask(taskStr)
	// tasknet.ExecNetworkCommand(taskStr)
}

func Event_BuildConnectionToTaskNetwork() {

	var n tasknet.NetworkSettings
	var c tasknet.NetworkCallbacks

	n.MyPeerID = AppSettings.ClientID
	n.ServerPort = AppSettings.ServerPort
	n.OnStatusUpdate = nil // s.OnStatusBarUpdateRequest
	n.BidTimeoutSeconds = 5
	n.AcceptedBidsPerTask = 3
	n.TaskReadyForProcessing = Event_TaskReadyForExecution
	n.DataFolder = AppSettings.DataFolder

	c.OnTaskReceived = nil // s.OnNewTaskReceived
	c.OnTaskApproved = nil //s.OnNetworkTaskApproval

	tasknet.Connect(n, c)

}

// This function is called when we have been selected to actually execute this task
// Depending on the task, it would be routed to different places. We can execute it
// it out of band and submit a result once done.
func Event_TaskReadyForExecution(cmd string) {
	fmt.Println("Task ready for execution: " + cmd)

}

func Event_TaskHasCompletedExecution(cmd string) {
	// We call EngineHasCompletedTask
}

// This is called anytime the tasks are updated
func Event_TaskAdded(tid string, taskText string) {

	// we switch thread to ui context
	glib.TimeoutAdd(10, func() bool {

		//if networkSettings.OnStatusUpdate != nil {
		//	networkSettings.OnStatusUpdate(str, section)
		//}

		tasks, _ := tasksAvailable.GetAllTasks()

		var items []string
		for i := 0; i < len(tasks); i++ {
			items = append(items, tasks[0].Command)
		}

		ui.TasksPanel.UpdateList(items)
		return false
	})

}
