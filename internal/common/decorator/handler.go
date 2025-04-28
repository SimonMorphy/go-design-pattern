package decorator

import (
	"context"
	"github.com/sirupsen/logrus"
)

// Handler <-* Impl Here *-> DecoratorPattern
type Handler[In, Out any] interface {
	Handle(ctx context.Context, query In) (Out, error)
}

type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, command C) (R, error)
}

type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

func ApplyHandlerDecorators[I, O any](
	handler Handler[I, O],
	logger *logrus.Entry,
	record MetricsRecord,
) Handler[I, O] {
	return LogDecorator[I, O]{
		logger: logger,
		base: MetricsDecorator[I, O]{
			base:     handler,
			recorder: record,
		},
	}
}
