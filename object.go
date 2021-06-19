package main

import "fmt"

type inflow interface{}
type outflow interface{}

type node interface {
	Draw()
	Mouse(at XY, buttons int)
	Rect() *XYWH
}

type tedstate struct {
	Winsize WH
	Objects []node
	Pos     XY
	focus   int
	hold    int
	start   XY
}

func (t *tedstate) Draw() {
	G.SetDrawColor(colx(FieldColor))
	G.Clear()
	for i := len(t.Objects) - 1; i >= 0; i-- {
		t.Objects[i].Draw()
	}
}

func (t *tedstate) Mouse(at XY, buttons int) {
	fmt.Println(at, buttons)
	//*
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
			// switch j := e.Widget.(type) {
			// case Toucher:
			// 	a.tag = j.Touch(At(at.X+frame.X, at.Y+frame.Y), which, how, e.XYWH)
			// 	if how == MousePress {
			// 		a.Top(e.Widget)
			// 		a.focus = len(a.store) - 1
			// 		a.hold = a.focus
			// 	}
			// }
			// break
		}
	}
	if !over && t.hold < 0 {
		t.focus = -1
	}
	if t.focus >= 0 && t.focus == t.hold {
		//if a.tag == "move" {
		// if how == MousePress {

		// }
		rc := t.Objects[t.hold].Rect()
		if buttons == MouseLeft {
			rc.X += at.X - t.start.X
			rc.Y += at.Y - t.start.Y
			t.start = at
		}
		if buttons == MouseRight {
			e := rc
			// wh := e.Widget.Constraint()
			dw := at.X - t.start.X
			dh := at.Y - t.start.Y

			// if e.W+dw < wh.W {
			// 	e.W = wh.W
			// } else {
			e.W += dw
			// }
			// if e.H+dh < wh.H {
			// 	e.H = wh.H
			// } else {
			e.H += dh
			// }
			t.start = at
			//}
		}
		//}
	}
	if buttons == 0 {
		t.hold = -1
	}
	//*/
}
