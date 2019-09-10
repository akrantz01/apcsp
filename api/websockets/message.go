package websockets

const (
	MessageAuthentication = iota
)

type BaseMessage struct {
	Type int `json:"type"`
}

type AuthenticationMessage struct {
	Type  int    `json:"type"`
	Token string `json:"token"`
}
