/*
 * (Almost) All javascript DOM events here in only copy
 * for the sake of optimization
 */

package events

import (
	"fmt"
	"riwo/internal/wm"
	"riwo/internal/wm/window"
	"syscall/js"
)

// New sets up global mouse event listeners.
func New() {
	addMouseMove()
	addMouseDown()
	addMouseUp()
	return
}

// addMouseMove
// registers mouse movement event in Document,
// includes assets and
func addMouseMove() {
	js.Global().Get("document").Call(
		"addEventListener",
		"mousemove",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// In "Move" mode
			if wm.IsDragging && wm.IsMovingMode && wm.GhostWindow.Truthy() {
				x := args[0].Get("clientX").Float() - wm.StartX
				y := args[0].Get("clientY").Float() - wm.StartY
				wm.GhostWindow.Get("style").Set("left", wm.Ftoa(x)+"px")
				wm.GhostWindow.Get("style").Set("top", wm.Ftoa(y)+"px")
				if wm.Verbose {
					wm.Print("Ghost window is moving.")
				}
			}

			// In "Resize" mode adjust ghost window size
			if wm.GhostWindow.Truthy() && wm.IsResizingMode && wm.IsResizingInit && wm.IsDragging && !wm.IsMovingMode {
				currentX := args[0].Get("clientX").Float()
				currentY := args[0].Get("clientY").Float()

				// Calculate and update ghost window size and position based on selection
				width := currentX - wm.StartX
				height := currentY - wm.StartY
				wm.GhostWindow.Get("style").Set("width", wm.Ftoa(width)+"px")
				wm.GhostWindow.Get("style").Set("height", wm.Ftoa(height)+"px")

				// Handle direction to ensure ghost window moves according to selection direction
				if width < 0 {
					wm.GhostWindow.Get("style").Set("left", wm.Ftoa(currentX)+"px")
					wm.GhostWindow.Get("style").Set("width", wm.Ftoa(-width)+"px")
				}
				if height < 0 {
					wm.GhostWindow.Get("style").Set("top", wm.Ftoa(currentY)+"px")
					wm.GhostWindow.Get("style").Set("height", wm.Ftoa(-height)+"px")
				}
				if wm.Verbose {
					wm.Print("Ghost window is resizing with freeform selection.")
				}
			}

			// In "New" mode adjusting selection
			if wm.GhostWindow.Truthy() && wm.IsNewMode && wm.IsDragging {
				currentX := args[0].Get("clientX").Float()
				currentY := args[0].Get("clientY").Float()

				// Calculate and update ghost window size and position based on selection
				width := currentX - wm.StartX
				height := currentY - wm.StartY
				wm.GhostWindow.Get("style").Set("width", wm.Ftoa(width)+"px")
				wm.GhostWindow.Get("style").Set("height", wm.Ftoa(height)+"px")

				// Handle direction to ensure ghost window moves according to selection direction
				if width < 0 {
					wm.GhostWindow.Get("style").Set("left", wm.Ftoa(currentX)+"px")
					wm.GhostWindow.Get("style").Set("width", wm.Ftoa(-width)+"px")
				}
				if height < 0 {
					wm.GhostWindow.Get("style").Set("top", wm.Ftoa(currentY)+"px")
					wm.GhostWindow.Get("style").Set("height", wm.Ftoa(-height)+"px")
				}

				if wm.Verbose {
					wm.Print("New window selection resizing.")
				}
			}

			return nil
		}))
}

