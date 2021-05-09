/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package tool

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"io/ioutil"
	"math"
	"math/rand"
	"strings"

	"awesomeProject1/model"
)

var MusicSheetList = new(model.DoubleList)
var CurrentMusicSheet = new(model.DoubleNode)
var TouchNoteList = make([]model.TouchNote, model.LINE)
var MusicNoteList = make([]model.DoubleList, model.LINE)
var AnimationText = new(model.DoubleList)
var MissCount = 0
var ScoreSum = 0 //总分
var FrameCount = 0
var FramActionList = new(model.DoubleList)
var SuccessiveNum = 0

var CurrentScreen = 0 //当前界面
const Title = 0
const InGAME = 1
const SongBox = 2

var OptionSelect = 0

const ChooseStartGame = 0
const ChooseMusicBox = 1
const ChooseCredit = 2

func UpdateDrawFrame() {
	switch CurrentScreen {
	case Title:
		DrawMenu()
	case InGAME:
		FlushGame()
		DrawGame()
	case SongBox:
		DrawSongBox()
	}
}

func InitGame(sheetFiles []string) {
	for k, _ := range sheetFiles {
		sheet := new(model.MusicSheet)
		err := LoadMusicSheetFromFile(sheet, sheetFiles[k])
		if err != nil {
			panic(err)
		}
		MusicSheetList.Append(sheet)
	}
	if len(sheetFiles) == 0 {
		sheet := new(model.MusicSheet)
		GetDefaultMusicSheet(sheet)
		MusicSheetList.Append(sheet)
	}

	// 初始化ListHead
	CurrentMusicSheet = MusicSheetList.Head
	// 初始化按钮
	for i := range TouchNoteList {
		TouchNoteList[i].Width = model.TOUCH_BLOCK_WIDTH
		TouchNoteList[i].Height = model.TOUCH_BLOCK_HEIGHT
		TouchNoteList[i].X = float32(model.LEFT_MARGIN + i*(model.TOUCH_BLOCK_WIDTH+model.MARGIN_BETWEEN_LINE))
		TouchNoteList[i].Y = model.SCREEN_HEIGHT - model.TOUCH_BLOCK_HEIGHT - model.TOUCH_BLOCK_MARGIN_BOTTOM
		TouchNoteList[i].Color = model.TOUCH_BLOCK_BAD_COLOR
		TouchNoteList[i].Key = model.KeyboardKey[i]
	}
}

// 获取下一个音符到MusicNoteList中
func getYOfNote(minHeight float32) float32 {
	if minHeight > 0 {
		return -model.MUSIC_NOTE_HEIGHT
	} else {
		return minHeight - model.MUSIC_NOTE_HEIGHT
	}
}
func generateNextNote() {
	// 1.找到当前最上面的音符
	var minHeight = float32(0x3f3f3f3f)
	for i := range MusicNoteList {
		node := new(model.DoubleNode)
		node = MusicNoteList[i].Head
		for node != nil {
			if minHeight > node.Data.(model.MusicNote).Y {
				minHeight = node.Data.(model.MusicNote).Y
			}
			node = node.Prev
		}
	}
	// 2. 当前需要产生新的音符
	for minHeight > 0 || minHeight >= -model.MUSIC_NOTE_HEIGHT {
		selectLine := rand.Intn(model.LINE)
		musicNote := model.MusicNote{
			X:      float32(model.LEFT_MARGIN + selectLine*(model.LINE_WIDTH+model.MARGIN_BETWEEN_LINE)),
			Y:      getYOfNote(minHeight),
			Width:  model.MUSIC_NOTE_WIDTH,
			Height: model.MUSIC_NOTE_HEIGHT,
			Color:  model.MUSIC_NOTE_INIT_COLOR[rand.Intn(5)],
		}
		MusicNoteList[selectLine].Append(musicNote)
		minHeight = musicNote.Y
	}
}

func getMissMusicNote() {
	node := new(model.DoubleNode)
	for i := range MusicNoteList {
		node = MusicNoteList[i].Head
		for node != nil {
			musicNote := node.Data.(model.MusicNote)
			// 超过即Miss
			if musicNote.Y > model.SCREEN_HEIGHT {
				MusicNoteList[i].Delete(node)
				MissCount += 1
				SuccessiveNum = 0
				RegisterAnimateText(FrameCount+1, &model.TextBox{
					X:         musicNote.X,
					Y:         model.SCREEN_HEIGHT - model.MISSED_FONT_SIZE,
					FontSize:  model.MISSED_FONT_SIZE,
					FontColor: model.MISSED_FONT_COLOR,
					Text:      "MISS",
				})
				break
			}
			node = node.Prev
		}
	}
}

