package tgmodel

type Location struct {
	Location Loc    `json:"location"`
	Address  string `json:"address"`
}

type Loc struct {
	Longitude            int `json:"longitude"`
	Latitude             int `json:"latitude"`
	HorizontalAccuracy   int `json:"horizontal_accuracy"`
	LivePeriod           int `json:"live_period"`
	Heading              int `json:"heading"`
	ProximityAlertRadius int `json:"proximity_alert_radius"`
}
