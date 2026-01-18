package router

import (
	"mova-backend/internal/database"
	"mova-backend/internal/handler"
	"mova-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(queries *database.Queries) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.GET("/", handler.HelloWorld)

	// Webhooks (no auth)
	clerkWebhookHandler := handler.NewClerkWebhookHandler(queries)
	r.POST("/webhooks/clerk", clerkWebhookHandler.HandleUserCreated)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.ClerkAuth())
	{
		userHandler := handler.NewUserHandler(queries)
		protected.GET("/me", userHandler.GetMe)
	}

	return r
}
