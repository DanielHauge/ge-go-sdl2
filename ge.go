package ge_go_sdl2

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	running bool
)

func Run(gui []View) {
	running = true

	registerViews(gui)
	sdlInit()
	defer destroyUI()
	renderGUI(gui[0])
	go handlePropertyChanges()
	handleEvents()
}

func sdlInit() {
	fmt.Println("Initializing SDL2")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	if err := ttf.Init(); err != nil {
		panic(err)
	}
}

func destroyUI() {
	running = false
	destroy()
	sdl.Quit()
	ttf.Quit()
}
