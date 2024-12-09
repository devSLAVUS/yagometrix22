package router

import (
	"github.com/devSLAVUS/yagometrix22/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *handlers.Handlers) *gin.Engine {
	r := gin.New()

	r.POST("/update/:type/:name/:value", handler.UpdateMetricHandler)
	r.GET("/value/:type/:name", handler.GetMetricValueHandler)
	r.GET("/", handler.GetMetricsHandler)

	return r
}
