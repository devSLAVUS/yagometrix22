package handlers

import (
	"net/http"
	"strconv"

	"github.com/devSLAVUS/yagometrix22/internal/server/storage"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	storage storage.Storage
}

func NewHandlers(storage storage.Storage) *Handlers {
	return &Handlers{storage: storage}
}

func (h *Handlers) UpdateMetricHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")

	if metricName == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metric name is required"})
		return
	}

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gauge value"})
			return
		}
		h.storage.UpdateGauge(metricName, value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter value"})
			return
		}
		h.storage.UpdateCounter(metricName, value)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
		return
	}

	c.String(http.StatusOK, "OK")
}

func (h *Handlers) GetMetricHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	switch metricType {
	case "gauge":
		value, ok := h.storage.GetGaugeValue(metricName)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gauge not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})
	case "counter":
		value, ok := h.storage.GetCounterValue(metricName)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Counter not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
	}
}

func (h *Handlers) GetAllMetricsHandler(c *gin.Context) {
	metrics := h.storage.GetAllMetrics()
	c.JSON(http.StatusOK, metrics)
}
