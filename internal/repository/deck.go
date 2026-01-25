package repository

import (
	"context"

	"github.com/google/uuid"

	"mova-backend/internal/database"
)

type DeckRepository struct {
	queries *database.Queries
}

func NewDeckRepository(queries *database.Queries) *DeckRepository {
	return &DeckRepository{queries: queries}
}

func (r *DeckRepository) Create(
	ctx context.Context, userID uuid.UUID, name, sourceLang, targetLang string,
) (database.Deck, error) {
	return r.queries.CreateDeck(ctx, database.CreateDeckParams{
		UserID:         userID,
		Name:           name,
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
	})
}

func (r *DeckRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]database.Deck, error) {
	return r.queries.GetDecksByUserID(ctx, userID)
}
