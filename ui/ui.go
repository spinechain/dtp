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
var commandPanel *PanelCommand
var historyPanel *PanelNetwork
var networkPanel *PanelNetwork

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
	commandPanel = new(PanelCommand)
	historyPanel = new(PanelNetwork)
	networkPanel = new(PanelNetwork)

	// Add the panels to the window
	mainWindow.AddPanel(commandPanel, "Commands")
	mainWindow.AddPanel(historyPanel, "History")
	mainWindow.AddPanel(networkPanel, "Network")

	mainWindow.SetIcon()

	mainWindow.Show()

	AddTreeviewItems()

	SwitchPanel("commands")
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

func SwitchPanel(str string) {

	if str == "commands" {
		networkPanel.Hide()
		historyPanel.Hide()
		commandPanel.Show()
	} else if str == "history" {
		networkPanel.Hide()
		historyPanel.Show()
		commandPanel.Hide()
	} else if str == "network" {
		networkPanel.Show()
		historyPanel.Hide()
		commandPanel.Hide()
	}

}

func AddTreeviewItems() {

	treeViewSidebar.addRow(treeViewSidebar.treeStore, "appendage", "Commands", "commands/commands")
	treeViewSidebar.addRow(treeViewSidebar.treeStore, "list", "History", "history/history")
	treeViewSidebar.addRow(treeViewSidebar.treeStore, "list", "Network", "network/network")

}
