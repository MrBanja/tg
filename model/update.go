package model

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID       int               `json:"message_id"`
	MessageThreadID *int              `json:"message_thread_id,omitempty"`
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

type User struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
	IsPremium    bool   `json:"is_premium,omitempty"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Type      string `json:"type"`
}

type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	CustomEmojiID string `json:"custom_emoji_id"`
}
