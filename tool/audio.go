package tool

import (
	"awesomeProject1/model"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var soundMap = make(map[string]rl.Sound)

const AUDIO_DIR = "audio"

// PlayMusicSheet Play the next music note of 'sheet' and advance the progress
func PlayMusicSheet(sheet *model.MusicSheet) {
	if sheet.CurrentNode != nil {
		playMusicNote(sheet.CurrentNode.Data)
		sheet.CurrentNode = sheet.CurrentNode.Prev
	}
}

// Play a sound of music note
func playMusicNote(data model.NodeObject) {
	note := data.(string)
	if s, ok := soundMap[note]; ok {
		rl.PlaySound(s)
	} else {
		sound := loadMusicNoteSound(note)
		soundMap[note] = sound
		rl.PlaySound(s)
	}
}

func loadMusicNoteSound(note string) rl.Sound {
	filePath := AUDIO_DIR + "/" + note + ".mp3"
	sound := rl.LoadSound(filePath)
	return sound
}
