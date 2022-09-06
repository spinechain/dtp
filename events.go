package main

import (
	"fmt"
	tasknet "spinedtp/tasknet"
	"spinedtp/ui"

	"github.com/gotk3/gotk3/glib"
)

func Event_SubmitTaskToNetwork(taskStr string) {

	SetNetworkSettings()

	tasknet.SendTaskToNetwork(taskStr)
	// tasknet.ExecNetworkCommand(taskStr)
}

func Event_BuildConnectionToTaskNetwork() {

	// TODO: Check what happens if this is called twice

	SetNetworkSettings()
	tasknet.Connect()

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

		tasks, _ := tasksAvailable.GetAllTasks()

		ui.TasksPanel.UpdateList(tasks)
		return false
	})

	tasknet.CheckForNewTasks()
}

func Event_StatusUpdate(txt string, section int) {
	// we switch thread to ui context
	glib.TimeoutAdd(10, func() bool {

		ui.UpdateStatusBar(txt, section)
		return false
	})
}

func Event_TaskWorkerAdded(tid string) {

	// we switch thread to ui context
	glib.TimeoutAdd(10, func() bool {

		taskWorkers, _ := taskWorkers.GetAllTaskWorkers()

		ui.NetworkPanel.UpdateList(taskWorkers)
		return false
	})
}

func Event_ClearTasksDB() {

	tasksAvailable.RemoveAllTasks()

	// we switch thread to ui context
	glib.TimeoutAdd(10, func() bool {

		tasks, _ := tasksAvailable.GetAllTasks()

		ui.TasksPanel.UpdateList(tasks)
		return false
	})
}
