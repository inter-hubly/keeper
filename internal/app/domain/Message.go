package domain

type Conversations struct {
	ProfileName string    `json:"profileName"`
	Messages    []Message `json:"messages"`
}

type Message struct {
	Status  []MessageStatus `json:"status,omitempty"`
	Id      string          `json:"id"`
	Text    string          `json:"text"`
	ToPhone string          `json:"toPhone"`
	IsOwner bool            `json:"isOwner"`
}

type MessageStatus struct {
	ReceivedAt float64 `json:"receivedAt"`
	Status     string  `json:"status"`
}
