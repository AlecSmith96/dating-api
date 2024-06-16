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

// DiscoverPotentialMatchesRequestBody represents the filters for the returned list of users
// @Description the request body for the discover endpoint
type DiscoverPotentialMatchesRequestBody struct {
	// PageInfo represents the filter information for the request
	PageInfo PageInfo `json:"pageInfo"`
}

// PageInfo represents the filters for the returned list of users
// @Description the filter information for the request
type PageInfo struct {
	// MinAge is the minimum age of any users returned in the list
	MinAge int `json:"minAge"`
	// MaxAge is the maximum age of any users returned in the list
	MaxAge int `json:"maxAge"`
	// PreferredGenders is an array of genders to include in the list
	PreferredGenders []string `json:"preferredGenders"`
}

// DiscoverPotentialMatchesResponseBody represents the response of the discover endpoint
// @Description the response body for the discover endpoint
type DiscoverPotentialMatchesResponseBody struct {
	// Users is the returned list of all users matching the filter criteria
	Users []UserResponseBody `json:"users"`
}

// UserResponseBody represents a user that is returned by the discover endpoint
// @Description a user matching the filter criteria
type UserResponseBody struct {
	// ID is the id of the user
	ID string `json:"id"`
	// Name is the name of the user
	Name string `json:"name"`
	// Gender is the gender of the user
	Gender string `json:"gender"`
	// Age is the age of the user
	Age int `json:"age"`
}

// NewDiscoverPotentialMatches get a filterable list of users
// @Summary Discover new users
// @Description Gets a filterable list of new users
// @Security BearerAuth
// @Tags users
// @Accept json
// @Produce json
// @Param user body DiscoverPotentialMatchesRequestBody true "Discover Potential Matches Request Body"
// @Success 200 {object} DiscoverPotentialMatchesResponseBody
// @Failure 400
// @Failure 500
// @Router /user/discover [get]
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
