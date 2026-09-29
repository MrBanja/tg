package model

import "encoding/json"

type Message struct {
	MessageID       int64             `json:"message_id"`
	MessageThreadID *int64            `json:"message_thread_id,omitempty"`
	ReplyToMessage  *Message          `json:"reply_to_message,omitempty"`
	IsTopicMessage  *bool             `json:"is_topic_message,omitempty"`
	From            *User             `json:"from,omitempty"`
	Chat            *Chat             `json:"chat,omitempty"`
	ForwardOrigin   *MessageOrigin    `json:"forward_origin,omitempty"`
	Date            int64             `json:"date"`
	Text            *string           `json:"text,omitempty"`
	Caption         *string           `json:"caption,omitempty"`
	RichMessage     *RichMessage      `json:"rich_message,omitempty"`
	Sticker         *Sticker          `json:"sticker,omitempty"`
	Photo           []*MessagePhoto   `json:"photo,omitempty"`
	Document        *MessageDocument  `json:"document,omitempty"`
	Animation       *MessageAnimation `json:"animation,omitempty"`
	Video           *MessageVideo     `json:"video,omitempty"`
	Entities        []MessageEntity   `json:"entities,omitempty"`
	MediaGroupID    *string           `json:"media_group_id,omitempty"`
	CaptionEntities []MessageEntity   `json:"caption_entities,omitempty"`
	Audio           *Audio            `json:"audio,omitempty"`
	Voice           *Voice            `json:"voice,omitempty"`
	Location        *Location         `json:"location,omitempty"`
	// Venue is set for shared venues; Location is also set in that case.
	Venue *Venue `json:"venue,omitempty"`
}

