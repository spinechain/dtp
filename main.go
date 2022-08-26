package main

import (
	"bufio"
	"fmt"
	"os"
	tasknet "spinedtp/tasknet"
	"spinedtp/ui"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	fmt.Println("Starting SpineChain DTP")

	LoadSettings()
	SaveSettings()

	if AppSettings.ShowUI {
		ui.Create()

		// Set the callback pressed when connect btn is pressed
		ui.OnConnectToNetwork = BuildConnectionToSpineNetwork
		ui.OnSubmitToNetworkButton = SubmitTaskToNetwork

		// Start the windowing thread
		gtk.Main()
	} else {
		fmt.Println("Spine running on the command line")
		fmt.Println("How many I help?")

		go func() {
			time.Sleep(5 * time.Second)

			// command line connect can happen here
			BuildConnectionToSpineNetwork()
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

func BuildConnectionToSpineNetwork() {

	var n tasknet.NetworkSettings
	var c tasknet.NetworkCallbacks

	n.MyPeerID = AppSettings.ClientID
	n.ServerPort = AppSettings.ServerPort
	n.OnStatusUpdate = nil // s.OnStatusBarUpdateRequest
	n.BidTimeoutSeconds = 5
	n.AcceptedBidsPerTask = 3
	n.TaskReadyForProcessing = nil // s.TaskReadyForProcessing
	n.DataFolder = AppSettings.DataFolder

	c.OnTaskReceived = nil // s.OnNewTaskReceived
	c.OnTaskApproved = nil //s.OnNetworkTaskApproval

	tasknet.Connect(n, c)

}

func Shutdown() {

	fmt.Println("Shutting down SpineChain...")
	SaveSettings()

	fmt.Println("Shut down complete.")

	os.Exit(1)
}
