package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type FnOnSubmitToNetwork func(string)

var OnSubmitToNetworkButton FnOnSubmitToNetwork
var OnConnectToNetwork func()
var OnClearTasksDb func()

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
		s, _ := CommandPanel.commandTextField.GetText()

		CommandPanel.AddToHistory("Submitted: " + s)
		OnSubmitToNetworkButton(s)

	}

}

func onBtnPlayClick() {

	fmt.Println("Connect Pressed.")

	if OnConnectToNetwork != nil {
		OnConnectToNetwork()
	}
}

func onBtnClearTasksClick() {
	dialog := gtk.MessageDialogNew(
		mainWindow.Window,                          //Specify the parent window
		gtk.DIALOG_MODAL,                           //Modal dialog
		gtk.MESSAGE_QUESTION,                       //Specify the dialog box type
		gtk.BUTTONS_YES_NO,                         //Default button
		"This will clear all tasks. Are you sure?") //Set content

	dialog.SetTitle("Clear Tasks?") //Dialog box setting title

	flag := dialog.Run() //Run dialog
	if flag == gtk.RESPONSE_YES {
		fmt.Println("Press yes")
		mainWindow.Window.QueueDraw()

		if OnClearTasksDb != nil {
			OnClearTasksDb()
		}
	} else if flag == gtk.RESPONSE_NO {
		fmt.Println("Press no")
	} else {
		fmt.Println("Press the close button")
	}

	dialog.Destroy() //Destroy the dialog

}
