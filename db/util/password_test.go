package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	pass := RandomString(6)

	hp, err := HashPassword(pass)
	require.NoError(t, err)
	require.NotEmpty(t, hp)

	whp, err := HashPassword(RandomString(73))
	require.Empty(t, whp)

	require.Equal(t, err, bcrypt.ErrPasswordTooLong)

	err = CompareHashPassword(pass, hp)
	require.NoError(t, err)

	wrongPass := RandomString(7)
	err = CompareHashPassword(wrongPass, hp)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}
