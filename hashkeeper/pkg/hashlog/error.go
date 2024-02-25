package hashlog

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func WithStackErrorf(format string, a ...any) error {
	return errors.WithStack(fmt.Errorf(format, a...))
}

func LogErrorWithStack(err error) *logrus.Entry {
	return logrus.WithField("errorWithStackTrace", fmt.Sprintf("%+v", err))
}
