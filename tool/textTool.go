/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package tool

import (
	"math"

	"awesomeProject1/model"
)

func RegisterAnimateText(frame int, text *model.TextBox) {
	// 1. 扩展到待展示内
	AnimationText.Append(text)
	// 2. 确定展示角度
	radians := math.Pi / 2
	// 3. 逐个注册帧
	for i := 0; i < model.ANIMATE_TEXT_DURATION; i++ {
		data := &model.FrameRegisterText{
			Text: AnimationText.Tail.Data.(*model.TextBox),
			Dx:   -2 * math.Cos(radians),
			Dy:   -2 * math.Sin(radians),
		}
		RegisterAtFrame(frame+i, data, ChangeTextPosition)
	}
	// 4. 删除文案
	data := AnimationText.Tail
	RegisterAtFrame(frame+model.ANIMATE_TEXT_DURATION, data, FrameDeleteText)
}

// 注册到某一帧中
func RegisterAtFrame(frame int, data interface{}, fun func(text interface{})) {
	if frame < FrameCount {
		return
	}
	registerAction := model.FrameRegisterAction{
		Function: fun,
		Frame:    frame,
		Data:     data,
	}
	FramActionList.Append(registerAction)
}

func ChangeTextPosition(text interface{}) {
	frameText := text.(*model.FrameRegisterText)
	frameText.Text.X += float32(frameText.Dx)
	frameText.Text.Y += float32(frameText.Dy)
}

func FrameDeleteText(text interface{}) {
	node := text.(*model.DoubleNode)
	AnimationText.Delete(node)
}

func ResetNodeColor(text interface{}) {
	node := text.(*model.SetTouchBlockColor)
	node.TouchNote.Color = node.Color
}
