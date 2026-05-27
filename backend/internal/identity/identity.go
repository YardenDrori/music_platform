package identity

import (
	"context"

	"github.com/google/uuid"
)

type contextUserID struct{}

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, contextUserID{}, id)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(contextUserID{}).(uuid.UUID)

	return id, ok
}
