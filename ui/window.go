package ui

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Window struct {
	mainBox   *gtk.Box
	winBox    *gtk.Box
	statusBox *gtk.Box
	panelBox  *gtk.Box
	Window    *gtk.Window
}

func (window *Window) CreateWindow() {

	// Setup the window
	gtk.Init(nil)

	var err error
	window.Window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	window.Window.SetTitle("SpineChain DTP")
	window.Window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	window.Window.SetDefaultSize(800, 600)
	window.Window.SetPosition(gtk.WIN_POS_CENTER)

	// Create Layout for toolbar vs other parts
	window.mainBox, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	window.Window.Add(window.mainBox)
}

func (window *Window) SetIcon() {
	window.Window.SetIcon(imageAppendage)
}

func (window *Window) AddTreeview(treeview *Treeview) {

	// Create Box for the Treeview vs content area
	window.winBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	window.mainBox.PackStart(window.winBox, true, true, 0)

	tv := treeview.Create()
	window.winBox.PackStart(tv, false, true, 0)

	window.panelBox, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	window.winBox.PackStart(window.panelBox, true, true, 0)
}

func (window *Window) Show() {
	window.Window.ShowAll()

}

func (window *Window) AddPanel(panel Panel, title string) {
	panelBox, _ := panel.Create(title)
	window.panelBox.PackStart(panelBox, false, false, 5)
}

func (window *Window) AddToolbar(toolbar *Toolbar) {

	tb := toolbar.Create()
	window.mainBox.PackStart(tb, false, false, 5)
}

func (window *Window) AddStatus(status *Status, pos uint) {

	status.Create()

	if window.statusBox == nil {
		window.statusBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		window.mainBox.PackStart(window.statusBox, false, true, 0)
	}

	if pos == 3 || pos == 0 {
		window.statusBox.PackEnd(status.StatusBar, false, false, 0)
	} else {
		window.statusBox.PackStart(status.StatusBar, false, false, 0)
	}

	/*
		if pos == 0 {
			window.statusBox.PackEnd(status.StatusBar, false, false, 0)
		} else {
			window.statusBox.PackStart(status.StatusBar, false, false, 0)
		}
	*/
}
