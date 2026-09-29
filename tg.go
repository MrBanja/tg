package tg

import "net/http"

var (
	token      string
	apiBaseURL = "https://api.telegram.org"
	httpClient = http.DefaultClient
)

func SetToken(t string) {
	token = t
}

// SetHTTPClient replaces the HTTP client used for Bot API requests
// (http.DefaultClient by default, which has no timeout).
func SetHTTPClient(c *http.Client) {
	if c == nil {
		c = http.DefaultClient
	}
	httpClient = c
}

// SetAPIBaseURL points the client to a different Bot API server, e.g. a local
// Bot API server or a test server. Defaults to https://api.telegram.org.
func SetAPIBaseURL(u string) {
	apiBaseURL = u
}

func ToOptional[T any](t T) *T {
	return &t
}
