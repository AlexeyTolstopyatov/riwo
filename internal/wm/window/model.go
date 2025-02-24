package window

import (
	"riwo/internal/controls"
	"syscall/js"
)

// Window
// Type `Window` manages single
// abstract window’s properties.
type Window struct {
	ID             int                     // For the most part unites DOM object and Go object
	Element        js.Value                // Connected DOM element.
	ContextEntries []controls.ContextEntry // "Move", "Resize", "Delete" and "Hide" are basic ones
}
