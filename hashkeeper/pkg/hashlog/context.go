package hashlog

import (
	"context"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

type requestID string

var _requestIDKey requestID = "requestID"

func AppendReqID(ctx context.Context, uuid uuid.UUID) context.Context {
	return context.WithValue(ctx, _requestIDKey, uuid)
}

func GetReqID(ctx context.Context) uuid.UUID {
	return ctx.Value(_requestIDKey).(uuid.UUID)
}

func LogReqID(ctx context.Context) *logrus.Entry {
	return logrus.WithField(string(_requestIDKey), ctx.Value(_requestIDKey))
}

func LogWithReqID(ent *logrus.Entry, ctx context.Context) *logrus.Entry {
	return ent.WithField(string(_requestIDKey), ctx.Value(_requestIDKey))
}
