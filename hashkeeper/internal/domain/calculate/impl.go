package calculate

import (
	"context"
	"hashkeeper/internal/domain/datahash"
	"hashkeeper/pkg/hashlog"

	"github.com/google/uuid"
)

// Calculate Request

type calculateRequestImpl struct {
	requestID uuid.UUID
	strings   []datahash.StringContent
}

var _ CalculateRequest = (*calculateRequestImpl)(nil)

func NewCalculateRequest(in []datahash.StringContent) CalculateRequest {
	return &calculateRequestImpl{
		requestID: uuid.New(),
		strings:   in,
	}
}

func (r calculateRequestImpl) RequestID() uuid.UUID {
	return r.requestID
}

func (r calculateRequestImpl) Strings() []datahash.StringContent {
	return r.strings
}

// Calculate Response

type calculateResponseImpl struct {
	requestID uuid.UUID
	hashes    []datahash.Hash
}

var _ CalculateResponse = (*calculateResponseImpl)(nil)

func NewCalculateResponse(requestID uuid.UUID, out []datahash.Hash) CalculateResponse {
	return &calculateResponseImpl{
		requestID: requestID,
		hashes:    out,
	}
}

func (r calculateResponseImpl) RequestID() uuid.UUID {
	return r.requestID
}

func (r calculateResponseImpl) Hashes() []datahash.Hash {
	return r.hashes
}

// Calculate Handler

type calculateHandlerImpl struct {
	maker datahash.HashMaker
	repo  datahash.HashRepository
}

var _ CalculateHandler = (*calculateHandlerImpl)(nil)

func NewCalculateHandler(maker datahash.HashMaker, repo datahash.HashRepository) CalculateHandler {
	return &calculateHandlerImpl{
		maker: maker,
		repo:  repo,
	}
}

func (h calculateHandlerImpl) Calculate(ctx context.Context, req CalculateRequest) (CalculateResponse, error) {
	calcRes, err := h.maker.Make(ctx, req.Strings())

	if err != nil {
		return nil, hashlog.WithStackErrorf("calculate make failed: %w", err)
	}

	err = h.repo.Store(ctx, calcRes)
	if err != nil {
		return nil, hashlog.WithStackErrorf("calculate store failed: %w", err)
	}

	storeRes, err := h.repo.FindByContent(ctx, calcRes)
	if err != nil {
		return nil, hashlog.WithStackErrorf("calculate find failed: %w", err)
	}

	resp := NewCalculateResponse(req.RequestID(), storeRes)

	return resp, nil
}
