package main

const (
	inappend = iota
	inupdate

	inmodes
)

type buf struct {
	Where     XYWH
	hold      bool
	inlet     node
	inletmode int
	outlets   map[node]struct{}
	Text      []rune
	Entry     *TedText
	Scrollpos uint
	Cursor    [2]uint
}

func (b *buf) Inlet() *node {
	return &b.inlet
}

func (b *buf) Limits() (WH, WH) {
	return Wt(32, 32), NoLimit()
}

func (b *buf) Outlets() *map[node]struct{} {
	return &b.outlets
}

func (b *buf) Draw() {
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

func (b *buf) Mouse(at XY, buttons int, delta int) int {
	if knobpos(b.Where).Inside(at) {
		if buttons == MouseLeft {
			return MoveMe
		}
		if buttons == MouseRight {
			return ResizeMe
		}
	}
	if inletpos(b.Where).Inside(at) {
		if buttons == delta && buttons == MouseLeft {
			b.inletmode = (b.inletmode + 1) % inmodes
		}
		return OverInlet
	}
	if outletpos(b.Where).Inside(at) {
		return OverOutlet
	}
	return 0
}

func (b *buf) TextInput(r rune) {
	b.Entry.TextInput(r)
}

func (b *buf) Rect() *XYWH {
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
