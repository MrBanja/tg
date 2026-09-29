package model

import "testing"

func TestReplyKeyboardMarkupMarshal(t *testing.T) {
	req := SendMessageRequest{
		ChatID: 1,
		Text:   "hi",
		ReplyMarkup: &ReplyKeyboardMarkup{
			Keyboard:              [][]KeyboardButton{{{Text: "Start"}, {Text: "Settings"}}},
			IsPersistent:          true,
			ResizeKeyboard:        true,
			InputFieldPlaceholder: "Ask me",
		},
	}
	want := `{"chat_id":1,"text":"hi","reply_markup":{"keyboard":[[{"text":"Start"},{"text":"Settings"}]],"is_persistent":true,"resize_keyboard":true,"input_field_placeholder":"Ask me"}}`
	if got := mustMarshal(t, req); got != want {
		t.Fatalf("unexpected json: %s", got)
	}
}

func TestReplyKeyboardRemoveMarshal(t *testing.T) {
	var m Markup = &ReplyKeyboardRemove{}
	if got, want := mustMarshal(t, m), `{"remove_keyboard":true}`; got != want {
		t.Fatalf("unexpected json: %s", got)
	}
	m = &ForceReply{InputFieldPlaceholder: "Your prompt"}
	if got, want := mustMarshal(t, m), `{"force_reply":true,"input_field_placeholder":"Your prompt"}`; got != want {
		t.Fatalf("unexpected json: %s", got)
	}
}

func TestInlineKeyboardButtonMarshalOmitsEmptyFields(t *testing.T) {
	data := "cb"
	btn := InlineKeyboardButton{Text: "Go", CallbackData: &data, Style: "primary"}
	if got, want := mustMarshal(t, btn), `{"text":"Go","style":"primary","callback_data":"cb"}`; got != want {
		t.Fatalf("unexpected json: %s", got)
	}
}
