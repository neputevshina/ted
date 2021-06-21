package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

/*
[Forwarded from иооператив "сучка" sticks and stones (sigma male grindset possessor)]
There's currently no method to retrieve the kerning for a pair of characters from SDL_ttf,
However correct kerning will be applied when a string of text is rendered instead of individual glyphs.
пидоры блядь.
а стринг рендерится в чём у нас???? правильно!!!! в utf-8. когда мне нужны руны.
сука, как же я ненавижу sdl
[Forwarded from иооператив "сучка" sticks and stones (sigma male grindset possessor)]
вообще я кажется щас изобрёл костыль
1 рендерим два символа как строку
2 рендерим два символа по отдельности
3 храним кеш разниц между двумя символами
4 когда надо нарисовать следующий символ, сверяем его с предыдущим, и применяем разницу
*/

/*
so.
type kerncache map[[2]rune]int
*/

// Glyph contains information about rendered rune
type Glyph struct {
	t  *sdl.Texture
	m  *ttf.GlyphMetrics
	cl sdl.Rect
}

// TedText is an elastic tabstop text box
type TedText struct {
	R          *sdl.Renderer
	Font       *ttf.Font
	GlyphCache map[rune]Glyph
	Where      XYWH
	Limit      bool
	Text       []rune
	Color      uint32
	//colors [][]uint
	//tabs [][]uint
}

// func (e *TedText) solvetabs() {

// }

func (e *TedText) mylittletypesetter() {
	f := e.Font
	characc := 0
	lineacc := f.LineSkip()

	rcache := e.GlyphCache
	for _, r := range e.Text {
		if _, k := rcache[r]; !k {
			m, err := f.GlyphMetrics(r)
			if err != nil {
				panic(err)
			}
			s, err := f.RenderGlyphBlended(r, rgba(e.Color))
			if err != nil {
				panic(err)
			}
			t, err := e.R.CreateTextureFromSurface(s)
			if err != nil {
				panic(err)
			}
			cl := s.ClipRect
			rcache[r] = Glyph{t, m, cl}
		}
		if r == '\n' {
			lineacc += f.LineSkip()
			characc = 0
			continue
		}
		t := rcache[r].t
		cl := rcache[r].cl
		e.R.Copy(
			t,
			&cl,
			(FromSDL(cl).Move(At(characc, lineacc).Wh(0, 0))).ToSDL(),
		)
		characc += rcache[r].m.Advance
	}
}

// Draw paints object to the screen
func (e *TedText) Draw() {
	e.mylittletypesetter()
}

// Mouse as in Drawer
func (e *TedText) Mouse(at XY, buttons int, delta int) int {
	return 0
}
