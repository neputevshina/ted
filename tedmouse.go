package main

func (t *tedstate) Mouse(at XY, buttons, delta int) {
	t.ov = t.hit(at)
	if t.ov != nil {
		t.code = t.ov.Mouse(at, buttons, delta)
	}
	// todo: fsa mouse input parser; see proton, pike's squeak
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
				{
					t.Objects = append(t.Objects, newbuf(Rect(at.X-100+4, at.Y-100+4, 100, 100)))
					t.hold = t.Objects[len(t.Objects)-1]
					t.hcode = MoveMe
					t.start = at
				}
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
