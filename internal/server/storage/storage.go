package storage

type Storage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGaugeValue(name string) (float64, bool)
	GetCounterValue(name string) (int64, bool)
	GetAllMetrics() map[string]interface{}
}

type MemStorage struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGauge(name string, value float64) {
	ms.GaugeMetrics[name] = value
}

func (ms *MemStorage) UpdateCounter(name string, value int64) {
	ms.CounterMetrics[name] += value
}

func (ms *MemStorage) GetGaugeValue(name string) (float64, bool) {
	value, ok := ms.GaugeMetrics[name]
	return value, ok
}

func (ms *MemStorage) GetCounterValue(name string) (int64, bool) {
	value, ok := ms.CounterMetrics[name]
	return value, ok
}

func (ms *MemStorage) GetAllMetrics() map[string]interface{} {
	return map[string]interface{}{
		"gauge":   ms.GaugeMetrics,
		"counter": ms.CounterMetrics,
	}
}
