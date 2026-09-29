package model

type Markup interface {
	__Markup()
}

var (
	_ Markup = (*InlineKeyboardMarkup)(nil)
	_ Markup = (*ReplyKeyboardMarkup)(nil)
	_ Markup = (*ReplyKeyboardRemove)(nil)
	_ Markup = (*ForceReply)(nil)
)

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

func (*InlineKeyboardMarkup) __Markup() {}

type InlineKeyboardButton struct {
	Text         string          `json:"text"`
	Style        string          `json:"style,omitempty"` // danger, success or primary
	URL          *string         `json:"url,omitempty"`
	CallbackData *string         `json:"callback_data,omitempty"`
	CopyText     *CopyTextButton `json:"copy_text,omitempty"`
	Disabled     *DisabledButton `json:"disabled,omitempty"`
}

// ReplyKeyboardMarkup is a custom keyboard with reply options.
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          bool               `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"` // 1-64 characters
	Selective             bool               `json:"selective,omitempty"`
}

func (*ReplyKeyboardMarkup) __Markup() {}

// KeyboardButton is one button of a reply keyboard. With only Text set,
// the text is sent as a message when the button is pressed.
type KeyboardButton struct {
	Text            string `json:"text"`
	Style           string `json:"style,omitempty"` // danger, success or primary
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

// ReplyKeyboardRemove removes the current custom keyboard.
type ReplyKeyboardRemove struct {
	Selective bool `json:"selective,omitempty"`
}

func (*ReplyKeyboardRemove) __Markup() {}

func (r *ReplyKeyboardRemove) MarshalJSON() ([]byte, error) {
	type alias ReplyKeyboardRemove
	return marshalWithTrueField("remove_keyboard", alias(*r))
}

// ForceReply shows a reply interface to the user.
type ForceReply struct {
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	Selective             bool   `json:"selective,omitempty"`
}

func (*ForceReply) __Markup() {}

func (r *ForceReply) MarshalJSON() ([]byte, error) {
	type alias ForceReply
	return marshalWithTrueField("force_reply", alias(*r))
}
