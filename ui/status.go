package ui

import "github.com/gotk3/gotk3/gtk"

const (
	StatusLeft = iota
	StatusRight
)

type Status struct {
	StatusBar *gtk.Statusbar
}

func (status *Status) Create() {

	status.StatusBar, _ = gtk.StatusbarNew()
	status.StatusBar.Push(status.StatusBar.GetContextId("Main"), "...")
}

func (status *Status) SetText(str string) {

	if status.StatusBar != nil {
		status.StatusBar.Pop(status.StatusBar.GetContextId("Main"))
		status.StatusBar.Push(status.StatusBar.GetContextId("Main"), str)
	}
}
