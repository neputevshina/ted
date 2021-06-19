package main

type button struct {
	Where XYWH
}

func (bt *button) Draw() {
	xy := bt.Where
	G.SetDrawColor(colx(BoxBorderColor))
	G.FillRect(xy.ToSDL())
	G.SetDrawColor(colx(BoxBgColor))
	G.FillRect(xy.Extrude(2).ToSDL())
}

func (bt *button) Mouse(at XY, buttons int) int {
	return buttons
}

func (bt *button) Rect() *XYWH {
	return &bt.Where
}
