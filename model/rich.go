package model

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// marshalTyped marshals v (which must encode to a JSON object) and injects
// a leading "type" field. It lets union members carry their discriminator
// without exposing a Type field that callers could get wrong.
func marshalTyped(typ string, v any) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(body) < 2 || body[0] != '{' {
		return nil, fmt.Errorf("marshal %s: expected JSON object, got %s", typ, body)
	}
	typeField, err := json.Marshal(typ)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	b.WriteString(`{"type":`)
	b.Write(typeField)
	if !bytes.Equal(body, []byte("{}")) {
		b.WriteByte(',')
		b.Write(body[1:])
	} else {
		b.WriteByte('}')
	}
	return b.Bytes(), nil
}

// marshalWithTrueField marshals v (a JSON object) and injects a leading
// `"<field>":true`, for objects whose field must always be true.
func marshalWithTrueField(field string, v any) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(body) < 2 || body[0] != '{' {
		return nil, fmt.Errorf("marshal %s: expected JSON object, got %s", field, body)
	}
	name, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	b.WriteByte('{')
	b.Write(name)
	b.WriteString(":true")
	if !bytes.Equal(body, []byte("{}")) {
		b.WriteByte(',')
		b.Write(body[1:])
	} else {
		b.WriteByte('}')
	}
	return b.Bytes(), nil
}

// ---------------------------------------------------------------------------
// RichText
// ---------------------------------------------------------------------------

// RichText is a rich formatted text: plain text (RichTextPlain), a list of
// rich texts (RichTexts) or one of the typed RichText* objects.
// Use RichTextRaw for types that are not modelled here.
type RichText interface {
	richText()
}

var (
	_ RichText = RichTextPlain("")
	_ RichText = RichTexts(nil)
	_ RichText = RichTextRaw(nil)
	_ RichText = RichTextBold{}
	_ RichText = RichTextItalic{}
	_ RichText = RichTextUnderline{}
	_ RichText = RichTextStrikethrough{}
	_ RichText = RichTextSpoiler{}
	_ RichText = RichTextMarked{}
	_ RichText = RichTextCode{}
	_ RichText = RichTextURL{}
	_ RichText = RichTextEmailAddress{}
	_ RichText = RichTextCustomEmoji{}
	_ RichText = RichTextButton{}
	_ RichText = RichTextAnchor{}
	_ RichText = RichTextAnchorLink{}
	_ RichText = RichTextReference{}
	_ RichText = RichTextReferenceLink{}
)

// RichTextPlain is plain text; it is serialized as a JSON string.
type RichTextPlain string

func (RichTextPlain) richText() {}

// RichTexts is a concatenation of rich texts; it is serialized as a JSON array.
type RichTexts []RichText

func (RichTexts) richText() {}

// RichTextRaw is an escape hatch for rich text types that are not modelled.
type RichTextRaw json.RawMessage

func (RichTextRaw) richText() {}

func (r RichTextRaw) MarshalJSON() ([]byte, error) {
	if len(r) == 0 {
		return []byte("null"), nil
	}
	return r, nil
}

type RichTextBold struct {
	Text RichText `json:"text"`
}

func (RichTextBold) richText() {}

func (r RichTextBold) MarshalJSON() ([]byte, error) {
	type alias RichTextBold
	return marshalTyped("bold", alias(r))
}

type RichTextItalic struct {
	Text RichText `json:"text"`
}

func (RichTextItalic) richText() {}

func (r RichTextItalic) MarshalJSON() ([]byte, error) {
	type alias RichTextItalic
	return marshalTyped("italic", alias(r))
}

type RichTextUnderline struct {
	Text RichText `json:"text"`
}

func (RichTextUnderline) richText() {}

func (r RichTextUnderline) MarshalJSON() ([]byte, error) {
	type alias RichTextUnderline
	return marshalTyped("underline", alias(r))
}

type RichTextStrikethrough struct {
	Text RichText `json:"text"`
}

func (RichTextStrikethrough) richText() {}

func (r RichTextStrikethrough) MarshalJSON() ([]byte, error) {
	type alias RichTextStrikethrough
	return marshalTyped("strikethrough", alias(r))
}

type RichTextSpoiler struct {
	Text RichText `json:"text"`
}

func (RichTextSpoiler) richText() {}

func (r RichTextSpoiler) MarshalJSON() ([]byte, error) {
	type alias RichTextSpoiler
	return marshalTyped("spoiler", alias(r))
}

type RichTextMarked struct {
	Text RichText `json:"text"`
}

func (RichTextMarked) richText() {}

func (r RichTextMarked) MarshalJSON() ([]byte, error) {
	type alias RichTextMarked
	return marshalTyped("marked", alias(r))
}

