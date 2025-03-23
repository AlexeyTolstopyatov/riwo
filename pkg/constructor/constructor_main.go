package constructor

import (
	"riwo/pkg/wm"
	"syscall/js"
)

// WindowContainer
// represents main <div> element
type WindowContainer struct {
	Styles  *WindowContainerStyles
	Content js.Value
}

type WindowContainerStyles struct {
	Width           string
	Height          string
	Display         string
	Justify         string
	Alignment       string
	BackgroundColor string
}

type WindowParameters struct {
	Container *WindowContainer
	Theme     string
}

// Add
// Appends new HTML element in page
func (wp *WindowParameters) Add(wnd *wm.Window) *WindowParameters {
	// make parent element
	container := js.Global().Call("createElement", "div")

	// join to black parade
	wnd.Element.Call("appendChild", container)
	return wp
}

// Render
// Is an Entry Point of Window constructor
// Based on den0602 code it draws common parts
// of Window using js calls.
//
// Theme used by default: Monochrome
func (wp *WindowParameters) Render(wnd *wm.Window) {

	if len(wp.Theme) == 0 {
		wp.Theme = "monochrome"
	}
	document := js.Global().Get("document")

	container := document.Call("createElement", "div")
	container.Get("style").Set("height", wp.Container.Styles.Height)
	container.Get("style").Set("display", wp.Container.Styles.Display)
	container.Get("style").Set("flexDirection", "column")
	container.Get("style").Set("justifyContent", wp.Container.Styles.Justify)
	container.Get("style").Set("alignItems", wp.Container.Styles.Alignment)
	container.Get("style").Set("backgroundColor", wm.GetColor[wp.Theme]["faded"])

	wnd.Element.Set("innerHTML", "")
	wnd.Element.Call("appendChild", container)
}
