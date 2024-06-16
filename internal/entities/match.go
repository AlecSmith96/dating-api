package entities

import "github.com/google/uuid"

type Match struct {
	ID            uuid.UUID `json:"id"`
	OwnerUserID   uuid.UUID `json:"ownerUserId"`
	MatchedUserID uuid.UUID `json:"matchedUserId"`
}
