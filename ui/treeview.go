package ui

import (
	"log"
	"strings"

	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// IDs to access the tree view columns by
const (
	COLUMN_FIRST = iota
	COLUMN_SECOND
	COLUMN_TYPE
)

type Treeview struct {
	treeStore *gtk.TreeStore
	treeView  *gtk.TreeView
}

//go:embed icons/red.png
var imgAppendageBuf []byte

//go:embed icons/center.png
var imgCenterBuf []byte

//go:embed icons/memory.png
var imgMemoryBuf []byte

//go:embed icons/wrench.png
var imgWrench []byte

//go:embed icons/list.png
var imgList []byte

// We want to have the icons in image files, hence need to individually
// specify the path to reach them. Here an environment variable is used.

// Icons Pixbuf representation
var (
	imageAppendage *gdk.Pixbuf = nil
	imageCenter    *gdk.Pixbuf = nil
	imageMemory    *gdk.Pixbuf = nil
	imageWrench    *gdk.Pixbuf = nil
	imageList      *gdk.Pixbuf = nil
)

func (tv *Treeview) Create() *gtk.TreeView {
	var treeView *gtk.TreeView
	treeView, tv.treeStore = tv.SetupTreeView()
	treeView.SetModel(tv.treeStore)

	var err error
	imageAppendage, err = gdk.PixbufNewFromBytesOnly(imgAppendageBuf)
	if err != nil {
		log.Fatal("Unable to load image:", err)
	}
	imageAppendage, _ = imageAppendage.ScaleSimple(16, 16, gdk.INTERP_BILINEAR)

	imageMemory, err = gdk.PixbufNewFromBytesOnly(imgMemoryBuf)
	if err != nil {
		log.Fatal("Unable to load image:", err)
	}
	imageMemory, _ = imageMemory.ScaleSimple(16, 16, gdk.INTERP_BILINEAR)

	imageCenter, err = gdk.PixbufNewFromBytesOnly(imgCenterBuf)
	if err != nil {
		log.Fatal("Unable to load image:", err)
	}
	imageCenter, _ = imageCenter.ScaleSimple(16, 16, gdk.INTERP_BILINEAR)

	imageWrench, err = gdk.PixbufNewFromBytesOnly(imgWrench)
	if err != nil {
		log.Fatal("Unable to load image:", err)
	}
	imageWrench, _ = imageWrench.ScaleSimple(16, 16, gdk.INTERP_BILINEAR)

	imageList, err = gdk.PixbufNewFromBytesOnly(imgList)
	if err != nil {
		log.Fatal("Unable to load image:", err)
	}
	imageList, _ = imageList.ScaleSimple(16, 16, gdk.INTERP_BILINEAR)

	return treeView
}

// Creates a tree view and the list store that holds its data
func (tv *Treeview) SetupTreeView() (*gtk.TreeView, *gtk.TreeStore) {
	var err error
	tv.treeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatal("Unable to create tree view:", err)
	}

	tv.treeView.SetSizeRequest(250, -1)

	tv.treeView.AppendColumn(tv.CreateImageColumn("", COLUMN_FIRST))
	tv.treeView.AppendColumn(tv.CreateTextColumn("Name", COLUMN_SECOND))

	// Creating a tree store. This is what holds the data that will be shown on our tree view.
	tv.treeStore, err = gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create tree store:", err)
	}

	selection, err := tv.treeView.GetSelection()
	if err != nil {
		log.Fatal("Could not get tree selection object.")
	}
	selection.SetMode(gtk.SELECTION_SINGLE)
	selection.Connect("changed", tv.TreeSelectionChangedCB)

	return tv.treeView, tv.treeStore
}

// Add a column to the tree view (during the initialization of the tree view)
// We need to distinct the type of data shown in either column.
func (tv *Treeview) CreateTextColumn(title string, id int) *gtk.TreeViewColumn {

	// In this column we want to show text, hence create a text renderer
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	// Tell the renderer where to pick input from. Text renderer understands
	// the "text" property.

	//column, err := gtk.TreeViewColumnNew()
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	return column
}

func (tv *Treeview) CreateImageColumn(title string, id int) *gtk.TreeViewColumn {
	// In this column we want to show image data from Pixbuf, hence
	// create a pixbuf renderer
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		log.Fatal("Unable to create pixbuf cell renderer:", err)
	}

	// Tell the renderer where to pick input from. Pixbuf renderer understands
	// the "pixbuf" property.
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	return column
}

// Handle selection
func (tv *Treeview) TreeSelectionChangedCB(selection *gtk.TreeSelection) {
	var iter *gtk.TreeIter
	var model gtk.ITreeModel
	var ok bool
	model, iter, ok = selection.GetSelected()

	if ok {

		tpath, err := model.(*gtk.TreeModel).GetPath(iter)
		if err != nil {
			log.Printf("treeSelectionChangedCB: Could not get path from model: %s\n", err)
			return
		}
		log.Printf("treeSelectionChangedCB: selected path: %s\n", tpath)

		val, err := tv.treeStore.GetValue(iter, COLUMN_TYPE)

		s, _ := val.GetString()
		log.Printf("treeSelectionChangedCB: selected path: %s\n", s)

		items := strings.Split(s, "/")
		if len(items) == 3 {
			type_name := items[0]
			center_name := items[1]
			action_name := items[2]
			log.Printf("Type=%s, Center=%s, Action=%s\n", type_name, center_name, action_name)

			OnTreeviewItemSelected(type_name, center_name, action_name)
		} else if len(items) == 2 {
			OnTreeviewItemSelected(items[0], items[1], "")
		}
	}
}

func (tv *Treeview) SelectItem(it *gtk.TreeIter) {

	selection, _ := tv.treeView.GetSelection()
	selection.SelectIter(it)
}

// Append a toplevel row to the tree store for the tree view
func (tv *Treeview) addRow(treeStore *gtk.TreeStore, texticon string, text string, dataType string) *gtk.TreeIter {
	return tv.addSubRow(treeStore, nil, texticon, text, dataType)
}

// Append a sub row to the tree store for the tree view
func (tv *Treeview) addSubRow(treeStore *gtk.TreeStore, iter *gtk.TreeIter, texticon string, text string, idPath string) *gtk.TreeIter {
	// Get an iterator for a new row at the end of the list store
	i := treeStore.Append(iter)
	// i := treeStore.InsertAfter(iter, nil)

	// Set the contents of the tree store row that the iterator represents
	var icon *gdk.Pixbuf = imageMemory

	if texticon == "memory" {
		icon = imageMemory
	} else if texticon == "center" {
		icon = imageCenter
	} else if texticon == "appendage" {
		icon = imageAppendage
	} else if texticon == "wrench" {
		icon = imageWrench
	} else if texticon == "list" {
		icon = imageList
	}

	err := treeStore.SetValue(i, COLUMN_FIRST, icon)
	if err != nil {
		log.Fatal("Unable set value:", err)
	}
	err = treeStore.SetValue(i, COLUMN_SECOND, text)
	if err != nil {
		log.Fatal("Unable set value:", err)
	}

	err = treeStore.SetValue(i, COLUMN_TYPE, idPath)
	if err != nil {
		log.Fatal("Unable set value:", err)
	}

	return i
}
