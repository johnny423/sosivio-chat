package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"sosivio-chat/internal/auth"
	"sosivio-chat/internal/config"
	"sosivio-chat/internal/listen"
	"time"
)

const componentKey = "component"

func main() {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	validator := auth.NewRESTClient("http://localhost:"+config.AuthPort, "client",
		logger.WithField(componentKey, "validator_client"))

	retriever := listen.NewMessagesRetriever(
		&listen.NSQListenerCreator{Addr: config.NSQLookupdAddr},
		&listen.SimpleClientCreator{Timeout: 10 * time.Millisecond})

	ChatFunc := listen.NewPoolingHandler(validator, retriever,
		logger.WithField(componentKey, "pool_handler"))

	http.HandleFunc("/chat", ChatFunc)
	err := http.ListenAndServe(":"+config.ChatPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
