package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenType byte

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Type      TokenType `json:"token_type"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, role string, duration time.Duration, tokenType TokenType) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// Username and role are included in the payload for identification
	payload := &Payload{
		ID:        tokenID,
		Type:      tokenType,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid(tokenType TokenType) error {
	// Check if the token type matches the expected type
	if payload.Type != tokenType {
		return ErrInvalidToken
	}

	// Check if the token has expired
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// GetExpirationTime returns the expiration time of the token payload
func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.ExpiredAt,
	}, nil
}

// GetIssuedAt returns the issued at time of the token payload
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

// GetNotBefore returns the not before time of the token payload
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

// GetIssuer returns the issuer of the token payload
func (payload *Payload) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject returns the subject of the token payload
func (payload *Payload) GetSubject() (string, error) {
	return "", nil
}

// GetAudience returns the audience claim of the token payload
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
