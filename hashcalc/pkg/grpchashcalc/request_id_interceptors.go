package grpchashcalc

import (
	"context"
	"hashkeeper/pkg/hashlog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _requestIDKey = "requestID"

func RequestIDClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// Append the requestID info
	id := hashlog.GetReqID(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, _requestIDKey, id.String())
	return streamer(ctx, desc, cc, method)
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func newWrappedStream(s grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	return &wrappedStream{s, ctx}
}

func RequestIDServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	idStr := metadata.ValueFromIncomingContext(ss.Context(), _requestIDKey)
	if len(idStr) > 0 {
		ctxReqID := hashlog.AppendReqID(ss.Context(), uuid.MustParse(idStr[0]))
		return handler(srv, newWrappedStream(ss, ctxReqID))
	}

	return handler(srv, ss)
}
