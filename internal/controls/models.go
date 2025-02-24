/*
 * Window struct for WM
 * and its functions + listeners
 */

package controls

// ContextEntry
// Define a named type
// for the context entry.
type ContextEntry struct {
	Name     string
	Callback func()
}