type RichTextCode struct {
	Text RichText `json:"text"`
}

func (RichTextCode) richText() {}

func (r RichTextCode) MarshalJSON() ([]byte, error) {
	type alias RichTextCode
	return marshalTyped("code", alias(r))
}

// RichTextURL is a text with a link.
type RichTextURL struct {
	Text RichText `json:"text"`
	URL  string   `json:"url"`
}

func (RichTextURL) richText() {}

func (r RichTextURL) MarshalJSON() ([]byte, error) {
	type alias RichTextURL
	return marshalTyped("url", alias(r))
}

type RichTextEmailAddress struct {
	Text         RichText `json:"text"`
	EmailAddress string   `json:"email_address"`
}

func (RichTextEmailAddress) richText() {}

func (r RichTextEmailAddress) MarshalJSON() ([]byte, error) {
	type alias RichTextEmailAddress
	return marshalTyped("email_address", alias(r))
}

type RichTextCustomEmoji struct {
	CustomEmojiID   string `json:"custom_emoji_id"`
	AlternativeText string `json:"alternative_text"`
}

func (RichTextCustomEmoji) richText() {}

func (r RichTextCustomEmoji) MarshalJSON() ([]byte, error) {
	type alias RichTextCustomEmoji
	return marshalTyped("custom_emoji", alias(r))
}

// RichTextButton is an inline button inside rich text.
type RichTextButton struct {
	Button RichMessageButton `json:"button"`
}

func (RichTextButton) richText() {}

func (r RichTextButton) MarshalJSON() ([]byte, error) {
	type alias RichTextButton
	return marshalTyped("button", alias(r))
}

type RichTextAnchor struct {
	Name string `json:"name"`
}

func (RichTextAnchor) richText() {}

func (r RichTextAnchor) MarshalJSON() ([]byte, error) {
	type alias RichTextAnchor
	return marshalTyped("anchor", alias(r))
}

// RichTextAnchorLink links to an anchor; an empty AnchorName links to the top of the message.
type RichTextAnchorLink struct {
	Text       RichText `json:"text"`
	AnchorName string   `json:"anchor_name"`
}

func (RichTextAnchorLink) richText() {}

func (r RichTextAnchorLink) MarshalJSON() ([]byte, error) {
	type alias RichTextAnchorLink
	return marshalTyped("anchor_link", alias(r))
}

type RichTextReference struct {
	Text RichText `json:"text"`
	Name string   `json:"name"`
}

func (RichTextReference) richText() {}

func (r RichTextReference) MarshalJSON() ([]byte, error) {
	type alias RichTextReference
	return marshalTyped("reference", alias(r))
}

type RichTextReferenceLink struct {
	Text          RichText `json:"text"`
	ReferenceName string   `json:"reference_name"`
}

func (RichTextReferenceLink) richText() {}

func (r RichTextReferenceLink) MarshalJSON() ([]byte, error) {
	type alias RichTextReferenceLink
	return marshalTyped("reference_link", alias(r))
}

// ---------------------------------------------------------------------------
// Buttons
// ---------------------------------------------------------------------------

// RichMessageButton is a button in a rich message. Exactly one of the fields
// other than Text and Style must be set.
type RichMessageButton struct {
	Text         RichText        `json:"text"`
	Style        string          `json:"style,omitempty"` // danger, success, primary or link
	URL          string          `json:"url,omitempty"`
	CallbackData string          `json:"callback_data,omitempty"`
	CopyText     *CopyTextButton `json:"copy_text,omitempty"`
	Disabled     *DisabledButton `json:"disabled,omitempty"`
}

type CopyTextButton struct {
	Text string `json:"text"`
}

// DisabledButton marks a button as disabled. It currently holds no information.
type DisabledButton struct{}

// ---------------------------------------------------------------------------
// InputRichBlock
// ---------------------------------------------------------------------------

// InputRichBlock is a block of a rich message to be sent.
// Use InputRichBlockRaw for block types that are not modelled here.
type InputRichBlock interface {
	inputRichBlock()
}

var (
	_ InputRichBlock = InputRichBlockRaw(nil)
	_ InputRichBlock = InputRichBlockParagraph{}
	_ InputRichBlock = InputRichBlockSectionHeading{}
	_ InputRichBlock = InputRichBlockPreformatted{}
	_ InputRichBlock = InputRichBlockFooter{}
	_ InputRichBlock = InputRichBlockDivider{}
	_ InputRichBlock = InputRichBlockAnchor{}
	_ InputRichBlock = InputRichBlockList{}
	_ InputRichBlock = InputRichBlockBlockQuotation{}
	_ InputRichBlock = InputRichBlockExpandableBlockQuotation{}
	_ InputRichBlock = InputRichBlockDetails{}
	_ InputRichBlock = InputRichBlockButtons{}
	_ InputRichBlock = InputRichBlockAudio{}
	_ InputRichBlock = InputRichBlockVoiceNote{}
	_ InputRichBlock = InputRichBlockPhoto{}
	_ InputRichBlock = InputRichBlockThinking{}
)

