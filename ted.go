package main

import "fmt"

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
	if t.NewBox.Rect().Inside(at) && t.hold < 0 {
		if buttons == MouseLeft && delta != 0 {
			t.Objects = append(t.Objects, newbuf(Rect(at.X-100+4, at.Y-100+4, 100, 100)))
			t.hold = len(t.Objects) - 1
			t.holdcode = MoveMe
			t.start = at
		}
		return
	}

	over := t.over(at)
	t.focus = over
	code := 0
	if over >= 0 {
		code = t.Objects[over].Mouse(at, buttons, delta)
	}

	if buttons == delta && buttons != 0 {
		t.start = at
		println("P")
		if over >= 0 && over != t.hold {
			t.Objects = append(
				append(t.Objects[:over], t.Objects[over+1:]...),
				t.Objects[over],
			)
			t.hold = len(t.Objects) - 1
			t.holdcode = code
		}
		if buttons == MouseLeft {
			if code == OverOutlet {
				t.doingconn = true
			}
		}
		if buttons == MouseRight {
			if code == OverOutlet {
				disconnect(t.Objects[over])
			}
		}
	}
	if buttons != 0 && delta == 0 && t.hold >= 0 {
		if t.doingconn {
			if buttons == MouseLeft {
				t.end = at
			}

		} else {
			dif := Dxy(at, t.start)
			fmt.Println(dif)
			r := t.Objects[t.hold].Rect()
			if t.holdcode == MoveMe {
				r.X += dif.X
				r.Y += dif.Y
				t.start = at
				fmt.Println(r)
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
				fmt.Println(r)
			}
		}
	}

	if buttons == 0 && delta != 0 {
		if delta == MouseLeft {
			if t.holdcode == OverOutlet && code == OverInlet && over >= 0 {
				connect(t.Objects[over], t.Objects[t.hold])
				t.doingconn = false
			}
		}
		t.hold = -1
		t.holdcode = 0
	}
}

func (t *tedstate) TextInput(r rune) {
	if t.focus >= 0 {
		t.Objects[t.focus].TextInput(r)
	}
}
