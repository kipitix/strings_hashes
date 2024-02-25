package datahash

import "fmt"

type hashImpl struct {
	id   HashID
	hash HashContent
}

// Hash interface

var _ Hash = (*hashImpl)(nil)

func NewHash(id HashID, hash HashContent) Hash {
	return &hashImpl{
		id:   id,
		hash: hash,
	}
}

func (h hashImpl) ID() HashID {
	return h.id
}

func (h hashImpl) Hash() HashContent {
	return h.hash
}

// Stringer interface

var _ fmt.Stringer = (*hashImpl)(nil)

func (h hashImpl) String() string {
	return fmt.Sprintf("id: %d, hash: %s", h.id, h.hash)
}
