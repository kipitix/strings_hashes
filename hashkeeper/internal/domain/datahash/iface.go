package datahash

import "context"

type HashID int
type HashContent string

type StringContent string

type Hash interface {
	ID() HashID
	Hash() HashContent
}

type HashMaker interface {
	Make(context.Context, []StringContent) ([]HashContent, error)
}

type HashRepository interface {
	Store(context.Context, []HashContent) error

	FindByContent(context.Context, []HashContent) ([]Hash, error)
	FindByID(context.Context, []HashID) ([]Hash, error)
}
