package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"sosivio-chat/internal/auth"
	"sosivio-chat/internal/config"
	"time"
)

const componentKey = "component"

func main() {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	service := auth.NewAuthService(
		&auth.DummyAuthenticator{},
		&auth.JwtTokenizer{JwtKey: []byte(config.TokenKey)},
		time.Minute*5,
	)
	SignIn := auth.NewSignInHandler(service, logger.WithField(componentKey, "signin_handler"))
	Validate := auth.NewValidateHandler(service, logger.WithField(componentKey, "validate_handler"))
	http.HandleFunc("/signin", SignIn)
	http.HandleFunc("/validate", Validate)
	log.Fatal(http.ListenAndServe(":"+config.AuthPort, nil))
}
