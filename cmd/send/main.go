package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"sosivio-chat/internal/auth"
	"sosivio-chat/internal/config"
	"sosivio-chat/internal/send"
)

const componentKey = "component"

func main() {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	validator := auth.NewRESTClient("http://localhost:"+config.AuthPort, "client",
		logger.WithField(componentKey, "validator_client"))

	f := &send.NsqSenderCreator{Addr: config.NSQAddr}

	http.HandleFunc("/send", send.NewSendHandler(f, validator,
		logger.WithField(componentKey, "send_handler")))

	err := http.ListenAndServe(":"+config.SendPort, nil)
	if err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}
