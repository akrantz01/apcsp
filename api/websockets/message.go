package websockets

const (
	MessageAuthentication = iota
	MessageNewMessage
)

type BaseMessage struct {
	Type int `json:"type"`
}

type AuthenticationMessage struct {
	Type  int    `json:"type"`
	Token string `json:"token"`
}

type NewMessage struct {
	Type        int    `json:"type"`
	Message     string `json:"message"`
	Chat        string `json:"chat"`
	Sender      string `json:"sender"`
	ContentType int    `json:"content-type"`
}
