package ge_go_sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	window *sdl.Window
	views  map[string]View
)

func registerViews(gui []View) {
	views = make(map[string]View)
	for _, v := range gui {
		views[v.Id] = v
	}
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
			renderView(surface, t)
			break
		case Container:
			renderContainer(surface, t)
			break
		case TextField:
			renderTextField(surface, t)
			break
		case Button:
			renderButton(surface, t)
			break
		case Text:
			renderText(surface, t)
			break
		}
	}

	window.UpdateSurface()
}

func destroy() {
	window.Destroy()
}

func renderContainer(surface *sdl.Surface, viewContainer Container) {

	renderThis := func(container *Container) {
		view := views[container.ViewId]
		view.X += container.X
		view.Y += container.Y
		renderView(surface, view)
	}
	renderThis(&viewContainer)

	updateChan := make(chan PropertyChange)
	updateChannels[viewContainer.Id] = updateChan
	go handleElementPropertyChanges(updateChan, &viewContainer, renderThis)

}

func renderView(surface *sdl.Surface, view View) {
	rect := sdl.Rect{X: view.X, Y: view.Y, W: view.W, H: view.H}
	borderRect := sdl.Rect{X: view.X - 1, Y: view.Y - 1, W: view.W + 2, H: view.H + 2}
	surface.FillRect(&borderRect, view.BorderColor)
	surface.FillRect(&rect, view.BgColor)
	for _, child := range view.Children {
		switch t := child.(type) {
		case View:
			t.X += view.X
			t.Y += view.Y
			renderView(surface, t)
			break
		case Container:
			t.X += view.X
			t.Y += view.Y
			renderContainer(surface, t)
			break
		case TextField:
			t.X += view.X
			t.Y += view.Y
			renderTextField(surface, t)
			break
		case Button:
			t.X += view.X
			t.Y += view.Y
			renderButton(surface, t)
			break
		case Text:
			t.X += view.X
			t.Y += view.Y
			renderText(surface, t)
			break
		}
	}

}

func renderButton(surface *sdl.Surface, btn Button) {

	renderThis := func(btn *Button) {
		borderRect := sdl.Rect{X: btn.X, Y: btn.Y, W: btn.W, H: btn.H}
		innerRect := sdl.Rect{X: btn.X + 1, Y: btn.Y + 1, W: btn.W - 2, H: btn.H - 2}
		surface.FillRect(&borderRect, btn.BorderColor)
		surface.FillRect(&innerRect, btn.BgColor)

		btn.ContentLabel.X = btn.X + 2
		btn.ContentLabel.Y = btn.Y + 1
		btn.ContentLabel.Content = btn.Content

		clickListeners[btn.Id] = clickListener{boundary: borderRect, onClickChan: btn.OnClick}

		renderText(surface, btn.ContentLabel)
	}
	renderThis(&btn)

	updateChan := make(chan PropertyChange)
	updateChannels[btn.Id] = updateChan
	go handleElementPropertyChanges(updateChan, &btn, renderThis)

}

func renderText(surface *sdl.Surface, txt Text) {
	red, green, blue, alpha := surface.At(int(txt.X), int(txt.Y)).RGBA()
	bgColor := sdl.MapRGBA(surface.Format, uint8(red), uint8(green), uint8(blue), uint8(alpha))
	var clearRect sdl.Rect
	clear := false

	renderThis := func(text *Text) {
		textFont, err := ttf.OpenFont(text.Font, text.Size)
		if err != nil {
			panic(err.Error())
		}
		defer textFont.Close()

		label, err := textFont.RenderUTF8Blended(text.Content, text.TextColor)
		if err != nil {
			return
		}
		defer label.Free()

		if clear {
			surface.FillRect(&clearRect, bgColor)
		}

		clearRect = sdl.Rect{X: text.X, Y: text.Y, W: label.ClipRect.W, H: label.ClipRect.H}
		clear = true
		err = label.Blit(nil, surface, &sdl.Rect{X: text.X, Y: text.Y, W: 0, H: 0})
	}

	renderThis(&txt)

	updateChan := make(chan PropertyChange)
	updateChannels[txt.Id] = updateChan
	go handleElementPropertyChanges(updateChan, &txt, renderThis)

}

func renderTextField(surface *sdl.Surface, txtField TextField) {

	renderThis := func() {
		borderRect := sdl.Rect{X: txtField.X, Y: txtField.Y, W: txtField.W, H: txtField.H}
		innerRect := sdl.Rect{X: txtField.X + 1, Y: txtField.Y + 1, W: txtField.W - 2, H: txtField.H - 2}
		surface.FillRect(&borderRect, txtField.BorderColor)
		surface.FillRect(&innerRect, txtField.BgColor)
		clickListeners[txtField.Id] = clickListener{boundary: borderRect, onClickChan: focusClick}
		txtLabel := Text{X: innerRect.X + 2, Y: innerRect.Y + 2, Size: txtField.Size, Font: txtField.Font, Content: txtField.Value, TextColor: txtField.TextColor}
		renderText(surface, txtLabel)
	}
	renderThis()

	onInput := make(chan string)
	onEdit := make(chan string)
	focusListeners[txtField.Id] = focusListener{onInput: onInput, onEdit: onEdit}

	go func() {
		for {
			select {
			case v := <-onInput:
				txtField.Value += v
			case v := <-onEdit:
				switch v {
				case "-":
					if len(txtField.Value) > 0 {
						txtField.Value = txtField.Value[0 : len(txtField.Value)-1]
					}
				}
			}
			renderThis()
			txtField.notifyChange()
			window.UpdateSurface()
		}
	}()

}
