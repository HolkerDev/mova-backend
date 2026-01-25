package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mova-backend/internal/middleware"
	"mova-backend/internal/service"
)

type DeckHandler struct {
	deckService *service.DeckService
}

func NewDeckHandler(deckService *service.DeckService) *DeckHandler {
	return &DeckHandler{deckService: deckService}
}

type CreateDeckRequest struct {
	Name           string `json:"name" binding:"required"`
	SourceLanguage string `json:"source_language" binding:"required"`
	TargetLanguage string `json:"target_language" binding:"required"`
}

type CreateDeckResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	CreatedAt      string `json:"created_at"`
}

type Language string

const (
	LangEnglish   Language = "en"
	LangGerman    Language = "de"
	LangRussian   Language = "ru"
	LangUkrainian Language = "uk"
	LangGreek     Language = "el"
	LangPolish    Language = "pl"
)

func (l Language) IsValid() bool {
	switch l {
	case LangEnglish, LangGerman, LangRussian, LangUkrainian, LangGreek, LangPolish:
		return true
	}
	return false
}

func (h *DeckHandler) CreateDeck(c *gin.Context) {
	user, ok := middleware.GetAuthUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateDeckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !Language(req.SourceLanguage).IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source_language"})
		return
	}
	if !Language(req.TargetLanguage).IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_language"})
		return
	}
	if req.SourceLanguage == req.TargetLanguage {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_language and target_language must be different"})
		return
	}

	deck, err := h.deckService.CreateDeck(
		c.Request.Context(), user.ID, req.Name, req.SourceLanguage, req.TargetLanguage,
	)
	if err != nil {
		middleware.Logger.Error("failed to create deck", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create deck"})
		return
	}

	c.JSON(http.StatusCreated, CreateDeckResponse{
		ID:             deck.ID.String(),
		Name:           deck.Name,
		SourceLanguage: deck.SourceLanguage,
		TargetLanguage: deck.TargetLanguage,
		CreatedAt:      deck.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	})
}
