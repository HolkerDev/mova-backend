package router

import (
	"mova-backend/internal/handler"
	"mova-backend/internal/middleware"
	"mova-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func Setup(
	userService *service.UserService, deckService *service.DeckService, clerkWebhookSecret string,
) (*gin.Engine, error) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// Public routes
	r.GET("/", handler.HelloWorld)

	// Webhooks (no auth, but verified with svix)
	clerkWebhookHandler, err := handler.NewClerkWebhookHandler(userService, clerkWebhookSecret)
	if err != nil {
		return nil, err
	}
	r.POST("/webhooks/clerk", clerkWebhookHandler.HandleUserCreated)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.ClerkAuth(userService))
	{
		userHandler := handler.NewUserHandler()
		protected.GET("/me", userHandler.GetMe)

		deckHandler := handler.NewDeckHandler(deckService)
		protected.POST("/decks", deckHandler.CreateDeck)
	}

	return r, nil
}