// InputRichBlockRaw is an escape hatch for block types that are not modelled.
type InputRichBlockRaw json.RawMessage

func (InputRichBlockRaw) inputRichBlock() {}

func (r InputRichBlockRaw) MarshalJSON() ([]byte, error) {
	if len(r) == 0 {
		return []byte("null"), nil
	}
	return r, nil
}

// RichBlockCaption is the caption of a rich block.
type RichBlockCaption struct {
	Text   RichText `json:"text"`
	Credit RichText `json:"credit,omitempty"`
}

// InputRichBlockParagraph corresponds to <p>.
type InputRichBlockParagraph struct {
	Text RichText `json:"text"`
}

func (InputRichBlockParagraph) inputRichBlock() {}

func (b InputRichBlockParagraph) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockParagraph
	return marshalTyped("paragraph", alias(b))
}

// InputRichBlockSectionHeading corresponds to <h1>…<h6>. Size is 1-6, 1 is the largest.
type InputRichBlockSectionHeading struct {
	Text RichText `json:"text"`
	Size int      `json:"size"`
}

func (InputRichBlockSectionHeading) inputRichBlock() {}

func (b InputRichBlockSectionHeading) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockSectionHeading
	return marshalTyped("heading", alias(b))
}

// InputRichBlockPreformatted corresponds to <pre><code>.
type InputRichBlockPreformatted struct {
	Text     RichText `json:"text"`
	Language string   `json:"language,omitempty"`
}

func (InputRichBlockPreformatted) inputRichBlock() {}

func (b InputRichBlockPreformatted) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockPreformatted
	return marshalTyped("pre", alias(b))
}

// InputRichBlockFooter corresponds to <footer>.
type InputRichBlockFooter struct {
	Text RichText `json:"text"`
}

func (InputRichBlockFooter) inputRichBlock() {}

func (b InputRichBlockFooter) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockFooter
	return marshalTyped("footer", alias(b))
}

// InputRichBlockDivider corresponds to <hr/>.
type InputRichBlockDivider struct{}

func (InputRichBlockDivider) inputRichBlock() {}

func (b InputRichBlockDivider) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockDivider
	return marshalTyped("divider", alias(b))
}

// InputRichBlockAnchor corresponds to <a name="…">.
type InputRichBlockAnchor struct {
	Name string `json:"name"`
}

func (InputRichBlockAnchor) inputRichBlock() {}

func (b InputRichBlockAnchor) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockAnchor
	return marshalTyped("anchor", alias(b))
}

// InputRichBlockList corresponds to <ul>/<ol>.
type InputRichBlockList struct {
	Items []InputRichBlockListItem `json:"items"`
}

func (InputRichBlockList) inputRichBlock() {}

func (b InputRichBlockList) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockList
	return marshalTyped("list", alias(b))
}

// InputRichBlockListItem is an item of a list to be sent.
type InputRichBlockListItem struct {
	Blocks      []InputRichBlock `json:"blocks"`
	HasCheckbox bool             `json:"has_checkbox,omitempty"`
	IsChecked   bool             `json:"is_checked,omitempty"`
	Value       int              `json:"value,omitempty"`
	Type        string           `json:"type,omitempty"` // a, A, i, I or 1 (ordered lists)
}

// InputRichBlockBlockQuotation corresponds to <blockquote>.
type InputRichBlockBlockQuotation struct {
	Blocks []InputRichBlock `json:"blocks"`
	Credit RichText         `json:"credit,omitempty"`
}

func (InputRichBlockBlockQuotation) inputRichBlock() {}

func (b InputRichBlockBlockQuotation) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockBlockQuotation
	return marshalTyped("blockquote", alias(b))
}

// InputRichBlockExpandableBlockQuotation corresponds to <blockquote expandable>.
type InputRichBlockExpandableBlockQuotation struct {
	Text   RichText `json:"text"`
	Credit RichText `json:"credit,omitempty"`
}

func (InputRichBlockExpandableBlockQuotation) inputRichBlock() {}

func (b InputRichBlockExpandableBlockQuotation) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockExpandableBlockQuotation
	return marshalTyped("expandable_blockquote", alias(b))
}

