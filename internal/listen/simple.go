package listen

import (
	"sosivio-chat/internal/chat"
	"time"
)

//SimpleClientCreator creates a SimpleClient instances
type SimpleClientCreator struct {
	Timeout time.Duration
}

//Create creates a new instance of SimpleClient
func (creator *SimpleClientCreator) Create() Client {
	return &SimpleClient{timeout: creator.Timeout}
}

//SimpleClient handles the extracting new messages from a listener for limited time.
type SimpleClient struct {
	receive chan *chat.Msg
	timeout time.Duration
}

//NewDispatcher creates a new dispatcher which will send any new message to the
//communication channel with the client
func (client *SimpleClient) NewDispatcher() Dispatcher {
	c := make(chan *chat.Msg)
	client.receive = c
	return &dispatcher{out: c}
}

//ExtractMessages extracting any messages arriving to the Client for limited amount of time
//and then gracefully wrapping up the connection.
func (client *SimpleClient) ExtractMessages(closeHandler func()) []*chat.Msg {
	var msgs []*chat.Msg

	defer func() {
		closeHandler()
		close(client.receive)
	}()

	for {
		select {
		case msg := <-client.receive:
			msgs = append(msgs, msg)
		case <-time.After(client.timeout):
			return msgs
		}
	}
}

//dispatcher is used as a simple directing new messages to the out channel
type dispatcher struct {
	out chan *chat.Msg
}

//Handle sends any new message to the out channel
func (d *dispatcher) Handle(msg *chat.Msg) error {
	d.out <- msg
	return nil
}
