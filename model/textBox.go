package model

import rl "github.com/gen2brain/raylib-go/raylib"

type TextBox struct {
	X         float32
	Y         float32
	FontSize  int
	Text      string
	FontColor rl.Color
}
