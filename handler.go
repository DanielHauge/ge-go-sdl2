package ge_go_sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
)

type clickListener struct {
	boundary    sdl.Rect
	onClickChan chan<- string
}

type focusListener struct {
	boundary sdl.Rect
	onInput  chan<- string
	onEdit   chan<- string
}

type position struct {
	x int32
	y int32
}

var (
	clickListeners map[string]clickListener
	focusListeners map[string]focusListener
	focusId        string
	focusClick     chan string
	focusInput     focusListener
	lastPress      position
	lastRelease    position
	pressed        bool
	inFocus        bool
)

func init() {
	clickListeners = make(map[string]clickListener)
	focusListeners = make(map[string]focusListener)

	focusClick = make(chan string)
}

func handleFocusCursor() {
	for running {
		focusId = <-focusClick
		value, ok := focusListeners[focusId]
		inFocus = ok
		if ok {
			focusInput = value
			sdl.StartTextInput()
		} else {
			sdl.StopTextInput()
		}
	}
}

func handleEvents() {
	go handleFocusCursor()
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
			case *sdl.TextInputEvent:
				if inFocus {
					text := t.GetText()
					focusInput.onInput <- text
				}
				break
			case *sdl.KeyboardEvent:
				handleKeyboardEvent(t)
				break
			default:
				break
			}
		}
	}
}

func handleMouseEvent(t *sdl.MouseButtonEvent) {
	if t.State == sdl.PRESSED {
		if t.Which == 0 && !pressed {
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
	if inFocus { // keycode = 8 -> backspace
		if keyCode == 8 {
			if t.Repeat > 0 {
				focusInput.onEdit <- "-"
			} else {
				if t.State == sdl.PRESSED {
					focusInput.onEdit <- "-"
				}
			}
		}
	}
}

func handleClick() {
	focusClick <- ""
	for id, listener := range clickListeners {
		if isWithinBoundary(&listener.boundary, &lastRelease) && isWithinBoundary(&listener.boundary, &lastPress) {
			listener.onClickChan <- id
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
