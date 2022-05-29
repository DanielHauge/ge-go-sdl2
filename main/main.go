package main

import (
	"fmt"
	"ge-go-sdl2"
)

func main() {
	fmt.Println("Setup test gui")
	var gui ge_go_sdl2.Window
	gui.X = 0
	gui.Y = 0
	gui.H = 600
	gui.W = 800
	gui.Title = "Test gui"

	btnLabel := ge_go_sdl2.Text{
		Label: "Button",
		Font:  "testfont",
		Size:  12,
	}

	onClickChannel := make(chan int)

	var btn ge_go_sdl2.Button
	btn.X = 20
	btn.Y = 20
	btn.H = 20
	btn.W = 60
	btn.Id = "superBtn"
	btn.OnClick = onClickChannel
	btn.BorderColor = 0xffff0000
	btn.BgColor = 0x000000000000
	btn.Content = btnLabel

	// ge_go_sdl2.Initialize(gui, nil)

	<-onClickChannel

}
