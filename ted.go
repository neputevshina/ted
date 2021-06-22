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
	// super ugly, but i don't know how to implement it other way
	// also maybe todo: pretty patchcord rendering a la max
	for i := range t.Objects {
		if i == t.focus {
			G.SetDrawColor(colx(0x0000ffff))
			G.FillRect(t.Objects[t.focus].Rect().Extrude(-2).ToSDL())
		}
		// debug
		// G.SetDrawColor(colx(0xff0000ff))
		// for l := range *t.Objects[i].Outlets() {
		// 	in := inletpos(*l.Rect()).Xy().Move(At(2, 2))
		// 	out := outletpos(*t.Objects[i].Rect()).Xy().Move(At(2, 2))
		// 	G.DrawLine(int32(in.X), int32(in.Y), int32(out.X), int32(out.Y))
		// }
		//
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

func (t *tedstate) outlets(at XY, buttons, delta int, over, code int) {
	if buttons == delta && buttons == MouseLeft &&
		code == OverOutlet {
		t.doingconn = true
		t.start = at
	}
	if buttons == delta && buttons == MouseRight &&
		code == OverOutlet {
		disconnect(t.Objects[over])
	}
	if t.doingconn {
		if buttons == MouseLeft && delta == 0 {
			t.end = at
		}
		if buttons == 0 && delta == MouseLeft {
			if t.holdcode == OverOutlet && code == OverInlet && over >= 0 {
				connect(t.Objects[t.hold], t.Objects[over])
				t.doingconn = false
			}
			t.hold = -1
			t.holdcode = 0
		}
	}
}

func (t *tedstate) move(at XY, buttons, delta int, over, code int) {
	if buttons == delta && buttons != 0 {
		if code == MoveMe || code == ResizeMe {
			t.start = at
		}
	}
	if buttons != 0 && delta == 0 && t.hold >= 0 {
		dif := Dxy(at, t.start)
		r := t.Objects[t.hold].Rect()
		if t.holdcode == MoveMe {
			r.X += dif.X
			r.Y += dif.Y
			t.start = at
		}
		if t.holdcode == ResizeMe {
			l, u := t.Objects[t.hold].Limits()
			if r.W+dif.X <= l.W || r.W+dif.X >= u.W {
				dif.X = 0
			}
			if r.H+dif.Y <= l.H || r.H+dif.Y >= u.H {
				dif.Y = 0
			}
			r.W += dif.X
			r.H += dif.Y
			t.start = at
		}
		return
	}
}

func (t *tedstate) mousefield(at XY, buttons, delta int) {
	over := t.over(at)
	code := 0
	if over >= 0 {
		code = t.Objects[over].Mouse(at, buttons, delta)
	}

	if buttons == delta && buttons != 0 {
		if over >= 0 {
			t.Objects = append(
				append(t.Objects[:over], t.Objects[over+1:]...),
				t.Objects[over],
			)
			t.hold = len(t.Objects) - 1
			t.holdcode = code
		}
	}
	t.move(at, buttons, delta, over, code)
	t.outlets(at, buttons, delta, over, code)
	if buttons == 0 && delta != 0 {

	}
	// click with the...

	// hold...
	// release over...

	// over := false
	// for i := len(t.Objects) - 1; i >= 0; i-- {
	// 	e := t.Objects[i]
	// 	if e.Rect().Inside(at) {
	// 		over = true
	// 		t.over = i
	// 		if t.hold >= 0 {
	// 			break
	// 		}
	// 		t.holdcode = e.Mouse(at, buttons, delta)
	// 		t.focus = i
	// 		t.start = at
	// 		if t.hold < 0 && buttons != 0 {
	// 			// to the top
	// 			t.Objects = append(append(t.Objects[:i], t.Objects[i+1:]...), e)
	// 			t.focus = len(t.Objects) - 1
	// 			t.hold = (t.focus)
	// 		}
	// 		break
	// 	}
	// }
	// if !over && t.hold < 0 {
	// 	t.focus = -1
	// }
	// if t.focus >= 0 && t.focus == t.hold {
	// 	if t.holdcode == OverKnob {
	// 		rc := t.Objects[t.hold].Rect()
	// 		if buttons == MouseLeft {
	// 			rc.X += at.X - t.start.X
	// 			rc.Y += at.Y - t.start.Y
	// 			t.start = at
	// 		}
	// 		if buttons == MouseRight {
	// 			e := rc
	// 			dw := at.X - t.start.X
	// 			dh := at.Y - t.start.Y
	// 			e.W += dw
	// 			e.H += dh
	// 			t.start = at
	// 		}
	// 	}
	// 	if t.holdcode == OverOutlet && delta != 0 {
	// 		t.doesconn = true
	// 	}

	// 	o := t.Objects[t.over]
	// 	hc2 := o.Mouse(at, buttons, delta)

	// 	if hc2 == OverInlet {
	// 		if t.doesconn {
	// 			h := t.Objects[t.hold]
	// 			if delta == MouseLeft && buttons == 0 {
	// 				connect(h, o)
	// 				t.doesconn = false
	// 			}
	// 		} else if delta == 0 && buttons == MouseRight {
	// 			delinlet(o)
	// 		}
	// 	}
	// }
	// if buttons == 0 {
	// 	t.hold = -1
	// 	t.doesconn = false
	// }
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

func (t *tedstate) Mouse(at XY, buttons, delta int) {
	if t.NewBox.Rect().Inside(at) && t.hold < 0 {
		if buttons == MouseLeft && delta != 0 {
			t.Objects = append(t.Objects, &buf{
				Where:   Rect(at.X-100+4, at.Y-100+4, 100, 100),
				outlets: make(map[node]struct{}, 10),
			})
			t.hold = len(t.Objects) - 1
			t.holdcode = MoveMe
			t.start = at
		}
	} else {
		t.mousefield(at, buttons, delta)
	}
}

func (t *tedstate) TextInput(r rune) {

}
