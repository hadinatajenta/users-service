package users

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts user-related endpoints.
func RegisterRoutes(r *gin.Engine, handler *Handler) {
	r.GET("/health", handler.Health)

	api := r.Group("/api/v1")
	{
		api.POST("/users", handler.CreateUser)
		api.GET("/users", handler.ListUsers)
		api.GET("/users/:id", handler.GetUserByID)
	}
}
