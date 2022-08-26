package ui

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Toolbar struct {
}

func (tb *Toolbar) Create() *gtk.Toolbar {

	toolbar, err := gtk.ToolbarNew()

	if err != nil {
		log.Fatal("Unable to create toolbar:", err)
	}
	toolbar.SetStyle(gtk.TOOLBAR_BOTH)

	img, err := gtk.ImageNew()
	if err != nil {
		log.Fatal("Unable to create image:", err)
	}
	img.SetFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	playTb, err := gtk.ToolButtonNew(img, "Connect")
	if err != nil {
		log.Fatal("Unable to create newTb:", err)
	}
	playTb.SetName("Play")
	toolbar.Insert(playTb, -1)

	img, _ = gtk.ImageNew()
	img.SetFromIconName("media-playback-stop", gtk.ICON_SIZE_BUTTON)
	stopTb, err := gtk.ToolButtonNew(img, "Disconnect")
	if err != nil {
		log.Fatal("Unable to create newTb:", err)
	}
	stopTb.SetName("Stop")
	toolbar.Insert(stopTb, -1)

	playTb.Connect("clicked", func() {
		onBtnPlayClick()
	})

	return toolbar
}
