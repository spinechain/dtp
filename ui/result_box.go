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
	image    *gtk.Image
	spinner  *gtk.Spinner
	text     string
}

func (b *ResultBox) SwitchToSpinner() {

	curChild, err := b.scroller.GetChild()
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	if curChild != nil {
		b.scroller.Remove(curChild)
	}

	if b.spinner == nil {
		spinner, _ := gtk.SpinnerNew()
		b.spinner = spinner
	}

	b.scroller.Add(b.spinner)

	if b.label != nil {
		b.box.Remove(b.label)
		b.label.SetText(b.text)
		b.box.PackStart(b.label, false, false, 0)
		b.box.ShowAll()
	}

	b.scroller.ShowAll()
	b.spinner.ShowAll()

	b.spinner.Start()

}

func (b *ResultBox) SwitchToImage() {

	curChild, err := b.scroller.GetChild()
	if err != nil {
		util.PrintRed(err.Error())
		return
	}

	if curChild != nil {
		b.scroller.Remove(curChild)
	}

	// Add the image to the scroller
	b.scroller.Add(b.image)

	// Add the label below the scrolled window
	if b.label != nil {
		b.box.Remove(b.label)
		b.label.SetText(b.text)
		b.box.PackStart(b.label, false, false, 0)
		b.box.ShowAll()
	}

}

func (b *ResultBox) AdaptToCircumstances(taskResult *TaskResult) {

	if taskResult.spinning {
		b.SwitchToSpinner()
	} else if b.image != nil {
		b.SwitchToImage()
	}

	/*
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

	/*
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
	*/
}
