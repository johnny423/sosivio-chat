package chat

type Msg struct {
	Body     string `json:"body"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Time     string `json:"time"`
}
