package main

import (
	"os"
	"unicode/utf8"

	_ "embed"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var window *sdl.Window

// G is a global window rendeder for an application
var G *sdl.Renderer
var et TedText

//go:embed field.go
var fieldgo []byte

// Gfont is a global application font
var Gfont *ttf.Font

var gcache *SpriteCache

func conv(a []byte) (u []rune) {
	u = make([]rune, 0, len(a)/2)
	for len(a) > 0 {
		r, s := utf8.DecodeRune(a)
		u = append(u, r)
		a = a[s:]
	}
	return
}

func init() {
	var err error
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	window, err = sdl.CreateWindow("ted", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	G, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(err)
	}
	if Gfont, err = ttf.OpenFont(`./Go-Regular.ttf`, 12); err != nil {
		panic(err)
	}
	gcache = NewSpriteCache(G, Gfont, rgba(0x000000ff))
	sdl.EnableScreenSaver()
	window.SetResizable(true)
	et = TedText{
		Text:        conv(fieldgo),
		Sel:         [2]int{9, 12},
		SpriteCache: gcache,
		Where:       Rect(20, 20, 800-20, 600-20),
		R:           G,
		Font:        Gfont,
		addlater:    'a',
		//dirty:       true,
	}
	gcache.Generate(et.Text)
}

func main() {
	go func() {
		for {
			G.SetDrawColor(colx(0xffffffff))
			G.Clear()
			et.Draw()
			G.Present()
			window.UpdateSurface()
			//fmt.Println(len(gcache))
		}
	}()
	// editor resizing, because fuck sdl
	sdl.AddEventWatchFunc(func(ee sdl.Event, d interface{}) bool {
		switch e := ee.(type) {
		case *sdl.WindowEvent:
			if e.Event == sdl.WINDOWEVENT_RESIZED {
				w, _ := sdl.GetWindowFromID(e.WindowID)
				if w == d.(*sdl.Window) && w == window {
					et.Where = At(0, 0).Wh(int(e.Data1), int(e.Data2))
				}
			}
		}
		return false
	}, window)
	eventloop()
	// for _, g := range gcache {
	// 	g.t.Destroy()
	// }
	window.Destroy()
	sdl.Quit()
	os.Exit(0)
}

var mouseprev int
var moused int

func eventloop() {
	wait := func() sdl.Event {
		return sdl.WaitEventTimeout(1000)
		//return sdl.PollEvent()
	}
	var event sdl.Event
	for {
	wt:
		event = wait()
		if event == nil {
			goto wt
		}

		switch event.(type) {
		case *sdl.QuitEvent:
			return
		}

		switch e := event.(type) {
		case *sdl.WindowEvent:
			if e.Type == sdl.WINDOWEVENT_SIZE_CHANGED {
				println("A")
				et.Where = At(0, 0).Wh(int(e.Data1), int(e.Data2))
			}
		case *sdl.TextInputEvent:
			et.TextInput(e.Text)
		case *sdl.MouseMotionEvent:
			s := int(e.State)
			moused = mouseprev ^ s
			et.Mouse(At(int(e.X), int(e.Y)), s, moused)
			mouseprev = int(e.State)
		case *sdl.MouseButtonEvent:
			ch := 0
			switch e.Button {
			case sdl.BUTTON_LEFT:
				ch = MouseLeft
			case sdl.BUTTON_MIDDLE:
				ch = MouseMiddle
			case sdl.BUTTON_RIGHT:
				ch = MouseRight
			}
			if e.State == sdl.RELEASED {
				et.Mouse(At(int(e.X), int(e.Y)), 0, ch)
			} else {
				et.Mouse(At(int(e.X), int(e.Y)), ch, ch)
			}
		case *sdl.KeyboardEvent:
			if e.State == sdl.PRESSED {
				switch e.Keysym.Sym {
				case sdl.K_RETURN:
					et.TextInput([32]byte{'\n'})
				case sdl.K_BACKSPACE:
					et.TextInput([32]byte{'\b'})
				}
			}
		}

	}
}
