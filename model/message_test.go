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
