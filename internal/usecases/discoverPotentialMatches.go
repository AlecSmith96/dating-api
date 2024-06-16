package usecases

import (
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type UserDiscoverer interface {
	DiscoverNewUsers(ownerUserID uuid.UUID, pageInfo entities.PageInfo) ([]entities.UserDiscovery, error)
}

type DiscoverPotentialMatchesRequestBody struct {
	PageInfo PageInfo `json:"pageInfo"`
}

type PageInfo struct {
	MinAge           int      `json:"minAge"`
	MaxAge           int      `json:"maxAge"`
	PreferredGenders []string `json:"preferredGenders"`
}

type DiscoverPotentialMatchesResponseBody struct {
	Users []UserResponseBody `json:"users"`
}

type UserResponseBody struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

func NewDiscoverPotentialMatches(discoverer UserDiscoverer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			slog.Error("unable to get userID from context")
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "unable to get users"})
			return
		}

		var request DiscoverPotentialMatchesRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Error("validating request body", "err", err)
			c.JSON(http.StatusBadRequest, entities.ErrorMessage{Message: "unable to validate request body"})
			return
		}

		userIDString := userID.(uuid.UUID)
		users, err := discoverer.DiscoverNewUsers(userIDString, entities.PageInfo{
			MinAge:           request.PageInfo.MinAge,
			MaxAge:           request.PageInfo.MaxAge,
			PreferredGenders: request.PageInfo.PreferredGenders,
		})
		if err != nil {
			slog.Error("getting users", "err", err)
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "unable to get users"})
			return
		}

		var returnedUsers []UserResponseBody
		for _, user := range users {
			returnedUsers = append(returnedUsers, UserResponseBody{
				ID:     user.ID.String(),
				Name:   user.Name,
				Gender: user.Gender,
				Age:    user.Age,
			})
		}

		c.JSON(http.StatusOK, DiscoverPotentialMatchesResponseBody{Users: returnedUsers})
	}
}
