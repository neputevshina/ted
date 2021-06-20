package main

type box struct {
	Where     XYWH
	inlet     node
	outlets   map[node]struct{}
	Lines     []string
	Scrollpos uint
	Cursor    [2]uint
}

func (b *box) Inlet() *node {
	return &b.inlet
}

func (b *box) Outlets() *map[node]struct{} {
	return &b.outlets
}

func (b *box) Draw() {
	xy := b.Where
	G.SetDrawColor(colx(BoxBorderColor))
	G.FillRect(xy.ToSDL())
	G.SetDrawColor(colx(BoxBgColor))
	G.FillRect(xy.Extrude(1).ToSDL())
	// inlet and outlet
	G.SetDrawColor(colx(BoxBorderColor))
	G.FillRect(inletpos(xy).ToSDL())
	G.FillRect(outletpos(xy).ToSDL())
	// knob
	G.FillRect(knobpos(xy).ToSDL())
}

func (b *box) Mouse(at XY, buttons int) int {
	if knobpos(b.Where).Inside(at) {
		return OverKnob
	}
	if inletpos(b.Where).Inside(at) {
		return OverInlet
	}
	if outletpos(b.Where).Inside(at) {
		return OverOutlet
	}
	return 0
}

// text input in modern oses is broken: why not just use ascii BS when you
// have to erase a symbol?
func (b *box) TextInput(text [32]byte) {

}

func (b *box) Rect() *XYWH {
	return &b.Where
}

const winout = 8
const hinout = 4

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

// func errletpos(xy XYWH) XYWH {
// 	return Rect(xx.X+xx.H+4*2, xx.Y, xx.W, xx.H)
// }

type cmd struct {
	Where  XYWH
	Inlet  node
	Outlet map[node]struct{}
	Errlet map[node]struct{}
	Cmd    string
}
