package ge_go_sdl2

import "github.com/veandco/go-sdl2/sdl"

type Window struct {
	Title    string
	X        int32
	Y        int32
	W        int32
	H        int32
	bgColor  uint32
	Flags    uint
	Children []interface{}
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

type Text struct {
	Label     string
	X         int32
	Y         int32
	Size      int
	Font      string
	TextColor sdl.Color
	Id        string
}

type Button struct {
	Content     Text
	X           int32
	Y           int32
	W           int32
	H           int32
	BorderColor uint32
	BgColor     uint32
	Id          string
	OnClick     chan int
}

type TextField struct {
	Value   string
	X       int32
	Y       int32
	W       int32
	H       int32
	Size    int
	Id      string
	OnInput chan int
}
