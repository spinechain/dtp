package ui

import (
	util "spinedtp/util"
	"strconv"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

type PanelCommand struct {
	cmdLabel         *gtk.Label
	commandTextField *gtk.Entry
	btn              *gtk.Button
	historyLabel     *gtk.Label
	commandBox       *gtk.Box
	resultGrid       *gtk.Grid
	panelFrames      []*gtk.Frame
	panelItems       []*gtk.Widget
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

	command.commandTextField.SetText("ping 8.8.8.8")

	// add the grid
	command.resultGrid, err = gtk.GridNew()
	if err != nil {
		util.PrintRed(err.Error())
		return nil, err
	}

	command.resultGrid.SetColumnHomogeneous(true)
	command.resultGrid.SetRowHomogeneous(true)

	// Add cell padding
	command.resultGrid.SetRowSpacing(20)
	command.resultGrid.SetColumnSpacing(20)

	command.commandBox.PackStart(command.resultGrid, false, false, padding)

	//	Create all the images
	for i := 0; i < 9; i++ {
		frm, err := gtk.FrameNew(strconv.Itoa(i))

		// Set the minimum height of the frame
		frm.SetSizeRequest(100, 200)

		if err != nil {
			util.PrintRed(err.Error())
			return nil, err
		}

		command.panelFrames = append(command.panelFrames, frm)
	}

	for i, pitem := range command.panelFrames {

		// get the row and column
		row := i / 3
		col := i % 3

		command.resultGrid.Attach(pitem, col, row, 1, 1)
	}

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

func (command *PanelCommand) AddResult(task string, mimeType string, data []byte) {

	if mimeType == "image/png" || mimeType == "image/jpeg" {
		// load the pixbuf from the data using the loader
		loader, err := gdk.PixbufLoaderNew()
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		// write the data to the loader
		_, err = loader.Write(data)
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		// close the loader
		err = loader.Close()
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		// get the pixbuf
		pixbuf, err := loader.GetPixbuf()
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		img, err := gtk.ImageNewFromPixbuf(pixbuf)
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		// Add to the panel items
		command.panelItems = append([]*gtk.Widget{img.ToWidget()}, command.panelItems...)

	} else if mimeType == "text/plain" {
		// create a label
		label, err := gtk.LabelNew(string(data))
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		label.SetHExpand(false)
		label.SetVExpand(false)
		label.SetMaxWidthChars(50)
		label.SetLineWrap(true)
		label.SetLineWrapMode(pango.WRAP_WORD_CHAR)
		label.SetLines(5)
		label.SetEllipsize(pango.ELLIPSIZE_MIDDLE)
		label.SetMarginTop(8)
		label.SetMarginBottom(8)
		label.SetMarginStart(8)
		label.SetMarginEnd(8)

		// Add to the panel items
		command.panelItems = append([]*gtk.Widget{label.ToWidget()}, command.panelItems...)

	}

	// delete the last item if more than 9
	if len(command.panelItems) > 9 {
		command.panelItems = command.panelItems[:len(command.panelItems)-1]
	}

	// Loop through the panel items and add them to the frames
	for i, pitem := range command.panelItems {

		// Remove existing children
		curChild, err := command.panelFrames[i].GetChild()
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		if curChild != nil {
			command.panelFrames[i].Remove(curChild)
		}

		// Add the new child
		command.panelFrames[i].Add(pitem)
	}

	command.commandBox.ShowAll()

}
