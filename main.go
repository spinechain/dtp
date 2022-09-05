package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	tasknet "spinedtp/tasknet"
	"spinedtp/taskworkers"
	"spinedtp/ui"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var tasksAvailable tasknet.Taskpool
var tasksCompleted tasknet.Taskpool
var taskWorkers taskworkers.TaskWorkers

func main() {
	fmt.Println("Starting SpineChain DTP")

	LoadSettings()

	tasknet.OpenTaskPool = &tasksAvailable

	if AppSettings.ShowUI {

		// Set the callback pressed when connect btn is pressed
		ui.OnConnectToNetwork = Event_BuildConnectionToTaskNetwork
		ui.OnSubmitToNetworkButton = Event_SubmitTaskToNetwork
		ui.OnClearTasksDb = Event_ClearTasksDB

		ui.Create()

		glib.TimeoutAdd(250, func() bool {

			go Start()
			return false
		})

		// Start the windowing thread
		gtk.Main()
	} else {
		fmt.Println("Spine running on the command line")
		fmt.Println("How many I help?")

		go Start()

		go func() {
			time.Sleep(5 * time.Second)

			// command line connect can happen here
			Event_BuildConnectionToTaskNetwork()
		}()

		// Read from the terminal and getting commands
		reader := bufio.NewReader(os.Stdin)

		for {
			text, _ := reader.ReadString('\n')
			if text == "q\n" {
				break
			} else {
				Event_SubmitTaskToNetwork(text)
				// fmt.Println("I don't understand")
			}
		}

	}

	Shutdown()
}

func Start() {
	SaveSettings()

	tasksAvailable.Start(filepath.Join(AppSettings.DataFolder, "tasks_available.db"), false)
	tasksAvailable.OnTaskAdded = Event_TaskAdded

	tasksCompleted.Start(filepath.Join(AppSettings.DataFolder, "tasks_done.db"), false)
	tasksCompleted.OnTaskAdded = Event_TaskAdded

	// load the existing items to listview
	Event_TaskAdded("-1", "")

	taskWorkers.Start(filepath.Join(AppSettings.DataFolder, "tasks_workers.db"), false)
	taskWorkers.OnTaskWorkersAdded = Event_TaskWorkerAdded
	taskWorkers.AddTaskWorker(GetMeAsTaskWorker())

}

// Return the info needed to make me a taskworker
func GetMeAsTaskWorker() *taskworkers.TaskWorker {
	var mtw taskworkers.TaskWorker
	mtw.Address = "127.0.0.1"
	mtw.Port = int(AppSettings.ServerPort)
	mtw.ID = AppSettings.ClientID
	return &mtw
}

func Shutdown() {

	fmt.Println("Shutting down SpineChain...")
	SaveSettings()

	// Close list of tasks
	tasksAvailable.Stop()
	tasksCompleted.Stop()

	// Close list of workers that execute tasks
	taskWorkers.Stop()

	fmt.Println("Shut down complete.")

	os.Exit(1)
}