func calculateScore(rectangleA, rectangleB model.Rectangle) int {
	if rectangleA.X+rectangleA.Width <= rectangleB.X || rectangleB.X+rectangleB.Width <= rectangleA.X {
		return 0
	}
	if rectangleA.Y+rectangleA.Height <= rectangleB.Y || rectangleB.Y+rectangleB.Height <= rectangleA.Y {
		return 0
	}
	x := math.Abs(math.Max(float64(rectangleA.X), float64(rectangleB.X)) - math.Min(float64(rectangleA.X+rectangleA.Width), float64(rectangleB.X+rectangleB.Width)))
	y := math.Abs(math.Max(float64(rectangleA.Y), float64(rectangleB.Y)) - math.Min(float64(rectangleA.Y+rectangleA.Height), float64(rectangleB.Y+rectangleB.Height)))
	present := x * y / math.Min(model.TOUCH_BLOCK_AREA, model.MUSIC_NOTE_AREA)
	// 低于最坏情况
	if present < model.TOUCH_BLOCK_BAD_TOLERANCE {
		return 0
	}
	if present > model.TOUCH_BLOCK_PERFECT_TOLERANCE {
		return 5
	}
	if present > model.TOUCH_BLOCK_OK_TOLERANCE {
		return 3
	}
	if present > model.TOUCH_BLOCK_BAD_TOLERANCE {
		return 1
	}
	return 0
}

// 增加得分
func addScore() {
	for i := range MusicNoteList {
		// 1. 判断该key是否按下去了
		touchRect := model.Rectangle{
			X:      TouchNoteList[i].X,
			Y:      TouchNoteList[i].Y,
			Width:  TouchNoteList[i].Width,
			Height: TouchNoteList[i].Height,
		}
		if rl.IsKeyPressed(TouchNoteList[i].Key) || judgeGamePadPressed(TouchNoteList[i].Key) {
			bingoSuccess := false
			scoreIncr := 0
			node := new(model.DoubleNode)
			node = MusicNoteList[i].Head
			if node != nil {
				musicNote := node.Data.(model.MusicNote)
				musicNoteRect := model.Rectangle{
					X:      musicNote.X,
					Y:      musicNote.Y,
					Width:  musicNote.Width,
					Height: musicNote.Height,
				}
				// 2.计算获得的分数
				scoreIncr = calculateScore(touchRect, musicNoteRect)
				if scoreIncr > 0 {
					bingoSuccess = true
					MusicNoteList[i].Delete(node)
					// 3. 根据分数选择颜色
					switch scoreIncr {
					case 5:
						TouchNoteList[i].Color = model.TOUCH_BLOCK_PERFECT_COLOR
					case 3:
						TouchNoteList[i].Color = model.TOUCH_BLOCK_OK_COLOR
					case 1:
						TouchNoteList[i].Color = model.TOUCH_BLOCK_BAD_COLOR
					}
					if scoreIncr == 5 {
						SuccessiveNum += 1
					} else {
						SuccessiveNum = 0
					}
					if SuccessiveNum >= model.PerfectLimit {
						RegisterAnimateText(FrameCount+5, &model.TextBox{
							X:         TouchNoteList[i].X,
							Y:         TouchNoteList[i].Y,
							FontSize:  model.SCORE_FONT_SIZE,
							Text:      model.PerfectWord,
							FontColor: model.GREAT_COLOR,
						})
					} else if SuccessiveNum >= model.GreatLimit {
						RegisterAnimateText(FrameCount+5, &model.TextBox{
							X:         TouchNoteList[i].X,
							Y:         TouchNoteList[i].Y,
							FontSize:  model.SCORE_FONT_SIZE,
							Text:      model.GreatWord,
							FontColor: model.GREAT_COLOR,
						})
					}
				}
			}
			// 4. 按成功
			if bingoSuccess {
				sheet := CurrentMusicSheet.Data.(*model.MusicSheet)
				if JudgeSheetEnded(sheet) {
					sheet.CurrentNode = sheet.List.Head
					if CurrentMusicSheet.Prev == nil {
						CurrentMusicSheet = MusicSheetList.Head
					} else {
						CurrentMusicSheet = CurrentMusicSheet.Prev
					}
				}
				PlayMusicSheet(sheet)
				msg := fmt.Sprintf("+%v", scoreIncr)
				//@Todo:显示msg
				RegisterAnimateText(FrameCount+1, &model.TextBox{
					X:         TouchNoteList[i].X,
					Y:         TouchNoteList[i].Y,
					FontSize:  model.SCORE_FONT_SIZE,
					Text:      msg,
					FontColor: model.SCORE_FONT_COLOR,
				})
				ScoreSum += scoreIncr

			} else {
				TouchNoteList[i].Color = model.TOUCH_BLOCK_MISTOUCH_COLOR
			}

			// 5. 10帧后改回颜色
			data := &model.SetTouchBlockColor{
				TouchNote: &TouchNoteList[i],
				Color:     model.TOUCH_BLOCK_INIT_COLOR,
			}
			RegisterAtFrame(FrameCount+10, data, ResetNodeColor)
		}
	}
}

