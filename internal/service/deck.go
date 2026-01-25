package service

import (
	"context"

	"github.com/google/uuid"

	"mova-backend/internal/database"
	"mova-backend/internal/repository"
)

type DeckService struct {
	deckRepo *repository.DeckRepository
}

func NewDeckService(deckRepo *repository.DeckRepository) *DeckService {
	return &DeckService{deckRepo: deckRepo}
}

func (s *DeckService) CreateDeck(
	ctx context.Context, userID uuid.UUID, name, sourceLang, targetLang string,
) (database.Deck, error) {
	return s.deckRepo.Create(ctx, userID, name, sourceLang, targetLang)
}

func (s *DeckService) ListDecks(ctx context.Context, userID uuid.UUID) ([]database.Deck, error) {
	return s.deckRepo.ListByUserID(ctx, userID)
}
