package find

import (
	"context"
	"hashkeeper/internal/domain/datahash"
	"hashkeeper/pkg/hashlog"

	"github.com/google/uuid"
)

// Find Request

type findRequestImpl struct {
	requestID uuid.UUID
	ids       []datahash.HashID
}

var _ FindRequest = (*findRequestImpl)(nil)

func NewFindRequest(in []datahash.HashID) FindRequest {
	return &findRequestImpl{
		requestID: uuid.New(),
		ids:       in,
	}
}

func (r findRequestImpl) RequestID() uuid.UUID {
	return r.requestID
}

func (r findRequestImpl) IDs() []datahash.HashID {
	return r.ids
}

// Find Response

type findResponseImpl struct {
	requestID uuid.UUID
	hashes    []datahash.Hash
}

var _ FindResponse = (*findResponseImpl)(nil)

func NewFindResponse(requestID uuid.UUID, out []datahash.Hash) FindResponse {
	return &findResponseImpl{
		requestID: requestID,
		hashes:    out,
	}
}

func (r findResponseImpl) RequestID() uuid.UUID {
	return r.requestID
}

func (r findResponseImpl) Hashes() []datahash.Hash {
	return r.hashes
}

// Find Handler

type findHandlerImpl struct {
	repo datahash.HashRepository
}

var _ FindHandler = (*findHandlerImpl)(nil)

func NewFindHandler(repo datahash.HashRepository) FindHandler {
	return &findHandlerImpl{
		repo: repo,
	}
}

func (h findHandlerImpl) Find(ctx context.Context, req FindRequest) (FindResponse, error) {
	restRes, err := h.repo.FindByID(ctx, req.IDs())
	if err != nil {
		return nil, hashlog.WithStackErrorf("find by id filed: %w", err)
	}

	resp := NewFindResponse(req.RequestID(), restRes)

	return resp, nil
}
