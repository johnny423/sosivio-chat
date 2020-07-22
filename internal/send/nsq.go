package send

import (
	"github.com/nsqio/go-nsq"
)

//NsqSenderCreator is a creator for a sender that uses a NSQ technology for sending data.
type NsqSenderCreator struct {
	Addr string
}

//Create creates a new NsqSender to send data to a specific user.
func (creator *NsqSenderCreator) Create(user string) Sender {
	return &NsqSender{
		addr:  creator.Addr,
		topic: user,
	}
}

//NsqSender is a Sender which sends data using NSQ connection.
type NsqSender struct {
	addr  string
	topic string
}

//Send sends a msg to an NSQ producer using a the user as a topic.
func (sender *NsqSender) Send(msg []byte) error {
	producer, err := sender.newProducer()
	if err != nil {
		return err
	}

	err = producer.Publish(sender.topic, msg)
	if err != nil {
		return err
	}

	producer.Stop()
	return nil
}

//newProducer initiates a new configured NSQ producer
func (sender *NsqSender) newProducer() (*nsq.Producer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(sender.addr, config)
	if err != nil {
		return nil, err
	}
	producer.SetLoggerLevel(nsq.LogLevelMax)
	return producer, err
}