func judgeGamePadXbox() bool {
	if rl.IsGamepadAvailable(rl.GamepadPlayer1) {
		//目前只支持Xbox手柄
		if strings.Contains(rl.GetGamepadName(rl.GamepadPlayer1), "Xbox") {
			return true
		}
		return false
	}
	return false
}

func judgeGamePadPressed(key int32) bool {
	if judgeGamePadXbox() && rl.IsGamepadButtonPressed(rl.GamepadPlayer1, model.GamePadXboxKey[key]) {
		return true
	} else {
		return false
	}
}

func speed() float32 {
	return float32(model.INIT_MUSIC_NOTE_SPEED + (FrameCount / 3000.0))
}

func updateNoteY() {
	for i := range MusicNoteList {
		node := new(model.DoubleNode)
		node = MusicNoteList[i].Head
		for node != nil {
			musicNote := node.Data.(model.MusicNote)
			musicNote.Y += speed()
			node.Data = musicNote
			node = node.Prev
		}
	}
}

func checkFrameAction() {
	// run previously registered functions
	node := FramActionList.Head
	for node != nil {
		action := node.Data.(model.FrameRegisterAction)
		if action.Frame == FrameCount {
			action.Function(action.Data)
			FramActionList.Delete(node)
			node = FramActionList.Head
			continue
		}
		node = node.Prev
	}
}

func FlushGame() {
	//if rl.GetGamepadButtonPressed() != -1 {
	//	fmt.Println(rl.GetGamepadButtonPressed())
	//}
	//if rl.IsGamepadButtonPressed(0, rl.GamepadXboxButtonA) || rl.IsGamepadButtonPressed(0, rl.GamepadXboxButtonX) || rl.IsGamepadButtonPressed(0, rl.GamepadXboxButtonY) || rl.IsGamepadButtonPressed(0, rl.GamepadXboxButtonB) {
	//	fmt.Println("find!!!")
	//}
	generateNextNote()
	getMissMusicNote()
	addScore()
	updateNoteY()
	FrameCount++
	checkFrameAction()
}

func DrawMenu() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	var fontSize0 int32 = 30
	var fontSize1 int32 = 30
	var fontSize2 int32 = 30

	switch OptionSelect {
	case 0:
		fontSize0 = 48
	case 1:
		fontSize1 = 48
	case 2:
		fontSize2 = 48
	}
	rl.DrawText("START GAME", 100, 100, fontSize0, model.TOUCH_BLOCK_FONT_COLOR)
	rl.DrawText("SONG BOX", 100, 100+64, fontSize1, model.TOUCH_BLOCK_FONT_COLOR)
	rl.DrawText("CREDIT", 100, 100+64+64, fontSize2, model.TOUCH_BLOCK_FONT_COLOR)

	if rl.IsKeyPressed(rl.KeyDown) {
		OptionSelect++
	} else if rl.IsKeyPressed(rl.KeyUp) {
		OptionSelect--
	}

	if OptionSelect < 0 {
		OptionSelect = 0
	}
	if OptionSelect > 2 {
		OptionSelect = 2
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		switch OptionSelect {
		case ChooseStartGame:
			CurrentScreen = InGAME
		case ChooseMusicBox:
			CurrentScreen = SongBox
		}
	}

	rl.EndDrawing()

}

var CurrentMusicName = ""
var CurrentChooseMusicIndex = 0

func DrawSongBox() {
	files, _ := ioutil.ReadDir("./sheet")
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	for idx, f := range files {
		if idx == CurrentChooseMusicIndex {
			rl.DrawText(f.Name(), 100, int32(100+(idx*30)), 24, rl.Red)
		} else {
			rl.DrawText(f.Name(), 100, int32(100+(idx*30)), 24, model.TOUCH_BLOCK_FONT_COLOR)
		}
	}

	if rl.IsKeyPressed(rl.KeyDown) && CurrentChooseMusicIndex < len(files)-1 {
		CurrentChooseMusicIndex++
	} else if rl.IsKeyPressed(rl.KeyUp) && CurrentChooseMusicIndex > 0 {
		CurrentChooseMusicIndex--
	}

	rl.EndDrawing()
}

