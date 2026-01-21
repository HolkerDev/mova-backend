package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	svix "github.com/svix/svix-webhooks/go"

	"mova-backend/internal/middleware"
	"mova-backend/internal/service"
)

type ClerkEmailAddress struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
}

type ClerkUserData struct {
	ID                    string              `json:"id"`
	PrimaryEmailAddressID string              `json:"primary_email_address_id"`
	EmailAddresses        []ClerkEmailAddress `json:"email_addresses"`
}

type ClerkWebhookPayload struct {
	Type string        `json:"type"`
	Data ClerkUserData `json:"data"`
}

type ClerkWebhookHandler struct {
	userService *service.UserService
	wh          *svix.Webhook
}

func NewClerkWebhookHandler(userService *service.UserService, webhookSecret string) (*ClerkWebhookHandler, error) {
	wh, err := svix.NewWebhook(webhookSecret)
	if err != nil {
		return nil, err
	}
	return &ClerkWebhookHandler{userService: userService, wh: wh}, nil
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

	err = h.wh.Verify(body, headers)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	var payload ClerkWebhookPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
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

	user, err := h.userService.CreateUser(c.Request.Context(), payload.Data.ID, email)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			middleware.Logger.Info("user already exists", "clerk_id", payload.Data.ID)
			c.JSON(http.StatusOK, gin.H{"message": "user already exists"})
			return
		}
		middleware.Logger.Error("failed to create user", "error", err)
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
