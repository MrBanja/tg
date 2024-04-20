package server

import "strings"

type Opt struct {
	addr     string
	whPath   string
	whSecret string

	tgToken string
	selfURL string
}

func Options() *Opt {
	return &Opt{
		addr:   ":8080",
		whPath: "/webhook",
	}
}

func (o *Opt) Addr(addr string) *Opt {
	o.addr = addr
	return o
}

func (o *Opt) WebhookPath(path string) *Opt {
	o.whPath = "/" + strings.TrimPrefix(path, "/")
	return o
}

func (o *Opt) WebhookSecret(secret string) *Opt {
	o.whSecret = secret
	return o
}

func (o *Opt) SetWebhook(tgToken string, selfURL string) *Opt {
	o.tgToken = tgToken
	o.selfURL = strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(selfURL, "https://"), "http://"), "/")
	return o
}
