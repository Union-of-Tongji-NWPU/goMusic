package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"musicDance/model"
	"musicDance/tool"
)

func main() {
	rl.InitAudioDevice()
	rl.InitWindow(model.SCREEN_WIDTH, model.SCREEN_HEIGHT, model.WINDOW_TITLE)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		tool.UpdateDrawFrame()
	}
	rl.CloseWindow()

}
