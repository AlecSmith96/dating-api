package usecases

import (
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type SwipeRegister interface {
	RegisterSwipe(ownerUserID, swipedUserID uuid.UUID, isPositivePreference bool) error
	IsMatch(ownerUserID, swipedUserID uuid.UUID) (*entities.Match, error)
}

type SwipeUserRequestBody struct {
	UserID     uuid.UUID `json:"userId" binding:"required"`
	Preference string    `json:"preference" binding:"required,oneof=YES yes NO no"`
}

type Result struct {
	Matched bool      `json:"matched"`
	MatchID uuid.UUID `json:"matchId"`
}

func NewSwipeUser(swipeRegister SwipeRegister) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			slog.Error("unable to get userID from context")
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "unable to get users"})
			return
		}
		userIDUUID := userID.(uuid.UUID)

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

		err = swipeRegister.RegisterSwipe(userIDUUID, request.UserID, isPositivePreference)
		if err != nil {
			slog.Error("registering swipe", "err", err)
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "an internal server error occurred"})
			return
		}

		match, err := swipeRegister.IsMatch(userIDUUID, request.UserID)
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
