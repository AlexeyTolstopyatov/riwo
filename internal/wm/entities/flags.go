package entities

// WindowManagerFlags
// describes all possible flags for
// checking system-state at the moment.
type WindowManagerFlags struct {
	IsDragging     bool
	IsMovingMode   bool
	IsResizingMode bool
	IsResizingInit bool
	JustSelected   bool
	IsDeleteMode   bool
	IsNewMode      bool
	IsHiding       bool
}
