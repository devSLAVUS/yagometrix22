package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Storage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetMetrics() map[string]interface{}
	GetGaugeValue(name string) (float64, error)
	GetCounterValue(name string) (int64, error)
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

func (ms *MemStorage) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"gauge":   ms.GaugeMetrics,
		"counter": ms.CounterMetrics,
	}
}
func (ms *MemStorage) GetGaugeValue(name string) (float64, error) {
	return ms.GaugeMetrics[name], nil
}
func (ms *MemStorage) GetCounterValue(name string) (int64, error) {
	return ms.CounterMetrics[name], nil

}

func UpdateHandler(storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("type")
		metricName := c.Param("name")
		metricValue := c.Param("value")

		if metricName == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Metric name is required"})
			return
		}

		if metricType != Gauge && metricType != Counter {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
			return
		}

		switch metricType {
		case Gauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gauge value"})
				return
			}
			storage.UpdateGauge(metricName, value)
		case Counter:
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter value"})
				return
			}
			storage.UpdateCounter(metricName, value)
		}

		c.String(http.StatusOK, "OK")

	}
}

func getMetricsHandler(storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, storage.GetMetrics())
	}
}

func getValueHandler(storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("type")
		metricName := c.Param("name")
		if metricType != Gauge && metricType != Counter {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
			return
		}
		switch metricType {
		case Gauge:
			value, err := storage.GetGaugeValue(metricName)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gauge name"})
				return
			}
			c.JSON(http.StatusOK, value)

		case Counter:
			value, err := storage.GetCounterValue(metricName)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter name"})
				return
			}
			c.JSON(http.StatusOK, value)
		}
	}
}

func main() {
	r := gin.Default()
	var storage Storage = NewMemStorage()

	r.POST("/update/:type/:name/:value", UpdateHandler(storage))
	r.GET("/value/:type/:name", getValueHandler(storage))
	r.GET("/", getMetricsHandler(storage))

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
