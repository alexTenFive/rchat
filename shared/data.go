package shared

import "time"

const (
	ServerPort         = 1212
	HANDSHAKING_REQREP = "handshaked"
)

type (
	Message struct {
		Time    time.Time `json:"time"`
		Name    string    `json:"name"`
		Message string    `json:"message"`
	}

	ClientRequest struct {
		Name       string `json:"name"`
		Data       string `json:"data"`
		ToUserName string `json:"to_user_name"`
	}

	TerminalData struct {
		Message string
	}
)
