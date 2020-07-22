package listen

import "sosivio-chat/internal/chat"

//NewMessagesRetriever creates a new instance of MessagesRetriever.
func NewMessagesRetriever(listenerCreator ListenerCreator, clientCreator ClientCreator) MessagesRetriever {
	return MessagesRetriever{
		listenerCreator: listenerCreator,
		clientCreator:   clientCreator,
	}
}

//MessagesRetriever mange the consuming messages from an async listener into a client.
//The MessagesRetriever using creators interface to create a new instance of a listener
//and a client for each request.
type MessagesRetriever struct {
	listenerCreator ListenerCreator
	clientCreator   ClientCreator
}

//GetUserMessages returns all the new messages that are waiting for the user.
func (mr *MessagesRetriever) GetUserMessages(user string) []*chat.Msg {

	//creates a listener to observe the messages for the user
	listener := mr.listenerCreator.Create(user)
	//create a client to extract the messages from the listener
	client := mr.clientCreator.Create()

	var msgs []*chat.Msg

	//connected the client to the listener by injecting the client's dispatcher into the listener
	err := listener.Init(client.NewDispatcher())
	if err != nil {
		return msgs
	}

	//starts the listener streaming into the dispatcher
	err = listener.Run()
	if err != nil {
		return msgs
	}

	//consuming the messages using the client, and ordering the client to stop the listener when it is done.
	msgs = client.ExtractMessages(func() { listener.Stop() })
	return msgs
}
