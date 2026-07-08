package model

import (
	"encoding/json"
	"testing"
)

func TestSetMessageReactionRequestMarshalIncludesEmptyReactionArray(t *testing.T) {
	data, err := json.Marshal(SetMessageReactionRequest{
		ChatID:    1,
		MessageID: 2,
		Reaction:  []ReactionTypeEmoji{},
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	want := `{"chat_id":1,"message_id":2,"reaction":[]}`
	if string(data) != want {
		t.Fatalf("unexpected json: %s", string(data))
	}
}

func TestSendMessageRequestMarshalIncludesLinkPreviewOptions(t *testing.T) {
	data, err := json.Marshal(SendMessageRequest{
		ChatID: 1,
		Text:   "hello",
		LinkPreviewOptions: &LinkPreviewOptions{
			IsDisabled: true,
		},
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	want := `{"chat_id":1,"text":"hello","link_preview_options":{"is_disabled":true}}`
	if string(data) != want {
		t.Fatalf("unexpected json: %s", string(data))
	}
}

func TestSendRichMessageRequestMarshal(t *testing.T) {
	data, err := json.Marshal(SendRichMessageRequest{
		ChatID: 1,
		RichMessage: InputRichMessage{
			HTML: "<p>hello</p>",
		},
		ReplyParameters: &ReplyParameters{MessageID: 2},
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	want := `{"chat_id":1,"rich_message":{"html":"\u003cp\u003ehello\u003c/p\u003e"},"reply_parameters":{"message_id":2}}`
	if string(data) != want {
		t.Fatalf("unexpected json: %s", string(data))
	}
}

func TestEditMessageTextRequestMarshalRichMessage(t *testing.T) {
	chatID := int64(1)
	messageID := int64(2)
	data, err := json.Marshal(EditMessageTextRequest{
		ChatID:    &chatID,
		MessageID: &messageID,
		RichMessage: &InputRichMessage{
			Markdown: "hello",
		},
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	want := `{"chat_id":1,"message_id":2,"rich_message":{"markdown":"hello"}}`
	if string(data) != want {
		t.Fatalf("unexpected json: %s", string(data))
	}
}
