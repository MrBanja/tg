package tgmodel

type Poll struct {
	ID                    string                  `json:"id"`
	Question              string                  `json:"question"`
	Options               []PollOption            `json:"options,omitempty"`
	TotalVoterCount       int                     `json:"total_voter_count"`
	IsClosed              bool                    `json:"is_closed"`
	IsAnonymous           bool                    `json:"is_anonymous"`
	Type                  string                  `json:"type"`
	AllowsMultipleAnswers bool                    `json:"allows_multiple_answers"`
	CorrectOptionID       int                     `json:"correct_option_id"`
	Explanation           string                  `json:"explanation"`
	ExplanationEntities   []PoolExplanationEntity `json:"explanation_entities,omitempty"`
	OpenPeriod            int                     `json:"open_period"`
	CloseDate             int                     `json:"close_date"`
}

type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

type PoolExplanationEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	URL      string `json:"url"`
	User     *User  `json:"user,omitempty"`
	Language string `json:"language"`
}
