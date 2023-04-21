package jwttoken

import "time"

// Maker is an interface for managing jwttoken
type Maker interface {
	CreateToken(username string, duration time.Duration) (string,*Payload, error)
	VerifyToken(token string) (*Payload, error)
}
