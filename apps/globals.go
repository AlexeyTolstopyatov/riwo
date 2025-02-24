/*
App globals, primarily AppRegistry
*/

package apps

import (
	"riwo/internal/controls"
)

// AppRegistry holds all available app functions.
// Each app should register itself (typically in its init function).
var AppRegistry = make(map[string]func(*controls.Window))
