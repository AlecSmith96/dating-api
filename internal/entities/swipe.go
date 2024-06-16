package entities

import "github.com/google/uuid"

type Swipe struct {
	ID                 uuid.UUID `json:"id"`
	OwnerUserID        uuid.UUID `json:"ownerUserID"`
	SwipedUserID       uuid.UUID `json:"swipedUserID"`
	PositivePreference bool      `json:"positivePreference"`
}
