package listen

import "sosivio-chat/internal/chat"

//Dispatcher handles the receiving of a new message.
type Dispatcher interface {
	Handle(*chat.Msg) error
}

//Listener is an interface for async messages observer.
type Listener interface {
	//Init sets up the listener with a dispatcher.
	Init(dispatcher Dispatcher) error
	//Run starts the asynchronously waiting for new messages and using the dispatcher for handling them.
	Run() error
	//Stop ends the waiting for new messages.
	Stop()
}

//ListenerCreator is an interface for creating new listener with will observe any messages sent to a user.
type ListenerCreator interface {
	Create(user string) Listener
}

//Client is an interface for collecting messages from async listener.
type Client interface {
	NewDispatcher() Dispatcher
	ExtractMessages(closeHandler func()) []*chat.Msg
}

//ClientCreator in an interface from creating a new Client.
type ClientCreator interface {
	Create() Client
}
