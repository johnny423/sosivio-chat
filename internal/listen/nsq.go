package listen

import (
	"bytes"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"sosivio-chat/internal/chat"
)

//NSQListenerCreator is a creator for new instances of NSQListener
type NSQListenerCreator struct {
	Addr string
}

//Create creates a new instance of NSQListener
func (creator *NSQListenerCreator) Create(user string) Listener {
	return &NSQListener{
		topic: user,
		addr:  creator.Addr,
		q:     nil,
	}
}

//NSQListener is an implementation of the listener interface using an NSQ queue.
type NSQListener struct {
	topic string
	addr  string
	q     *nsq.Consumer
}

//Init sets up a new nsq consumer and inject the dispatcher into it.
func (listener *NSQListener) Init(dispatcher Dispatcher) error {
	q, err := listener.newConsumer()
	if err != nil {
		return err
	}
	q.AddHandler(listener.handlerCreator(dispatcher))
	listener.q = q
	return nil
}

func (listener *NSQListener) newConsumer() (*nsq.Consumer, error) {
	config := nsq.NewConfig()
	q, err := nsq.NewConsumer(listener.topic, "ch", config)
	if err != nil {
		return nil, err
	}
	q.SetLoggerLevel(nsq.LogLevelMax)
	return q, err

}

//Run connected the consumer to the NSQLookupd and by doing so starts consuming messages from the NSQ.
func (listener *NSQListener) Run() error {
	err := listener.q.ConnectToNSQLookupd(listener.addr)
	return err
}

//Stop ends the consuming from the NSQ.
func (listener *NSQListener) Stop() {
	listener.q.Stop()
	listener.q = nil
}

//handlerCreator create a new nsq.Handler from given dispatcher.
func (listener *NSQListener) handlerCreator(dispatcher Dispatcher) nsq.Handler {
	handler := func(message *nsq.Message) error {
		var msg chat.Msg
		err := json.NewDecoder(bytes.NewReader(message.Body)).Decode(&msg)
		if err != nil {
			return err
		}
		err = dispatcher.Handle(&msg)
		return err
	}
	return nsq.HandlerFunc(handler)
}
