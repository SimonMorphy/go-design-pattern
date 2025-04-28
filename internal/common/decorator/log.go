package decorator

import (
	"context"
	"github.com/sirupsen/logrus"
)

type LogDecorator[I, O any] struct {
	logger *logrus.Entry
	base   Handler[I, O]
}

func (l LogDecorator[I, O]) Handle(ctx context.Context, in I) (out O, err error) {
	logger := logrus.StandardLogger()
	logrus.Debug("Handler Executing...")
	defer func() {
		if err == nil {
			logger.Info("Handler handled Successfully")
		} else {
			logger.Error("Failed to Handle...")
		}
	}()
	return l.base.Handle(ctx, in)
}
