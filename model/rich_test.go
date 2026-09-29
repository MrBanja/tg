package model

import (
	"encoding/json"
	"testing"
)

func mustMarshal(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	return string(data)
}

func TestInputRichMessageBlocksMarshal(t *testing.T) {
	msg := InputRichMessage{
		Blocks: []InputRichBlock{
			InputRichBlockSectionHeading{Text: RichTextPlain("Title"), Size: 2},
			InputRichBlockVoiceNote{VoiceNote: InputMediaVoiceNote{Media: AttachURI("voice"), Duration: 12}},
			InputRichBlockDetails{
				Summary: RichTextPlain("Text"),
				Blocks: []InputRichBlock{
					InputRichBlockParagraph{Text: RichTexts{
						RichTextPlain("see "),
						RichTextURL{Text: RichTextPlain("[1]"), URL: "https://example.com"},
					}},
				},
			},
			InputRichBlockDivider{},
			InputRichBlockFooter{Text: RichTextItalic{Text: RichTextPlain("Sources")}},
			InputRichBlockButtons{Buttons: []RichMessageButton{
				{Text: RichTextPlain("More"), CallbackData: "fu:1"},
			}},
		},
	}

	want := `{"blocks":[` +
		`{"type":"heading","text":"Title","size":2},` +
		`{"type":"voice_note","voice_note":{"type":"voice_note","media":"attach://voice","duration":12}},` +
		`{"type":"details","summary":"Text","blocks":[{"type":"paragraph","text":["see ",{"type":"url","text":"[1]","url":"https://example.com"}]}]},` +
		`{"type":"divider"},` +
		`{"type":"footer","text":{"type":"italic","text":"Sources"}},` +
		`{"type":"buttons","buttons":[{"text":"More","callback_data":"fu:1"}]}` +
		`]}`
	if got := mustMarshal(t, msg); got != want {
		t.Fatalf("unexpected json:\n got: %s\nwant: %s", got, want)
	}
}

func TestInputRichBlockThinkingAndRawMarshal(t *testing.T) {
	blocks := []InputRichBlock{
		InputRichBlockThinking{Text: RichTextPlain("Searching…")},
		InputRichBlockRaw(`{"type":"map","location":{"latitude":1,"longitude":2},"zoom":10}`),
	}
	want := `[{"type":"thinking","text":"Searching…"},{"type":"map","location":{"latitude":1,"longitude":2},"zoom":10}]`
	if got := mustMarshal(t, blocks); got != want {
		t.Fatalf("unexpected json: %s", got)
	}
}

func TestInputRichMessageMediaMarshal(t *testing.T) {
	msg := InputRichMessage{
		Markdown: "![](tg://audio?id=a1)",
		Media: []InputRichMessageMedia{
			{ID: "a1", Media: InputMediaAudio{Media: AttachURI("a1"), Title: "Guide", Performer: "Bot"}},
		},
	}
	want := `{"markdown":"![](tg://audio?id=a1)","media":[{"id":"a1","media":{"type":"audio","media":"attach://a1","performer":"Bot","title":"Guide"}}]}`
	if got := mustMarshal(t, msg); got != want {
		t.Fatalf("unexpected json: %s", got)
	}
}

func TestSendRichMessageDraftRequestMarshalCanStop(t *testing.T) {
	req := SendRichMessageDraftRequest{
		ChatID:      1,
		DraftID:     7,
		RichMessage: InputRichMessage{Markdown: "thinking"},
		CanStop:     true,
	}
	want := `{"chat_id":1,"draft_id":7,"rich_message":{"markdown":"thinking"},"can_stop":true}`
	if got := mustMarshal(t, req); got != want {
		t.Fatalf("unexpected json: %s", got)
	}
}

func TestSendMessageDraftRequestMarshalKeepsEmptyText(t *testing.T) {
	want := `{"chat_id":1,"draft_id":2,"text":""}`
	if got := mustMarshal(t, SendMessageDraftRequest{ChatID: 1, DraftID: 2}); got != want {
		t.Fatalf("unexpected json: %s", got)
	}
}
