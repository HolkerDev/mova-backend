package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mova-backend/internal/database"
)

type ClerkEmailAddress struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
}

type ClerkUserData struct {
	ID                    string              `json:"id"`
	EmailAddresses        []ClerkEmailAddress `json:"email_addresses"`
	PrimaryEmailAddressID string              `json:"primary_email_address_id"`
}

type ClerkWebhookPayload struct {
	Data ClerkUserData `json:"data"`
	Type string        `json:"type"`
}

type ClerkWebhookHandler struct {
	queries *database.Queries
}

func NewClerkWebhookHandler(queries *database.Queries) *ClerkWebhookHandler {
	return &ClerkWebhookHandler{queries: queries}
}

func (h *ClerkWebhookHandler) HandleUserCreated(c *gin.Context) {
	var payload ClerkWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if payload.Type != "user.created" {
		c.JSON(http.StatusOK, gin.H{"message": "event ignored"})
		return
	}

	email := getPrimaryEmail(payload.Data)
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no primary email found"})
		return
	}

	user, err := h.queries.CreateUser(c.Request.Context(), database.CreateUserParams{
		ClerkID: payload.Data.ID,
		Email:   email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "clerk_id": user.ClerkID})
}

func getPrimaryEmail(data ClerkUserData) string {
	for _, email := range data.EmailAddresses {
		if email.ID == data.PrimaryEmailAddressID {
			return email.EmailAddress
		}
	}
	return ""
}