// Voice represents a voice note.
type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// Audio represents an audio file to be treated as music.
type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer,omitempty"`
	Title        string `json:"title,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// Location represents a point on the map.
type Location struct {
	Latitude             float64  `json:"latitude"`
	Longitude            float64  `json:"longitude"`
	HorizontalAccuracy   *float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           *int     `json:"live_period,omitempty"`
	Heading              *int     `json:"heading,omitempty"`
	ProximityAlertRadius *int     `json:"proximity_alert_radius,omitempty"`
}

// Venue represents a venue.
type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareID    string   `json:"foursquare_id,omitempty"`
	FoursquareType  string   `json:"foursquare_type,omitempty"`
	GooglePlaceID   string   `json:"google_place_id,omitempty"`
	GooglePlaceType string   `json:"google_place_type,omitempty"`
}

// File represents a file ready to be downloaded with DownloadFile.
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
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

// MessageOrigin describes the origin of a forwarded message; only the fields
// shared by all origin variants (user, hidden_user, chat, channel) are mapped.
type MessageOrigin struct {
	Type string `json:"type"`
	Date int64  `json:"date"`
}

type MessageEntity struct {
	Type           string `json:"type"`
	Offset         int    `json:"offset"`
	Length         int    `json:"length"`
	URL            string `json:"url,omitempty"`
	User           *User  `json:"user,omitempty"`
	Language       string `json:"language,omitempty"`
	CustomEmojiID  string `json:"custom_emoji_id,omitempty"`
	UnixTime       int64  `json:"unix_time,omitempty"`
	DateTimeFormat string `json:"date_time_format,omitempty"`
}

type ReplyParameters struct {
	MessageID                int64  `json:"message_id"`
	ChatID                   *int64 `json:"chat_id,omitempty"`
	AllowSendingWithoutReply bool   `json:"allow_sending_without_reply,omitempty"`
	Quote                    string `json:"quote,omitempty"`
	QuoteParseMode           string `json:"quote_parse_mode,omitempty"`
	QuotePosition            *int   `json:"quote_position,omitempty"`
}

type DeleteMessageRequest struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int64 `json:"message_id"`
}

type SendMessageRequest struct {
	ChatID              int64               `json:"chat_id"`
	MessageThreadID     *int64              `json:"message_thread_id,omitempty"`
	Text                string              `json:"text,omitempty"`
	ReplyParameters     *ReplyParameters    `json:"reply_parameters,omitempty"`
	ReplyMarkup         Markup              `json:"reply_markup,omitempty"`
	ParseMode           string              `json:"parse_mode,omitempty"`
	LinkPreviewOptions  *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	Entities            []MessageEntity     `json:"entities,omitempty"`
	DisableNotification bool                `json:"disable_notification,omitempty"`
	ProtectContent      bool                `json:"protect_content,omitempty"`
}

// SendMessageDraftRequest streams a partial plain-text message to a private chat.
// An empty Text shows a "Thinking…" placeholder. DraftID must be non-zero.
type SendMessageDraftRequest struct {
	ChatID          int64           `json:"chat_id"`
	MessageThreadID *int64          `json:"message_thread_id,omitempty"`
	DraftID         int64           `json:"draft_id"`
	Text            string          `json:"text"`
	ParseMode       string          `json:"parse_mode,omitempty"`
	Entities        []MessageEntity `json:"entities,omitempty"`
	// CanStop shows a stop button; pressing it produces a stopped_message_generation update.
	CanStop    bool `json:"can_stop,omitempty"`
	KeepOnStop bool `json:"keep_on_stop,omitempty"`
}

type SendRichMessageRequest struct {
	BusinessConnectionID  string           `json:"business_connection_id,omitempty"`
	ChatID                int64            `json:"chat_id"`
	MessageThreadID       *int64           `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int64            `json:"direct_messages_topic_id,omitempty"`
	RichMessage           InputRichMessage `json:"rich_message"`
	DisableNotification   bool             `json:"disable_notification,omitempty"`
	ProtectContent        bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool             `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string           `json:"message_effect_id,omitempty"`
	ReplyParameters       *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup           Markup           `json:"reply_markup,omitempty"`
}

type SendRichMessageDraftRequest struct {
	ChatID          int64            `json:"chat_id"`
	MessageThreadID *int64           `json:"message_thread_id,omitempty"`
	DraftID         int64            `json:"draft_id"`
	RichMessage     InputRichMessage `json:"rich_message"`
	// CanStop shows a stop button; pressing it produces a stopped_message_generation update.
	CanStop    bool `json:"can_stop,omitempty"`
	KeepOnStop bool `json:"keep_on_stop,omitempty"`
}

// InputRichMessage describes a rich message to be sent.
// Exactly one of HTML, Markdown or Blocks must be used.
type InputRichMessage struct {
	Blocks              []InputRichBlock        `json:"blocks,omitempty"`
	HTML                string                  `json:"html,omitempty"`
	Markdown            string                  `json:"markdown,omitempty"`
	Media               []InputRichMessageMedia `json:"media,omitempty"`
	IsRTL               bool                    `json:"is_rtl,omitempty"`
	SkipEntityDetection bool                    `json:"skip_entity_detection,omitempty"`
}

type LinkPreviewOptions struct {
	IsDisabled bool `json:"is_disabled,omitempty"`
}

type SetMessageReactionRequest struct {
	ChatID    int64               `json:"chat_id"`
	MessageID int64               `json:"message_id"`
	Reaction  []ReactionTypeEmoji `json:"reaction"`
	IsBig     *bool               `json:"is_big,omitempty"`
}

type EditMessageTextRequest struct {
	// Unique identifier of the business connection on behalf of which the message to be edited was sent
	BusinessConnectionID *string `json:"business_connection_id,omitempty"`
	// Required if inline_message_id is not specified. Unique identifier for the target chat or username of the target channel
	ChatID *int64 `json:"chat_id,omitempty"`
	// Required if inline_message_id is not specified. Identifier of the message to edit
	MessageID *int64 `json:"message_id,omitempty"`
	// Required if chat_id and message_id are not specified. Identifier of the inline message
	InlineMessageID    *string             `json:"inline_message_id,omitempty"`
	Text               string              `json:"text,omitempty"`
	ParseMode          string              `json:"parse_mode,omitempty"`
	Entities           []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	RichMessage        *InputRichMessage   `json:"rich_message,omitempty"`
	ReplyMarkup        Markup              `json:"reply_markup,omitempty"`
}

type RichMessage struct {
	Blocks []json.RawMessage `json:"blocks"`
	IsRTL  bool              `json:"is_rtl,omitempty"`
}

type ReactionTypeEmojiType string

const (
	Emoji ReactionTypeEmojiType = "emoji"
)

type ReactionTypeEmoji struct {
	Type  ReactionTypeEmojiType `json:"type"`
	Emoji string                `json:"emoji"`
}

// SendVoiceRequest sends a voice note. Voice is a file_id or an HTTP URL;
// leave it empty and use SendVoiceUpload to upload a new file.
type SendVoiceRequest struct {
	ChatID              int64            `json:"chat_id"`
	MessageThreadID     *int64           `json:"message_thread_id,omitempty"`
	Voice               string           `json:"voice,omitempty"`
	Caption             string           `json:"caption,omitempty"`
	ParseMode           string           `json:"parse_mode,omitempty"`
	CaptionEntities     []MessageEntity  `json:"caption_entities,omitempty"`
	Duration            int              `json:"duration,omitempty"`
	DisableNotification bool             `json:"disable_notification,omitempty"`
	ProtectContent      bool             `json:"protect_content,omitempty"`
	ReplyParameters     *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup         Markup           `json:"reply_markup,omitempty"`
}

type ChatAction string

const (
	ChatActionTyping          ChatAction = "typing"
	ChatActionUploadPhoto     ChatAction = "upload_photo"
	ChatActionRecordVoice     ChatAction = "record_voice"
	ChatActionUploadVoice     ChatAction = "upload_voice"
	ChatActionUploadDocument  ChatAction = "upload_document"
	ChatActionChooseSticker   ChatAction = "choose_sticker"
	ChatActionFindLocation    ChatAction = "find_location"
	ChatActionRecordVideo     ChatAction = "record_video"
	ChatActionUploadVideo     ChatAction = "upload_video"
	ChatActionRecordVideoNote ChatAction = "record_video_note"
	ChatActionUploadVideoNote ChatAction = "upload_video_note"
)

type SendChatActionRequest struct {
	ChatID          int64      `json:"chat_id"`
	MessageThreadID *int64     `json:"message_thread_id,omitempty"`
	Action          ChatAction `json:"action"`
}

type AnswerCallbackQueryRequest struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// BotCommandScope describes the scope of bot commands. Type is one of default,
// all_private_chats, all_group_chats, all_chat_administrators, chat,
// chat_administrators or chat_member; ChatID and UserID are set when the type requires them.
type BotCommandScope struct {
	Type   string `json:"type"`
	ChatID int64  `json:"chat_id,omitempty"`
	UserID int64  `json:"user_id,omitempty"`
}

type SetMyCommandsRequest struct {
	Commands     []BotCommand     `json:"commands"`
	Scope        *BotCommandScope `json:"scope,omitempty"`
	LanguageCode string           `json:"language_code,omitempty"`
}

type GetFileRequest struct {
	FileID string `json:"file_id"`
}

// EditMessageReplyMarkupRequest edits only the inline keyboard of a message.
// A nil ReplyMarkup removes the keyboard.
type EditMessageReplyMarkupRequest struct {
	ChatID          *int64                `json:"chat_id,omitempty"`
	MessageID       *int64                `json:"message_id,omitempty"`
	InlineMessageID *string               `json:"inline_message_id,omitempty"`
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}
