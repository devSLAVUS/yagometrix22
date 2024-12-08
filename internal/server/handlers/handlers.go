package handlers

import (
	"net/http"
	"strconv"

	"github.com/devSLAVUS/yagometrix22/internal/server/storage"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	storage *storage.MemStorage
}

func NewHandlers(store *storage.MemStorage) *Handlers {
	return &Handlers{storage: store}
}

func (h *Handlers) UpdateMetricHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")

	if metricType == "gauge" {
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gauge value"})
			return
		}
		h.storage.UpdateGauge(metricName, value)
	} else if metricType == "counter" {
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter value"})
			return
		}
		h.storage.UpdateCounter(metricName, value)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
		return
	}

	c.String(http.StatusOK, "OK")
}

func (h *Handlers) GetMetricsHandler(c *gin.Context) {
	metrics := h.storage.GetAllMetrics()
	c.JSON(http.StatusOK, metrics)
}

func (h *Handlers) GetMetricValueHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	if metricType != "gauge" && metricType != "counter" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
		return
	}

	metric, ok := h.storage.GetMetric(metricName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": metric.GetValue()})
}
