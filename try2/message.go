package reservoir

/*
	Message exists as a way to integrate packets in socket connections.
	When a socket connection with the worker is established, PING/PONG will be
	sent often; and when work needs to be done, the following will happen:

	HANDSHAKE -> worker
	HANDSHAKE -> host
	MSG_TYPE -> worker
	MSG_OK or MSG_NO -> host

	if MSG_OK {
		MSG_DATA -> worker
		MSG_OK or MSG_NO -> host

		if MSG_NO {
			(something went terribly wrong)
		}
	}

	HANDSHAKE -> worker
	HANDSHAKE -> host

	The handshakes at the end are merely used to symbolise the end of a transaction.
	They're just there for the politeness.

	Also, we're sending the packets as JSON so that we can do some good with the info.
	See https://gist.github.com/faried/239744!

*/

const (
	MSG_PING = iota
	MSG_PONG
	MSG_HANDSHAKE
	MSG_OK
	MSG_NO
	MSG_TYPE
	MSG_DATA
)

type Message struct {
	msgtype byte   `json:"type"`
	msgdest	string `json:"dest"`
	message string `json:"message"`
}

type Dispatchable interface {
	SendMessage(msg *Message) bool
	Ping() bool
}

type Processable interface {
	ProcessMessage(msg *Message)
}
