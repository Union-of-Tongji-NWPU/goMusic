package tool

import (
	"io/ioutil"
	"os"
	"strings"

	"awesomeProject1/model"
)

const DEFAULT_MUSIC_SHEET = "c4 c4 g4 g4 a4 a4 a4 a4 g4 f4 f4 e4 e4 d4 d4 c4 g4 g4 g4 f4 f4 e4 e4 e4 d4 c4 g4 g4 g4 f4 f4 f4 f4 e4 e4 e4 d4"

func LoadMusicSheetFromString(sheet *model.MusicSheet, str string) {
	if sheet == nil {
		sheet = new(model.MusicSheet)
	}
	musicSheetList := strings.Fields(str)
	for _, v := range musicSheetList {
		sheet.List.Append(v)
	}
}

func LoadMusicSheetFromFile(sheet *model.MusicSheet, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	bytes, err := (ioutil.ReadAll(f))
	if err != nil {
		return err
	}
	LoadMusicSheetFromString(sheet, string(bytes))
	return nil
}

func JudgeSheetEnded(sheet *model.MusicSheet) bool{
	return sheet.CurrentNode == nil
}

func GetDefaultMusicSheet(sheet *model.MusicSheet){
	LoadMusicSheetFromString(sheet,DEFAULT_MUSIC_SHEET)
}