package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mova-backend/internal/database"
	"mova-backend/internal/middleware"
)

type UserHandler struct {
	queries *database.Queries
}

func NewUserHandler(queries *database.Queries) *UserHandler {
	return &UserHandler{queries: queries}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.queries.GetUserByClerkID(c.Request.Context(), clerkUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"clerk_id": user.ClerkID,
		"email":    user.Email,
	})
}
