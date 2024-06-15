package drivers

import (
	"errors"
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/AlecSmith96/dating-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// TokenAuthMiddleware is a custom middleware function that processes the provided JWT in the Authorization header of
// the request. If the JWT is valid, it sets the userID value in the requests context and parses it to the usecase.
func TokenAuthMiddleware(jwtProcessor usecases.JwtProcessor) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeaderValue := c.GetHeader("Authorization")
		jwt := strings.Split(authHeaderValue, " ")

		if len(jwt) != 2 || (jwt[0] != "Bearer" && jwt[0] != "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entities.ErrorMessage{Message: "invalid jwt"})
			return
		}

		if jwt[1] == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entities.ErrorMessage{Message: "jwt is missing"})
			return
		}

		userID, err := jwtProcessor.ValidateJwtForUser(jwt[1])
		if err != nil {
			if errors.Is(err, entities.ErrJwtExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, entities.ErrorMessage{Message: "jwt is expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, entities.ErrorMessage{Message: "invalid jwt"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func NewRouter(
	userCreator usecases.UserCreator,
	userAuthenticator usecases.UserAuthenticator,
	jwtProcessor usecases.JwtProcessor,
	userDiscoverer usecases.UserDiscoverer,
) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/dating-api/v1")
	{
		v1.POST("/login", usecases.NewLoginUser(userAuthenticator))

		protected := v1.Group("/user", TokenAuthMiddleware(jwtProcessor))
		{
			protected.POST("/create", usecases.NewCreateUser(userCreator))
			protected.GET("/discover", usecases.NewDiscoverPotentialMatches(userDiscoverer))
		}
	}

	return r
}
