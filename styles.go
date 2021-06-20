package main

// Colors
const (
	BoxBorderColor   = 0x000000ff
	BoxBgColor       = 0xfffaffff
	FieldColor       = 0xddddddff
	InletAppendColor = 0x000000ff
	InletUpdateColor = 0x0000ffff
)

// Sizes
const (
	BoxBorderWidth           = 1
	BoxKnobsSize             = 8
	AnywhereFocusBorderWidth = 2
	NodeHeight               = 16
)

func colx(i uint32) (uint8, uint8, uint8, uint8) {
	a := uint8(i >> 24 & 0xff)
	r := uint8(i >> 16 & 0xff)
	g := uint8(i >> 8 & 0xff)
	b := uint8(i & 0xff)
	return a, r, g, b
}

const winout = 8
const hinout = 8

func inletpos(xy XYWH) XYWH {
	//return Rect(xy.X+4, xy.Y, 8, 3)
	return Rect(xy.X, xy.Y, winout, hinout)
}

func outletpos(xy XYWH) XYWH {
	//return Rect(xy.X+4, xy.Y+xy.H-3, 8, 3)
	return Rect(xy.X, xy.Y+xy.H-hinout, winout, hinout)
}

func knobpos(xy XYWH) XYWH {
	bx := BoxKnobsSize
	return Rect(xy.X+xy.W-bx, xy.Y+xy.H-bx, bx, bx)
}
