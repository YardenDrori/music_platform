package summaries

import "github.com/google/uuid"

type UserSummary struct {
	ID            uuid.UUID
	Username      string
	ProfilePicKey *uuid.UUID
}
