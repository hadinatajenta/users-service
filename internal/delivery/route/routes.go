package route

import (
	"github.com/gin-gonic/gin"

	handler "users-service/internal/delivery/http"
)

func Register(r *gin.Engine, userHandler *handler.UserHandler) {
	r.GET("/health", userHandler.Health)

	api := r.Group("/api/v1")
	{
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users", userHandler.ListUsers)
		api.GET("/users/:id", userHandler.GetUserByID)
	}
}
