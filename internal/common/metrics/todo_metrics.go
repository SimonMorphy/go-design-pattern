package metrics

import "github.com/sirupsen/logrus"

type TodoMetrics struct {
}

func NewTodoMetrics() *TodoMetrics {
	return &TodoMetrics{}
}

func (t TodoMetrics) Inc(k string, v int) {
	logrus.Info(k, ":", v)
}
