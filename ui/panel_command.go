package ui

import (
	"fmt"
	"math/rand"
	"os"
	"spinedtp/tasknet"
	util "spinedtp/util"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

type TaskResult struct {
	task       *tasknet.Task
	box        *ResultBox // The box that this result is in
	mimeType   string
	data       []byte
	filePath   string
	spinning   bool
	image      *gtk.Image
	statusText string
}

type PanelCommand struct {
	cmdLabel         *gtk.Label
	commandTextField *gtk.Entry
	btn              *gtk.Button
	historyLabel     *gtk.Label
	commandBox       *gtk.Box
	resultGrid       *gtk.Grid
	resultBoxes      []*ResultBox
	taskResults      []*TaskResult
	maxResultBoxes   int
}

var testUIFlipBit bool = false

func (command *PanelCommand) Create(title string) (*gtk.Box, error) {

	var err error
	command.commandBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	command.maxResultBoxes = 9

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

	//	Create all the result boxes
	for i := 0; i < command.maxResultBoxes; i++ {

		var resultBox ResultBox
		resultBox.Create()

		command.resultBoxes = append(command.resultBoxes, &resultBox)
	}

	// Add all the boxes to the grid
	for i, rbox := range command.resultBoxes {

		// get the row and column
		row := i / 3
		col := i % 3

		command.resultGrid.Attach(rbox.box, col, row, 1, 1)
	}

	return command.commandBox, err
}

func (command *PanelCommand) TestUI() {

	// generate random ID
	val := rand.Intn(100000)

	// val to string
	s := fmt.Sprint(val)

	var t tasknet.Task
	t.ID = s
	t.Command = "draw a picture of a happy spaceship"

	command.PrepareForNewResult(&t)

	// Delay a bit, then load the image
	glib.TimeoutAdd(3250, func() bool {

		// load sample image
		if testUIFlipBit {
			data, _ := os.ReadFile("assets/test1.jpg")
			command.AddResult(&t, "image/jpeg", data)
			testUIFlipBit = false
		} else {
			data, _ := os.ReadFile("assets/test2.jpg")
			command.AddResult(&t, "image/jpeg", data)
			testUIFlipBit = true
		}

		/*
			glib.TimeoutAdd(3250, func() bool {
				var t tasknet.Task
				t.ID = "124"
				t.Command = "another one"

				command.PrepareForNewResult(&t)
				command.UpdateTasks(nil)

				return false
			})
		*/

		return false
	})

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

// Shift the boxes by one to the right so the new result is in the first pos
func (command *PanelCommand) PrepareForNewResult(task *tasknet.Task) {

	var newResult TaskResult
	newResult.task = task
	newResult.spinning = true
	newResult.statusText = "Loading..."

	command.taskResults = append([]*TaskResult{&newResult}, command.taskResults...)

	for i, result := range command.taskResults {

		if i == command.maxResultBoxes {
			break
		}

		result.box = command.resultBoxes[i]

		// make the result box show the main thing
		result.box.AdaptToCircumstances(result)
	}

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
	taskResult.spinning = false
	taskResult.statusText = "Done"

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
		taskResult.image = img

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

		// TODO, put this label in the scroller
		taskResult.statusText = s
	}

	command.UpdateTask(task)

}

func (command *PanelCommand) UpdateTasks(tasks []*tasknet.Task) {

	if tasks == nil {
		// Loop over all results
		for _, result := range command.taskResults {
			command.UpdateTask(result.task)
		}

	} else {
		for _, task := range tasks {
			command.UpdateTask(task)
		}
	}

}

func (command *PanelCommand) UpdateTask(task *tasknet.Task) {

	// Find the task result

	// Loop through the panel items and add them to the frames
	for i, result := range command.taskResults {

		if result.task.ID == task.ID {
			command.taskResults[i].task = task

			// we only show limited results on ui
			if i == command.maxResultBoxes {
				break
			}

			command.taskResults[i].box.AdaptToCircumstances(command.taskResults[i])
		}
	}

	command.commandBox.ShowAll()
}
