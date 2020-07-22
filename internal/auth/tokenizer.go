package auth

import (
	"time"
)

//Tokenizer is an interface for handling token related procedures
type Tokenizer interface {
	CreateToken(credentials Credentials, expiration time.Time) (string, error)
	ValidateToken(token string) (string, error)
}
