package router

import (
	"mova-backend/internal/database"
	"mova-backend/internal/handler"
	"mova-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(queries *database.Queries, clerkWebhookSecret string) (*gin.Engine, error) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// Public routes
	r.GET("/", handler.HelloWorld)

	// Webhooks (no auth, but verified with svix)
	clerkWebhookHandler, err := handler.NewClerkWebhookHandler(queries, clerkWebhookSecret)
	if err != nil {
		return nil, err
	}
	r.POST("/webhooks/clerk", clerkWebhookHandler.HandleUserCreated)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.ClerkAuth(queries))
	{
		userHandler := handler.NewUserHandler()
		protected.GET("/me", userHandler.GetMe)
	}

	return r, nil
}
