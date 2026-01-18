package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	svix "github.com/svix/svix-webhooks/go"

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
	wh      *svix.Webhook
}

func NewClerkWebhookHandler(queries *database.Queries, webhookSecret string) (*ClerkWebhookHandler, error) {
	wh, err := svix.NewWebhook(webhookSecret)
	if err != nil {
		return nil, err
	}
	return &ClerkWebhookHandler{queries: queries, wh: wh}, nil
}

func (h *ClerkWebhookHandler) HandleUserCreated(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	headers := http.Header{}
	headers.Set("svix-id", c.GetHeader("svix-id"))
	headers.Set("svix-timestamp", c.GetHeader("svix-timestamp"))
	headers.Set("svix-signature", c.GetHeader("svix-signature"))

	if err := h.wh.Verify(body, headers); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	var payload ClerkWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
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
