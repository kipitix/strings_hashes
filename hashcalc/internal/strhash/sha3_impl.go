package strhash

import (
	"context"
	"hashkeeper/pkg/hashlog"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"

	"go.uber.org/atomic"
)

type ModeSHA3 int

const (
	ModeUndef ModeSHA3 = iota
	Mode224
	Mode256
	Mode384
	Mode512
)

type Sha3StrHashCfg struct {
	Mode            ModeSHA3
	BuffersCapacity int
}

type sha3StrHashMakerImpl struct {
	calcFunc func([]byte) []byte
	inChan   chan InItem
	outChan  chan OutItem

	canCloseOutChan atomic.Bool
	inProcessCount  atomic.Int32

	mut sync.Mutex
}

func NewSha3StrHashMaker(cfg Sha3StrHashCfg) (StrHashMaker, error) {
	res := &sha3StrHashMakerImpl{
		inChan:  make(chan InItem, cfg.BuffersCapacity),
		outChan: make(chan OutItem, cfg.BuffersCapacity),
	}

	switch cfg.Mode {
	case Mode224:
		res.calcFunc = calc224
	case Mode256:
		res.calcFunc = calc256
	case Mode384:
		res.calcFunc = calc384
	case Mode512:
		res.calcFunc = calc512
	default:
		return nil, errors.WithStack(hashlog.WithStackErrorf("creation of sha3 string to hash maker failed: wrong mode: %d", cfg.Mode))
	}

	return res, nil
}

var _ StrHashMaker = (*sha3StrHashMakerImpl)(nil)

func (mi *sha3StrHashMakerImpl) Run(ctx context.Context) error {
	if mi.calcFunc == nil {
		return hashlog.WithStackErrorf("hash maker is not inited")
	}

	for {
		select {
		case inItem, sent := <-mi.inChan:
			if !sent {
				mi.canCloseOutChan.Store(true)
				return nil
			}

			mi.inProcessCount.Add(1)
			go func(in InItem) {
				outData := mi.calcFunc(in.Data)
				mi.outChan <- OutItem{
					Index: in.Index,
					Hash:  outData,
				}

				mi.mut.Lock()
				mi.inProcessCount.Sub(1)
				if mi.inProcessCount.Load() == 0 && mi.canCloseOutChan.Load() {
					close(mi.outChan)
				}
				mi.mut.Unlock()

			}(inItem)

		case <-ctx.Done():
			mi.canCloseOutChan.Store(true)
			return hashlog.WithStackErrorf("finished by context, reason: %w", ctx.Err())
		}
	}
}

func (mi *sha3StrHashMakerImpl) InChan() chan<- InItem {
	mi.mut.Lock()
	defer mi.mut.Unlock()

	return mi.inChan
}

func (mi *sha3StrHashMakerImpl) OutChan() <-chan OutItem {
	mi.mut.Lock()
	defer mi.mut.Unlock()

	return mi.outChan
}

func calc224(in []byte) []byte {
	out := sha3.Sum224(in)
	return out[:]
}

func calc256(in []byte) []byte {
	out := sha3.Sum256(in)
	return out[:]
}

func calc384(in []byte) []byte {
	out := sha3.Sum384(in)
	return out[:]
}

func calc512(in []byte) []byte {
	out := sha3.Sum512(in)
	return out[:]
}
