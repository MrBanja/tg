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
