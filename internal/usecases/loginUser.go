package usecases

import (
	"errors"
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type UserAuthenticator interface {
	LoginUser(email string, password string) (*entities.User, error)
	IssueJWT(userID uuid.UUID) (*entities.Token, error)
}

type LoginUserRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponseBody struct {
	Token string `json:"token"`
}

func NewLoginUser(userAuthenticator UserAuthenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request LoginUserRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Error("binding request body", "err", err)
			c.JSON(http.StatusBadRequest, entities.ErrorMessage{Message: err.Error()})
			return
		}

		user, err := userAuthenticator.LoginUser(request.Email, request.Password)
		if err != nil {
			if errors.Is(err, entities.ErrUserNotFound) {
				slog.Error("user not found for parsed details")
				c.JSON(http.StatusUnauthorized, entities.ErrorMessage{Message: "incorrect email or password"})
				return
			}
			slog.Error("authenticating user login", "err", err)
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "unable to login user"})
			return
		}

		token, err := userAuthenticator.IssueJWT(user.ID)
		if err != nil {
			slog.Error("issuing user JWT", "err", err)
			c.JSON(http.StatusInternalServerError, entities.ErrorMessage{Message: "unable to login user"})
			return
		}

		c.JSON(http.StatusOK, LoginUserResponseBody{
			Token: token.Value,
		})
	}
}
