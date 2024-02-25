package calculate

import (
	"context"

	"github.com/google/uuid"

	"hashkeeper/internal/domain/datahash"
)

type CalculateRequest interface {
	RequestID() uuid.UUID
	Strings() []datahash.StringContent
}

type CalculateResponse interface {
	RequestID() uuid.UUID
	Hashes() []datahash.Hash
}

type CalculateHandler interface {
	Calculate(context.Context, CalculateRequest) (CalculateResponse, error)
}
