package main

type box struct {
	Where     XYWH
	Inlet     inflow
	Outlet    []outflow
	Text      []string
	Scrollpos uint
}

func inletpos(xy XYWH) XYWH {
	return Rect(xy.X+4, xy.Y, 8, 3)
}

func outletpos(xy XYWH) XYWH {
	return Rect(xy.X+4, xy.Y+xy.H-3, 8, 3)
}

func knobpos(xy XYWH) XYWH {
	bx := BoxKnobsSize
	return Rect(xy.X+xy.W-bx, xy.Y+xy.H-bx, bx, bx)
}

// func errletpos(xy XYWH) XYWH {
// 	return Rect(xx.X+xx.H+4*2, xx.Y, xx.W, xx.H)
// }

type cmd struct {
	Where  XYWH
	Inlet  inflow
	Outlet []outflow
	Errlet []outflow
	Cmd    string
}
