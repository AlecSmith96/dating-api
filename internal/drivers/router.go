package drivers

import (
	"github.com/AlecSmith96/dating-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userCreator usecases.UserCreator,
	userAuthenticator usecases.UserAuthenticator,
) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/dating-api/v1")
	{
		v1.POST("/user/create", usecases.NewCreateUser(userCreator))
		v1.POST("/login", usecases.NewLoginUser(userAuthenticator))
	}

	return r
}
