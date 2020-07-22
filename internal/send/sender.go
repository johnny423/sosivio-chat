package send

import (
	"encoding/json"
	"sosivio-chat/internal/chat"
)

//Creator is an interface for creating a sender to a specific user
type Creator interface {
	Create(user string) Sender
}

//Sender is an interface for sending data
type Sender interface {
	Send(msg []byte) error
}

//SendMsg handle the sending of msg using the sender interface.
//Currently uses only a json as a serializer.
func SendMsg(msg chat.Msg, sender Sender) error {
	var b []byte
	b, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	err = sender.Send(b)
	if err != nil {
		return err
	}
	return nil
}
