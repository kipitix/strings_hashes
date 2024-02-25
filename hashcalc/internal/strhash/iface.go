package strhash

import "context"

type StrHashMaker interface {
	Run(context.Context) error

	InChan() chan<- InItem
	OutChan() <-chan OutItem
}

type InItem struct {
	Index int
	Data  []byte
}

type OutItem struct {
	Index int
	Hash  []byte
}
