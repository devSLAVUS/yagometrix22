package router

import (
	"github.com/devSLAVUS/yagometrix22/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

// NewRouter создает новый маршрутизатор с зарегистрированными маршрутами.
func NewRouter(handler *handlers.Handlers) *gin.Engine {
	r := gin.New()

	// Регистрация маршрутов
	r.POST("/update/:type/:name/:value", handler.UpdateMetricHandler)
	r.GET("/value/:type/:name", handler.GetMetricHandler)
	r.GET("/", handler.GetAllMetricsHandler)

	return r
}
