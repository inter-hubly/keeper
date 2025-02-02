package domain

type Conversations struct {
	WhatsAppProfileName string    `json:"whatsAppProfileName"`
	LocalProfileName    string    `json:"localProfileName"`
	Messages            []Message `json:"messages"`
}

type Message struct {
	Status  []MessageStatus `json:"status,omitempty"`
	Id      string          `json:"id"`
	Text    string          `json:"text"`
	IsOwner bool            `json:"isOwner"`
}

type MessageStatus struct {
	ReceivedAt float64 `json:"receivedAt"`
	Status     string  `json:"status"`
}
