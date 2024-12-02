package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	storage := NewMemStorage()
	storage.UpdateGauge("TestGauge", 123.45)

	assert.Equal(t, 123.45, storage.GaugeMetrics["TestGauge"], "bad gauge")
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	storage := NewMemStorage()
	storage.UpdateCounter("TestCounter", 10)
	storage.UpdateCounter("TestCounter", 5)

	assert.Equal(t, int64(15), storage.CounterMetrics["TestCounter"], "bad counter")
}

func TestUpdateHandler(t *testing.T) {

	storage := NewMemStorage()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/update/:type/:name/:value", updHandler(storage))

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/TestGauge/123.45", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 123.45, storage.GaugeMetrics["TestGauge"])

	req = httptest.NewRequest(http.MethodPost, "/update/counter/TestCounter/10", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), storage.CounterMetrics["TestCounter"])

	req = httptest.NewRequest(http.MethodPost, "/update/invalid/TestMetric/123", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMetricsHandler(t *testing.T) {
	storage := NewMemStorage()
	storage.UpdateGauge("TestGauge", 123.45)
	storage.UpdateCounter("TestCounter", 10)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", getMetricsHandler(storage))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := `{"gauge":{"TestGauge":123.45},"counter":{"TestCounter":10}}`
	assert.JSONEq(t, expectedBody, rec.Body.String(), "bad response")
}
