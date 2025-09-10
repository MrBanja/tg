package tg

import "errors"

var (
	ErrWebhookAlreadySetToDiffAddress = errors.New("webhook is already set to different url")
)
