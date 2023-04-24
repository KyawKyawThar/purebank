package db

import (
	"context"
	"database/sql"
	"fmt"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	Users        Users
	VerifyEmails VerifyEmails
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {

	var result VerifyEmailTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmails, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})

		if err != nil {
			fmt.Println("VerifyEmailTx-1", err)
			return err
		}

		result.Users, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmails.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		fmt.Println("VerifyEmailTx-2", err)

		return err

	})
	fmt.Println("VerifyEmailTx-3", err)
	return result, err
}
