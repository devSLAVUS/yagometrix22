package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devSLAVUS/yagometrix22/internal/server/handlers"
	"github.com/devSLAVUS/yagometrix22/internal/server/router"
	"github.com/devSLAVUS/yagometrix22/internal/server/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	store := storage.NewMemStorage()
	store.UpdateGauge("TestGauge", 123.45)

	metrics := store.GetAllMetrics()
	t.Logf("Metrics: %+v", metrics) // Для отладки

	value, exists := metrics["TestGauge"]
	if !exists {
		t.Fatalf("Gauge metric 'TestGauge' does not exist")
	}

	assert.Equal(t, 123.45, value.(float64), "bad gauge value")
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	store := storage.NewMemStorage()
	store.UpdateCounter("TestCounter", 10)
	store.UpdateCounter("TestCounter", 5)

	metrics := store.GetAllMetrics()
	t.Logf("Metrics: %+v", metrics) // Для отладки

	value, exists := metrics["TestCounter"]
	if !exists {
		t.Fatalf("Counter metric 'TestCounter' does not exist")
	}

	assert.Equal(t, int64(15), value.(int64), "bad counter value")
}

func TestUpdateMetricHandler(t *testing.T) {
	store := storage.NewMemStorage()
	handler := handlers.NewHandlers(store)

	gin.SetMode(gin.TestMode)
	r := router.NewRouter(handler)

	// Test for updating gauge metric
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/TestGauge/123.45", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 123.45, store.GetAllMetrics()["TestGauge"].(float64))

	// Test for updating counter metric
	req = httptest.NewRequest(http.MethodPost, "/update/counter/TestCounter/10", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), store.GetAllMetrics()["TestCounter"].(int64))

	// Test for invalid metric type
	req = httptest.NewRequest(http.MethodPost, "/update/invalid/TestMetric/123", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMetricsHandler(t *testing.T) {
	store := storage.NewMemStorage()
	store.UpdateGauge("TestGauge", 123.45)
	store.UpdateCounter("TestCounter", 10)

	handler := handlers.NewHandlers(store)
	gin.SetMode(gin.TestMode)
	r := router.NewRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := `{"TestGauge":123.45,"TestCounter":10}`
	assert.JSONEq(t, expectedBody, rec.Body.String(), "bad response")
}
