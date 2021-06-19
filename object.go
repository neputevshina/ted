package main

type inflow interface{}
type outflow interface{}

type drawer interface {
	Draw()
	Mouse(at XY, buttons int) int
	Rect() *XYWH
}

type node interface {
	drawer
}

type tedstate struct {
	Winsize  WH
	Objects  []node
	Pos      XY
	focus    int
	hold     int
	holdcode int
	start    XY
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
	for i := len(t.Objects) - 1; i >= 0; i-- {
		t.Objects[i].Draw()
	}
	// ugly, but will work
	b := &t.NewBox
	*b.Rect() = newbrect(t.Winsize)
	b.Draw()
}

func mousefield(t *tedstate, at XY, buttons int) {
	over := false
	for i := len(t.Objects) - 1; i >= 0; i-- {
		e := t.Objects[i]
		if e.Rect().Inside(at) {
			over = true
			if t.hold >= 0 {
				break
			}
			t.focus = i
			t.start = at
			t.holdcode = e.Mouse(at, buttons)
			if t.hold <= 0 {
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
		if t.holdcode == 1 {
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
	}
	if buttons == 0 {
		t.hold = -1
	}
}

func (t *tedstate) Mouse(at XY, buttons int) {
	//fmt.Println(at, buttons)
	if t.NewBox.Rect().Inside(at) {
		t.NewBox.Mouse(at, buttons)
	} else {
		mousefield(t, at, buttons)
	}
}
