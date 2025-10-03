package model

type Markup interface {
	__Markup()
}

var (
	_ Markup = (*InlineKeyboardMarkup)(nil)
)

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

func (*InlineKeyboardMarkup) __Markup() {}

type InlineKeyboardButton struct {
	Text         string  `json:"text"`
	CallbackData *string `json:"callback_data,omitempty"`
}
