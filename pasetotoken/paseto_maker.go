package pasetotoken

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20"
	"time"
)

// PasetoMaker is a JSON Paseto Token Maker
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

// NewPasetoMaker create a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {

	if len(symmetricKey) != chacha20.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20.KeySize)
	}

	return &PasetoMaker{symmetricKey: paseto.NewV4SymmetricKey()}, nil
}

// CreateToken is a paseto createToken
func (p PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", nil, err
	}

	token := paseto.NewToken()
	token.SetString("ID", payload.ID.String())
	token.SetIssuer(payload.Username)

	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.Expiration)

	v4token := token.V4Encrypt(p.symmetricKey, nil)

	return v4token, payload, nil

}

// VerifiedToken checks if the token is valid or not
func (p PasetoMaker) VerifiedToken(token string) (*Payload, error) {

	parser := paseto.NewParser()

	t, err := parser.ParseV4Local(p.symmetricKey, token, nil)

	if err != nil {

		return nil, ErrInvalidToken
	}

	id, err := t.GetString("ID")
	user, err := t.GetIssuer()
	issueAt, err := t.GetIssuedAt()
	expired, err := t.GetExpiration()

	uuid, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:         uuid,
		Username:   user,
		IssuedAt:   issueAt,
		Expiration: expired,
	}

	return payload, nil
}
