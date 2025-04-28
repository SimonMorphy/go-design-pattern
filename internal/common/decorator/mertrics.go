package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MetricsRecord interface {
	Inc(key string, val int)
}
type MetricsDecorator[I, O any] struct {
	base     Handler[I, O]
	recorder MetricsRecord
}

func (m MetricsDecorator[I, O]) Handle(ctx context.Context, in I) (result O, err error) {
	start := time.Now()
	defer func() {
		end := time.Since(start)
		res := fmt.Sprintf("%s 's durationo", strings.Split(fmt.Sprintf("%T", in), ".")[1])
		m.recorder.Inc(res, int(end.Seconds()))
		if err == nil {
			m.recorder.Inc(res+"success", 1)
		} else {
			m.recorder.Inc(res+"failed", 1)
		}
	}()
	return m.base.Handle(ctx, in)
}
