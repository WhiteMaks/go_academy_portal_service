package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenExpired = NewTokenError("token has expired")
	ErrTokenInvalid = NewTokenError("token is invalid")
)

type TokenError struct {
	Message string
}

func NewTokenError(message string) *TokenError {
	return &TokenError{
		Message: message,
	}
}

func (e *TokenError) Error() string {
	return e.Message
}

type Payload struct {
	Username  string
	Role      string
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewTokenPayload(username string, role string, duration time.Duration) *Payload {
	return &Payload{
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrTokenExpired
	}

	return nil
}

type TokenMaker interface {
	CreateToken(username string, role string) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type JWTTokenMaker struct {
	secretKey string
	duration  time.Duration
}

func NewJWTTokenMaker(secretKey string, duration time.Duration) TokenMaker {
	return &JWTTokenMaker{
		secretKey: secretKey,
		duration:  duration,
	}
}

func (maker *JWTTokenMaker) CreateToken(username string, role string) (string, error) {
	payload := NewTokenPayload(username, role, maker.duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTTokenMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrTokenInvalid
		}

		return []byte(maker.secretKey), nil
	}

	claims, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if errors.Is(ve.Inner, ErrTokenExpired) {
				return nil, ErrTokenExpired
			}
		}

		return nil, ErrTokenInvalid
	}

	payload, ok := claims.Claims.(*Payload)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return payload, nil
}
