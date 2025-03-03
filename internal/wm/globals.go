/*
 * Global vars for WM
 * and common utils
 */

package wm

import (
	"riwo/internal/wm/window"
	"strconv"
	"syscall/js"
)

// bad idea. Needs to be more structurised.

var (
	ContextMenu      js.Value
	ContextMenuHides []js.Value
	StartY           float64
	StartX           float64
	AllWindows       map[string]*window.Window // All Go Windows

	// WindowCount
	// Counter for creating multiple windows with unique z-index
	WindowCount int
	// HighestZIndex
	// Track the highest z-index for bringing windows to front
	HighestZIndex int = 10
	// Verbose
	// holds flag of logger usage.
	Verbose bool = false
	// ColorSchemes
	// Represents loaded color schemes
	ColorSchemes = map[string]map[string]string{
		"monochrome": {
			"faded":  "#ffffff",
			"normal": "#777777",
			"vivid":  "#000000",
		},
		"red": {
			"faded":  "#ffeaea",
			"normal": "#df9595",
			"vivid":  "#bb5d5d",
		},
		"green": {
			"faded":  "#eaffea",
			"normal": "#88cc88",
			"vivid":  "#448844",
		},
		"blue": {
			"faded":  "#c0eaff",
			"normal": "#00aaff",
			"vivid":  "#0088cc",
		},
		"yellow": {
			"faded":  "#ffffea",
			"normal": "#eeee9e",
			"vivid":  "#99994c",
		},
		"aqua": {
			"faded":  "#eaffff",
			"normal": "#9eeeee",
			"vivid":  "#8888cc",
		},
		"gray": {
			"faded":  "#eeeeee",
			"normal": "#cccccc",
			"vivid":  "#888888",
		},
	}
)

func Print(value string) {
	js.Global().Get("console").Call("log", value)
}
func Ftoa(value float64) string {
	return strconv.FormatFloat(value, 'f', 6, 64)
}
