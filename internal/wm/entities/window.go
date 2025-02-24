package entities

import (
	"riwo/internal/wm/window"
	"syscall/js"
)

type WindowFaces struct {
	Window      *window.Window // Active Window structure
	JsWindow    js.Value       // Active JS window structure
	GhostWindow js.Value
}
