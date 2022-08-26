package ui

import (
	"fmt"
	"net/url"

	"github.com/gotk3/gotk3/gtk"
)

type FnOnSubmitToNetwork func(string)

var OnSubmitToNetworkButton FnOnSubmitToNetwork
var OnConnectToNetworkButton func()

func onBtnClearDBClick(btn *gtk.ToolButton, item *gtk.ToolButton) {

	// p, _ := btn.GetParentWindow()
	//New message dialog, select dialog
	dialog := gtk.MessageDialogNew(
		mainWindow.Window,    //Specify the parent window
		gtk.DIALOG_MODAL,     //Modal dialog
		gtk.MESSAGE_QUESTION, //Specify the dialog box type
		gtk.BUTTONS_YES_NO,   //Default button
		"This will clear your current db. Are you sure?") //Set content

	dialog.SetTitle("Clear DB?") //Dialog box setting title

	flag := dialog.Run() //Run dialog
	if flag == gtk.RESPONSE_YES {
		fmt.Println("Press yes")
		mainWindow.Window.QueueDraw()
	} else if flag == gtk.RESPONSE_NO {
		fmt.Println("Press no")
	} else {
		fmt.Println("Press the close button")
	}

	dialog.Destroy() //Destroy the dialog
}

func OnTreeviewItemSelected(itemType string, center string, action string) {

	SwitchPanel(itemType)

	mainWindow.Window.QueueDraw()

}

func onBtnSubmitToNetworkClick() {
	if OnSubmitToNetworkButton != nil {
		s, _ := commandPanel.commandTextField.GetText()

		_, err := url.ParseRequestURI(s)
		if err != nil {
			commandPanel.AddToHistory("Invalid URL: " + s)
		} else {
			commandPanel.AddToHistory("Submitted: " + s)
			OnSubmitToNetworkButton(s)
		}

	}

}

func onBtnPlayClick() {

	fmt.Println("Connect Pressed.")

	if OnConnectToNetworkButton != nil {
		OnConnectToNetworkButton()
	}
}