func DrawGame() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// --- all lines' sides ---
	for line := 0; line < model.LINE; line++ {
		// left side
		rl.DrawLineEx(
			rl.Vector2{
				X: float32(model.LEFT_MARGIN + (line * model.LINE_WIDTH) + (line * model.MARGIN_BETWEEN_LINE))},
			rl.Vector2{
				X: float32(model.LEFT_MARGIN + (line * model.LINE_WIDTH) + (line * model.MARGIN_BETWEEN_LINE)),
				Y: model.SCREEN_HEIGHT - model.TOUCH_BLOCK_MARGIN_BOTTOM,
			},
			model.SIDE_LINE_WIDTH, model.SIDE_LINE_COLOR)
		// right side
		rl.DrawLineEx(
			rl.Vector2{
				X: float32(model.LEFT_MARGIN + model.LINE_WIDTH + (line * model.LINE_WIDTH) + (line * model.MARGIN_BETWEEN_LINE))},
			rl.Vector2{
				X: float32(model.LEFT_MARGIN + model.LINE_WIDTH + (line * model.LINE_WIDTH) + (line * model.MARGIN_BETWEEN_LINE)),
				Y: model.SCREEN_HEIGHT - model.TOUCH_BLOCK_MARGIN_BOTTOM,
			},
			model.SIDE_LINE_WIDTH, model.SIDE_LINE_COLOR)
	}
	// --- music nodes ---
	for line := 0; line < model.LINE; line++ {
		for node := MusicNoteList[line].Head; node != nil; node = node.Prev {
			//fmt.Print(strconv.Itoa(int(node.Data.(model.MusicNote).Y)))
			rl.DrawRectangle(int32(node.Data.(model.MusicNote).X),
				int32(node.Data.(model.MusicNote).Y),
				int32(node.Data.(model.MusicNote).Width),
				int32(node.Data.(model.MusicNote).Height),
				node.Data.(model.MusicNote).Color)
		}
	}

	// --- touch blocks ---
	for line := 0; line < model.LINE; line++ {
		rl.DrawRectangle(int32(TouchNoteList[line].X),
			int32(TouchNoteList[line].Y),
			int32(TouchNoteList[line].Width),
			int32(TouchNoteList[line].Height),
			TouchNoteList[line].Color)

		// text on touch block
		text := TouchNoteList[line].Key
		rl.DrawText(string(text), int32(TouchNoteList[line].X+model.TOUCH_BLOCK_FONT_SIZE/2.0),
			int32(TouchNoteList[line].Y+model.TOUCH_BLOCK_FONT_SIZE/2.0),
			model.TOUCH_BLOCK_FONT_SIZE,
			model.TOUCH_BLOCK_FONT_COLOR)
		if judgeGamePadXbox() {
			rl.DrawCircle(int32(TouchNoteList[line].X+model.TOUCH_BLOCK_FONT_SIZE*3),
				int32(TouchNoteList[line].Y+TouchNoteList[line].Height/2.0), TouchNoteList[line].Height/2.5, rl.Pink)
			rl.DrawText(model.GamePadXboxKeyLetter[model.GamePadXboxKey[TouchNoteList[line].Key]],
				int32(TouchNoteList[line].X+model.TOUCH_BLOCK_FONT_SIZE*3-model.TOUCH_BLOCK_FONT_SIZE/3),
				int32(TouchNoteList[line].Y+model.TOUCH_BLOCK_FONT_SIZE/2.0), model.TOUCH_BLOCK_FONT_SIZE,
				model.TOUCH_BLOCK_FONT_COLOR)
		}
	}

	// --- text animations ---
	for node := AnimationText.Head; node != nil; node = node.Prev {
		textBox := node.Data.(*model.TextBox)
		rl.DrawText(textBox.Text, int32(textBox.X), int32(textBox.Y), int32(textBox.FontSize), textBox.FontColor)
	}

	// --- text of upper-left corner ---
	rl.DrawText(fmt.Sprintf("SCORE: %d", ScoreSum), model.UI_MARGIN, model.UI_MARGIN, model.UI_FONT_SIZE, model.LIGHTGRAY)
	rl.DrawText(fmt.Sprintf("MISS: %d", MissCount), model.UI_MARGIN, model.UI_MARGIN+model.UI_FONT_SIZE, model.UI_FONT_SIZE, model.UI_FONT_COLOR)
	rl.DrawText(fmt.Sprintf("SPEED: %.1f", speed()), model.UI_MARGIN, model.UI_MARGIN+model.UI_FONT_SIZE*2, model.UI_FONT_SIZE,
		model.UI_FONT_COLOR)

	rl.EndDrawing()

}
