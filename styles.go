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
	FontSize                 = 12
	BoxBorderWidth           = 1
	BoxKnobsSize             = 8
	AnywhereFocusBorderWidth = 2
	NodeHeight               = 16
	TextNewlineWidth         = 6 // px
)

const winout = BoxKnobsSize
const hinout = BoxKnobsSize

func inletpos(xy XYWH) XYWH {
	return Rect(xy.X, xy.Y, winout, hinout)
}

func outletpos(xy XYWH) XYWH {
	return Rect(xy.X, xy.Y+xy.H-hinout, winout, hinout)
}

func knobpos(xy XYWH) XYWH {
	bx := BoxKnobsSize
	return Rect(xy.X+xy.W-bx, xy.Y+xy.H-bx, bx, bx)
}
