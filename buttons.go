package main

type button struct {
	Where      XYWH
	PressLeft  int
	PressRight int
	Hold       int
	Release    int
}

func (bt *button) Draw() {
	xy := bt.Where
	G.SetDrawColor(colx(BoxBorderColor))
	G.FillRect(xy.ToSDL())
	G.SetDrawColor(colx(BoxBgColor))
	G.FillRect(xy.Extrude(2).ToSDL())
}

func (bt *button) Mouse(at XY, buttons, delta int) int {
	if buttons == delta {
		if delta == MouseLeft {
			return bt.PressLeft
		}
		if delta == MouseRight {
			return bt.PressRight
		}
	}
	if buttons != 0 && delta == 0 {
		return bt.Hold
	}
	if buttons == 0 && delta != 0 {
		return bt.Release
	}
	return 0
}

func (bt *button) Rect() *XYWH {
	return &bt.Where
}

func (bt *button) Limits() (l, u WH) {
	return Wt(0, 0), NoLimit()
}
