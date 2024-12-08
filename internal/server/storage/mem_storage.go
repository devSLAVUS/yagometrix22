package storage

type MemStorage struct {
	metrics map[string]Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]Metric),
	}
}

func (ms *MemStorage) UpdateGauge(name string, value float64) {
	if metric, ok := ms.metrics[name]; ok {
		_ = metric.UpdateValue(value)
	} else {
		g := Gauge(value)
		ms.metrics[name] = &g
	}
}

func (ms *MemStorage) UpdateCounter(name string, value int64) {
	if metric, ok := ms.metrics[name]; ok {
		_ = metric.UpdateValue(value)
	} else {
		c := Counter(value)
		ms.metrics[name] = &c
	}
}

func (ms *MemStorage) GetAllMetrics() map[string]interface{} {
	result := make(map[string]interface{})
	for name, metric := range ms.metrics {
		result[name] = metric.GetValue()
	}
	return result
}

func (ms *MemStorage) GetMetric(name string) (Metric, bool) {
	metric, ok := ms.metrics[name]
	return metric, ok
}
