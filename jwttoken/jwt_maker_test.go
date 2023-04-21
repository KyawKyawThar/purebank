package jwttoken

import (
	"github.com/golang-jwt/jwt/v5"

	//"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"purebank/db/util"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner(6)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	jwtToken, payload, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)
	require.NotEmpty(t, payload)

	token, err := maker.VerifyToken(jwtToken)

	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)

	require.NotZero(t, token.ID)
	require.Equal(t, username, payload.Issuer)

	require.WithinDuration(t, issuedAt, token.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiredAt, token.ExpiresAt.Time, time.Second)

}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner(6)
	jwtToken, payload, err := maker.CreateToken(username, -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, jwtToken)

	token, err := maker.VerifyToken(jwtToken)

	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, token)
}

func TestInvalidToken(t *testing.T) {

	payload, err := NewPayload(util.RandomOwner(6), time.Minute)
	require.NoError(t, err)

	//
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	//
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	//
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorInvalidToken.Error())
	require.Nil(t, payload)

}
