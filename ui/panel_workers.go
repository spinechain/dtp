package ui

import (
	"log"
	"spinedtp/taskworkers"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type PanelWorkers struct {
	cmdLabel *gtk.Label
	list     *gtk.TreeView
	frame    *gtk.Frame
	store    *gtk.ListStore
}

// IDs to access the tree view columns by
const (
	COLUMN_WORKERS_ID = iota
	COLUMN_WORKERS_ADDR
	COLUMN_WORKERS_PORT
	COLUMN_WORKERS_TASKSDONE
	COLUMN_WORKERS_REPUTATION
)

// Create a new panel that lists out all tasks
func (tasks *PanelWorkers) Create(title string) (*gtk.Box, error) {

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

	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("ID", COLUMN_WORKERS_ID))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Address", COLUMN_WORKERS_ADDR))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Port", COLUMN_WORKERS_PORT))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Tasks Done", COLUMN_WORKERS_TASKSDONE))
	tasks.list.AppendColumn(tasks.CreateTasksTextColumn("Reputation", COLUMN_WORKERS_REPUTATION))

	tasks.frame.Add(tasks.list)
	tasks.frame.SetMarginStart(20)
	tasks.frame.SetMarginEnd(20)

	tasks.list.SetHeadersVisible(true)

	box.PackStart(tasks.frame, true, true, 25)

	tasks.store = tasks.initList(tasks.list)

	return box, err
}

// Add a column to the tree view (during the initialization of the tree view)
// We need to distinct the type of data shown in either column.
func (tasks *PanelWorkers) CreateTasksTextColumn(title string, id int) *gtk.TreeViewColumn {

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

func (tasks *PanelWorkers) UpdateList(list []*taskworkers.TaskWorker) {

	if tasks.store != nil {
		tasks.store.Clear()
		for _, item := range list {
			iter := tasks.store.Append()
			tasks.store.SetValue(iter, COLUMN_WORKERS_ID, item.ID)
			tasks.store.SetValue(iter, COLUMN_WORKERS_ADDR, item.Address)
			tasks.store.SetValue(iter, COLUMN_WORKERS_PORT, item.Port)
			tasks.store.SetValue(iter, COLUMN_WORKERS_TASKSDONE, item.TasksDone)
			tasks.store.SetValue(iter, COLUMN_WORKERS_REPUTATION, item.Reputation)
		}
	}
}

func (tasks *PanelWorkers) initList(list *gtk.TreeView) *gtk.ListStore {
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

	store, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal(err)
	}
	list.SetModel(store)

	return store
}

func (tasks *PanelWorkers) Destroy() {

}

func (tasks *PanelWorkers) Show() {
	tasks.cmdLabel.ShowAll()
	tasks.list.ShowAll()
	tasks.frame.ShowAll()
}

func (tasks *PanelWorkers) Hide() {
	tasks.cmdLabel.Hide()
	tasks.list.Hide()
	tasks.frame.Hide()
}
