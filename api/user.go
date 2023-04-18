package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	db "purebank/db/sqlc"
	"purebank/db/util"
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

type userResponse struct {
	UserName          string    `json:"username"`
	Email             string    `json:"email" `
	FirstName         string    `json:"first_name"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createUser(c *gin.Context) {

	fmt.Println("createUser")
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

	arg := db.CreateUserParams{
		Username:  req.UserName,
		Password:  hashPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
	}

	user, err := s.store.CreateUser(c, arg)

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

	res := userResponse{
		UserName:          user.Username,
		Email:             user.Email,
		FirstName:         user.FirstName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

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

	res := userResponse{
		UserName:          user.Username,
		Email:             user.Email,
		FirstName:         user.FirstName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}
