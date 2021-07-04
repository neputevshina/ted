package main

import (
	"log"
	"os"
	"runtime/pprof"
	"unicode/utf8"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var window *sdl.Window

// G is the global window rendeder for an application
var G *sdl.Renderer
var ted tedstate

var gcache *FontSprites

// Gfont is the global application font
var Gfont *ttf.Font

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
	if Gfont, err = ttf.OpenFont(`./Go-Regular.ttf`, FontSize); err != nil {
		panic(err)
	}
	gcache = NewFontSprites(G, Gfont, rgba(0x000000ff))
	sdl.EnableScreenSaver()
	window.SetResizable(true)
	createted()
}

func main() {
	f, err := os.Create("ted.profile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	go func() {
		for {
			ted.Draw()
			G.Present()
			window.UpdateSurface()
		}
	}()
	// editor resizing, because fuck sdl
	sdl.AddEventWatchFunc(func(ee sdl.Event, d interface{}) bool {
		if e, k := ee.(*sdl.WindowEvent); k {
			if e.Event == sdl.WINDOWEVENT_RESIZED {
				w, _ := sdl.GetWindowFromID(e.WindowID)
				if w == d.(*sdl.Window) && w == window {
					ted.Where = Wt(int(e.Data1), int(e.Data2))
				}
			}
		}
		return false
	}, window)
	eventloop()
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

		if _, k := event.(*sdl.QuitEvent); k {
			pprof.StopCPUProfile()
			return
		}

		switch e := event.(type) {
		case *sdl.TextInputEvent:
			r, _ := utf8.DecodeRune(e.Text[:])
			ted.TextInput(r)

		case *sdl.MouseMotionEvent:
			s := int(e.State)
			moused = mouseprev ^ s
			ted.Mouse(At(int(e.X), int(e.Y)), s, moused)
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
				ted.Mouse(At(int(e.X), int(e.Y)), 0, ch)
			} else {
				ted.Mouse(At(int(e.X), int(e.Y)), ch, ch)
			}
			mouseprev = ch

		case *sdl.KeyboardEvent:
			if e.State == sdl.PRESSED {
				switch e.Keysym.Sym {
				case sdl.K_RETURN:
					ted.TextInput('\n')
				case sdl.K_BACKSPACE:
					ted.TextInput('\b')
				case sdl.K_DELETE:
					ted.TextInput('\x7f')
				case sdl.K_HOME:
					ted.TextInput('\x01')
				case sdl.K_END:
					ted.TextInput('\x05')
				case sdl.K_UP:
					ted.TextInput('\x11')
				case sdl.K_DOWN:
					ted.TextInput('\x12')
				case sdl.K_LEFT:
					ted.TextInput('\x13')
				case sdl.K_RIGHT:
					ted.TextInput('\x14')
				}
			}

		}

	}
}
