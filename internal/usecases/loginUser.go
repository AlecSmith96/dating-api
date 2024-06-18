package usecases

import (
	"errors"
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/userAuthenticator.go  . "UserAuthenticator"
type UserAuthenticator interface {
	LoginUser(email string, password string) (*entities.User, error)
	IssueJWT(userID uuid.UUID) (*entities.Token, error)
}

// LoginUserRequestBody represents the login credentials for the user
// @Description the login information for the user
type LoginUserRequestBody struct {
	// Email represents the email of the user to log in as
	Email string `json:"email" binding:"required"`
	// Password represents the password of the user to log in as
	Password string `json:"password" binding:"required"`
}

// LoginUserResponseBody represents the Bearer token to use in authenticated requests
// @Description the newly issued JWT for the logged in user
type LoginUserResponseBody struct {
	// Token represents the JWT issued for the logged in user
	Token string `json:"token"`
}

// NewLoginUser logs in a user
// @Summary Login a user
// @Description Logs in a user with the provided credentials
// @Tags users
// @Accept json
// @Produce json
// @Param user body LoginUserRequestBody true "Login User Request Body"
// @Success 200 {object} LoginUserResponseBody
// @Failure 400
// @Failure 500
// @Router /login [post]
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
