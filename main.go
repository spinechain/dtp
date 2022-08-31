package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	tasknet "spinedtp/tasknet"
	taskpool "spinedtp/taskpool"
	"spinedtp/ui"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

var tasksAvailable taskpool.Taskpool
var tasksCompleted taskpool.Taskpool

func main() {
	fmt.Println("Starting SpineChain DTP")

	LoadSettings()
	SaveSettings()

	tasksAvailable.Start(filepath.Join(AppSettings.DataFolder, "tasks_available.db"), true)
	// tasksCompleted.Start(filepath.Join(AppSettings.DataFolder, "tasks_done.db"), true)

	tasksAvailable.AddTask("1", "test")
	tasksAvailable.GetAllTasks()

	if AppSettings.ShowUI {
		ui.Create()

		// Set the callback pressed when connect btn is pressed
		ui.OnConnectToNetwork = BuildConnectionToTaskNetwork
		ui.OnSubmitToNetworkButton = SubmitTaskToNetwork

		// Start the windowing thread
		gtk.Main()
	} else {
		fmt.Println("Spine running on the command line")
		fmt.Println("How many I help?")

		go func() {
			time.Sleep(5 * time.Second)

			// command line connect can happen here
			BuildConnectionToTaskNetwork()
		}()

		// Read from the terminal and getting commands
		reader := bufio.NewReader(os.Stdin)

		for {
			text, _ := reader.ReadString('\n')
			if text == "q\n" {
				break
			} else {
				SubmitTaskToNetwork(text)
				// fmt.Println("I don't understand")
			}
		}

	}

	Shutdown()
}

func SubmitTaskToNetwork(taskStr string) {
	tasknet.ExecNetworkCommand(taskStr)
}

func BuildConnectionToTaskNetwork() {

	var n tasknet.NetworkSettings
	var c tasknet.NetworkCallbacks

	n.MyPeerID = AppSettings.ClientID
	n.ServerPort = AppSettings.ServerPort
	n.OnStatusUpdate = nil // s.OnStatusBarUpdateRequest
	n.BidTimeoutSeconds = 5
	n.AcceptedBidsPerTask = 3
	n.TaskReadyForProcessing = TaskReadyForExecution
	n.DataFolder = AppSettings.DataFolder

	c.OnTaskReceived = nil // s.OnNewTaskReceived
	c.OnTaskApproved = nil //s.OnNetworkTaskApproval

	tasknet.Connect(n, c)

}

// This function is called when we have been selected to actually execute this task
// Depending on the task, it would be routed to different places. We can execute it
// it out of band and submit a result once done.
func TaskReadyForExecution(cmd string) {
	fmt.Println("Task ready for execution: " + cmd)

}

func TaskHasCompletedExecution(cmd string) {
	// We call EngineHasCompletedTask
}

func Shutdown() {

	fmt.Println("Shutting down SpineChain...")
	SaveSettings()

	tasksAvailable.Stop()
	tasksCompleted.Stop()

	fmt.Println("Shut down complete.")

	os.Exit(1)
}
