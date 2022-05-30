package ge_go_sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
)

type uiElement interface {
	View | Container | Text | Button | TextField
}

type View struct {
	Id          string
	X           int32
	Y           int32
	W           int32
	H           int32
	BgColor     uint32
	BorderColor uint32
	Flags       uint
	Children    []interface{}
}

type Container struct {
	Id     string
	ViewId string
	X      int32
	Y      int32
}

func (e *Container) ChangeView(viewId string) {
	PropertyChangeChannel <- PropertyChange{Id: e.Id, Name: "ViewId", Value: viewId}
}

type Text struct {
	Content   string
	X         int32
	Y         int32
	Size      int
	Font      string
	TextColor sdl.Color
	Id        string
}

type Button struct {
	ContentLabel Text
	Content      string
	X            int32
	Y            int32
	W            int32
	H            int32
	BorderColor  uint32
	BgColor      uint32
	Id           string
	OnClick      chan string
}

func (e *Button) SetContent(content string) {
	PropertyChangeChannel <- PropertyChange{Id: e.Id, Name: "Content", Value: content}
}

type TextField struct {
	Value       string
	OnChanged   chan string
	X           int32
	Y           int32
	W           int32
	H           int32
	BorderColor uint32
	BgColor     uint32
	Size        int
	Font        string
	TextColor   sdl.Color
	Id          string
}

func (e *TextField) SetText(text string) {
	PropertyChangeChannel <- PropertyChange{Id: e.Id, Name: "Value", Value: text}
}

func (e *TextField) notifyChange() {
	valueToChange := e.Value
	go func() {
		e.OnChanged <- valueToChange
	}()
}

func HandleAsCallbackArg[T interface{}](channel chan T, callback func(T)) {
	go func() {
		for {
			callback(<-channel)
		}
	}()
}

func HandleAsCallback(channel chan string, callback func()) {
	go func() {
		for {
			<-channel
			callback()
		}
	}()
}
