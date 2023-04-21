package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	db "purebank/db/sqlc"
	"purebank/db/util"
	"purebank/worker"
	"time"
)

type createUserRequest struct {
	UserName  string `json:"username" binding:"required,alphanum"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
}

type logInUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type logInUserResponse struct {
	SessionID           uuid.UUID `json:"session_id"`
	AccessToken         string    `json:"access_token"`
	AccessTokenExpired  time.Time `json:"access_token_expired"`
	RefreshToken        string    `json:"refresh_token"`
	RefreshTokenExpired time.Time `json:"refresh_token_expired"`
	User                userResponse
}

type userResponse struct {
	UserName          string    `json:"username"`
	Email             string    `json:"email" `
	FirstName         string    `json:"first_name"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.Users) userResponse {

	return userResponse{
		UserName:          user.Username,
		Email:             user.Email,
		FirstName:         user.FirstName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) createUser(c *gin.Context) {

	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPassword, err := util.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:  req.UserName,
			Password:  hashPassword,
			Email:     req.Email,
			FirstName: req.FirstName,
		},
		AfterCreate: func(users db.Users) error {
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second), //delay
				asynq.Queue(worker.CriticalQueue),
			}

			taskPayloadEmail := &worker.PayloadSendVerifyEmail{Username: users.Username}
			return s.taskdistributor.DistributorSendVerifyEmail(c, taskPayloadEmail, opts...)

		},
	}

	//user, err := s.store.CreateUserTx(c, arg)

	txResult, err := s.store.CreateUserTx(c, arg)

	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok {

			switch pgErr.Routine {
			case "_bt_check_unique":

				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Send verify email to user

	res := newUserResponse(txResult.Users)

	c.JSON(http.StatusOK, res)

}

func (s *Server) loginUser(c *gin.Context) {

	var req logInUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(c, req.UserName)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CompareHashPassword(req.Password, user.Password)

	if err != nil {

		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.maker.CreateToken(user.Username, s.config.TokenDuration)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.maker.CreateToken(user.Username, s.config.RefreshTokenDuration)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlock:      false,
		ExpiredAt:    refreshPayload.Expiration,
	}

	session, err := s.store.CreateSession(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := logInUserResponse{
		SessionID:           session.ID,
		AccessToken:         accessToken,
		AccessTokenExpired:  accessPayload.Expiration,
		RefreshToken:        refreshToken,
		RefreshTokenExpired: refreshPayload.Expiration,
		User:                newUserResponse(user),
	}

	c.JSON(http.StatusOK, res)
}
