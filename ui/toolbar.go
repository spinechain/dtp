package ui

import (
	"log"

	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

//go:embed icons/connect.png
var imgConnectBuf []byte

//go:embed icons/disconnect.png
var imgDisconnectBuf []byte

var (
	imageConnect    *gdk.Pixbuf = nil
	imageDisconnect *gdk.Pixbuf = nil
)

type Toolbar struct {
}

func (tb *Toolbar) Create() *gtk.Toolbar {

	toolbar, err := gtk.ToolbarNew()

	if err != nil {
		log.Fatal("Unable to create toolbar:", err)
	}
	toolbar.SetStyle(gtk.TOOLBAR_BOTH)

	imageConnect, err = gdk.PixbufNewFromBytesOnly(imgConnectBuf)
	if err != nil {
		log.Fatal("Unable to load image:", err)
	}
	imageConnect, _ = imageConnect.ScaleSimple(16, 16, gdk.INTERP_BILINEAR)

	imageDisconnect, err = gdk.PixbufNewFromBytesOnly(imgDisconnectBuf)
	if err != nil {
		log.Fatal("Unable to load image:", err)
	}
	imageDisconnect, _ = imageDisconnect.ScaleSimple(16, 16, gdk.INTERP_BILINEAR)

	img, err := gtk.ImageNew()
	if err != nil {
		log.Fatal("Unable to create image:", err)
	}
	img.SetFromPixbuf(imageConnect)

	// .SetFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	playTb, err := gtk.ToolButtonNew(img, "Connect")
	if err != nil {
		log.Fatal("Unable to create newTb:", err)
	}
	playTb.SetName("Play")
	toolbar.Insert(playTb, -1)

	img, _ = gtk.ImageNew()
	// img.SetFromIconName("media-playback-stop", gtk.ICON_SIZE_BUTTON)
	img.SetFromPixbuf(imageDisconnect)
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
