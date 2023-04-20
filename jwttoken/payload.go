package jwttoken

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	ErrorExpiredToken = errors.New("jwttoken is expired")

	ErrorInvalidToken = errors.New("jwttoken is invalid")
)

// Payload contains the payload data of the token
type Payload struct {
	jwt.RegisteredClaims
}

// NewPayload creates a new jwttoken payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	token, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}
	//IssuedAt
	payload := &Payload{
		jwt.RegisteredClaims{
			ID:        token.String(),
			Issuer:    username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return payload, nil
}
