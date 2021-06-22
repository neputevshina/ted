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

type Bang struct{}

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
	R           *sdl.Renderer //
	Font        *ttf.Font     //
	SpriteCache *SpriteCache  //
	Where       XYWH          // Position on screen
	Limit       bool          // Cull text to its hitbox
	Oneliner    bool          // Display only first line and ignore newlines input
	Text        *[]rune       //
	Color       uint32        //
	Sel         [2]int        // Selection
	addlater    rune          //
	Wakeup      chan Bang     // Bang here wakes box up from drawing stasis
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
	for _, r := range *e.Text {
		//e.SpriteCache.Update(r)
		if r == '\n' {
			// proper oneliner won't have \ns
			if e.Oneliner {
				return
			}
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
	esel := e.Sel
	if esel[1] < esel[0] {
		esel[0], esel[1] = esel[1], esel[0]
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

	if esel[1] >= len(*e.Text) {
		esel[1] = len(*e.Text)
	}

	inside := false
	for i, r := range *e.Text {
		// cache is already there, skip
		if i == esel[0] {
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
		if i == esel[1] {
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

	if esel[1] == len(*e.Text) {
		e.R.SetDrawColor(colx(0x000000ff))
		e.R.FillRect(Rect(line.X+line.W-1, line.Y, 1, line.H).ToSDL())
	}

}

// Draw paints object to the screen
func (e *TedText) Draw() {
	//if e.dirty { // TODO
	// sooooo, with this shitty »hack« we will update only incoming runes! yay!!
	// and this speeds up almost nothing.
	// fucking add blitting cache, for god's sake
	e.SpriteCache.Update(e.addlater)
	e.mylittletypesetter()
	//e.dirty = false
	//}
	e.paintsel()
}

// actually, this func is part of the style
func measline(font *ttf.Font, at XY) int {
	zero := at.Y
	fh := font.Height()
	rem := 0
	// if zero%fh > fh/2 {
	// 	rem = 1
	// }
	return zero/fh + rem
}

func (e *TedText) measchar(at XY) (j int) {
	where := *e.Text
	glyphs := e.SpriteCache.Cache

	// skip lines
	atline := measline(e.Font, at)
	if e.Oneliner {
		atline = 0
	}

	for i, r := range where {
		if atline == 0 {
			j = i
			break
		}
		if r == '\n' {
			atline--
		}
	}
	zero := at.X //.Move(at).X

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

// Rect as in Drawer
func (e *TedText) Rect() *XYWH {
	return &e.Where
}

// Limits as in Drawer
func (e *TedText) Limits() (lower WH, upper WH) {
	return Wt(e.Font.Height(), e.Font.Height()), Wt(-1, -1)
}

// Mouse as in Drawer
func (e *TedText) Mouse(at XY, buttons int, delta int) int {
	at = at.Move(At(-e.Where.X, -e.Where.Y))
	measure := e.measchar(at)
	if buttons != 0 {
		if buttons == delta {
			e.Sel[0] = measure
			e.Sel[1] = measure
		}
		if delta == 0 {
			if measure <= e.Sel[0] {
				e.Sel[0] = measure
			} else {
				e.Sel[1] = measure
			}

		}
	}
	return 1
}

func (e *TedText) TextInput(b [32]byte) {
	r, _ := utf8.DecodeRune(b[:])
	t := *e.Text
	e.addlater = r
	if e.Oneliner && r == '\n' {
		return
	}
	if e.Sel[0] > e.Sel[1] {
		e.Sel[0], e.Sel[1] = e.Sel[1], e.Sel[0]
	}
	x0, x1 := e.Sel[0], e.Sel[1]
	if r == '\x7f' {
		if x0 == x1 {
			x0++
			x1++
		}
		r = '\b'
	}
	if r == '\b' {

		if x0 == x1 && x0 > 0 {
			t = append(t[:x0-1], t[x0:]...)
			e.Sel[0] = x0 - 1
		} else {
			t = append(t[:x0], t[x1:]...)
		}
	} else {
		t = append(t[:x0], append([]rune{r}, t[x1:]...)...)
		e.Sel[0]++
	}
	e.Sel[1] = e.Sel[0]
	*e.Text = t
}
