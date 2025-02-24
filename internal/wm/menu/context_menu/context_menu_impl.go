/*
 * Context Menu, its callbacks and listeners
 */

package context_menu

import (
	"riwo/internal/wm"
	"strconv"
	"syscall/js"
)

// New
// creates and sets up the global context menu.
// It adds both default options and (if available)
// custom entries from the current window.
func New() {
	document := js.Global().Get("document")
	body := document.Get("body")

	// Create the context menu container.
	wm.ContextMenu = document.Call("createElement", "div")
	wm.ContextMenu.Set("id", "contextMenu")
	menuStyle := "position: absolute; display: none; background-color: #EEFFEE; border: solid #8BCE8B; padding: 0; text-align: center;"
	wm.ContextMenu.Set("style", menuStyle)
	body.Call("appendChild", wm.ContextMenu)

	// Pre-create default options.
	moveOption := NewMenuItem("Move")
	newOption := NewMenuItem("New")
	resizeOption := NewMenuItem("Resize")
	deleteOption := NewMenuItem("Delete")
	hideOption := NewMenuItem("Hide")

	// Set up default options event handlers.
	moveOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			wm.JustSelected = true
			wm.IsMovingMode = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			wm.ContextMenu.Get("style").Set("display", "none")
			if wm.Verbose {
				wm.Print("Move mode activated.")
			}
		}
		return nil
	}))
	resizeOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			wm.JustSelected = true
			wm.IsResizingMode = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			wm.ContextMenu.Get("style").Set("display", "none")
			if wm.Verbose {
				wm.Print("Resize mode activated.")
			}
		}
		return nil
	}))
	deleteOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			wm.JustSelected = true
			wm.IsDeleteMode = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			wm.ContextMenu.Get("style").Set("display", "none")
			if wm.Verbose {
				wm.Print("Delete mode activated.")
			}
		}
		return nil
	}))
	hideOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			wm.JustSelected = true
			wm.IsHiding = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			wm.ContextMenu.Get("style").Set("display", "none")

			if wm.Verbose {
				wm.Print("Hide mode activated.")
			}
		}
		return nil
	}))

	// Background-specific options
	newOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			wm.JustSelected = true
			wm.IsNewMode = true
			wm.IsDragging = false
			wm.StartX = 0
			wm.StartY = 0
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			wm.ContextMenu.Get("style").Set("display", "none")
			if wm.Verbose {
				wm.Print("New mode activated. Select an area to create a window.")
			}
		}
		return nil
	}))

	// Omit browser's context menu
	body.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")
		//justSelected = false
		return nil
	}))
	// And call ours
	body.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")
		if (args[0].Get("button").Int() == 2) && !wm.JustSelected && !wm.IsMovingMode && !wm.IsResizingMode && !wm.IsDeleteMode && !wm.IsNewMode && !wm.IsHiding {
			// Clear all previous menu items.
			for wm.ContextMenu.Get("firstChild").Truthy() {
				wm.ContextMenu.Call("removeChild", wm.ContextMenu.Get("firstChild"))
			}
			// Add default options.
			wm.ContextMenu.Call("appendChild", moveOption)
			wm.ContextMenu.Call("appendChild", newOption)
			wm.ContextMenu.Call("appendChild", resizeOption)
			wm.ContextMenu.Call("appendChild", deleteOption)
			wm.ContextMenu.Call("appendChild", hideOption)
			// Clear active window if deleted (or something)
			if (wm.Window != nil) && (wm.Window.ID == -1) {
				wm.Window = nil
			}
			// Append custom context menu entries from the active window, if any.
			if wm.JsWindow.Truthy() && wm.Window != nil && wm.Window.ContextEntries != nil {
				for _, customOption := range wm.Window.ContextEntries {
					opt := NewMenuItem(customOption.Name)
					opt.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						args[0].Call("preventDefault")
						args[0].Call("stopPropagation")
						wm.JustSelected = true
						// Execute default callback for this option.
						customOption.Callback()
						wm.ContextMenu.Get("style").Set("display", "none") // hide menu after click
						wm.JustSelected = false
						if wm.Verbose {
							wm.Print("Custom option " + customOption.Name + " called")
						}
						return nil
					}))
					wm.ContextMenu.Call("appendChild", opt)
				}
			}
			// Append hidden windows' unhides, if any.
			if len(wm.ContextMenuHides) > 0 {
				for _, hiddenWindowButton := range wm.ContextMenuHides {
					wm.ContextMenu.Call("appendChild", hiddenWindowButton)
				}
			}

			// Position and show the menu.
			wm.ContextMenu.Get("style").Set("z-index", strconv.Itoa(wm.HighestZIndex+10))
			wm.ContextMenu.Get("style").Set("left", strconv.Itoa(args[0].Get("clientX").Int())+"px")
			wm.ContextMenu.Get("style").Set("top", strconv.Itoa(args[0].Get("clientY").Int())+"px")
			wm.ContextMenu.Get("style").Set("display", "block")
		}
		wm.JustSelected = false
		return nil
	}))
}

// NewMenuItem
// creates a new context menu option element.
func NewMenuItem(optionText string) js.Value {
	document := js.Global().Get("document")
	option := document.Call("createElement", "div")
	option.Set("innerText", optionText)
	option.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
	option.Get("style").Set("padding", "10px")
	option.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		option.Get("style").Set("background-color", "#418941")
		return nil
	}))
	option.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		option.Get("style").Set("background-color", "#EEFFEE")
		return nil
	}))
	return option
}

// DeleteMenuItem
// removes a given menu option from the context menu.
func DeleteMenuItem(option js.Value) {
	document := js.Global().Get("document")
	menu := document.Call("getElementById", "contextMenu")
	menu.Call("removeChild", option)
}
