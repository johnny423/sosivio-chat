package auth

import "time"

//Credentials contains the information needed to authenticate a user
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

//Authenticator returns true if the credentials are valid, false otherwise
type Authenticator interface {
	Authenticate(credentials Credentials) bool
}

//Signer creates a token from the credentials
type Signer interface {
	SignIn(credentials Credentials) (string, time.Time, error)
}

//Validator checks if a token is valid and if so return the username of the token
type Validator interface {
	Validate(token string) (string, error)
}
