package ge_go_sdl2

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type clickListener struct {
	boundary    sdl.Rect
	onClickChan chan<- int
}

type keyListener struct {
	key     sdl.Keysym
	onPress chan<- int
}

type position struct {
	x int32
	y int32
}

var (
	clickListeners map[string]clickListener
	lastPress      position
	lastRelease    position
	pressed        bool
)

func init() {
	clickListeners = make(map[string]clickListener)
}

func handleEvents() {
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.MouseButtonEvent:
				handleMouseEvent(t)
				break
			case *sdl.KeyboardEvent:
				handleKeyboardEvent(t)
				break
			}
		}
	}
}

func handleMouseEvent(t *sdl.MouseButtonEvent) {
	if t.State == sdl.PRESSED {
		if t.Which == 1 && !pressed {
			lastPress = position{x: t.X, y: t.Y}
			pressed = true
		}
	} else {
		pressed = false
		lastRelease = position{x: t.X, y: t.Y}
		go handleClick()
	}
}

func handleKeyboardEvent(t *sdl.KeyboardEvent) {
	keyCode := t.Keysym.Sym
	keys := ""

	// Modifier keys
	switch t.Keysym.Mod {
	case sdl.KMOD_LALT:
		keys += "Left Alt"
	case sdl.KMOD_LCTRL:
		keys += "Left Control"
	case sdl.KMOD_LSHIFT:
		keys += "Left Shift"
	case sdl.KMOD_LGUI:
		keys += "Left Meta or Windows key"
	case sdl.KMOD_RALT:
		keys += "Right Alt"
	case sdl.KMOD_RCTRL:
		keys += "Right Control"
	case sdl.KMOD_RSHIFT:
		keys += "Right Shift"
	case sdl.KMOD_RGUI:
		keys += "Right Meta or Windows key"
	case sdl.KMOD_NUM:
		keys += "Num Lock"
	case sdl.KMOD_CAPS:
		keys += "Caps Lock"
	case sdl.KMOD_MODE:
		keys += "AltGr Key"
	}

	if keyCode < 10000 {
		if keys != "" {
			keys += " + "
		}

		// If the key is held down, this will fire
		if t.Repeat > 0 {
			keys += string(keyCode) + " repeating"
		} else {
			if t.State == sdl.RELEASED {
				keys += string(keyCode) + " released"
			} else if t.State == sdl.PRESSED {
				keys += string(keyCode) + " pressed"
			}
		}

	}

	if keys != "" {
		fmt.Println(keys)
	}
}

func handleClick() {
	for _, listener := range clickListeners {
		if isWithinBoundary(&listener.boundary, &lastRelease) && isWithinBoundary(&listener.boundary, &lastPress) {
			listener.onClickChan <- 1
		}
	}
}

func isWithinBoundary(rect *sdl.Rect, pos *position) bool {
	if pos.x < rect.X ||
		pos.y < rect.Y ||
		pos.x > rect.X+rect.W ||
		pos.y > rect.Y+rect.H {
		return false
	}
	return true
}

func onKeyPress() {
	// Let UI Elements in focus react to key presses.
}
