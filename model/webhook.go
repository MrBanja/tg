package model

type WebhookInfo struct {
	URL                  string `json:"url"`
	HasCustomCertificate bool   `json:"has_custom_certificate"`
	PendingUpdateCount   int    `json:"pending_update_count"`
}

type SetWebhookRequest struct {
	URL                string   `json:"url"`                            // HTTPS URL to send updates to. Use an empty string to remove webhook integration
	IPAddress          string   `json:"ip_address,omitempty"`           // The fixed IP address for webhook requests instead of resolving through DNS
	MaxConnections     int      `json:"max_connections,omitempty"`      // Maximum allowed simultaneous HTTPS connections (1-100), defaults to 40
	AllowedUpdates     []string `json:"allowed_updates,omitempty"`      // List of update types to receive
	DropPendingUpdates bool     `json:"drop_pending_updates,omitempty"` // Pass true to drop all pending updates
	SecretToken        string   `json:"secret_token,omitempty"`         // Secret token (1-256 chars: A-Z, a-z, 0-9, _, -) sent in header
}

type DeleteWebhookRequest struct {
	DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
}
