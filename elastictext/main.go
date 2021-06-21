package main

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var window *sdl.Window

// G is a global window rendeder for an application
var G *sdl.Renderer
var et TedText

// Gfoint is a global application font
var Gfoint *ttf.Font

var gcache = make(map[rune]Glyph)

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
	if Gfoint, err = ttf.OpenFont(`./Go-Regular.ttf`, 12); err != nil {
		panic(err)
	}
	sdl.EnableScreenSaver()
	window.SetResizable(true)
	et = TedText{
		Text: []rune(`пошёл нахуй	шалава ебаная
хахахаха я наебал тебя

amogus


nigger
`),
		GlyphCache: gcache,
		Where:      Rect(0, 0, 800, 600),
		R:          G,
		Font:       Gfoint,
	}
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

func eventloop() {
	wait := func() sdl.Event {
		return sdl.WaitEventTimeout(1000)
		//return sdl.PollEvent()
	}
	var e sdl.Event
	for {
	wt:
		e = wait()
		if e == nil {
			goto wt
		}

		switch e.(type) {
		case *sdl.QuitEvent:
			return
		}

		switch j := e.(type) {
		case *sdl.WindowEvent:
			if j.Type == sdl.WINDOWEVENT_SIZE_CHANGED {
				println("A")
				et.Where = At(0, 0).Wh(int(j.Data1), int(j.Data2))
			}
		case *sdl.TextInputEvent:
			// ted.TextInput(j.Text)
		case *sdl.KeyboardEvent:
			//keyboard(e.(*sdl.KeyboardEvent))
		case *sdl.MouseMotionEvent:
			et.Mouse(At(int(j.X), int(j.Y)), int(j.State), 0)
		case *sdl.MouseButtonEvent:
			ch := 0
			switch j.Button {
			case sdl.BUTTON_LEFT:
				ch = MouseLeft
			case sdl.BUTTON_MIDDLE:
				ch = MouseMiddle
			case sdl.BUTTON_RIGHT:
				ch = MouseRight
			}
			if j.State == sdl.RELEASED {
				ch = 0
			}
			et.Mouse(At(int(j.X), int(j.Y)), ch, 0)
		}

	}
}
