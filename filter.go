package tg

import (
	"context"
	"strings"

	"github.com/mrbanja/tg/model"
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
