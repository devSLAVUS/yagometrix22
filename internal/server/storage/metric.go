package storage

import "fmt"

type Metric interface {
	GetValue() interface{}
	UpdateValue(value interface{}) error
}

type Gauge float64

func (g *Gauge) GetValue() interface{} {
	return float64(*g)
}

func (g *Gauge) UpdateValue(value interface{}) error {
	v, ok := value.(float64)
	if !ok {
		return fmt.Errorf("invalid value type for Gauge, expected float64")
	}
	*g = Gauge(v)
	return nil
}

type Counter int64

func (c *Counter) GetValue() interface{} {
	return int64(*c)
}

func (c *Counter) UpdateValue(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return fmt.Errorf("invalid value type for Counter, expected int64")
	}
	*c += Counter(v)
	return nil
}
