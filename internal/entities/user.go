package entities

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID
	Email       string
	Password    string
	Name        string
	Gender      string
	DateOfBirth time.Time
}

func (u *User) GetAge() int {
	now := time.Now()
	age := now.Year() - u.DateOfBirth.Year()

	if now.YearDay() < u.DateOfBirth.YearDay() {
		age--
	}

	return age
}
