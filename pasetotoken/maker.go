package pasetotoken

import "time"

// Maker is an interface for managing paseto token
type Maker interface {

	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	//VerifiedToken checks if the token is valid or not
	VerifiedToken(token string) (*Payload, error)
}
