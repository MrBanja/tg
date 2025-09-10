package tg

import (
	"context"
	"strings"

	"github.com/mrbanja/tg/v2/model"
)

type Filter func(ctx context.Context, u *model.Update) bool

// MessageFilter should be used with FilterMessage func
type MessageFilter func(ctx context.Context, u *model.Message) bool

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

	return func(ctx context.Context, u *model.Update) bool {
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

func FilterOr(filters ...Filter) Filter {
	return func(ctx context.Context, u *model.Update) bool {
		for _, f := range filters {
			if f(ctx, u) {
				return true
			}
		}
		return false
	}
}
func FilterMessage(filters ...MessageFilter) Filter {
	return func(ctx context.Context, u *model.Update) bool {
		if u.Message == nil {
			return false
		}
		for _, f := range filters {
			if !f(ctx, u.Message) {
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
