package ui

import "github.com/gotk3/gotk3/gtk"

// This is the Appendage type which defines all the methods
// that an appendage has to support
type Panel interface {

	// Turn on everything and start the threads. There is no
	// guarantee that by the time it comes on that the threads
	// are up and running
	Create(title string) (*gtk.Box, error)

	// Shut down all threads
	Destroy()

	// Hide all widgets
	Hide()

	// Show all widgets
	Show()
}
