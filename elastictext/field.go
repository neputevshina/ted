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

type glyph struct {
	s *sdl.Surface
	m *ttf.GlyphMetrics
}

// ElasticText is an elastic tabstop text box
type ElasticText struct {
	R    *sdl.Renderer
	Font *ttf.Font
	//GlyphCache map[rune]glyph
	Where XYWH
	Limit bool
	Text  []rune
	Color uint32
	//colors [][]uint
	//tabs [][]uint
}

// func (e *ElasticText) solvetabs() {

// }

// Draw is a my little typesetter
func (e *ElasticText) Draw() {
	f := e.Font
	characc := 0
	lineacc := f.LineSkip()

	// maybe global cache? but i don't think it will be necessary
	// in the case of single field. oh, i get it
	// todo.
	rcache := make(map[rune]glyph)
	defer func() {
		for _, c := range rcache {
			c.s.Free()
		}
	}()
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
			rcache[r] = glyph{s, m}
		}

		characc += rcache[r].m.Advance
		if r == '\n' {
			lineacc += f.LineSkip()
		}

	}

}

// Mouse as in Drawer
func (e *ElasticText) Mouse(at XY, buttons int, delta int) int {
	return 0
}
