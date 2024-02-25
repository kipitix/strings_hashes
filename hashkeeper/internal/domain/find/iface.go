package find

import (
	"context"
	"hashkeeper/internal/domain/datahash"

	"github.com/google/uuid"
)

type FindRequest interface {
	RequestID() uuid.UUID
	IDs() []datahash.HashID
}

type FindResponse interface {
	RequestID() uuid.UUID
	Hashes() []datahash.Hash
}

type FindHandler interface {
	Find(context.Context, FindRequest) (FindResponse, error)
}
