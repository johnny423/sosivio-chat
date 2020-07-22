package auth

import (
	"errors"
	"time"
)

const AuthenticateError = "authentication error"

//Service handles authentication calls
type Service struct {
	authenticator Authenticator
	tokenizer     Tokenizer
	expiration    time.Duration
}

//NewAuthService creates a new instance of auth.Service
func NewAuthService(authenticator Authenticator, tokenizer Tokenizer, expiration time.Duration) *Service {
	return &Service{
		authenticator: authenticator,
		tokenizer:     tokenizer,
		expiration:    expiration,
	}
}

//SignIn authenticate the credentials and creates a token from them
func (s *Service) SignIn(credentials Credentials) (string, time.Time, error) {
	if !s.authenticator.Authenticate(credentials) {
		return "", time.Time{}, errors.New(AuthenticateError)
	}

	expirationTime := time.Now().Add(s.expiration)
	token, err := s.tokenizer.CreateToken(credentials, expirationTime)
	if err != nil {
		return "", time.Time{}, errors.New("token creation error")
	}
	return token, expirationTime, nil
}

//Validate checks if a given token is valid and if so returns the token's username, otherwise returns error
func (s *Service) Validate(token string) (string, error) {
	return s.tokenizer.ValidateToken(token)
}
