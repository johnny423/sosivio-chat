package listen

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sosivio-chat/internal/auth"
	"sosivio-chat/internal/config"
)

//NewPoolingHandler creates a new http handler that checks if the user is authenticated  if so it gets the user's waiting messages
func NewPoolingHandler(validator auth.Validator, handler MessagesRetriever, logger log.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received new pooling request")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		cookie, err := r.Cookie(config.TokenName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Warn("request with no token")
			return
		}

		username, err := validator.Validate(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Warn("validate error\t", err)
			return
		}
		logger.Info("validated user\t", username)

		msgs := handler.GetUserMessages(username)
		if len(msgs) == 0 {
			w.WriteHeader(http.StatusNoContent)
			logger.Warn("no messages found for user:\t", username)
			return
		}
		logger.Debug(fmt.Sprintf("found %v messages for user:\t%s", len(msgs), username))

		content, err := json.Marshal(msgs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Warn("err while marshaling the messages:\t", err)
			return
		}

		_, err = w.Write(content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Warn("err while writing the messages to the respond:\t", err)
			return
		}

	}
}
