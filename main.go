package main

import (
	"bufio"
	"fmt"
	"os"
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

		// Start the windowing thread
		gtk.Main()
	} else {
		fmt.Println("Spine running on the command line")
		fmt.Println("How many I help?")

		go func() {
			time.Sleep(5 * time.Second)

			// command line connect can happen here
		}()

		// Read from the terminal and getting commands
		reader := bufio.NewReader(os.Stdin)

		for {
			text, _ := reader.ReadString('\n')
			if text == "q\n" {
				break
			} else {
				fmt.Println("I don't understand")
			}
		}

	}

	Shutdown()
}

func Shutdown() {

	fmt.Println("Shutting down SpineChain...")
	SaveSettings()

	fmt.Println("Shut down complete.")

	os.Exit(1)
}
