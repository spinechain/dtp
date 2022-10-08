package ui

import (
	util "spinedtp/util"

	"github.com/gotk3/gotk3/gtk"
)

// One of the square result boxes
type ResultBox struct {
	box      *gtk.Box
	label    *gtk.Label
	scroller *gtk.ScrolledWindow
	spinner  *gtk.Spinner
}

func (b *ResultBox) Create() error {

	// Create a new box
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	b.box = box

	// Create new scrolled window
	frm, _ := b.MakeScrolledWindow()

	b.scroller = frm

	// Add the frame to the box
	box.PackStart(frm, true, true, 0)

	// Create a new label
	label, err := gtk.LabelNew("")
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	b.label = label

	return nil
}

func (b *ResultBox) SwitchToSpinner(taskResult *TaskResult) {

	// Remove the image from the scroller
	curChild, err := b.scroller.GetChild()
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	if curChild != nil {
		b.scroller.Remove(curChild)
	}

	// Create a new spinner
	if b.spinner == nil {
		spinner, _ := gtk.SpinnerNew()
		b.spinner = spinner
	}

	// Add the spinner to the scroller
	b.scroller.Add(b.spinner)

	// Add the label below the scrolled window
	if b.label != nil {
		b.box.Remove(b.label)
		b.label.SetText(taskResult.statusText)
		b.box.PackStart(b.label, false, false, 0)
		b.box.ShowAll()
	}

	// Show everything
	b.scroller.ShowAll()
	b.spinner.ShowAll()

	// Start the spinner
	b.spinner.Start()

}

func (b *ResultBox) MakeScrolledWindow() (*gtk.ScrolledWindow, error) {
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

func (b *ResultBox) SwitchToImage(taskResult *TaskResult) {

	// Remove any child of the scroler
	curChild, err := b.scroller.GetChild()
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	if curChild != nil {
		b.scroller.Remove(curChild)
	}

	// Add the image to the scroller
	b.scroller.Add(taskResult.image)

	// Add the label below the scrolled window
	if b.label != nil {
		b.box.Remove(b.label)
		b.label.SetText(taskResult.statusText)
		b.box.PackStart(b.label, false, false, 0)
		b.box.ShowAll()
	}

}

func (b *ResultBox) AdaptToCircumstances(taskResult *TaskResult) {

	// make scrolled window if nil
	if b.scroller == nil {
		frm, _ := b.MakeScrolledWindow()
		b.scroller = frm

		// add
		b.box.PackStart(frm, true, true, 0)
	}

	if taskResult.spinning {
		b.SwitchToSpinner(taskResult)
	} else if taskResult.image != nil {
		b.SwitchToImage(taskResult)
	}
}
