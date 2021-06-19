package main

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var window *sdl.Window
var g *sdl.Renderer
var font *ttf.Font

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
	g, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(err)
	}
	if font, err = ttf.OpenFont(`./Go-Regular.ttf`, 10); err != nil {
		panic(err)
	}
	sdl.EnableScreenSaver()
}

func main() {
	go func() {
		//predraw()
		for {
			draw()
			g.Present()
			window.UpdateSurface()
		}
	}()
	//s := make(chan os.Signal, 0)
	//signal.Notify(s, os.Interrupt)
	eventloop()
	window.Destroy()
	sdl.Quit()
	os.Exit(0)
	//<-s
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

		switch e.(type) {
		case *sdl.KeyboardEvent:
			//keyboard(e.(*sdl.KeyboardEvent))
		case *sdl.MouseButtonEvent:
			//mouse(e.(*sdl.MouseButtonEvent))
		case *sdl.MouseMotionEvent:
			//mousemove(e.(*sdl.MouseMotionEvent))
		}

	}
}

func draw() {

}
