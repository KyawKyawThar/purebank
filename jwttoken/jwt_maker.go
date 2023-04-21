package jwttoken

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"

	//"github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretKeyLength = 32

// JWTMaker is a JSON Web Token Maker
type JWTMaker struct {
	secretkey string
}

// NewJWTMaker is a creation new JWT Token
func NewJWTMaker(secretkey string) (Maker, error) {

	if len(secretkey) < minSecretKeyLength {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeyLength)
	}

	return &JWTMaker{secretkey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (j JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", nil, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(j.secretkey))
	if err != nil {
		return "", nil, err
	}
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (j JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {

			return nil, ErrorInvalidToken
		}

		return []byte(j.secretkey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc, jwt.WithLeeway(5*time.Second))

	if err != nil {

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrorExpiredToken
		}

		return nil, ErrorInvalidToken
	}

	payload := jwtToken.Claims.(*Payload)

	if jwtToken.Valid {
		return payload, nil
	} else {
		return nil, err
	}
}
