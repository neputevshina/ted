package main

import (
	"unicode/utf8"

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

// Sprite contains information about rendered rune
type Sprite struct {
	t  *sdl.Texture
	m  ttf.GlyphMetrics
	cl sdl.Rect
	// todo: add kerning cache
}

// SpriteCache is a storage of prerendered glyphs of certain size and color.
// Don't even try to update cache not in drawing context.
// todo: multicolored?
type SpriteCache struct {
	Color sdl.Color
	R     *sdl.Renderer
	Font  *ttf.Font
	Cache map[rune]Sprite
}

// NewSpriteCache returns ready-to-work glyph cache.
func NewSpriteCache(r *sdl.Renderer, f *ttf.Font, c sdl.Color) *SpriteCache {
	return &SpriteCache{
		Color: c,
		R:     r,
		Font:  f,
		Cache: make(map[rune]Sprite, 128),
	}
}

// TedText is an elastic tabstop text box
type TedText struct {
	R           *sdl.Renderer
	Font        *ttf.Font
	SpriteCache *SpriteCache
	Where       XYWH
	Limit       bool
	Text        []rune
	Color       uint32
	Selection   [2]int
	//colors [][]uint
	//tabs [][]uint
}

const (
	TextNewlineWidth = 6 // px
)

// func (e *TedText) solvetabs() {

// }

func (s *SpriteCache) Generate(text []rune) {
	for _, r := range text {
		s.Update(r)
	}
}

func (sc *SpriteCache) Update(r rune) {
	if _, k := sc.Cache[r]; !k {
		m, err := sc.Font.GlyphMetrics(r)
		if err != nil {
			panic(err)
		}
		s, err := sc.Font.RenderGlyphBlended(r, sc.Color)
		if err != nil {
			panic(err)
		}
		t, err := sc.R.CreateTextureFromSurface(s)
		if err != nil {
			panic(err)
		}
		cl := s.ClipRect
		sc.Cache[r] = Sprite{t, *m, cl}
	}
}

func (e *TedText) mylittletypesetter() {
	f := e.Font
	characc := 0
	lineacc := 0

	rcache := e.SpriteCache.Cache
	for _, r := range e.Text {
		e.SpriteCache.Update(r)
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
			(FromSDL(cl).Move(e.Where).Move(At(characc, lineacc).Wh(0, 0))).ToSDL(),
		)
		characc += rcache[r].m.Advance
	}
}

func (e *TedText) paintsel() {
	if e.Selection[1] < e.Selection[0] {
		e.Selection[0], e.Selection[1] = e.Selection[1], e.Selection[0]
	}

	fh := e.Font.Height()
	fl := e.Font.LineSkip()
	characc := 0
	lineacc := 0

	var oldbm sdl.BlendMode
	_ = e.R.GetDrawBlendMode(&oldbm)
	e.R.SetDrawBlendMode(sdl.BLENDMODE_MOD)
	defer e.R.SetDrawBlendMode(oldbm)

	var line XYWH

	if e.Selection[1] >= len(e.Text) {
		e.Selection[1] = len(e.Text)
	}

	inside := false
	for i, r := range e.Text {
		// cache is already there, skip
		if i == e.Selection[0] {
			inside = true
			line = Rect(characc, lineacc, 0, fh).Move(e.Where)

			e.R.SetDrawColor(colx(0x000000ff))
			e.R.FillRect(Rect(line.X, line.Y, 1, line.H).ToSDL())

			characc = 0
		}
		if r == '\n' {
			characc += TextNewlineWidth
			line.W = characc
			if inside {
				e.R.SetDrawColor(colx(0xffff00ff))
				e.R.FillRect(line.ToSDL())
			}
			lineacc += fl
			characc = 0
			line = Rect(characc, lineacc, 0, fh).Move(e.Where)
			continue
		}
		if i == e.Selection[1] {
			inside = false
			line.W = characc
			e.R.SetDrawColor(colx(0xffff00ff))
			e.R.FillRect(line.ToSDL())

			e.R.SetDrawColor(colx(0x000000ff))
			e.R.FillRect(Rect(line.X+line.W, line.Y, 1, line.H).ToSDL())
			// last time
			characc += e.SpriteCache.Cache[r].m.Advance

			break
		}
		characc += e.SpriteCache.Cache[r].m.Advance
	}

	if e.Selection[1] == len(e.Text) {
		e.R.SetDrawColor(colx(0x000000ff))
		e.R.FillRect(Rect(line.X+line.W-1, line.Y, 1, line.H).ToSDL())
	}

}

// Draw paints object to the screen
func (e *TedText) Draw() {
	e.mylittletypesetter()
	e.paintsel()
}

// actually, this func is part of the style
func measline(font *ttf.Font, base XYWH, at XY) int {
	zero := base.Xy().Move(at).Y
	fh := font.Height()
	rem := 0
	// if zero%fh > fh/2 {
	// 	rem = 1
	// }
	return zero/fh + rem
}

func measchar(where []rune, glyphs map[rune]Sprite, font *ttf.Font, base XYWH, at XY) (j int) {
	// skip lines
	atline := measline(font, base, at)
	for i, r := range where {
		if atline == 0 {
			j = i
			break
		}
		if r == '\n' {
			atline--
		}
	}
	zero := base.Xy().Move(at).X

	zero -= glyphs[where[0]].m.Advance / 2
	characc := 0
	last := 0
	for i, r := range where {
		abs := zero - characc
		curr := glyphs[r].m.Advance

		if r == '\n' {
			// we need a way to touch the newline symbol
			curr = TextNewlineWidth
		}
		if abs <= last {
			j = i + 1
			return
		}
		if i < j {
			continue
		}
		last = curr
		characc += last
	}
	// not found, return last
	return len(where) - 1
}

// Mouse as in Drawer
func (e *TedText) Mouse(at XY, buttons int, delta int) int {
	//at = at.Move(At(-e.Where.X, -e.Where.Y))
	measure := measchar(e.Text, e.SpriteCache.Cache, e.Font, e.Where, at)
	if buttons == MouseLeft && delta == MouseLeft {
		e.Selection[0] = measure
	}
	if buttons == MouseLeft && delta == 0 &&
		(measure >= e.Selection[0] || measure <= e.Selection[1]) {
		e.Selection[1] = measure
	}
	if buttons == 0 && delta == MouseLeft {
		if measure == e.Selection[0] {
			e.Selection[1] = e.Selection[0]
		}
	}
	return 0
}

func (e *TedText) TextInput(b [32]byte) {
	r, _ := utf8.DecodeRune(b[:])
	x0, x1 := e.Selection[0], e.Selection[1]
	e.Text = append(e.Text[:x0], append([]rune{r}, e.Text[x1:]...)...)
	e.Selection[0]++
	e.Selection[1] = e.Selection[0]
}