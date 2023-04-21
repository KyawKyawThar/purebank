package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type sessionRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type sessionResponse struct {
	AccessToken         string    `json:"access_token" binding:"required"`
	AccessTokenDuration time.Time `json:"access_token_duration"`
}

func (s *Server) renewAccessToken(c *gin.Context) {

	var req sessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := s.maker.VerifiedToken(req.RefreshToken)
	//
	if err != nil {
		fmt.Println("code running here")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := s.store.GetSession(c, refreshPayload.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlock {
		err := fmt.Errorf("blocked session")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err = fmt.Errorf("incorrect session user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		fmt.Println("session expired")
		err = fmt.Errorf("expired session")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.maker.CreateToken(refreshPayload.Username, s.config.TokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := sessionResponse{
		AccessToken:         accessToken,
		AccessTokenDuration: accessPayload.Expiration,
	}

	c.JSON(http.StatusOK, res)
}
