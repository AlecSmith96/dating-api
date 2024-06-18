package usecases

import (
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/userCreator.go  . "UserCreator"
type UserCreator interface {
	CreateUser(user *entities.User) (*entities.User, error)
}

// CreateUserResponseBody represents the newly generated user object
// @Description Response body for the newly created user
type CreateUserResponseBody struct {

	// ID the generated id for the user
	ID string `json:"id"`
	// Email the generated email for the user
	Email string `json:"email"`
	// Password the generated password for the user
	Password string `json:"password"`
	// Name the generated name for the user
	Name string `json:"name"`
	// Gender the generated gender for the user
	Gender string `json:"gender"`
	// Age the generated age for the user
	Age int `json:"age"`
	// Location the generated location for the user
	Location Location `json:"location"`
}

type Location struct {
	// Latitude the generated latitude for the users location
	Latitude float64 `json:"latitude"`
	// Longitude the generated longitude for the users location
	Longitude float64 `json:"longitude"`
}

// @BasePath /dating-api/v1

// NewCreateUser generates a new user record
// @Summary Create a new user
// @Description Generates a new user record based on fake data
// @Security BearerAuth
// @Tags users
// @Produce json
// @Success 200 {object} CreateUserResponseBody
// @Failure 400
// @Failure 500
// @Router /user/create [post]
func NewCreateUser(userCreator UserCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		newUser := &entities.User{
			Email:       gofakeit.Email(),
			Password:    gofakeit.Password(true, true, true, true, true, 15),
			Name:        gofakeit.Name(),
			Gender:      gofakeit.Gender(),
			DateOfBirth: gofakeit.Date(),
			Location: entities.Location{
				Latitude:  gofakeit.Address().Latitude,
				Longitude: gofakeit.Address().Longitude,
			},
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
			Location: Location{
				Latitude:  user.Location.Latitude,
				Longitude: user.Location.Longitude,
			},
		})
	}
}
