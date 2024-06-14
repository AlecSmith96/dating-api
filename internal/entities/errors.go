package entities

import "errors"

var (
	ErrUserNotFound = errors.New("user not found for parsed details")
)

type ErrorMessage struct {
	Message string `json:"message"`
}
