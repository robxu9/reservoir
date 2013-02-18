package reservoir

type Message struct {
	msgtype byte   `json:"type"`
	message string `json:"message"`
}
