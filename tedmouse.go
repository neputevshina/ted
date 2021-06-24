package main

func (t *tedstate) hit(at XY, top bool) drawer {
	for i := len(t.Objects) - 1; i >= 0; i-- {
		e := t.Objects[i]
		if e.Rect().Inside(at) {
			// fix this
			// if top {
			// 	serch := func(n node) int {
			// 		for i, m := range t.Objects {
			// 			if n == m {
			// 				return i
			// 			}
			// 		}
			// 		return len(t.Objects)
			// 	}
			// 	ap := func(i int) {
			// 		t.Objects = append(
			// 			append(t.Objects[:i], t.Objects[i+1:]...),
			// 			t.Objects[i],
			// 		)
			// 	}
			// 	ovi := serch(t.ov.(node))
			// 	ap(ovi)
			// 	holdi := serch(t.hold.(node))
			// 	ap(holdi)
			// 	ap(i)
			// 	return t.Objects[len(t.Objects)-1]
			// }
			return e
		}
	}
	if t.NewBox.Rect().Inside(at) {
		return &t.NewBox
	}
	return nil
}

func (t *tedstate) Mouse(at XY, buttons, delta int) {
	t.ov = t.hit(at, buttons == delta && delta != 0 && t.hold == nil)
	if t.ov != nil {
		t.code = t.ov.Mouse(at, buttons, delta)
	}
	// todo: fsa-based mouse input parser; see proton, pike's squeak
	if buttons == delta && delta != 0 { // press
		if t.code == MoveMe || t.code == ResizeMe {
			t.start = at
			t.end = at
			t.hold = t.ov
			t.hcode = t.code
		}
		if t.code == OverOutlet {
			if buttons == MouseLeft {
				t.start = at
				t.end = at
				t.hold = t.ov
				t.hcode = t.code
			}
			if buttons == MouseRight {
				switch ov := t.ov.(type) {
				case node:
					disconnect(ov)
				}
			}
		}
		if t.code == 10 {
			if t.hold == nil {
				t.Objects = append(t.Objects, newbuf(Rect(at.X-100+BoxKnobsSize/24, at.Y-100+BoxKnobsSize/2, 100, 100)))
				t.hold = t.Objects[len(t.Objects)-1]
				t.hcode = MoveMe
				t.start = at
			}
		}
		if t.code == 11 {
			if t.hold == nil {
				c := newcmd(XYWH{})
				lo, _ := c.Limits()
				c.Where = Rect(at.X-lo.W+BoxKnobsSize/2, at.Y-lo.H+BoxKnobsSize/2, lo.W, lo.H)
				t.Objects = append(t.Objects, c)
				t.hold = t.Objects[len(t.Objects)-1]
				t.hcode = MoveMe
				t.start = at
			}
		}

	} else if buttons != 0 && delta == 0 { // hold
		if t.hcode == MoveMe {
			if t.hold != nil {
				dif := Dxy(at, t.start)
				r := t.hold.Rect()
				r.X += dif.X
				r.Y += dif.Y
				t.start = at
			}
		}
		if t.hcode == ResizeMe {
			if t.hold != nil {
				dif := Dxy(at, t.start)
				r := t.hold.Rect()
				l, u := t.hold.Limits()
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
		}
		if t.hcode == OverOutlet {
			t.end = at
		}

	} else if buttons == 0 && delta != 0 { // release
		if t.hcode == OverOutlet {
			switch in := t.ov.(type) {
			case node:
				switch out := t.hold.(type) {
				case node:
					connect(out, in)
				}
			}
		}
		t.hold = nil
		t.hcode = 0
	}

}
