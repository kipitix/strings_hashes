package grpccalc

import (
	"context"
	"errors"
	"fmt"
	"hashkeeper/internal/domain/datahash"
	"hashkeeper/pkg/hashlog"
	"io"
	"time"

	"hashcalc/pkg/grpchashcalc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcHashMaker interface {
	datahash.HashMaker

	Dial(context.Context) error
	Close() error
}

type hashMakerImpl struct {
	address        string
	dialTimeout    time.Duration
	requestTimeout time.Duration
	connection     *grpc.ClientConn
	client         grpchashcalc.HashCalcClient
}

var _ GrpcHashMaker = (*hashMakerImpl)(nil)

type HashMakerCfg struct {
	HashCalcServerAddress        string `arg:"--hash-calc-server-address,env:HASH_CALC_SERVER_ADDRESS" default:"localhost:50051" help:"Address and port of gRPC hash calc server"`
	HashCalcServerDialTimeout    string `arg:"--hash-calc-server-dial-timeout,env:HASH_CALC_SERVER_DIAL_TIMEOUT" default:"5s" help:"Dial to server timeout"`
	HashCalcServerRequestTimeout string `arg:"--hash-calc-server-request-timeout,env:HASH_CALC_SERVER_REQUEST_TIMEOUT" default:"60s" help:"Request to server timeout"`
}

func NewHashMaker(cfg HashMakerCfg) (GrpcHashMaker, error) {
	var res hashMakerImpl

	res.address = cfg.HashCalcServerAddress

	if dur, err := time.ParseDuration(cfg.HashCalcServerDialTimeout); err != nil {
		return nil, hashlog.WithStackErrorf("wrong dial timeout: %w", err)
	} else {
		res.dialTimeout = dur
	}

	if dur, err := time.ParseDuration(cfg.HashCalcServerDialTimeout); err != nil {
		return nil, hashlog.WithStackErrorf("wrong request timeout: %w", err)
	} else {
		res.requestTimeout = dur
	}

	return &res, nil
}

func (hm hashMakerImpl) Make(ctx context.Context, in []datahash.StringContent) ([]datahash.HashContent, error) {
	hashlog.LogReqID(ctx).Debug("making of hashes from strings started")
	hashlog.LogReqID(ctx).WithField("input", in).Trace("received strings")

	// Prepare
	calcCtx, calcCancel := context.WithTimeout(ctx, hm.requestTimeout)
	defer calcCancel()

	stream, err := hm.client.Calc(calcCtx)
	if err != nil {
		return nil, hashlog.WithStackErrorf("can`t create request stream: %w", err)
	}

	// Send data
	for index, str := range in {
		reqItem := grpchashcalc.InItem{
			Index: uint32(index),
			Data:  []byte(str),
		}

		err = stream.Send(&reqItem)
		if err != nil {
			return nil, hashlog.WithStackErrorf("can`t send request item: %w", err)
		}
	}

	// Finish sending
	if err := stream.CloseSend(); err != nil {
		return nil, hashlog.WithStackErrorf("can`t close sending: %w", err)
	}

	// Start receiving answer
	result := make([]datahash.HashContent, len(in))

	for {
		respItem, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, hashlog.WithStackErrorf("can`t receive response item: %w", err)
		}

		result[int(respItem.Index)] = datahash.HashContent(fmt.Sprintf("%x", respItem.Hash))
	}

	hashlog.LogReqID(ctx).WithField("hashes", result).Trace("result hashes")
	hashlog.LogReqID(ctx).Debug("making of hashes from strings finished successfully")

	return result, nil
}

func (hm *hashMakerImpl) Dial(ctx context.Context) error {
	connCtx, connCancel := context.WithTimeout(ctx, hm.dialTimeout)
	defer connCancel()

	if conn, err := grpc.DialContext(
		connCtx,
		hm.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpchashcalc.RequestIDClientInterceptor),
		grpc.WithBlock()); err != nil {
		return hashlog.WithStackErrorf("can`t connect to hash calc server: %w", err)
	} else {
		hm.connection = conn
		hm.client = grpchashcalc.NewHashCalcClient(conn)
	}

	return nil
}

func (hm *hashMakerImpl) Close() error {
	err := hm.connection.Close()

	hm.connection = nil
	hm.client = nil

	return err
}
