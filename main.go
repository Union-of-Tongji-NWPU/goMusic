package main

import (
	"awesomeProject1/model"
	"awesomeProject1/tool"
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitAudioDevice()
	rl.InitWindow(model.SCREEN_WIDTH, model.SCREEN_HEIGHT, "sample:music")
	//test
	tool.InitGame([]string{"7-years-lukas-graham.txt"})

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		tool.FlushGame()
		tool.DrawGame()
	}
	rl.CloseWindow()

}
