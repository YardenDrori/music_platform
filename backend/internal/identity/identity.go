package identity

import (
	"context"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

type contextUserID struct{}

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, contextUserID{}, id)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(contextUserID{}).(uuid.UUID)
	if !ok {
		return uuid.Nil, apperrors.NewErrInternal("").
			WithInternal("identity missing or formatted incorrectly in context")
	}
	return id, nil
}
