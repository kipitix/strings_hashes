package grpchandlers

import (
	"errors"
	"fmt"
	"hashkeeper/pkg/hashlog"
	"io"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	log "github.com/sirupsen/logrus"

	"hashcalc/internal/strhash"
	"hashcalc/pkg/grpchashcalc"
)

const (
	_hashMakerMode            = strhash.Mode256
	_hashMakerBuffersCapacity = 1000
)

type HashCalcServerImpl struct {
	grpchashcalc.UnimplementedHashCalcServer
}

var _ grpchashcalc.HashCalcServer = (*HashCalcServerImpl)(nil)

func (hcs *HashCalcServerImpl) Calc(stream grpchashcalc.HashCalc_CalcServer) error {
	hashlog.LogReqID(stream.Context()).Debug("hash calculations started")

	// Create hash maker
	hashMaker, err := strhash.NewSha3StrHashMaker(strhash.Sha3StrHashCfg{Mode: _hashMakerMode, BuffersCapacity: _hashMakerBuffersCapacity})
	if err != nil {
		return hashlog.WithStackErrorf("can`t init hash maker: %w", err)
	}

	// Start hash maker
	go func() {
		err := hashMaker.Run(stream.Context())
		if err != nil {
			hashlog.LogErrorWithStack(err).Error("calc failed")
		}
	}()

	hashlog.LogReqID(stream.Context()).Trace("hash maker running")

	// Create wait group to count received and send items
	var wg sync.WaitGroup

	// Process output data
	go func() {
		// OutChan must be closed when context canceled or data finished well
		for outData := range hashMaker.OutChan() {
			outItem := grpchashcalc.OutItem{
				Index: uint32(outData.Index),
				Hash:  outData.Hash,
			}

			err = stream.Send(&outItem)
			if err != nil {
				hashlog.LogErrorWithStack(err).Error("can`t send hash")
			}

			hashlog.LogReqID(stream.Context()).WithField("index", outItem.Index).WithField("hash", outItem.Hash).Trace("output sent")

			wg.Done()
		}
	}()

	// Process input
	for {
		inItem, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// Data receiving finished
			close(hashMaker.InChan())
			break
		} else if err != nil {
			hashlog.LogErrorWithStack(err).Error("can`t read data from input stream")
			return err
		}

		inData := strhash.InItem{
			Index: int(inItem.Index),
			Data:  inItem.Data,
		}

		wg.Add(1)

		hashMaker.InChan() <- inData

		hashlog.LogReqID(stream.Context()).WithField("index", inItem.Index).WithField("hash", inItem.Data).Trace("input received")
	}

	wg.Wait()

	hashlog.LogReqID(stream.Context()).Debug("hash calculations finished successfully")

	return nil
}

// TODO make or remove interceptor
func LogFileSizeInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Trace("Converting images to PDF...")

	rw := newRecvWrapper(ss)

	// Handle external func
	if err := handler(srv, rw); err != nil {
		hashlog.LogErrorWithStack(err).Error("Can`t send PDF data")
		return err
	}

	message := fmt.Sprintf("Count of received images: %d\n", len(rw.ContentSize))

	for k, v := range rw.ContentSize {
		message += fmt.Sprintf("Index: %d; Size: %d\n", k, v)
	}

	log.Trace(message)

	return nil
}

type recvWrapper struct {
	grpc.ServerStream
	ContentSize map[uint64]int
}

func newRecvWrapper(ss grpc.ServerStream) *recvWrapper {
	rw := &recvWrapper{ServerStream: ss}
	rw.ContentSize = make(map[uint64]int)
	return rw
}

func (rw *recvWrapper) RecvMsg(m interface{}) error {
	if err := rw.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	if p, ok := m.(proto.Message); ok {
		var index uint64
		var size int
		p.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
			switch fd.Name() {
			case "Content":
				size = len(v.Bytes())
			case "Index":
				index = v.Uint()
			}
			return true
		})

		rw.ContentSize[index] += size
	}

	return nil
}
