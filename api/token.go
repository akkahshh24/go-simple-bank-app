package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/akkahshh24/go-simple-bank-app/db/sqlc"
	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	// Bind the request to the renewAccessTokenRequest struct
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If binding fails, return a bad request error
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify the refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken, token.TokenTypeRefreshToken)
	if err != nil {
		// If the token verification fails, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Get the session associated with the refresh token
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		// If the session is not found, return a not found error
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// If there's another error, return an internal server error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check if the session is blocked
	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		// If the session is blocked, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check if the session username matches the refresh token username
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		// If the session username does not match, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check if the session refresh token matches the request refresh token
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		// If the session refresh token does not match, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check if the session has expired
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		// If the session has expired, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Create a new access token for the user using the same username and role from the refresh payload
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		refreshPayload.Role,
		server.config.AccessTokenDuration,
		token.TokenTypeAccessToken,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
