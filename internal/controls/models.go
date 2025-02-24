/*
 * Window struct for WM
 * and its functions + listeners
 */

package controls

import (
	"syscall/js"
)

// ContextEntry
// Define a named type
// for the context entry.
type ContextEntry struct {
	Name     string
	Callback func()
}

// Window
// Type `Window` manages single
// abstract window’s properties.
type Window struct {
	ID             int            // For the most part unites DOM object and Go object
	Element        js.Value       // Connected DOM element.
	ContextEntries []ContextEntry // Tho "Move", "Resize", "Delete" and "Hide" are basic ones
}
