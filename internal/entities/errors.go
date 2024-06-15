package entities

import "errors"

var (
	ErrUserNotFound = errors.New("user not found for parsed details")
	ErrJwtExpired   = errors.New("jwt is expired")
)

type ErrorMessage struct {
	Message string `json:"message"`
}
