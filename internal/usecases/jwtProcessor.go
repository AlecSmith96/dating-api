package usecases

import "github.com/google/uuid"

//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/jwtProcessor.go  . "JwtProcessor"
type JwtProcessor interface {
	ValidateJwtForUser(tokenValue string) (uuid.UUID, error)
}
