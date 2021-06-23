package main

type tedstate struct {
	Where     WH
	Objects   []node
	Pos       XY
	focus     int
	hold      int
	holdcode  int
	prev      int
	doingconn bool
	start     XY
	end       XY
	NewBox    button
}

func newbrect(winsize WH) XYWH {
	m := 4
	w := 24
	h := 24
	x, y := winsize.W-w-m, winsize.H-h-m
	return Rect(x, y, w, h)
}

func (t *tedstate) Draw() {
	G.SetDrawColor(colx(FieldColor))
	G.Clear()
	for i := range t.Objects {
		// connections
		G.SetDrawColor(colx(BoxBorderColor))
		if o := (*t.Objects[i].Inlet()); o != nil {
			ou := outletpos(*o.Rect()).Center()
			in := inletpos(*t.Objects[i].Rect()).Center()
			G.DrawLine(int32(ou.X), int32(ou.Y), int32(in.X), int32(in.Y))
		}
	}
	for i := range t.Objects {
		if i == t.focus {
			G.SetDrawColor(colx(0x0000ffff))
			G.FillRect(t.Objects[t.focus].Rect().Extrude(-2).ToSDL())
		}
		t.Objects[i].Draw()
	}
	// ugly, but will work
	if t.hold < 0 {
		*t.NewBox.Rect() = newbrect(t.Where)
		t.NewBox.Draw()
	}
	if t.doingconn {
		G.DrawLine(int32(t.start.X), int32(t.start.Y), int32(t.end.X), int32(t.end.Y))
	}
}

func (t *tedstate) over(at XY) int {
	for i := len(t.Objects) - 1; i >= 0; i-- {
		e := t.Objects[i]
		if e.Rect().Inside(at) {
			return i
		}
	}
	return -1
}

var str string
var ostr string

func connect(out node, in node) {
	(*in.Outlets())[out] = struct{}{}
	disconnect(out)
	*out.Inlet() = in
}

func disconnect(o node) {
	if *o.Inlet() != nil {
		if _, has := (*(*o.Inlet()).Outlets())[o]; has {
			delete(*((*o.Inlet()).Outlets()), o)
		}
		*o.Inlet() = nil
	}
}

func (t *tedstate) Mouse(at XY, buttons, delta int) {

}

func (t *tedstate) TextInput(r rune) {
	if t.focus >= 0 {
		t.Objects[t.focus].TextInput(r)
	}
}
