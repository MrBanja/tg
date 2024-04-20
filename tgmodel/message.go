package tgmodel

type Message struct {
	MessageID       int         `json:"message_id"`
	From            *From       `json:"from,omitempty"`
	SenderChat      *Chat       `json:"sender_chat,omitempty"`
	Date            int         `json:"date"`
	Chat            *Chat       `json:"chat,omitempty"`
	ViaBot          *User       `json:"via_bot,omitempty"`
	EditDate        int         `json:"edit_date"`
	MediaGroupID    string      `json:"media_group_id"`
	AuthorSignature string      `json:"author_signature"`
	Text            string      `json:"text"`
	Entities        []Entity    `json:"entities,omitempty"`
	Animation       *Animation  `json:"animation,omitempty"`
	Audio           *Audio      `json:"audio,omitempty"`
	Document        *Document   `json:"document,omitempty"`
	Photo           []PhotoSize `json:"photo,omitempty"`
	Sticker         *Sticker    `json:"sticker,omitempty"`
	Video           *Video      `json:"video,omitempty"`
	VideoNote       *VideoNote  `json:"video_note,omitempty"`
	Voice           *Voice      `json:"voice,omitempty"`
	Caption         string      `json:"caption"`
	CaptionEntities []Entity    `json:"caption_entities"`
	Contact         *Contact    `json:"contact,omitempty"`
	Dice            *Dice       `json:"dice,omitempty"`
	Poll            *Poll       `json:"poll,omitempty"`
	ReplyMarkup     *Markup     `json:"reply_markup,omitempty"`
}

type SendMessage struct {
	ChatID                   *int     `json:"chat_id,omitempty"`
	Text                     *string  `json:"text,omitempty"`
	ParseMode                *string  `json:"parse_mode,omitempty"`
	Entities                 []Entity `json:"entities,omitempty"`
	DisableWebPagePreview    *bool    `json:"disable_web_page_preview,omitempty"`
	DisableNotification      *bool    `json:"disable_notification,omitempty"`
	ReplyToMessageID         *int     `json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply *bool    `json:"allow_sending_without_reply,omitempty"`
	ReplyMarkup              *Markup  `json:"reply_markup,omitempty"`
}

type Entity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url"`
	User   User   `json:"user"`

	Language string `json:"language"`
}

type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`

	CanJoinGroups           bool `json:"can_join_groups"`
	CanReadAllGroupMessages bool `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool `json:"supports_inline_queries"`
}

type Chat struct {
	ID               int          `json:"id"`
	Type             string       `json:"type"`
	Title            string       `json:"title"`
	Username         string       `json:"username"`
	FirstName        string       `json:"first_name"`
	LastName         string       `json:"last_name"`
	Photo            *Photo       `json:"photo,omitempty"`
	Bio              string       `json:"bio"`
	Description      string       `json:"description"`
	InviteLink       string       `json:"invite_link"`
	Permissions      *Permissions `json:"permissions,omitempty"`
	SlowModeDelay    int          `json:"slow_mode_delay"`
	StickerSetName   string       `json:"sticker_set_name"`
	CanSetStickerSet bool         `json:"can_set_sticker_set"`
	LinkedChatID     int          `json:"linked_chat_id"`
	Location         *Location    `json:"location"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int    `json:"user_id"`
	Vcard       string `json:"vcard"`
}
