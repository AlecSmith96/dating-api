package usecases

import (
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/swipeRegister.go  . "SwipeRegister"
type SwipeRegister interface {
	RegisterSwipe(ownerUserID, swipedUserID uuid.UUID, isPositivePreference bool) error
	IsMatch(ownerUserID, swipedUserID uuid.UUID) (*entities.Match, error)
}

// SwipeUserRequestBody represents the swipe result on a user
// @Description the swipe result on a user
type SwipeUserRequestBody struct {
	// UserID the id of the user that is swiped on
	UserID uuid.UUID `json:"userId" binding:"required"`
	// Preference the preference for the user
	Preference string `json:"preference" binding:"required,oneof=YES yes NO no"`
}

// SwipeUserResponseBody represents the result of the swipe
// @Description the result of the swipe, if theres a match it returns the matchID
type SwipeUserResponseBody struct {
	// Results the result of the swipe
	Results Result `json:"results"`
}

// Result represents the information of the swipe result
// @Description the information of the swipe result
type Result struct {
	// Matched whether the swipe resulted in a match
	Matched bool `json:"matched"`
	// MatchID the id of the match if the swipe resulted in one
	MatchID uuid.UUID `json:"matchId"`
}

// NewSwipeUser swipe on a user
// @Summary Swipe on a user
// @Description Provides a swipe result on a user
// @Security BearerAuth
// @Tags users
// @Accept json
// @Produce json
// @Param user body SwipeUserRequestBody true "Swipe User Request Body"
// @Success 200 {object} SwipeUserResponseBody
// @Failure 400
// @Failure 500
// @Router /user/swipe [post]
func NewSwipeUser(swipeRegister SwipeRegister) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			slog.Error("unable to get userID from context")
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "unable to get users"})
			return
		}
		requestingUserID := userID.(uuid.UUID)

		var request SwipeUserRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Error("validating request", "err", err)
			c.JSON(http.StatusBadRequest, entities.ErrorMessage{Message: "invalid request body"})
			return
		}

		var isPositivePreference bool
		switch request.Preference {
		case "YES":
			isPositivePreference = true
		case "yes":
			isPositivePreference = true
		case "NO":
			isPositivePreference = false
		case "no":
			isPositivePreference = false
		default:
			c.JSON(http.StatusBadRequest, entities.ErrorMessage{
				Message: "invalid preference",
			})
			return
		}

		err = swipeRegister.RegisterSwipe(requestingUserID, request.UserID, isPositivePreference)
		if err != nil {
			slog.Error("registering swipe", "err", err)
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "an internal server error occurred"})
			return
		}

		match, err := swipeRegister.IsMatch(requestingUserID, request.UserID)
		if err != nil {
			slog.Error("checking for match", "err", err)
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "an internal server error occurred"})
			return
		}

		isMatch := match != nil
		results := map[string]interface{}{
			"matched": isMatch,
		}

		if isMatch {
			results = map[string]interface{}{
				"matched": isMatch,
				"matchId": request.UserID,
			}
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"results": results,
		})
	}
}
