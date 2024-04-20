package tgmodel

type Markup struct {
	InlineKeyboard [][]InlineKeyBoardButton `json:"inline_keyboard"`
}

type InlineKeyBoardButton struct {
	Text                         string    `json:"text"`
	URL                          string    `json:"url"`
	LoginURL                     *LoginURL `json:"login_url,omitempty"`
	CallbackData                 string    `json:"callback_data"`
	SwitchInlineQuery            string    `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string    `json:"switch_inline_query_current_chat"`
	Pay                          bool      `json:"pay"`
}

type LoginURL struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text"`
	BotUsername        string `json:"bot_username"`
	RequestWriteAccess bool   `json:"request_write_access"`
}
