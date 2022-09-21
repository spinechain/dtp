package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	tasknet "spinedtp/tasknet"
	"spinedtp/tasktypes"
	"spinedtp/taskworkers"
	"spinedtp/ui"
	"spinedtp/util"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/default_peers.txt
var default_peers string

var tasksAvailable tasknet.Taskpool
var tasksCompleted tasknet.Taskpool
var taskWorkers taskworkers.TaskWorkers
var TaskForSubmissionAvailable chan int

func main() {
	fmt.Println("Starting SpineChain DTP")

	settingsFile := LoadSettings()

	tasknet.OpenTaskPool = &tasksAvailable

	// Create shared channel for task execution to talk to the network
	TaskForSubmissionAvailable = make(chan int)
	tasknet.TaskForSubmissionAvailable = TaskForSubmissionAvailable
	tasktypes.TaskForSubmissionAvailable = TaskForSubmissionAvailable

	SetNetworkSettings()

	if AppSettings.ShowUI {

		// Set the callback pressed when connect btn is pressed
		ui.OnConnectToNetwork = Event_BuildConnectionToTaskNetwork
		ui.OnSubmitToNetworkButton = Event_SubmitTaskToNetwork
		ui.OnClearTasksDb = Event_ClearTasksDB

		err := ui.Create()
		if err != nil {
			SaveSettings()
			util.PrintRed("Running in UI mode, but could not create UI. Please change ShowUI to false in settings file: " + settingsFile)
			return
		}

		glib.TimeoutAdd(250, func() bool {

			go Start()

			return false
		})

		// Start the windowing thread
		gtk.Main()
	} else {
		fmt.Println("Spine running on the command line")
		fmt.Println("How may I help?")

		util.PrintLocalIPAddresses()

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

	tasknet.LoadPeerTable()

	UpdateInfoStatusBar()

	// Add this node as a peer. Will not be needed in future. Good for testing
	// To add local peer, put it in the default_peers.txt file

	// var peer tasknet.Peer
	// peer.Address = AppSettings.ListenAddress
	// peer.ConnectPort = int(AppSettings.ServerPort)
	// peer.ID = AppSettings.ClientID
	// tasknet.AddToPeerTable(&peer)

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

	tasknet.SavePeerTable()

	fmt.Println("Shut down complete.")

	os.Exit(1)
}
