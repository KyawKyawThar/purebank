package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"purebank/db/util"
	"testing"
	"time"
)

var ctx = context.Background()

func CreateTestUser(t *testing.T) Users {
	randomPass := util.RandomString(6)
	hp, err := util.HashPassword(randomPass)
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:  util.RandomOwner(8),
		Password:  hp,
		Email:     util.RandomEmail(5),
		FirstName: util.RandomOwner(6),
	}

	user, err := testQueries.CreateUser(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FirstName, user.FirstName)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestUser(t *testing.T) {
	CreateTestUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateTestUser(t)

	user2, err := testQueries.GetUser(ctx, user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.FirstName, user2.FirstName)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
