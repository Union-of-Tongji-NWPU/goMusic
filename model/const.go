package model

import rl "github.com/gen2brain/raylib-go/raylib"

const WINDOW_TITLE = "Music Tiles"
const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 800
const LINE = 4
const LINE_WIDTH = 100
const MAX_MUSIC_NOTE_EACH_LINE = 100
const MARGIN_BETWEEN_LINE = 10
const MUSIC_NOTE_WIDTH = LINE_WIDTH
const MUSIC_NOTE_HEIGHT = 160

var MUSIC_NOTE_INIT_COLOR = rl.DarkGray
var TOUCH_BLOCK_INIT_COLOR = rl.LightGray
var TOUCH_BLOCK_FONT_COLOR = rl.DarkGray
var TOUCH_BLOCK_PERFECT_COLOR = rl.Green
var TOUCH_BLOCK_OK_COLOR = rl.SkyBlue
var TOUCH_BLOCK_BAD_COLOR = rl.Orange
var TOUCH_BLOCK_MISTOUCH_COLOR = rl.Red
var LIGHTGRAY = rl.LightGray
var UI_FONT_COLOR = rl.LightGray
var MISSED_FONT_COLOR = rl.DarkGray
var SCORE_FONT_COLOR = rl.Gray
var SIDE_LINE_COLOR = rl.LightGray

const MUSIC_NOTE_AREA = MUSIC_NOTE_WIDTH * MUSIC_NOTE_HEIGHT

const INIT_MUSIC_NOTE_SPEED = 3 // 滑块下落速度

// 按扭相关位置定义
const TOUCH_BLOCK_WIDTH = LINE_WIDTH
const TOUCH_BLOCK_HEIGHT = 40
const TOUCH_BLOCK_AREA = TOUCH_BLOCK_WIDTH * TOUCH_BLOCK_HEIGHT

const TOUCH_BLOCK_MARGIN_BOTTOM = 20

// 命中滑块分数
const TOUCH_BLOCK_PERFECT_TOLERANCE = 0.9
const TOUCH_BLOCK_OK_TOLERANCE = 0.6
const TOUCH_BLOCK_BAD_TOLERANCE = 0.1
const TOUCH_BLOCK_FONT_SIZE = 20

const SIDE_LINE_WIDTH = 1
const UI_FONT_SIZE = 20
const MISSED_FONT_SIZE = 20
const SCORE_FONT_SIZE = 22
const UI_MARGIN = 5
const ANIMATE_TEXT_DURATION = 30 // 动画停留的帧数

const LEFT_MARGIN = (SCREEN_WIDTH - (LINE_WIDTH * LINE) - (MARGIN_BETWEEN_LINE * (LINE - 1))) / 2

const RIGHT_MARGIN = (SCREEN_WIDTH - (LINE_WIDTH * LINE) - (MARGIN_BETWEEN_LINE * (LINE - 1))) / 2

//KEY_A               = 65,
//KEY_B               = 66,
//KEY_C               = 67,
//KEY_D               = 68,
var KeyboardKey = [5]int{65, 66, 67, 68}
