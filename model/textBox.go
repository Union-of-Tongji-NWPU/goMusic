package model

import rl "github.com/gen2brain/raylib-go/raylib"

type TextBox struct {
	X         float32
	Y         float32
	FontSize  int
	Text      string
	FontColor rl.Color
}

type FrameRegisterText struct {
	Text *TextBox
	Dx   float64
	Dy   float64
}

type FrameRegisterAction struct {
	Function func(text interface{})
	Frame    int
	Data     interface{}
}

type SetTouchBlockColor struct {
	TouchNote *TouchNote
	Color     rl.Color
}
