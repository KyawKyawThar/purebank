package jwttoken

import "time"

// Maker is an interface for managing jwttoken
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
