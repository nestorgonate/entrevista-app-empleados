package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter arma el motor HTTP y monta cada handler bajo /api/v1.
func NewRouter(employee *EmployeeHandler, task *TaskHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	employee.RegisterRoutes(v1)
	task.RegisterRoutes(v1)

	return router
}
