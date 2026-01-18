package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mova-backend/internal/middleware"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	user, ok := middleware.GetAuthUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"clerk_id": user.ClerkID,
		"email":    user.Email,
	})
}
