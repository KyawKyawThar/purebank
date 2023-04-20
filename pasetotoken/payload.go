package pasetotoken

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("this token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	IssuedAt   time.Time `json:"issued_at"`
	Expiration time.Time `json:"expiration"`
}

// NewPayload creates a new jwttoken payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:         tokenID,
		Username:   username,
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(duration),
	}

	return payload, err
}

//// Valid checks if the token payload is valid or not
//func (p *Payload) Valid() error {
//
//	if time.Now().After(p.Expiration) {
//		return ErrExpiredToken
//	}
//	return nil
//}
