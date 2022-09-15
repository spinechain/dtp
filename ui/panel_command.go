package ui

import (
	"io/ioutil"
	util "spinedtp/util"

	"github.com/gotk3/gotk3/gtk"
)

type PanelCommand struct {
	cmdLabel         *gtk.Label
	commandTextField *gtk.Entry
	btn              *gtk.Button
	historyLabel     *gtk.Label
	commandBox       *gtk.Box
	images           []*gtk.Image
}

func (command *PanelCommand) Create(title string) (*gtk.Box, error) {

	var err error
	command.commandBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	command.cmdLabel, _ = gtk.LabelNew("Enter the command you would like to send to the network:")
	command.commandTextField, _ = gtk.EntryNew()
	command.btn, _ = gtk.ButtonNew()
	command.historyLabel, _ = gtk.LabelNew("")

	command.btn.SetLabel("Send")
	command.btn.Connect("clicked", onBtnSubmitToNetworkClick)

	var padding uint = 5
	command.commandBox.PackStart(command.cmdLabel, false, false, padding)
	command.commandBox.PackStart(command.commandTextField, false, false, padding)
	command.commandBox.PackStart(command.btn, false, false, padding)
	command.commandBox.PackStart(command.historyLabel, false, false, padding*3)

	command.commandBox.SetMarginStart(20)
	command.commandBox.SetMarginEnd(20)

	command.commandTextField.SetText("draw a picture of a happy spaceship")

	// init the images
	command.images = make([]*gtk.Image, 1)

	return command.commandBox, err
}

func (command *PanelCommand) Destroy() {

}

func (command *PanelCommand) AddToHistory(s string) {
	cs, _ := command.historyLabel.GetText()

	command.historyLabel.SetText(s + "\n" + cs)

}

func (command *PanelCommand) Show() {
	command.cmdLabel.ShowAll()
	command.commandTextField.ShowAll()
	command.btn.ShowAll()
	command.historyLabel.ShowAll()
	command.commandBox.ShowAll()
}

func (command *PanelCommand) Hide() {
	command.cmdLabel.Hide()
	command.commandTextField.Hide()
	command.btn.Hide()
	command.historyLabel.Hide()
	command.commandBox.Hide()
}

func (command *PanelCommand) AddResult(task string, key string, data []byte) {

	// Write the tt.Submission to disk
	err := ioutil.WriteFile("task_submissio.jpeg", data, 0644)
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	command.images[0], err = gtk.ImageNewFromFile("task_submissio.jpeg")
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	// add the image to the panel
	command.commandBox.PackStart(command.images[0], false, false, 5)
	command.commandBox.ShowAll()

}
