package ui

import (
	"log"
	"spinedtp/tasknet"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type PanelTasks struct {
	cmdLabel *gtk.Label
	list     *gtk.TreeView
	frame    *gtk.Frame
	store    *gtk.ListStore
	btnClear *gtk.Button
}

// IDs to access the tree view columns by
const (
	COLUMN_TASK_FIRST = iota
	COLUMN_TASK_SECOND
	COLUMN_TASK_THIRD
	COLUMN_TASK_FOURTH
	COLUMN_TASK_FIFTH
	COLUMN_TASK_SIXTH
	COLUMN_TASK_SEVENTH
	COLUMN_TASK_EIGTH
)

// Create a new panel that lists out all tasks
func (tasks *PanelTasks) Create(title string) (*gtk.Box, error) {

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}

	tasks.cmdLabel, _ = gtk.LabelNew(title)

	box.PackStart(tasks.cmdLabel, false, false, 5)

	tasks.frame, _ = gtk.FrameNew(title)

	tasks.list, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatal(err)
	}

	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Task", COLUMN_TASK_FIRST))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("ID", COLUMN_TASK_SECOND))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Owner", COLUMN_TASK_THIRD))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Fee", COLUMN_TASK_FOURTH))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Reward", COLUMN_TASK_FIFTH))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Created", COLUMN_TASK_SIXTH))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Synced", COLUMN_TASK_SEVENTH))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Height", COLUMN_TASK_EIGTH))

	tasks.frame.Add(tasks.list)
	tasks.frame.SetMarginStart(20)
	tasks.frame.SetMarginEnd(20)

	tasks.btnClear, _ = gtk.ButtonNew()
	tasks.btnClear.SetLabel("Clear Tasks")
	tasks.btnClear.Connect("clicked", onBtnClearTasksClick)

	tasks.list.SetHeadersVisible(true)

	box.PackStart(tasks.frame, true, true, 25)
	box.PackStart(tasks.btnClear, false, false, 25)

	tasks.store = tasks.initList(tasks.list)

	return box, err
}

// Add a column to the tree view (during the initialization of the tree view)
// We need to distinct the type of data shown in either column.
func (tasks *PanelTasks) CreateTasksTextColumn(title string, id int) *gtk.TreeViewColumn {

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

func (tasks *PanelTasks) UpdateList(list []*tasknet.Task) {

	if tasks.store != nil {
		tasks.store.Clear()
		for _, item := range list {
			iter := tasks.store.Append()
			tasks.store.SetValue(iter, 0, item.Command)
			tasks.store.SetValue(iter, 1, item.ID)
			tasks.store.SetValue(iter, 2, item.TaskOwnerID)
			tasks.store.SetValue(iter, 3, item.Fee)
			tasks.store.SetValue(iter, 4, item.Reward)
			tasks.store.SetValue(iter, 5, item.Created)
			tasks.store.SetValue(iter, 6, item.FullyPropagated)
			tasks.store.SetValue(iter, 7, item.Index)
		}
	}
}

func (tasks *PanelTasks) initList(list *gtk.TreeView) *gtk.ListStore {
	/*
		cellRenderer, err := gtk.CellRendererTextNew()
		if err != nil {
			log.Fatal(err)
		}


			column, err := gtk.TreeViewColumnNewWithAttribute("List Items", cellRenderer, "text", 0)
			if err != nil {
				log.Fatal(err)
			}
			list.AppendColumn(column)
	*/

	store, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal(err)
	}
	list.SetModel(store)

	return store
}

func (tasks *PanelTasks) Destroy() {

}

func (tasks *PanelTasks) Show() {
	tasks.cmdLabel.ShowAll()
	tasks.list.ShowAll()
	tasks.frame.ShowAll()
	tasks.btnClear.ShowAll()
}

func (tasks *PanelTasks) Hide() {
	tasks.cmdLabel.Hide()
	tasks.list.Hide()
	tasks.frame.Hide()
	tasks.btnClear.Hide()
}
