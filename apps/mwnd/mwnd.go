package mwnd

import (
	"riwo/apps"
	"riwo/pkg/constructor"
	"riwo/pkg/wm"
)

func init() {
	apps.AppRegistry["AppMyWindow"] = AppMyWindow
}

// AppMyWindow
// Describes handlers and other backend of application
// and Hyper-text markup registration
func AppMyWindow(window *wm.Window) {
	mainContainerStyles := constructor.WindowContainerStyles{
		Height:          "300",
		Display:         "flex",
		Justify:         "center",
		Alignment:       "center",
		BackgroundColor: "faded",
	}

	mainContainer := constructor.WindowContainer{
		Styles: &mainContainerStyles,
	}

	parameters := constructor.WindowParameters{
		Theme:     "aqua", // builtin
		Container: &mainContainer,
	}

	parameters.Render(window)
}
