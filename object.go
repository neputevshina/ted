package main

// type inflow interface{}
// type outflow interface{}

type drawer interface {
	Draw()
	Mouse(at XY, buttons int, delta int) int
	Rect() *XYWH
	Limits() (lower WH, upper WH)
}

type node interface {
	drawer
	Inlet() *node
	Outlets() *map[node]struct{}
}

const (
	MouseLeft   = 1 // sdl.ButtonLMask()
	MouseMiddle = 2
	MouseRight  = 4
)

// drawer.Touch messages
const (
	_ = iota
	OverKnob
	OverInlet
	OverOutlet
	// OverTextA
	// OverTextB
)
