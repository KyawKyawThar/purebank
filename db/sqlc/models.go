// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"time"
)

type Users struct {
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}