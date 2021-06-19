package main

type inflow interface{}
type outflow interface{}

type drawer interface {
	Draw()
	Mouse(at XY, buttons int) int
	Rect() *XYWH
}

type node interface {
	drawer
}
