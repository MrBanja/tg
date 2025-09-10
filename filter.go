package tg

import (
	"context"
	"strings"

	"github.com/mrbanja/tg/v2/model"
)

type Filter func(ctx context.Context, u *model.Update) bool

func FilterCommand(cmd string) Filter {
	cmd = "/" + strings.ToLower(strings.TrimPrefix(cmd, "/"))

	return func(ctx context.Context, u *model.Update) bool {
		if u.Message == nil {
			return false
		}
		if len(u.Message.Entities) == 0 {
			return false
		}
		if u.Message.Text == "" {
			return false
		}
		for _, e := range u.Message.Entities {
			if e.Type == "bot_command" && e.Offset == 0 && u.Message.Text[e.Offset:e.Offset+e.Length] == cmd {
				return true
			}
		}
		return false
	}
}
