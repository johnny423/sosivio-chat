package auth

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sosivio-chat/internal/config"
)

//NewSignInHandler creates SignIn http handler which receives credentials and validates
//them and create a new token
func NewSignInHandler(service Signer, logger log.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received requests")

		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		var credentials Credentials
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			logger.Warn("bad request with error ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Info("requests from user - ", credentials.Username)

		token, expiration, err := service.SignIn(credentials)
		if err != nil {
			if err.Error() == AuthenticateError {
				logger.Warn("authentication failed for user - ", credentials.Username)
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				logger.Error("error while trying to authenticate ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		logger.Info(fmt.Sprintf("setting up token for "+
			"user %s for %v time", credentials.Username, expiration))
		http.SetCookie(w, &http.Cookie{
			Name:    config.TokenName,
			Value:   token,
			Expires: expiration,
		})
	}

}

//NewValidateHandler creates a Validate handler which receive a token and checks if
//it is a valid and if so it returns the username
func NewValidateHandler(Validator Validator, logger log.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("receive request")

		var tokenMsg TokenMsg
		err := json.NewDecoder(r.Body).Decode(&tokenMsg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Error("parse error")
			return
		}

		username, err := Validator.Validate(tokenMsg.Token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Warn("validate error", err)
			return
		}

		logger.Debug("got from token the username - ", username)
		result, err := json.Marshal(&ValidateResult{UserName: username})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("marshaling request error", err)
			return
		}

		_, err = w.Write(result)
		if err != nil {
			logger.Error("result writing error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

}
