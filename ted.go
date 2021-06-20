package main

import "fmt"

type tedstate struct {
	Winsize  WH
	Objects  []node
	Pos      XY
	focus    int
	hold     int
	over     int
	holdcode int
	prev     int
	delta    int
	doesconn bool
	start    XY
	current  XY
	NewBox   button
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
		// super-ugly, but i don't know how to implement it other way
		// also maybe todo: pretty patchcord rendering a la max
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
		*t.NewBox.Rect() = newbrect(t.Winsize)
		t.NewBox.Draw()
	}
	if t.doesconn {
		G.DrawLine(int32(t.start.X), int32(t.start.Y), int32(t.current.X), int32(t.current.Y))
	}
}

func mousefield(t *tedstate, at XY, buttons int) {
	over := false
	for i := len(t.Objects) - 1; i >= 0; i-- {
		e := t.Objects[i]
		if e.Rect().Inside(at) {
			over = true
			t.over = i
			if t.hold >= 0 {
				break
			}
			t.holdcode = e.Mouse(at, buttons, t.delta)
			t.focus = i
			t.start = at
			if t.hold < 0 && buttons != 0 {
				// to the top
				t.Objects = append(append(t.Objects[:i], t.Objects[i+1:]...), e)
				t.focus = len(t.Objects) - 1
				t.hold = (t.focus)
			}
			break
		}
	}
	if !over && t.hold < 0 {
		t.focus = -1
	}
	if t.focus >= 0 && t.focus == t.hold {
		if t.holdcode == OverKnob {
			rc := t.Objects[t.hold].Rect()
			if buttons == MouseLeft {
				rc.X += at.X - t.start.X
				rc.Y += at.Y - t.start.Y
				t.start = at
			}
			if buttons == MouseRight {
				e := rc
				dw := at.X - t.start.X
				dh := at.Y - t.start.Y
				e.W += dw
				e.H += dh
				t.start = at
			}
		}
		if t.holdcode == OverOutlet && t.delta != 0 {
			t.doesconn = true
		}

		o := t.Objects[t.over]
		hc2 := o.Mouse(at, buttons, t.delta)

		if hc2 == OverInlet {
			if t.doesconn {
				h := t.Objects[t.hold]
				if t.delta == MouseLeft && buttons == 0 {
					connect(h, o)
					t.doesconn = false
				}
			} else if t.delta == 0 && buttons == MouseRight {
				delinlet(o)
			}
		}
	}
	if buttons == 0 {
		t.hold = -1
		t.doesconn = false
	}
}

func connect(out node, in node) {
	(*out.Outlets())[in] = struct{}{}
	delinlet(in)
	*in.Inlet() = out
}

func delinlet(o node) {
	if *o.Inlet() != nil {
		if _, has := (*(*o.Inlet()).Outlets())[o]; has {
			delete(*((*o.Inlet()).Outlets()), o)
		}
		*o.Inlet() = nil
	}
}

func (t *tedstate) Mouse(at XY, buttons int) {
	fmt.Println(at, buttons, t.prev, t.delta)
	t.current = at
	t.delta = t.prev ^ buttons
	if t.NewBox.Rect().Inside(at) && t.hold < 0 {
		if buttons == MouseLeft && t.delta != 0 {
			t.Objects = append(t.Objects, &box{
				Where:   Rect(at.X-100+4, at.Y-100+4, 100, 100),
				outlets: make(map[node]struct{}, 10),
			})
			t.start = at
			t.hold = len(t.Objects) - 1
			t.focus = t.hold
			t.holdcode = t.Objects[t.focus].Mouse(at, buttons, t.delta)
		}
	} else {
		mousefield(t, at, buttons)
	}
	t.prev = buttons
}
