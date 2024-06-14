package usecases

import (
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UserCreator interface {
	CreateUser(user *entities.User) (*entities.User, error)
}

type CreateUserResponseBody struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func NewCreateUser(userCreator UserCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		newUser := &entities.User{
			Email:       gofakeit.Email(),
			Password:    gofakeit.Password(true, true, true, true, true, 15),
			Name:        gofakeit.Name(),
			Gender:      gofakeit.Gender(),
			DateOfBirth: gofakeit.Date(),
		}

		user, err := userCreator.CreateUser(newUser)
		if err != nil {
			slog.Error("creating new user", "err", err)
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, CreateUserResponseBody{
			ID:       user.ID.String(),
			Email:    user.Email,
			Password: user.Password,
			Name:     user.Name,
			Gender:   user.Gender,
			Age:      user.GetAge(),
		})
	}
}
