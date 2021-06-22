package main

import "github.com/veandco/go-sdl2/sdl"

// XY is a point
type XY struct {
	X, Y int
}

// Move adds values of y to the values of x
func (x XY) Move(y XY) XY {
	return At(x.X+y.X, x.Y+y.Y)
}

// Val explodes the struct into values
func (x XY) Val() (int, int) {
	return x.X, x.Y
}

func FromSDL(r sdl.Rect) XYWH {
	return Rect(int(r.X), int(r.Y), int(r.W), int(r.H))
}

// Val explodes the struct into values
func (w WH) Val() (int, int) {
	return w.W, w.H
}

// Wt adds W and H components to the XY to yield an XYWH
func (x XY) Wt(w int, h int) XYWH {
	return XYWH{x.X, x.Y, w, h}
}

func (x XY) Wh() WH {
	return WH{x.X, x.Y}
}

// Dxy returns difference of two xys
func Dxy(xy XY, at XY) XY {
	return At(xy.X-at.X, xy.Y-at.Y)
}

// WH is a point intended for storing width and height
type WH struct {
	W, H int
}

// XYWH is a rectangle
type XYWH struct {
	X, Y, W, H int
}

// Extrude creates a new rect by adding a margin to the all sides
func (x XYWH) Extrude(length int) XYWH {
	return Rect(x.X+length, x.Y+length, x.W-length*2, x.H-length*2)
}

func (xy XYWH) Center() XY {
	return At((2*xy.X+xy.W)/2, (2*xy.Y+xy.H)/2)
}

// Inside checks if point and rect intersects
func (x XYWH) Inside(xy XY) bool {
	return xy.X >= x.X &&
		xy.Y >= x.Y &&
		xy.X <= x.X+x.W &&
		xy.Y <= x.Y+x.H
}

// Move translates x by y
func (x XYWH) Move(y XY) XYWH {
	return Rect(x.X+y.X, x.Y+y.Y, x.W, x.H)
}

// ToSDL converts rect to an SDL-edible rect
func (x XYWH) ToSDL() *sdl.Rect {
	return &sdl.Rect{X: int32(x.X), Y: int32(x.Y), W: int32(x.W), H: int32(x.H)}
}

// Xy returns X and Y components of a rect
func (x XYWH) Xy() XY {
	return At(x.X, x.Y)
}

// Wh returns W and H components of a rect
func (x XYWH) Wh() WH {
	return Wt(x.W, x.H)
}

// Val expodes rect to four values
func (x XYWH) Val() (int, int, int, int) {
	return x.X, x.Y, x.W, x.H
}

// At is a constructor for XY
func At(x, y int) XY {
	return XY{X: x, Y: y}
}

func NoLimit() WH {
	const max = int(^uint(0) >> 1)
	return WH{max, max}
}

// Wt is a constructor for Wt
func Wt(w, h int) WH {
	return WH{W: w, H: h}
}

// Rect is a constructor for XYWH
func Rect(x, y, w, h int) XYWH {
	return XYWH{X: x, Y: y, W: w, H: h}
}
