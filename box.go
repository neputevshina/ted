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
}

func newbuf(where XYWH) *buf {
	b := &buf{
		Where:   where,
		outlets: make(map[node]struct{}),
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

type cmd struct {
	Where  XYWH
	inlet  node
	Entry  *TedText
	outlet map[node]struct{}
	errlet map[node]struct{}
	Cmd    []rune
}

func newcmd(where XYWH) *cmd {
	c := &cmd{
		Where:  where,
		outlet: make(map[node]struct{}),
		errlet: make(map[node]struct{}),
		Cmd:    make([]rune, 0),
	}
	c.Entry = NewTedText(&c.Cmd, G, gcache, true, false)
	lo, up := c.Limits()
	if lo.W > up.W {
		up = lo
	}
	c.Where.W, c.Where.H = up.Val()
	return c
}

func linemeasure(text []rune, font *FontSprites) (px int) {
	for _, r := range text {
		if r == '\n' {
			break
		}
		px += font.Cache[r].m.Advance
	}
	return
}

func (c *cmd) Inlet() *node {
	return &c.inlet
}

func (c *cmd) Limits() (WH, WH) {
	lim := Wt(linemeasure(c.Cmd, c.Entry.Sprites)+2*BoxKnobsSize, 3*FontSize/2)
	return Wt(72, 3*FontSize/2), lim
}

func (c *cmd) Outlets() *map[node]struct{} {
	return &c.outlet
}

func (c *cmd) Draw() {
	xy := c.Where
	G.SetDrawColor(colx(BoxBorderColor))
	G.FillRect(xy.ToSDL())
	G.SetDrawColor(colx(BoxBgColor))
	G.FillRect(xy.Extrude(1).ToSDL())
	G.SetDrawColor(colx(InletAppendColor))
	G.FillRect(inletpos(xy).ToSDL())
	G.SetDrawColor(colx(BoxBorderColor))
	// outlet
	G.FillRect(outletpos(xy).ToSDL())
	// knob
	G.FillRect(knobpos(xy).ToSDL())
	c.Entry.Draw()
}

func (c *cmd) Mouse(at XY, buttons int, delta int) int {
	// todo: ugly
	x, y, w, h := c.Where.Val()
	c.Entry.Where = Rect(x+BoxKnobsSize, y, w-BoxKnobsSize*2, h)
	if knobpos(c.Where).Inside(at) {
		if buttons == MouseLeft {
			return MoveMe
		}
	}
	if inletpos(c.Where).Inside(at) {
		return OverInlet
	}
	if outletpos(c.Where).Inside(at) {
		return OverOutlet
	}
	c.Entry.Mouse(at, buttons, delta)
	return 0
}

func (c *cmd) TextInput(r rune) {
	c.Entry.TextInput(r)
	lo, up := c.Limits()
	if lo.W > up.W {
		up = lo
	}
	c.Where.W, c.Where.H = up.Val()
}

func (c *cmd) Rect() *XYWH {
	return &c.Where
}
