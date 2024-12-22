package domain

type SingleMessage struct {
	Id      string          `json:"id"`
	Text    string          `json:"text"`
	IsOwner bool            `json:"isOwner"`
	Status  []MessageStatus `json:"status"`
}

type MessageStatus struct {
	ReceivedAt float64 `json:"receivedAt"`
	Status     string  `json:"status"`
}
