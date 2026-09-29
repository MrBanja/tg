package model

import (
	"encoding/json"
	"testing"
)

func TestUpdateUnmarshalNewFields(t *testing.T) {
	raw := `{"update_id":1,"message":{"message_id":5,"date":1,"chat":{"id":9,"type":"private"},
		"media_group_id":"g1","caption":"hi","caption_entities":[{"type":"bold","offset":0,"length":2}],
		"voice":{"file_id":"v","file_unique_id":"vu","duration":3,"mime_type":"audio/ogg","file_size":100},
		"venue":{"location":{"latitude":59.94,"longitude":30.31},"title":"Hermitage","address":"Palace Sq"},
		"location":{"latitude":59.94,"longitude":30.31,"horizontal_accuracy":12.5}}}`
	var u Update
	if err := json.Unmarshal([]byte(raw), &u); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	m := u.Message
	if m.MediaGroupID == nil || *m.MediaGroupID != "g1" {
		t.Fatalf("media_group_id not decoded: %+v", m.MediaGroupID)
	}
	if len(m.CaptionEntities) != 1 || m.CaptionEntities[0].Type != "bold" {
		t.Fatalf("caption_entities not decoded: %+v", m.CaptionEntities)
	}
	if m.Voice == nil || m.Voice.FileID != "v" || m.Voice.Duration != 3 {
		t.Fatalf("voice not decoded: %+v", m.Voice)
	}
	if m.Venue == nil || m.Venue.Title != "Hermitage" || m.Venue.Location.Latitude != 59.94 {
		t.Fatalf("venue not decoded: %+v", m.Venue)
	}
	if m.Location == nil || m.Location.HorizontalAccuracy == nil || *m.Location.HorizontalAccuracy != 12.5 {
		t.Fatalf("location not decoded: %+v", m.Location)
	}
}

func TestUpdateUnmarshalStoppedMessageGeneration(t *testing.T) {
	raw := `{"update_id":2,"stopped_message_generation":{"chat":{"id":9,"type":"private"},"draft_id":42}}`
	var u Update
	if err := json.Unmarshal([]byte(raw), &u); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	s := u.StoppedMessageGeneration
	if s == nil || s.DraftID != 42 || s.Chat.ID != 9 {
		t.Fatalf("stopped_message_generation not decoded: %+v", s)
	}
}
