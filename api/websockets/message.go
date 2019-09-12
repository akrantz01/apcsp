package websockets

const (
	MessageAuthentication = iota
	MessageReceive
	MessageSent
)

type BaseMessage struct {
	Type int `json:"type"`
}

type AuthenticationMessage struct {
	Type  int    `json:"type"`
	Token string `json:"token"`
}

type ReceiveMessage struct {
	Type        int    `json:"type"`
	Message     string `json:"message"`
	Chat        string `json:"chat"`
	Sender      string `json:"sender"`
	ContentType int    `json:"content-type"`
}

type SentMessage struct {
	Type        int    `json:"type"`
	Chat        string `json:"chat"`
	Message     string `json:"message"`
	Filename    string `json:"filename"`
	ContentType string `json:"content-type"`
}
