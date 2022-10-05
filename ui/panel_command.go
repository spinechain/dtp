package ui

import (
	"spinedtp/tasknet"
	util "spinedtp/util"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

type TaskResult struct {
	task     *tasknet.Task
	widget   *gtk.Widget
	mimeType string
	data     []byte
	filePath string
}

type PanelCommand struct {
	cmdLabel         *gtk.Label
	commandTextField *gtk.Entry
	btn              *gtk.Button
	historyLabel     *gtk.Label
	commandBox       *gtk.Box
	resultGrid       *gtk.Grid
	panelBoxes       []*gtk.Box
	taskResults      []*TaskResult
	Spinning         bool
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

	// Add cell padding
	command.resultGrid.SetRowSpacing(20)
	command.resultGrid.SetColumnSpacing(20)

	command.commandBox.PackStart(command.resultGrid, false, false, padding)

	//	Create all the images
	for i := 0; i < 9; i++ {

		// Create a new box
		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		if err != nil {
			util.PrintRed(err.Error())
			return nil, err
		}

		// Create new scrolled window
		frm, _ := command.MakeScrolledWindow()

		// Add the frame to the box
		box.PackStart(frm, true, true, 0)

		command.panelBoxes = append(command.panelBoxes, box)
	}

	for i, pitem := range command.panelBoxes {

		// get the row and column
		row := i / 3
		col := i % 3

		command.resultGrid.Attach(pitem, col, row, 1, 1)
	}

	return command.commandBox, err
}

func (command *PanelCommand) MakeScrolledWindow() (*gtk.ScrolledWindow, error) {
	frm, err := gtk.ScrolledWindowNew(nil, nil)

	// Add frame to the scrolled window
	frm.SetShadowType(gtk.SHADOW_ETCHED_IN)

	// Set the minimum height of the frame
	frm.SetSizeRequest(100, 200)

	if err != nil {
		util.PrintRed(err.Error())
		return nil, err
	}

	return frm, nil
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

func (command *PanelCommand) PrepareForNewResult(task *tasknet.Task) {

	var newResult TaskResult
	newResult.task = task

	// Create a new spinner
	spinner, err := gtk.SpinnerNew()
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	newResult.widget = spinner.ToWidget()

	command.taskResults = append([]*TaskResult{&newResult}, command.taskResults...)
	command.Spinning = true

	// Start the spinner
	spinner.Start()
	command.UpdateTask(task)
}

func (command *PanelCommand) AddResult(task *tasknet.Task, mimeType string, data []byte) {

	var taskResult *TaskResult

	// Loop over all task result and identify the one with this task
	for _, result := range command.taskResults {

		if result.task.ID == task.ID {
			taskResult = result
			break
		}
	}

	if taskResult == nil {
		util.PrintRed("Could not find task result for task " + task.ID)
		return
	}

	taskResult.mimeType = mimeType
	taskResult.data = data

	/*
		if command.Spinning {
			// delete the first item in command.panelItems
			command.panelItems = command.panelItems[1:]
			command.Spinning = false
		}
	*/

	if mimeType == "image/png" || mimeType == "image/jpeg" {
		// load the pixbuf from the data using the loader
		loader, err := gdk.PixbufLoaderNew()
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		// write the data to the loader
		_, err = loader.Write(taskResult.data)
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

		// Scale the image
		width := 250
		height := 250

		// maintain aspect ratio
		if pixbuf.GetWidth() > pixbuf.GetHeight() {
			height = int(float64(pixbuf.GetHeight()) / float64(pixbuf.GetWidth()) * float64(width))
		} else {
			width = int(float64(pixbuf.GetWidth()) / float64(pixbuf.GetHeight()) * float64(height))
		}

		pixbuf, _ = pixbuf.ScaleSimple(width, height, gdk.INTERP_BILINEAR)

		img, err := gtk.ImageNewFromPixbuf(pixbuf)
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		// Add to the panel items
		taskResult.widget = img.ToWidget()

	} else if mimeType == "text/plain" {

		// convert to asccii
		// data = bytes.Replace(data, []byte{0x0}, []byte{0x20}, -1)

		s, _ := util.ToAscii(string(taskResult.data))

		// create a label
		label, err := gtk.LabelNew(s)
		if err != nil {
			util.PrintRed(err.Error())
			return
		}

		label.SetHExpand(false)
		label.SetVExpand(false)
		label.SetMaxWidthChars(50)
		// label.SetLineWrap(true)
		label.SetLineWrapMode(pango.WRAP_WORD_CHAR)
		label.SetLines(5)
		label.SetEllipsize(pango.ELLIPSIZE_MIDDLE)
		label.SetMarginTop(8)
		label.SetMarginBottom(8)
		label.SetMarginStart(8)
		label.SetMarginEnd(8)

		taskResult.widget = label.ToWidget()
	}

	command.UpdateTask(task)

}

func (command *PanelCommand) UpdateTasks(tasks []*tasknet.Task) {

	for _, task := range tasks {
		command.UpdateTask(task)
	}

}

func (command *PanelCommand) UpdateTask(task *tasknet.Task) {

	// Find the task result

	// Loop through the panel items and add them to the frames
	for i, result := range command.taskResults {

		if result.task.ID == task.ID {
			command.taskResults[i].task = task

			// we only show 9 results for now
			if i == 9 {
				break
			}

			// remove all children of the box
			command.panelBoxes[i].GetChildren().Foreach(func(child interface{}) {

				item := child.(gtk.IWidget)
				// item.GetChild()

				command.panelBoxes[i].Remove(item)
			})

			/*
				// Remove existing children
				curChild, err := command.panelFrames[i].GetChild()
				if err != nil {
					util.PrintRed(err.Error())
					return
				}

				if curChild != nil {
					command.panelFrames[i].Remove(curChild)
				}
			*/

			// This section is for full failure
			if result.task.LocalWorkProviderStatus == tasknet.StatusTimeout {

				// Write timeout on the widget
				label, err := gtk.LabelNew("Timeout")
				if err != nil {
					util.PrintRed(err.Error())
					return
				}

				label.SetHExpand(false)
				label.SetVExpand(false)
				label.SetMarginTop(8)
				label.SetMarginBottom(8)
				label.SetMarginStart(8)
				label.SetMarginEnd(8)

				command.panelBoxes[i].Add(label.ToWidget())

			} else {

				// Create new scrolled window
				frm, _ := command.MakeScrolledWindow()
				frm.Add(result.widget)

				// Add the new child
				command.panelBoxes[i].PackStart(frm, true, true, 0)

				if result.task.LocalWorkProviderStatus == tasknet.StatusWaitingForExecution {
					label, err := gtk.LabelNew("Executing...")
					if err != nil {
						util.PrintRed(err.Error())
						return
					}
					label.SetHExpand(false)
					label.SetVExpand(false)
					label.SetMarginTop(8)
					label.SetMarginBottom(8)
					label.SetMarginStart(8)
					label.SetMarginEnd(8)

					command.panelBoxes[i].Add(label)
				}
			}

			break
		}
	}

	command.commandBox.ShowAll()
}
