package main

import (
	"fmt"
	ge_go_sdl2 "ge-go-sdl2"
	"github.com/veandco/go-sdl2/sdl"
)

type name struct {
	gg string
}

func main() {

	test()
}

func test() {
	fmt.Println("Setup test gui")

	btnLabel := ge_go_sdl2.Text{
		Label: "Button",
		Font:  "./asset/test.ttf",
		Size:  12,
	}

	label := ge_go_sdl2.Text{
		Label: "Test header",
		Font:  "./asset/test.ttf",
		Size:  14,
		X:     200,
		Y:     10,
	}

	onClickChannel := make(chan string)

	go func() {
		for {
			<-onClickChannel
			fmt.Println("Clicked")
		}
	}()

	var btn ge_go_sdl2.Button
	btn.X = 20
	btn.Y = 20
	btn.H = 20
	btn.W = 60
	btn.Id = "superBtn"
	btn.OnClick = onClickChannel
	btn.BorderColor = 0xffff0000
	btn.BgColor = 0xffffff00
	btn.Content = btnLabel

	txtField := ge_go_sdl2.TextField{
		X:           10,
		Y:           10,
		H:           80,
		W:           200,
		Id:          "txtField",
		Size:        12,
		Value:       "text field",
		BgColor:     0xffffffff,
		BorderColor: 0x00000000,
		Font:        "./asset/test.ttf",
	}

	var viewChildren []interface{}
	viewChildren = append(viewChildren, txtField)

	view := ge_go_sdl2.View{
		X:           20,
		Y:           150,
		H:           300,
		W:           600,
		Id:          "view",
		BgColor:     0xeeeeeeee,
		BorderColor: 0xffffff00,
		Children:    viewChildren,
	}

	var children []interface{}

	children = append(children, btn)
	children = append(children, label)
	children = append(children, view)

	var gui = ge_go_sdl2.View{
		Id:       "Title",
		X:        0,
		Y:        0,
		W:        800,
		H:        600,
		Children: children,
		BgColor:  0xffffffff,
	}

	ge_go_sdl2.Run(gui, nil, nil)
}

func example() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
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

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
