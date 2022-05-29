package ge_go_sdl2

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

var running = false

func Initialize(gui Window, additionalViews []View) chan<- PropertyChange {
	fmt.Println("Initializing GUI")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	go handleEvents()
	RenderGUI(gui)
	running = true

	return nil
}

func DestroyUI() {
	running = false
}

func initializeGUI() {
	defer sdl.Quit()

}
