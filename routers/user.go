package routers

import (
	. "gin_projects/handlers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {

	// User routes
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", GetUsers)
		userGroup.POST("/", CreateUser)
		userGroup.GET("/:id", GetUser)
		userGroup.PUT("/:id", UpdateUser)
		userGroup.DELETE("/:id", DeleteUser)
	}
}
