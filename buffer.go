package main

import "io"

const (
	inappend = iota
	inupdate

	inmodes
)

var _ = node(&buf{})

type buf struct {
	Where XYWH

	inletmode int
	inlet     node
	outlets   map[node]struct{}
	sellets   map[node]struct{}

	in  io.ReadCloser
	out io.WriteCloser
	sel io.WriteCloser

	Text      []rune
	Entry     *TedText
	Scrollpos uint
}

func newbuf(where XYWH) *buf {
	b := &buf{
		Where:   where,
		outlets: make(map[node]struct{}),
		sellets: make(map[node]struct{}),
		Text:    make([]rune, 0, 100),
	}
	e := NewTedText(&b.Text, G, gcache, false, false)
	b.Entry = e
	x, y, w, h := where.Val()
	e.Where = Rect(x, y+BoxKnobsSize, w, h-BoxKnobsSize*2)
	return b
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

func boxdraw(b node) {
	xy := *b.Rect()
	G.SetDrawColor(colx(BoxBorderColor))
	G.FillRect(xy.ToSDL())
	G.SetDrawColor(colx(BoxBgColor))
	G.FillRect(xy.Extrude(1).ToSDL())
	// inlet
	switch b := b.(type) {
	case *buf:
		switch b.inletmode {
		case inappend:
			G.SetDrawColor(colx(InletAppendColor))
		case inupdate:
			G.SetDrawColor(colx(InletUpdateColor))
		}
	default:
		G.SetDrawColor(colx(BoxBorderColor))
	}
	G.FillRect(inletpos(xy).ToSDL())
	G.SetDrawColor(colx(BoxBorderColor))
	// outlet
	G.FillRect(outletpos(xy).ToSDL())
	// knob
	G.FillRect(knobpos(xy).ToSDL())
	// killer
	if ted.killmode && ted.ov == b {
		G.SetDrawColor(colx(0xff0000ff))
	}
	G.FillRect(killerpos(xy).ToSDL())
}

func (b *buf) Draw() {
	boxdraw(b)
	b.Entry.Draw()
}

func (b *buf) Mouse(at XY, buttons int, delta int) int {
	// todo: ugly
	x, y, w, h := b.Where.Val()
	b.Entry.Where = Rect(x+BoxKnobsSize, y+BoxKnobsSize, w-BoxKnobsSize*2, h-BoxKnobsSize*2)
	if knobpos(b.Where).Inside(at) {
		if buttons == MouseLeft {
			return MoveMe
		}
		if buttons == MouseRight {
			return ResizeMe
		}
	}
	if killerpos(b.Where).Inside(at) {
		return OverKiller
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
	b.Entry.Mouse(at, buttons, delta)
	return 0
}

func (b *buf) TextInput(r rune) {
	b.Entry.TextInput(r)
}

func (b *buf) Rect() *XYWH {
	return &b.Where
}

func (b *buf) Input() *io.ReadCloser {
	return &b.in
}

func (b *buf) Primary() *io.WriteCloser {
	return &b.out
}

func (b *buf) Secondary() *io.WriteCloser {
	return &b.sel
}

func (b *buf) Errlets() *map[node]struct{} {
	return &b.sellets
}