// addMouseUp
// Registers MouseUp event
// (when mouse button-down event firstly stops handling)
func addMouseUp() {
	js.Global().Get("document").Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// In "move" stop dragging (Teleport to ghost)
		if wm.IsMovingMode && wm.IsDragging {
			args[0].Call("stopPropagation")
			wm.IsDragging = false
			wm.IsMovingMode = false
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")

			// Move the window to the ghost's position
			if wm.GhostWindow.Truthy() && wm.JsWindow.Truthy() {
				wm.JsWindow.Get("style").Set("left", wm.GhostWindow.Get("style").Get("left"))
				wm.JsWindow.Get("style").Set("top", wm.GhostWindow.Get("style").Get("top"))
				wm.GhostWindow.Call("remove") // Delete ghost window
				wm.GhostWindow = js.Null()    // Reset ghost window reference
			}
			wm.JustSelected = false

			if wm.Verbose {
				wm.Print("Dragging ended and window teleported to ghost position.")
			}
		}

		// In "Resize" stop selecting (Teleport to ghost)
		if wm.GhostWindow.Truthy() && wm.JsWindow.Truthy() && wm.IsResizingMode && wm.IsResizingInit && wm.IsDragging && !wm.IsMovingMode {
			args[0].Call("stopPropagation")
			wm.IsResizingMode = false
			wm.IsResizingInit = false
			wm.IsDragging = false

			// Replace all dimensions with ghost's ones
			wm.JsWindow.Get("style").Set("left", wm.GhostWindow.Get("style").Get("left"))
			wm.JsWindow.Get("style").Set("top", wm.GhostWindow.Get("style").Get("top"))
			wm.JsWindow.Get("style").Set("width", wm.GhostWindow.Get("style").Get("width"))
			wm.JsWindow.Get("style").Set("height", wm.GhostWindow.Get("style").Get("height"))
			wm.GhostWindow.Call("remove") // Delete the ghost window
			wm.GhostWindow = js.Null()    // Reset ghost window reference

			// Reset cursor
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
			wm.JustSelected = false
			if wm.Verbose {
				wm.Print("Resizing completed and window resized to match selection.")
			}
		}

		// In "New" make new window
		if wm.IsNewMode && wm.IsDragging {
			args[0].Call("stopPropagation")
			wm.IsNewMode = false
			wm.IsDragging = false
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")

			// Create a new window at the ghost window's position and size
			if wm.GhostWindow.Truthy() {
				x := wm.GhostWindow.Get("style").Get("left").String()
				y := wm.GhostWindow.Get("style").Get("top").String()
				width := wm.GhostWindow.Get("style").Get("width").String()
				height := wm.GhostWindow.Get("style").Get("height").String()
				wm.GhostWindow.Call("remove")
				wm.GhostWindow = js.Null()
				neuwindow := window.New(x, y, width, height, "")
				wm.JustSelected = false
				if wm.Verbose {
					wm.Print("New window created at selected area.")
				}
				// Intends existance of default app
				AppDefault := js.Global().Get("LaunchDefault")
				go AppDefault.Invoke(neuwindow.ID)
				if wm.Verbose {
					wm.Print("LaunchDefault attached to the new window")
				}
			}
		}
		return nil
	}))
}

// addMouseDown
// Registers Mouse button-down event
// (invokes when mouse key is under pressure.)
func addMouseDown() {
	js.Global().Get("document").Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		// In "Resize" Second RMB after choosing window
		// Also "New"'s initial area selecting
		if (wm.IsResizingMode || wm.IsNewMode) && args[0].Get("button").Int() == 2 && (wm.IsResizingInit || wm.IsNewMode) {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			wm.JustSelected = true

			// Second RMB hold - Start resizing by creating a selection anywhere
			if wm.Verbose {
				wm.Print("Second right-click: Resizing initiated.")
			}
			wm.IsDragging = true

			// Reset selection start position
			wm.StartX = args[0].Get("clientX").Float()
			wm.StartY = args[0].Get("clientY").Float()

			// Create ghost window for resizing
			wm.GhostWindow = js.Global().Get("document").Call("createElement", "div")
			wm.GhostWindow.Set("style", fmt.Sprintf("position: absolute; z-index: %d", wm.HighestZIndex+1)+"; border: solid 2px #FF0000;")
			wm.GhostWindow.Get("style").Set("left", wm.Ftoa(wm.StartX)+"px")
			wm.GhostWindow.Get("style").Set("top", wm.Ftoa(wm.StartY)+"px")
			js.Global().Get("document").Get("body").Call("appendChild", wm.GhostWindow)

			if wm.Verbose {
				wm.Print("Resizing initiated with freeform selection.")
			}
		}

		// Cancel menu and reset all modes on leftclick.
		if args[0].Get("button").Int() == 0 {
			if wm.ContextMenu.Type() == js.TypeUndefined {
				wm.Print("ContextMenu is [undefined]")
				// Or handle the error appropriately
			} else {
				wm.ContextMenu.Get("style").Set("display", "none")
			}

			if wm.GhostWindow.Truthy() && wm.IsDragging {
				wm.GhostWindow.Call("remove")
				wm.GhostWindow = js.Null()
			}

			wm.IsDragging = false
			wm.IsMovingMode = false
			wm.IsResizingMode = false
			wm.IsResizingInit = false
			wm.JustSelected = false
			wm.IsDeleteMode = false
			wm.IsNewMode = false
			wm.IsHiding = false

			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
		}
		return nil
	}))
}
