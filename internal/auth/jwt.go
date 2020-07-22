package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//JwtTokenizer is a type for creating and validating tokens using the JSON Web Tokens protocol
type JwtTokenizer struct {
	JwtKey []byte
}

//CreateToken creates a new jwt token from given credentials
func (jt *JwtTokenizer) CreateToken(credentials Credentials, expiration time.Time) (string, error) {
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jt.JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//ValidateToken checks if a jwt token is valid and if so returns the username of the token's user
func (jt *JwtTokenizer) ValidateToken(token string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jt.JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("invalid")
		}
		return "", err
	}

	if tkn == nil || !tkn.Valid {
		return "", errors.New("invalid")

	}

	return claims.Username, nil
}

//Claims add the username for a standard jwt claim
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
