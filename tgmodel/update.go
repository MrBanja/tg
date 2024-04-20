package tgmodel

type Update struct {
	ID            int64          `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type CallbackQuery struct {
	ID           string   `json:"id"`
	From         From     `json:"from"`
	Message      *Message `json:"message"`
	ChatInstance string   `json:"chat_instance"`
	Data         string   `json:"data"`
}
