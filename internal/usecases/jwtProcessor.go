package usecases

import "github.com/google/uuid"

type JwtProcessor interface {
	ValidateJwtForUser(tokenValue string) (uuid.UUID, error)
}
