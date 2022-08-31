package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	taskpool "spinedtp/taskpool"
	"spinedtp/taskworkers"
	"spinedtp/ui"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

var tasksAvailable taskpool.Taskpool
var tasksCompleted taskpool.Taskpool
var taskWorkers taskworkers.TaskWorkers

func main() {
	fmt.Println("Starting SpineChain DTP")

	LoadSettings()
	SaveSettings()

	tasksAvailable.Start(filepath.Join(AppSettings.DataFolder, "tasks_available.db"), true)
	tasksAvailable.OnTaskAdded = Event_TaskAdded

	tasksCompleted.Start(filepath.Join(AppSettings.DataFolder, "tasks_done.db"), true)
	tasksCompleted.OnTaskAdded = Event_TaskAdded

	tasksAvailable.AddTask("1", "test")
	tasksAvailable.GetAllTasks()

	taskWorkers.Start(filepath.Join(AppSettings.DataFolder, "tasks_workers.db"), true)

	if AppSettings.ShowUI {

		// Set the callback pressed when connect btn is pressed
		ui.OnConnectToNetwork = Event_BuildConnectionToTaskNetwork
		ui.OnSubmitToNetworkButton = Event_SubmitTaskToNetwork

		ui.Create()

		// Start the windowing thread
		gtk.Main()
	} else {
		fmt.Println("Spine running on the command line")
		fmt.Println("How many I help?")

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

func Shutdown() {

	fmt.Println("Shutting down SpineChain...")
	SaveSettings()

	tasksAvailable.Stop()
	tasksCompleted.Stop()

	taskWorkers.Stop()

	fmt.Println("Shut down complete.")

	os.Exit(1)
}
