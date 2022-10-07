package main

import (
	tasknet "spinedtp/tasknet"
	"spinedtp/ui"

	"github.com/gotk3/gotk3/glib"
)

func Event_SubmitTaskToNetwork(taskStr string) *tasknet.Task {

	SetNetworkSettings()

	if taskStr == "test" {
		TestUI()
		return nil
	}

	return tasknet.SendTaskToNetwork(taskStr)
}

func Event_BuildConnectionToTaskNetwork() {

	// TODO: Check what happens if this is called twice

	SetNetworkSettings()
	tasknet.Connect()

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

		ui.WorkersPanel.UpdateList(taskWorkers)
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

func Event_TaskResultReceived(task *tasknet.Task, mimeType string, data []byte) {

	// we switch thread to ui context
	glib.TimeoutAdd(10, func() bool {

		ui.CommandPanel.AddResult(task, mimeType, data)
		return false
	})
}

func Event_TaskStatusUpdate() {

	// we switch thread to ui context
	glib.TimeoutAdd(10, func() bool {

		tasks, _ := tasksAvailable.GetAllTasks()

		ui.TasksPanel.UpdateList(tasks)
		ui.CommandPanel.UpdateTasks(tasks)
		return false
	})
}
