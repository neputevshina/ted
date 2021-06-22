package main

const (
	inappend = iota
	inupdate

	inmodes
)

type bufer struct {
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

func (b *bufer) Inlet() *node {
	return &b.inlet
}

func (b *bufer) Limits() (WH, WH) {
	return Wt(32, 32), Wt(-1, -1)
}

func (b *bufer) Outlets() *map[node]struct{} {
	return &b.outlets
}

func (b *bufer) Draw() {
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

func (b *bufer) Mouse(at XY, buttons int, delta int) int {
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

func (b *bufer) TextInput(text [32]byte) {
	b.Entry.TextInput(text)
}

func (b *bufer) Rect() *XYWH {
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
