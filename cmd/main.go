package main

import (
	"riwo/apps"
	"riwo/wm"
	"strconv"
	"syscall/js"
)

// Logging
// Switching verbose mode and informates about in Console.
func Logging(this js.Value, args []js.Value) interface{} {
	wm.Verbose = !wm.Verbose
	if wm.Verbose {
		wm.Print("Logging is now ON")
	} else {
		wm.Print("Logging is now OFF")
	}
	return nil
}

// LaunchDefault
// Callback for window-manager. Settings up perferences by default :D
func LaunchDefault(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		// Arguments count more or less than one. (expected WID)
		return "Expected WID only!"
	}
	jsNum := args[0] // Get the js.Value argument

	if jsNum.Type() != js.TypeNumber { // Check if it's a number
		return "Argument must be a uint32"
	}
	num := jsNum.Int() // Convert js.Value to Go int

	fetchedWindow, ok := wm.AllWindows[strconv.Itoa(num)]
	if !ok {
		// I'm really not okay (trust me)
		if wm.Verbose {
			wm.Print("Couldn't start AppDefault on window " + strconv.Itoa(num))
		}
		return nil
	}
	apps.AppDefault(fetchedWindow)
	return nil
}

func main() {
	c := make(chan struct{})

	wm.Print(`
Great, You've found yourself in the console
Then you are likely to want to know this:
- Click LMB to cancel any action
- Press RMB to open context menu
- Select option by pressing RMB
- "New" will open another window after you
  make a selection with RMB
- Select state wants RMB click ("Delete", "Resize")
  or hold ("Move") on desired window
For logging there are:
+ Logging()
`)

	// Logging toggler
	js.Global().Set("Logging", js.FuncOf(Logging))

	wm.AllWindows = make(map[string]*wm.Window)
	wm.ContextMenuHides = make([]js.Value, 0)

	// Set default app for window
	js.Global().Set("LaunchDefault", js.FuncOf(LaunchDefault))
	// Essential for context menu's "New"

	// Window manager core
	wm.InitializeContextMenu()
	wm.InitializeGlobalMouseEvents()

	<-c
}
