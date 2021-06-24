package main

type tedstate struct {
	Where   WH
	Objects []node
	NewBox  button

	ov    drawer
	code  int
	hold  drawer
	hcode int

	start XY
	end   XY
}

func createted() {
	ted = tedstate{
		Where: Wt(800, 600),
		NewBox: button{
			PressLeft:  11,
			PressRight: 10,
		},
	}

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
		if t.Objects[i] == t.ov {
			G.SetDrawColor(colx(0x0000ffff))
			G.FillRect(t.ov.Rect().Extrude(-2).ToSDL())
		}
		G.SetClipRect(t.Objects[i].Rect().ToSDL())
		t.Objects[i].Draw()
		G.SetClipRect(nil)
	}
	// connections
	for i := range t.Objects {
		G.SetDrawColor(colx(BoxBorderColor))
		if o := (*t.Objects[i].Inlet()); o != nil {
			ou := outletpos(*o.Rect()).Center()
			in := inletpos(*t.Objects[i].Rect()).Center()
			G.DrawLine(int32(ou.X), int32(ou.Y), int32(in.X), int32(in.Y))
		}
	}
	// ugly, but will work
	if t.hold == nil {
		*t.NewBox.Rect() = newbrect(t.Where)
		t.NewBox.Draw()
	}
	if t.hcode == OverOutlet {
		G.DrawLine(int32(t.start.X), int32(t.start.Y), int32(t.end.X), int32(t.end.Y))
	}
}

func connect(out node, in node) {
	(*out.Outlets())[in] = struct{}{}
	disconnect(in)
	*in.Inlet() = out
}

func disconnect(o node) {
	if *o.Inlet() != nil {
		if _, has := (*(*o.Inlet()).Outlets())[o]; has {
			delete(*((*o.Inlet()).Outlets()), o)
		}
		*o.Inlet() = nil
	}
}

func (t *tedstate) TextInput(r rune) {
	if t.ov != nil {
		switch ov := t.ov.(type) {
		case node:
			ov.TextInput(r)
		}
	}
}
