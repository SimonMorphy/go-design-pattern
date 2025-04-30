package metrics

type TodoMetrics struct {
}

func NewTodoMetrics() *TodoMetrics {
	return &TodoMetrics{}
}

func (t TodoMetrics) Inc(k string, v int) {
	//logrus.Infof(k, ":", v)
}