// InputRichBlockDetails corresponds to <details>; collapsed unless IsOpen.
type InputRichBlockDetails struct {
	Summary RichText         `json:"summary"`
	Blocks  []InputRichBlock `json:"blocks"`
	IsOpen  bool             `json:"is_open,omitempty"`
}

func (InputRichBlockDetails) inputRichBlock() {}

func (b InputRichBlockDetails) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockDetails
	return marshalTyped("details", alias(b))
}

// InputRichBlockButtons is a row of 1-8 buttons (<tg-button-row>).
type InputRichBlockButtons struct {
	Buttons []RichMessageButton `json:"buttons"`
	Align   string              `json:"align,omitempty"` // left, center or right
}

func (InputRichBlockButtons) inputRichBlock() {}

func (b InputRichBlockButtons) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockButtons
	return marshalTyped("buttons", alias(b))
}

// InputRichBlockAudio is a music file block (<audio>).
type InputRichBlockAudio struct {
	Audio   InputMediaAudio   `json:"audio"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

func (InputRichBlockAudio) inputRichBlock() {}

func (b InputRichBlockAudio) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockAudio
	return marshalTyped("audio", alias(b))
}

// InputRichBlockVoiceNote is a voice note block (<audio>).
type InputRichBlockVoiceNote struct {
	VoiceNote InputMediaVoiceNote `json:"voice_note"`
	Caption   *RichBlockCaption   `json:"caption,omitempty"`
}

func (InputRichBlockVoiceNote) inputRichBlock() {}

func (b InputRichBlockVoiceNote) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockVoiceNote
	return marshalTyped("voice_note", alias(b))
}

// InputRichBlockPhoto is a photo block (<img>).
type InputRichBlockPhoto struct {
	Photo   InputMediaPhoto   `json:"photo"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

func (InputRichBlockPhoto) inputRichBlock() {}

func (b InputRichBlockPhoto) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockPhoto
	return marshalTyped("photo", alias(b))
}

// InputRichBlockThinking is a "Thinking…" placeholder (<tg-thinking>).
// It may be used only in sendRichMessageDraft.
type InputRichBlockThinking struct {
	Text RichText `json:"text"`
}

func (InputRichBlockThinking) inputRichBlock() {}

func (b InputRichBlockThinking) MarshalJSON() ([]byte, error) {
	type alias InputRichBlockThinking
	return marshalTyped("thinking", alias(b))
}

// ---------------------------------------------------------------------------
// InputMedia
// ---------------------------------------------------------------------------

// InputMedia is a media element to be sent.
type InputMedia interface {
	inputMedia()
}

var (
	_ InputMedia = InputMediaAudio{}
	_ InputMedia = InputMediaVoiceNote{}
	_ InputMedia = InputMediaPhoto{}
)

// InputMediaAudio is an audio file to be treated as music.
// Media is a file_id, an HTTP URL or "attach://<name>" for a multipart upload.
type InputMediaAudio struct {
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Performer string `json:"performer,omitempty"`
	Title     string `json:"title,omitempty"`
}

func (InputMediaAudio) inputMedia() {}

func (m InputMediaAudio) MarshalJSON() ([]byte, error) {
	type alias InputMediaAudio
	return marshalTyped("audio", alias(m))
}

// InputMediaVoiceNote is a voice message file.
// Media is a file_id, an HTTP URL or "attach://<name>" for a multipart upload.
type InputMediaVoiceNote struct {
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
	Duration  int    `json:"duration,omitempty"`
}

func (InputMediaVoiceNote) inputMedia() {}

func (m InputMediaVoiceNote) MarshalJSON() ([]byte, error) {
	type alias InputMediaVoiceNote
	return marshalTyped("voice_note", alias(m))
}

// InputMediaPhoto is a photo.
// Media is a file_id, an HTTP URL or "attach://<name>" for a multipart upload.
type InputMediaPhoto struct {
	Media      string `json:"media"`
	Caption    string `json:"caption,omitempty"`
	ParseMode  string `json:"parse_mode,omitempty"`
	HasSpoiler bool   `json:"has_spoiler,omitempty"`
}

func (InputMediaPhoto) inputMedia() {}

func (m InputMediaPhoto) MarshalJSON() ([]byte, error) {
	type alias InputMediaPhoto
	return marshalTyped("photo", alias(m))
}

// InputRichMessageMedia is a media element referenced from rich markdown/html
// via tg://photo?id=, tg://video?id=, tg://document?id= or tg://audio?id= links.
type InputRichMessageMedia struct {
	ID    string     `json:"id"`
	Media InputMedia `json:"media"`
}

// AttachURI returns the "attach://<name>" reference for a multipart upload part.
func AttachURI(name string) string {
	return "attach://" + name
}
