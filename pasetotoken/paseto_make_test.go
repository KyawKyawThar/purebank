package pasetotoken

import (
	"github.com/stretchr/testify/require"
	"purebank/db/util"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {

	maker, err := NewPasetoMake(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner(6)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifiedToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.Expiration, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {

	maker, err := NewPasetoMake(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(6), -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifiedToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
