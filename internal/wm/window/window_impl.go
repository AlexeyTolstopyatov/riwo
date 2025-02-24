package window

//
// A "window" package represents modeling and implementation
// of Window entry.
//

import (
	"riwo/internal/controls"
	"riwo/internal/wm"
	"riwo/internal/wm/menu/context_menu"
	"strconv"
	"syscall/js"
)

// New
// NewWindow creates a new Window,
// sets up its DOM element, and returns a pointer to it.
func New(x, y, width, height, content string) *controls.Window {
	wm.WindowCount++
	id := wm.WindowCount

	document := js.Global().Get("document")
	body := document.Get("body")

	// Create the DOM element for the window.
	winElem := document.Call("createElement", "div")
	style := "overflow: hidden; position: absolute; z-index: " + strconv.Itoa(wm.HighestZIndex) + "; left: " + x +
		"; top: " + y + "; width: " + width + "; height: " + height +
		"; background-color: #f0f0f0; border: solid #55AAAA; padding: 0;"
	winElem.Set("style", style)
	winElem.Set("innerHTML", content)
	winElem.Set("id", strconv.Itoa(id)) // Assing shared ID

	body.Call("appendChild", winElem)

	// Logging
	if wm.Verbose {
		wm.Print("Generated window's ID (wid) is \"" +
			strconv.Itoa(id) + "\"")
	}

	neuwindow := &controls.Window{
		ID:      id,
		Element: winElem,
		// No custom ContextEntries
	}
	wm.Window = neuwindow
	wm.JsWindow = winElem
	wm.AllWindows[strconv.Itoa(neuwindow.ID)] = neuwindow

	// Bring to front when clicked
	winElem.Call(
		"addEventListener",
		"mousedown",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			wm.Window = neuwindow
			wm.JsWindow = winElem
			// Right-click (RMB) on the window to
			// select it for resizing, second right-click
			// activates resizing
			if wm.IsResizingMode && !wm.IsResizingInit && args[0].Get("button").Int() == 2 {
				// First RMB hold - Select the window for resizing
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")
				wm.JustSelected = true
				if wm.Verbose {
					wm.Print("First right-click: Window selected for resizing.")
				}
				wm.HighestZIndex++
				winElem.Get("style").Set("z-index", strconv.Itoa(wm.HighestZIndex))
				wm.IsResizingInit = true
				js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-selection.svg) 12 12, auto")
			}
			// Mouse down event for selecting
			// and dragging the window (click brings it to front)
			if !wm.IsResizingInit {
				wm.HighestZIndex++
				winElem.Get("style").Set("z-index", strconv.Itoa(wm.HighestZIndex))
				if wm.Verbose {
					wm.Print("Window brought to front.")
				}

				if wm.IsMovingMode && args[0].Get("button").Int() == 2 {
					args[0].Call("preventDefault")
					args[0].Call("stopPropagation")
					//JustSelected = true
					wm.StartX = args[0].Get("clientX").Float() - winElem.Get("offsetLeft").Float()
					wm.StartY = args[0].Get("clientY").Float() - winElem.Get("offsetTop").Float()
					wm.IsDragging = true
					// Create ghost window
					wm.GhostWindow = document.Call("createElement", "div")
					rect := winElem.Call("getBoundingClientRect")
					width := rect.Get("width").Float()
					height := rect.Get("height").Float()
					// Ensure ghost window is above everything during drag
					wm.GhostWindow.Set("style", "position: absolute; z-index: "+strconv.Itoa(wm.HighestZIndex+1)+"; width: "+wm.Ftoa(width)+"px; height: "+wm.Ftoa(height)+"px; border: solid 2px #FF0000; cursor: url(assets/cursor-drag.svg) 12 12, auto;")
					wm.GhostWindow.Get("style").Set("left", wm.Ftoa(winElem.Get("offsetLeft").Float())+"px")
					wm.GhostWindow.Get("style").Set("top", wm.Ftoa(winElem.Get("offsetTop").Float())+"px")
					body.Call("appendChild", wm.GhostWindow)
					wm.JustSelected = true
					if wm.Verbose {
						wm.Print("Dragging initiated with ghost window.")
					}
				}
				if wm.IsHiding && args[0].Get("button").Int() == 2 {
					// Hide window
					args[0].Call("preventDefault")
					args[0].Call("stopPropagation")
					wm.JustSelected = true
					wm.IsHiding = false

					hiddenWindowOption := context_menu.NewMenuItem("wid " + strconv.Itoa(neuwindow.ID))
					if winElem.Get("title").String() != "" {
						hiddenWindowOption = context_menu.NewMenuItem(winElem.Get("title").String())
					}

					hiddenWindowOption.Set("id", "menuopt"+strconv.Itoa(neuwindow.ID))
					// Unhide option activation
					hiddenWindowOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						if args[0].Get("button").Int() == 2 {
							args[0].Call("preventDefault")
							args[0].Call("stopPropagation")
							wm.JustSelected = true
							context_menu.DeleteMenuItem(hiddenWindowOption)
							winElem.Get("style").Set("display", "block")
							wm.ContextMenu.Get("style").Set("display", "none")
							// Delete by value
							for index, value := range wm.ContextMenuHides {
								if value.Get("id").String() == hiddenWindowOption.Get("id").String() {
									wm.ContextMenuHides = append(wm.ContextMenuHides[:index], wm.ContextMenuHides[index+1:]...)
								}
							}
							wm.JustSelected = false
							if wm.Verbose {
								wm.Print("Unhide activated.")
							}
						}
						return nil
					}))
					wm.ContextMenuHides = append(wm.ContextMenuHides, hiddenWindowOption)
					winElem.Get("style").Set("display", "none")
					js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
					wm.JustSelected = false
					if wm.Verbose {
						wm.Print("WID " + strconv.Itoa(neuwindow.ID) + " hidden")
					}

				}
			}
			// Right-click (RMB) deletes the window in delete mode
			if wm.IsDeleteMode && args[0].Get("button").Int() == 2 {
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")
				wm.JustSelected = true
				Delete(neuwindow)
				wm.IsDeleteMode = false
				js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
				wm.JustSelected = false

				if wm.Verbose {
					wm.Print("Window deleted.")
				}
			}
			return nil
		}))

	return neuwindow
}

// Move
// Sets pos and dimensions for the window.
// (Actually useless). {So, then f___ you}.
func (w *controls.Window) Move(newX, newY, newWidth, newHeight string) {
	w.Element.Get("style").Set("left", newX)
	w.Element.Get("style").Set("top", newY)
	w.Element.Get("style").Set("width", newWidth)
	w.Element.Get("style").Set("height", newHeight)
}

// Delete
// Deletes the window from DOM and Go.
func Delete(wnd *controls.Window) {
	// huh... C++ в жопе заиграл.
	wnd.ID = -1
	wnd.Element.Call("remove")
	delete(wm.AllWindows, strconv.Itoa(wnd.ID))
	return
}
