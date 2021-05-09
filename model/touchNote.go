package model

import rl "github.com/gen2brain/raylib-go/raylib"

type TouchNote struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
	Color  rl.Color
	Key    int32
}
