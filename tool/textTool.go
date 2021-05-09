/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package tool

import (
	"awesomeProject1/model"
)

// 显示工具
//void register_animate_text(unsigned int frame, const char* text, int fontsize, Color fontcolor, float start_x,
//float start_y)
//{
//TextBox textbox = { .x = start_x, .y = start_y, .fontsize = fontsize, .text = strdup(text), .fontcolor = fontcolor };
//
//list_append_tail(&animation_texts, &textbox, sizeof(TextBox));
//
//float radians = randf(PI / 2 - PI / 6, PI / 2 + PI / 6); // 60 ~ 120 degrees
//for (int i = 0; i < ANIMATE_TEXT_DURATION; ++i) {
//__frame_register_animate_text_t data = {
//.textbox = animation_texts.tail->data, .dx = -2 * cosf(radians), .dy = -2 * sinf(radians)
//};
//register_at_frame(frame + i, (void (*)(void*))__frame_register_animate_text, &data,
//sizeof(__frame_register_animate_text_t));
//}
//
//__frame_register_drop_animate_text_t _data = { .node = animation_texts.tail };
//register_at_frame(frame + ANIMATE_TEXT_DURATION, (void (*)(void*))__frame_register_drop_animate_text, &_data,
//sizeof(__frame_register_drop_animate_text_t));
//}

func RegisterAnimateText(frame int,text model.TextBox){

}