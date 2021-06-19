package main

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var window *sdl.Window

// G is a global window rendeder for an application
var G *sdl.Renderer
var ted tedstate

// var font *ttf.Font

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
	// if font, err = ttf.OpenFont(`./Go-Regular.ttf`, 10); err != nil {
	// 	panic(err)
	// }
	sdl.EnableScreenSaver()
	window.SetResizable(true)
	ted = tedstate{
		Winsize: Wt(800, 600),
		focus:   -1,
		hold:    -1,
		NewBox:  button{},
	}
}

func main() {
	go func() {
		for {
			ted.Draw()
			G.Present()
			window.UpdateSurface()
		}
	}()
	// editor resizing, because fuck sdl
	sdl.AddEventWatchFunc(func(ee sdl.Event, d interface{}) bool {
		switch e := ee.(type) {
		case *sdl.WindowEvent:
			if e.Event == sdl.WINDOWEVENT_RESIZED {
				w, _ := sdl.GetWindowFromID(e.WindowID)
				if w == d.(*sdl.Window) && w == window {
					ted.Winsize = Wt(int(e.Data1), int(e.Data2))
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
				ted.Winsize = Wt(int(j.Data1), int(j.Data2))
			}
		case *sdl.TextInputEvent:
			// ted.TextInput(j.Text)
		case *sdl.KeyboardEvent:
			//keyboard(e.(*sdl.KeyboardEvent))
		case *sdl.MouseMotionEvent:
			ted.Mouse(At(int(j.X), int(j.Y)), int(j.State))
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
			ted.Mouse(At(int(j.X), int(j.Y)), ch)
		}

	}
}
