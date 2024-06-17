package entities

import (
	"github.com/google/uuid"
	"time"
)

// UserDiscovery is a struct representing a users as it appears in the discover endpoint
type UserDiscovery struct {
	ID          uuid.UUID
	Email       string
	Password    string
	Name        string
	Gender      string
	DateOfBirth time.Time
	Location    Location
	Age         int
}
