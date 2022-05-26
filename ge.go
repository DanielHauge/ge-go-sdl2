package ge_go_sdl2

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

var running = false

func Initialize() {

	log.Println("Initializing GUI")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	running = true
	go renderGui() // Keep running given a UI.
}

func DestroyUI() {
	running = false
}

func renderGui() {
	defer sdl.Quit()
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	rect := sdl.Rect{0, 0, 200, 200}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					fmt.Println("Mouse", t.Which, "button", t.Button, "pressed at", t.X, t.Y)
				} else {
					fmt.Println("Mouse", t.Which, "button", t.Button, "released at", t.X, t.Y)
				}
				// Handle Click
				break
			}
		}
	}

}
