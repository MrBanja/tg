package model

type Message struct {
	MessageID       int64             `json:"message_id"`
	MessageThreadID *int64            `json:"message_thread_id,omitempty"`
	ReplyToMessage  *Message          `json:"reply_to_message,omitempty"`
	IsTopicMessage  *bool             `json:"is_topic_message,omitempty"`
	From            *User             `json:"from,omitempty"`
	Chat            *Chat             `json:"chat,omitempty"`
	Date            int64             `json:"date"`
	Text            *string           `json:"text,omitempty"`
	Sticker         *Sticker          `json:"sticker,omitempty"`
	Photo           []*MessagePhoto   `json:"photo,omitempty"`
	Document        *MessageDocument  `json:"document,omitempty"`
	Animation       *MessageAnimation `json:"animation,omitempty"`
	Video           *MessageVideo     `json:"video,omitempty"`
	Entities        []MessageEntity   `json:"entities,omitempty"`
}

type MessageVideo struct {
	Duration     int    `json:"duration"`
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	FileSize     int    `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	Height       int    `json:"height"`
	MimeType     string `json:"mime_type"`
	Width        int    `json:"width"`
}

type MessageAnimation struct {
	Duration     int    `json:"duration"`
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	FileSize     int    `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	Height       int    `json:"height"`
	MimeType     string `json:"mime_type"`
	Width        int    `json:"width"`
}

type MessageDocument struct {
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	FileSize     int    `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	MimeType     string `json:"mime_type"`
}

type MessagePhoto struct {
	FileID       string `json:"file_id"`
	FileSize     int    `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
}

type Sticker struct {
	Emoji        string `json:"emoji"`
	FileID       string `json:"file_id"`
	FileSize     int    `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	Height       int    `json:"height"`
	IsAnimated   bool   `json:"is_animated"`
	IsVideo      bool   `json:"is_video"`
	SetName      string `json:"set_name"`
	Type         string `json:"type"`
	Width        int    `json:"width"`
}

type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	CustomEmojiID string `json:"custom_emoji_id"`
}

type ReplyParameters struct {
	MessageID int64 `json:"message_id"`
}

type DeleteMessageRequest struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int64 `json:"message_id"`
}

type SendMessageRequest struct {
	ChatID          int64            `json:"chat_id"`
	MessageThreadID *int64           `json:"message_thread_id,omitempty"`
	Text            string           `json:"text,omitempty"`
	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup     Markup           `json:"reply_markup,omitempty"`
	ParseMode       string           `json:"parse_mode,omitempty"`
}

type SetMessageReactionRequest struct {
	ChatID    int64               `json:"chat_id"`
	MessageID int64               `json:"message_id"`
	Reaction  []ReactionTypeEmoji `json:"reaction,omitempty"`
	IsBig     *bool               `json:"is_big,omitempty"`
}

type ReactionTypeEmojiType string

const (
	Emoji ReactionTypeEmojiType = "emoji"
)

type ReactionTypeEmoji struct {
	Type  ReactionTypeEmojiType `json:"type"`
	Emoji string                `json:"emoji"`
}
