package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// type inflow interface{}
// type outflow interface{}

type drawable interface {
	Draw()
	Mouse(at XY, buttons int, delta int) int
	Rect() *XYWH
	Limits() (lower WH, upper WH)
}

type node interface {
	drawable
	Inlet() *node
	// Play(in *io.PipeReader, out *io.PipeWriter, err *io.PipeWriter)
	Outlets() *map[node]struct{}
	TextInput(r rune)
}

func colx(i uint32) (uint8, uint8, uint8, uint8) {
	a := uint8(i >> 24 & 0xff)
	r := uint8(i >> 16 & 0xff)
	g := uint8(i >> 8 & 0xff)
	b := uint8(i & 0xff)
	return a, r, g, b
}

func rgba(i uint32) sdl.Color {
	r, g, b, a := colx(i)
	return sdl.Color{r, g, b, a}
}

// Mouse states
const (
	MouseLeft   = 1 // sdl.ButtonLMask()
	MouseMiddle = 2
	MouseRight  = 4
)

// drawer.Touch messages
const (
	_ = iota
	MoveMe
	ResizeMe
	OverInlet
	OverOutlet
	OverKiller
	// OverTextA
	// OverTextB
)
