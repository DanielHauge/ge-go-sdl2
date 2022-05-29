package ge_go_sdl2

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	window *sdl.Window
)

func init() {

}

func renderGUI(wnd View) {
	var err error
	if wnd.X == 0 && wnd.Y == 0 {
		window, err = sdl.CreateWindow(wnd.Id, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, wnd.W, wnd.H, sdl.WINDOW_SHOWN)
	} else {
		window, err = sdl.CreateWindow(wnd.Id, wnd.X, wnd.Y, wnd.W, wnd.H, sdl.WINDOW_SHOWN)
	}
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: wnd.W, H: wnd.H}, wnd.BgColor)

	for _, child := range wnd.Children {
		switch t := child.(type) {
		case View:
			renderView(surface, t, false)
			break
		case TextField:
			renderTextField(surface, t, false)
			break
		case Button:
			renderButton(surface, t, false)
			break
		case Text:
			renderText(surface, t, false)
			break
		}
	}

	window.UpdateSurface()
}

func destroy() {
	window.Destroy()
}

func renderView(surface *sdl.Surface, view View, update bool) {
	rect := sdl.Rect{X: view.X, Y: view.Y, W: view.W, H: view.H}
	borderRect := sdl.Rect{X: view.X - 1, Y: view.Y - 1, W: view.W + 2, H: view.H + 2}
	surface.FillRect(&borderRect, view.BorderColor)
	surface.FillRect(&rect, view.BgColor)
	for _, child := range view.Children {
		switch t := child.(type) {
		case View:
			t.X += view.X
			t.Y += view.Y
			renderView(surface, t, false)
			break
		case TextField:
			t.X += view.X
			t.Y += view.Y
			renderTextField(surface, t, false)
			break
		case Button:
			t.X += view.X
			t.Y += view.Y
			renderButton(surface, t, false)
			break
		case Text:
			t.X += view.X
			t.Y += view.Y
			renderText(surface, t, false)
			break
		}
	}
	if update {
		window.UpdateSurface()
	}
}

func renderButton(surface *sdl.Surface, btn Button, update bool) {
	borderRect := sdl.Rect{X: btn.X, Y: btn.Y, W: btn.W, H: btn.H}
	innerRect := sdl.Rect{X: btn.X + 1, Y: btn.Y + 1, W: btn.W - 2, H: btn.H - 2}
	surface.FillRect(&borderRect, btn.BorderColor)
	surface.FillRect(&innerRect, btn.BgColor)

	btn.Content.X = btn.X + 2
	btn.Content.Y = btn.Y + 1

	clickListeners[btn.Id] = clickListener{boundary: borderRect, onClickChan: btn.OnClick}

	renderText(surface, btn.Content, false)

	if update {
		window.UpdateSurface()
	}
}

func renderText(surface *sdl.Surface, text Text, update bool) {
	textFont, err := ttf.OpenFont(text.Font, text.Size)
	if err != nil {
		panic(err.Error())
	}
	defer textFont.Close()

	label, err := textFont.RenderUTF8Blended(text.Label, text.TextColor)
	if err != nil {
		return
	}
	defer label.Free()

	err = label.Blit(nil, surface, &sdl.Rect{X: text.X, Y: text.Y, W: 0, H: 0})

	if update {
		window.UpdateSurface()
	}

}

func renderTextField(surface *sdl.Surface, txtField TextField, update bool) {

	renderThis := func() {
		borderRect := sdl.Rect{X: txtField.X, Y: txtField.Y, W: txtField.W, H: txtField.H}
		innerRect := sdl.Rect{X: txtField.X + 1, Y: txtField.Y + 1, W: txtField.W - 2, H: txtField.H - 2}
		surface.FillRect(&borderRect, txtField.BorderColor)
		surface.FillRect(&innerRect, txtField.BgColor)
		clickListeners[txtField.Id] = clickListener{boundary: borderRect, onClickChan: focusClick}
		txtLabel := Text{X: innerRect.X + 2, Y: innerRect.Y + 2, Size: txtField.Size, Font: txtField.Font, Label: txtField.Value, TextColor: txtField.TextColor}
		renderText(surface, txtLabel, false)
	}
	renderThis()

	onInput := make(chan string)
	focusListeners[txtField.Id] = focusListener{onInput: onInput}

	go func() {
		for {
			v := <-onInput
			fmt.Println("found: ", v)
			txtField.Value += v
			renderThis()
			window.UpdateSurface()
		}

	}()

	if update {
		window.UpdateSurface()
	}
}
