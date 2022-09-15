package ui

import (
	"io/ioutil"
	util "spinedtp/util"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type PanelCommand struct {
	cmdLabel         *gtk.Label
	commandTextField *gtk.Entry
	btn              *gtk.Button
	historyLabel     *gtk.Label
	commandBox       *gtk.Box
	resultGrid       *gtk.Grid
	images           []gtk.Image
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

	// add the grid
	command.resultGrid, err = gtk.GridNew()
	if err != nil {
		util.PrintRed(err.Error())
		return nil, err
	}

	command.resultGrid.SetColumnHomogeneous(true)
	command.resultGrid.SetRowHomogeneous(true)

	command.commandBox.PackStart(command.resultGrid, false, false, padding)

	// init the images
	// command.images = make([]*gtk.Image, 1)

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

	// load the image to pixbuf
	pixbuf, err := gdk.PixbufNewFromFile("task_submissio.jpeg")
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	img, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	// append image to the list, limit to 6 items
	command.images = append(command.images, *img)
	if len(command.images) > 6 {
		command.images = command.images[1:]
	}

	// Set each image to a grid cell, limit is 6x6
	for i, img2 := range command.images {
		command.resultGrid.Attach(&img2, i%3, i/3, 1, 1)
	}

	//for i, img2 := range command.images {
	//	command.resultGrid.Attach(img2, i, 0, 1, 1)
	//}

	// add the image to the grid
	//command.resultGrid.Attach(command.images[0], 0, 0, 1, 1)

	// add the image to the panel
	//command.commandBox.PackStart(command.images[0], false, false, 5)
	command.commandBox.ShowAll()

}
