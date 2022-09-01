package ui

import (
	util "spinedtp/util"
)

// The main window
var mainWindow *Window

// The sidebar treeview
var treeViewSidebar *Treeview

// Status bar at the bottom
var statusBar1 *Status
var statusBar2 *Status

// The panel
var CommandPanel *PanelCommand
var HistoryPanel *PanelNetwork
var NetworkPanel *PanelNetwork
var TasksPanel *PanelTasks

// Create all the windows
func Create() {

	// Create main window
	mainWindow = new(Window)
	mainWindow.CreateWindow()

	// Create top toolbar
	toolbar := new(Toolbar)

	// Add to window
	mainWindow.AddToolbar(toolbar)

	treeViewSidebar = new(Treeview)
	mainWindow.AddTreeview(treeViewSidebar)

	// Create the status bars at the bottom and add to the window
	statusBar1 = new(Status)
	statusBar2 = new(Status)

	mainWindow.AddStatus(statusBar1, 0)
	mainWindow.AddStatus(statusBar2, 1)

	// Create the panels on the right that are selected when a treeview item is clicked
	CommandPanel = new(PanelCommand)
	HistoryPanel = new(PanelNetwork)
	NetworkPanel = new(PanelNetwork)
	TasksPanel = new(PanelTasks)

	// Add the panels to the window
	mainWindow.AddPanel(CommandPanel, "Commands")
	mainWindow.AddPanel(HistoryPanel, "History")
	mainWindow.AddPanel(NetworkPanel, "Workers")
	mainWindow.AddPanel(TasksPanel, "Tasks")

	mainWindow.SetIcon()

	mainWindow.Show()

	// Add all the treeview items to switch panels
	AddTreeviewItems()

	// Select the default
	SwitchPanel("commands")
}

// Configures all the things we want in the treeview on the left
func AddTreeviewItems() {
	treeViewSidebar.addRow(treeViewSidebar.treeStore, "appendage", "Commands", "commands/commands")
	treeViewSidebar.addRow(treeViewSidebar.treeStore, "list", "History", "history/history")
	treeViewSidebar.addRow(treeViewSidebar.treeStore, "list", "Workers", "network/network")
	treeViewSidebar.addRow(treeViewSidebar.treeStore, "list", "Tasks", "tasks/tasks")
}

// Can be called to update the status bar at the bottom of the window
func UpdateStatusBar(s string, section int) {

	if statusBar1 == nil || statusBar2 == nil {

		if section == 0 {
			util.PrintPurple(s)
		} else {
			util.PrintYellow(s)
		}

		return
	}

	if section == 0 {
		statusBar1.SetText(s)
	} else {
		statusBar2.SetText(s)
	}

}

// Switches to the panel we want to show. Hides the others.
func SwitchPanel(str string) {

	if str == "commands" {
		NetworkPanel.Hide()
		HistoryPanel.Hide()
		CommandPanel.Show()
		TasksPanel.Hide()
	} else if str == "history" {
		NetworkPanel.Hide()
		HistoryPanel.Show()
		CommandPanel.Hide()
		TasksPanel.Hide()
	} else if str == "network" {
		NetworkPanel.Show()
		HistoryPanel.Hide()
		CommandPanel.Hide()
		TasksPanel.Hide()
	} else if str == "tasks" {
		NetworkPanel.Hide()
		HistoryPanel.Hide()
		CommandPanel.Hide()
		TasksPanel.Show()
	}

}
