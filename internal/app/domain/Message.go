package domain

type SingleMessage struct {
	Id      string          `json:"id"`
	Text    string          `json:"text"`
	IsOwner bool            `json:"isOwner"`
	ToPhone string          `json:"toPhone"`
	Status  []MessageStatus `json:"status,omitempty"`
}

type MessageStatus struct {
	ReceivedAt float64 `json:"receivedAt"`
	Status     string  `json:"status"`
}
