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
			PressLeft: 10,
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
		// connections
		G.SetDrawColor(colx(BoxBorderColor))
		if o := (*t.Objects[i].Inlet()); o != nil {
			ou := outletpos(*o.Rect()).Center()
			in := inletpos(*t.Objects[i].Rect()).Center()
			G.DrawLine(int32(ou.X), int32(ou.Y), int32(in.X), int32(in.Y))
		}
	}
	for i := range t.Objects {
		// if i == t.focus {
		// 	G.SetDrawColor(colx(0x0000ffff))
		// 	G.FillRect(t.Objects[t.focus].Rect().Extrude(-2).ToSDL())
		// }
		t.Objects[i].Draw()
	}
	// // ugly, but will work
	if t.hold == nil {
		*t.NewBox.Rect() = newbrect(t.Where)
		t.NewBox.Draw()
	}
	// if t.doingconn {
	// 	G.DrawLine(int32(t.start.X), int32(t.start.Y), int32(t.end.X), int32(t.end.Y))
	// }
}

func (t *tedstate) hit(at XY) drawer {
	for i := len(t.Objects) - 1; i >= 0; i-- {
		e := t.Objects[i]
		if e.Rect().Inside(at) {
			return e
		}
	}
	if t.NewBox.Rect().Inside(at) {
		return &t.NewBox
	}
	return nil
}

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

func (t *tedstate) TextInput(r rune) {
	if t.ov != nil {
		switch ov := t.ov.(type) {
		case node:
			ov.TextInput(r)
		}
	}
}
