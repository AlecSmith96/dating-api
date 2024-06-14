package entities

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID       string
	UserID   uuid.UUID
	Value    string
	IssuedAt time.Time
}
