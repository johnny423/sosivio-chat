package send

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sosivio-chat/internal/auth"
	"sosivio-chat/internal/chat"
	"sosivio-chat/internal/config"
)

//NewSendHandler creates a new http handler func which handles messages sending requests.
//The send handler validate the sender user and creates a sender from the sender user to the receiver user.
func NewSendHandler(senderCreator Creator, validator auth.Validator, logger log.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received request")

		var msg chat.Msg
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Warn("error decoding request's body ", err)
			return
		}

		token, err := r.Cookie(config.TokenName)
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				logger.Warn("received request with no token")
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username, err := validator.Validate(token.Value)
		if err != nil {
			logger.Warn("could not validate token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if username != msg.Sender {
			logger.Warn(fmt.Sprintf("username - %s, is unauthorized"+
				" to send message as - %s ", username, msg.Sender))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sender := senderCreator.Create(msg.Receiver)
		err = SendMsg(msg, sender)
		if err != nil {
			logger.Error("error while sending message", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}
