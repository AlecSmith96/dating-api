package usecases

import (
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type UserDiscoverer interface {
	DiscoverNewUsers(ownerUserID uuid.UUID) ([]entities.User, error)
}

type DiscoverPotentialMatchesResponseBody struct {
	Users []UserResponseBody `json:"users"`
}

type UserResponseBody struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func NewDiscoverPotentialMatches(discoverer UserDiscoverer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			slog.Error("unable to get userID from context")
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "unable to get users"})
			return
		}

		userIDString := userID.(uuid.UUID)
		users, err := discoverer.DiscoverNewUsers(userIDString)
		if err != nil {
			slog.Error("getting users", "err", err)
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "unable to get users"})
			return
		}

		var returnedUsers []UserResponseBody
		for _, user := range users {
			returnedUsers = append(returnedUsers, UserResponseBody{
				ID:       user.ID.String(),
				Email:    user.Email,
				Password: user.Password,
				Name:     user.Name,
				Gender:   user.Gender,
				Age:      user.GetAge(),
			})
		}

		c.JSON(http.StatusOK, DiscoverPotentialMatchesResponseBody{Users: returnedUsers})
	}
}
