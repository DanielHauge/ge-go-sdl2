package main

import (
	"fmt"
	ge_go_sdl2 "github.com/DanielHauge/ge-go-sdl2"
)

func main() {
	test()
}

func test() {
	fmt.Println("Setup test gui")

	btnLabel := ge_go_sdl2.Text{
		Id:      "btnLaben",
		Content: "Button",
		Font:    "./asset/test.ttf",
		Size:    12,
	}

	label := ge_go_sdl2.Text{
		Id:      "labelTest",
		Content: "Test header",
		Font:    "./asset/test.ttf",
		Size:    14,
		X:       200,
		Y:       10,
	}

	onClickChannel := make(chan string)
	onNewValueChannel := make(chan string)

	var btn ge_go_sdl2.Button
	btn.X = 20
	btn.Y = 20
	btn.H = 20
	btn.W = 60
	btn.Id = "superBtn"
	btn.OnClick = onClickChannel
	btn.BorderColor = 0xffff0000
	btn.BgColor = 0xffffff00
	btn.Content = "Button"
	btn.ContentLabel = btnLabel

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
		OnChanged:   onNewValueChannel,
	}

	gg := func() {
		fmt.Println("Clicked -> changing stuff")
		ge_go_sdl2.PropertyChangeChannel <- ge_go_sdl2.PropertyChange{Id: "labelTest", Name: "Content", Value: txtField.Value}
	}

	ge_go_sdl2.HandleAsCallback(onClickChannel, gg)

	go func() {
		for {
			newTxt := <-onNewValueChannel
			ge_go_sdl2.PropertyChangeChannel <- ge_go_sdl2.PropertyChange{Id: "labelTest", Name: "Content", Value: newTxt}
		}
	}()

	var viewChildren []interface{}
	viewChildren = append(viewChildren, txtField)

	view := ge_go_sdl2.View{
		X:           0,
		Y:           0,
		H:           300,
		W:           600,
		Id:          "view",
		BgColor:     0xeeeeeeee,
		BorderColor: 0xffffff00,
		Children:    viewChildren,
	}

	view2 := ge_go_sdl2.View{
		X:           0,
		Y:           0,
		H:           300,
		W:           600,
		Id:          "view2",
		BgColor:     0xeeeeeeee,
		BorderColor: 0xffffff00,
	}

	viewRef := ge_go_sdl2.Container{
		ViewId: "view",
		X:      20,
		Y:      150,
	}

	var children []interface{}

	children = append(children, btn)
	children = append(children, label)
	children = append(children, viewRef)

	var gui = ge_go_sdl2.View{
		Id:       "Title",
		X:        0,
		Y:        0,
		W:        800,
		H:        600,
		Children: children,
		BgColor:  0xffffffff,
	}

	var views []ge_go_sdl2.View
	views = append(views, gui)
	views = append(views, view)
	views = append(views, view2)

	ge_go_sdl2.Run(views)

}
