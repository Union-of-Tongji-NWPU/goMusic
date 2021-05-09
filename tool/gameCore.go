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
	"math"
	"math/rand"

	"awesomeProject1/model"
)

var MusicSheetList = new(model.DoubleList)
var CurrentMusicSheet = new(model.DoubleNode)
var TouchNoteList = make([]model.TouchNote, model.LINE)
var MusicNoteList = make([]*model.DoubleList, model.LINE)
var AnimationText = new(model.DoubleList)
var MissCount = 0
var ScoreSum = 0 //总分
var FrameCount = 0

func InitGame(sheetFiles []string) {
	for k, _ := range sheetFiles {
		sheet := new(model.MusicSheet)
		LoadMusicSheetFromFile(sheet, sheetFiles[k])
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
		TouchNoteList[i].X = float32(model.LLEFT_MARGIN + i*(model.TOUCH_BLOCK_WIDTH+model.MARGIN_BETWEEN_LINE))
		TouchNoteList[i].Y = model.SCREEN_HEIGHT - model.TOUCH_BLOCK_HEIGHT - model.TOUCH_BLOCK_MARGIN_BOTTOM
		TouchNoteList[i].Color = model.TOUCH_BLOCK_BAD_COLOR
		TouchNoteList[i].Key = int32(model.KeyboardKey[i])
	}
}

// 获取下一个音符到MusicNoteList中
func getYOfNote(minHeight float32) float32 {
	if minHeight > 0 {
		return -model.MUSIC_NOTE_HEIGHT
	} else {
		return (minHeight - model.MUSIC_NOTE_HEIGHT)
	}
}
func generateNextNote() {
	// 1.找到当前最上面的音符
	var minHeight = float32(0x3f3f3f3f)
	for i := range MusicNoteList {
		node := new(model.DoubleNode)
		node = MusicNoteList[i].Head
		for node != nil {
			if minHeight < node.Data.(model.MusicNote).Y {
				minHeight = node.Data.(model.MusicNote).Y
			}
			node = node.Prev
		}
	}
	// 2. 当前需要产生新的音符
	for minHeight > 0 || minHeight >= -model.MUSIC_NOTE_HEIGHT {
		selectLine := rand.Intn(model.LINE)
		musicNote := model.MusicNote{
			X:      float32(model.LLEFT_MARGIN + selectLine*(model.LINE_WIDTH+model.MARGIN_BETWEEN_LINE)),
			Y:      getYOfNote(minHeight),
			Width:  model.MUSIC_NOTE_WIDTH,
			Height: model.MUSIC_NOTE_HEIGHT,
			Color:  model.MUSIC_NOTE_INIT_COLOR,
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
				//@Todo： 显示Miss
				break
			}
			node = node.Prev
		}
	}
}

// if ((a.x + a.width <= b.x || a.y + a.height <= b.y) || (b.x + b.width <= a.x || b.y + b.height <= a.y))
//        return 0.0f;
//    else
//        return fabsf((MAX(a.x, b.x) - MIN(a.x + a.width, b.x + b.width)) * (MAX(a.y, b.y) - MIN(a.y + a.height, b.y + b.height)));
func calculateScore(rectangleA, rectangleB model.Rectangle) int {
	if rectangleA.X+rectangleA.Width <= rectangleB.X || rectangleB.X+rectangleB.Width <= rectangleA.X {
		return 0
	}
	if rectangleA.Y+rectangleA.Height <= rectangleB.Y || rectangleB.Y+rectangleB.Height <= rectangleB.Y {
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
		if rl.IsKeyPressed(TouchNoteList[i].Key) {
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
				scoreIncr := calculateScore(touchRect, musicNoteRect)
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
					case 0:
						TouchNoteList[i].Color = model.TOUCH_BLOCK_MISTOUCH_COLOR

					}
				}
			}
			// 4. 按成功
			if bingoSuccess {
				sheet := CurrentMusicSheet.Data.(model.MusicSheet)
				if JudgeSheetEnded(&sheet) {
					sheet.CurrentNode = sheet.List.Head
					if CurrentMusicSheet.Prev == nil {
						CurrentMusicSheet = MusicSheetList.Head
					} else {
						CurrentMusicSheet = CurrentMusicSheet.Prev
					}
				}
				//PlayMusic @XieQiJun
				msg := fmt.Sprintf("+%v分", scoreIncr)
				//@Todo:显示msg
				ScoreSum += scoreIncr
			}
		}
	}
}

func speed()float32{

		return float32(model.INIT_MUSIC_NOTE_SPEED + (FrameCount / 3000.0))
}

func updateNoteY(){
	for i := range MusicNoteList{
		node := new(model.DoubleNode)
		node = MusicNoteList[i].Head
		for node != nil{
			musicNote := node.Data.(model.MusicNote)
			musicNote.Y += speed()
			node = node.Prev
		}
	}
}

func FlushGame() {
	generateNextNote()
	getMissMusicNote()
	addScore()
	updateNoteY()
	FrameCount++
}