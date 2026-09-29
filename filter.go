package tg

import (
	"context"
	"strings"

	"github.com/mrbanja/tg/v4/model"
)

type Filter func(ctx context.Context, req *Request) bool

// MessageFilter should be used with FilterMessage func
type MessageFilter func(ctx context.Context, m *model.Message) bool

func FilterReplyMessage(replyMessageFilters ...MessageFilter) MessageFilter {
	return func(ctx context.Context, m *model.Message) bool {
		if m.ReplyToMessage == nil {
			return false
		}

		for _, f := range replyMessageFilters {
			if !f(ctx, m.ReplyToMessage) {
				return false
			}
		}
		return true
	}
}

func FilterCommand(cmd string) Filter {
	cmd = "/" + strings.ToLower(strings.TrimPrefix(cmd, "/"))

	return func(ctx context.Context, req *Request) bool {
		u := req.Update
		if u.Message == nil {
			return false
		}
		if len(u.Message.Entities) == 0 {
			return false
		}
		if u.Message.Text == nil {
			return false
		}
		for _, e := range u.Message.Entities {
			if e.Type == "bot_command" && e.Offset == 0 && (*u.Message.Text)[e.Offset:e.Offset+e.Length] == cmd {
				return true
			}
		}
		return false
	}
}

// FilterCommandWithOptionalBotName filters command with optional bot tag: /text@test_bot
func FilterCommandWithOptionalBotName(cmd string, botTag string) Filter {
	return func(ctx context.Context, req *Request) bool {
		u := req.Update
		if u.Message == nil {
			return false
		}
		if len(u.Message.Entities) == 0 {
			return false
		}
		if u.Message.Text == nil {
			return false
		}
		text := *u.Message.Text

		for _, e := range u.Message.Entities {
			if e.Type == "bot_command" &&
				e.Offset == 0 &&
				((text)[e.Offset:e.Offset+e.Length] == cmd ||
					(text)[e.Offset:e.Offset+e.Length] == cmd+botTag) {
				return true
			}
		}
		return false
	}
}

func FilterOr(filters ...Filter) Filter {
	return func(ctx context.Context, req *Request) bool {
		for _, f := range filters {
			if f(ctx, req) {
				return true
			}
		}
		return false
	}
}
func FilterMessage(filters ...MessageFilter) Filter {
	return func(ctx context.Context, req *Request) bool {
		if req.Update.Message == nil {
			return false
		}
		for _, f := range filters {
			if !f(ctx, req.Update.Message) {
				return false
			}
		}
		return true
	}
}

func FilterMessageOr(filters ...MessageFilter) MessageFilter {
	return func(ctx context.Context, m *model.Message) bool {
		for _, f := range filters {
			if f(ctx, m) {
				return true
			}
		}
		return false
	}
}

func FilterContainsMedia(ctx context.Context, m *model.Message) bool {
	return FilterMessageOr(
		FilterContainsPhoto,
		FilterContainsSticker,
		FilterContainsDocument,
		FilterContainsAnimation,
		FilterContainsVideo,
	)(ctx, m)
}

func FilterContainsPhoto(_ context.Context, m *model.Message) bool {
	return len(m.Photo) != 0
}

func FilterContainsSticker(_ context.Context, m *model.Message) bool {
	return m.Sticker != nil
}

func FilterContainsDocument(_ context.Context, m *model.Message) bool {
	return m.Document != nil
}

func FilterContainsAnimation(_ context.Context, m *model.Message) bool {
	return m.Animation != nil
}

func FilterContainsVideo(_ context.Context, m *model.Message) bool {
	return m.Video != nil
}

// FilterCallbackQuery matches callback queries whose data passes all filters.
func FilterCallbackQuery(filters ...func(ctx context.Context, q *model.CallbackQuery) bool) Filter {
	return func(ctx context.Context, req *Request) bool {
		q := req.Update.CallbackQuery
		if q == nil {
			return false
		}
		for _, f := range filters {
			if !f(ctx, q) {
				return false
			}
		}
		return true
	}
}

// FilterCallbackDataPrefix matches callback queries whose data starts with prefix.
func FilterCallbackDataPrefix(prefix string) Filter {
	return FilterCallbackQuery(func(_ context.Context, q *model.CallbackQuery) bool {
		return strings.HasPrefix(q.Data, prefix)
	})
}

// FilterStoppedMessageGeneration matches stopped_message_generation updates.
func FilterStoppedMessageGeneration(_ context.Context, req *Request) bool {
	return req.Update.StoppedMessageGeneration != nil
}

// FilterPrivateChat matches updates that come from a private chat.
func FilterPrivateChat(_ context.Context, req *Request) bool {
	u := req.Update
	switch {
	case u.Message != nil:
		return u.Message.Chat != nil && u.Message.Chat.Type == "private"
	case u.CallbackQuery != nil:
		m := u.CallbackQuery.Message
		return m != nil && m.Chat != nil && m.Chat.Type == "private"
	case u.StoppedMessageGeneration != nil:
		return u.StoppedMessageGeneration.Chat.Type == "private"
	}
	return false
}

// FilterText matches messages with text that is not a bot command.
func FilterText(_ context.Context, m *model.Message) bool {
	if m.Text == nil {
		return false
	}
	for _, e := range m.Entities {
		if e.Type == "bot_command" && e.Offset == 0 {
			return false
		}
	}
	return true
}

func FilterContainsVoice(_ context.Context, m *model.Message) bool {
	return m.Voice != nil
}

func FilterContainsAudio(_ context.Context, m *model.Message) bool {
	return m.Audio != nil
}

// FilterContainsLocation matches shared locations, including venues.
func FilterContainsLocation(_ context.Context, m *model.Message) bool {
	return m.Location != nil || m.Venue != nil
}

func FilterContainsVenue(_ context.Context, m *model.Message) bool {
	return m.Venue != nil
}

// FilterMediaGroup matches messages that belong to a media group (album).
func FilterMediaGroup(_ context.Context, m *model.Message) bool {
	return m.MediaGroupID != nil
}
