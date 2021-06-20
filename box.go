package main

const (
	inappend = iota
	inupdate

	inmodes
)

type box struct {
	Where     XYWH
	hold      bool
	inlet     node
	inletmode int
	outlets   map[node]struct{}
	Lines     []string
	Scrollpos uint
	Cursor    [2]uint
}

func (b *box) Inlet() *node {
	return &b.inlet
}

func (b *box) Limits() (WH, WH) {
	return Wt(32, 32), Wt(-1, -1)
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
	// inlet
	switch b.inletmode {
	case inappend:
		G.SetDrawColor(colx(InletAppendColor))
	case inupdate:
		G.SetDrawColor(colx(InletUpdateColor))
	}
	G.FillRect(inletpos(xy).ToSDL())
	G.SetDrawColor(colx(BoxBorderColor))
	// outlet
	G.FillRect(outletpos(xy).ToSDL())
	// knob
	G.FillRect(knobpos(xy).ToSDL())
}

func (b *box) Mouse(at XY, buttons int, delta int) int {
	if knobpos(b.Where).Inside(at) {
		return OverKnob
	}
	if inletpos(b.Where).Inside(at) {
		if buttons == MouseLeft && delta != 0 {
			println("yes")
			b.inletmode = (b.inletmode + 1) % inmodes
			b.hold = true
			return 0
		}
		return OverInlet
	}
	if outletpos(b.Where).Inside(at) {
		return OverOutlet
	}
	b.hold = false
	return 0
}

// text input in modern oses is broken: why not just use ascii BS when you
// have to erase a symbol?
func (b *box) TextInput(text [32]byte) {

}

func (b *box) Rect() *XYWH {
	return &b.Where
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
