package main

const (
	BoxBorderColor = 0x000000ff
	BoxBgColor     = 0xfffaffff
	FieldColor     = 0xddddddff
)

const (
	BoxBorderWidth           = 1
	BoxKnobsSize             = 8
	AnywhereFocusBorderWidth = 2
	NodeHeight               = 16
)

const (
	MouseLeft   = 1 // sdl.ButtonLMask()
	MouseMiddle = 2
	MouseRight  = 4
)

const (
	_ = iota
	OverKnob
	OverInlet
	OverOutlet
)

func colx(i uint32) (uint8, uint8, uint8, uint8) {
	a := uint8(i >> 24 & 0xff)
	r := uint8(i >> 16 & 0xff)
	g := uint8(i >> 8 & 0xff)
	b := uint8(i & 0xff)
	return a, r, g, b
}
