package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"entrevista/core/domain/dto"
	"entrevista/core/port"
)

type EmployeeHandler struct {
	svc     port.EmployeeService
	taskSvc port.TaskService
}

// NewEmployeeHandler recibe los casos de uso por inyeccion de dependencias.
// taskSvc se usa solo para listar las tareas de un empleado.
func NewEmployeeHandler(svc port.EmployeeService, taskSvc port.TaskService) *EmployeeHandler {
	return &EmployeeHandler{svc: svc, taskSvc: taskSvc}
}

// RegisterRoutes monta las rutas de empleados sobre el grupo recibido.
func (h *EmployeeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	employees := rg.Group("/employees")
	{
		employees.POST("", h.Create)
		employees.GET("", h.List)
		employees.GET("/:id", h.GetByID)
		employees.PUT("/:id", h.Update)
		employees.DELETE("/:id", h.Delete)
		employees.GET("/:id/tasks", h.ListTasks)
	}
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}

	employee, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, employee)
}

func (h *EmployeeHandler) List(c *gin.Context) {
	employees, err := h.svc.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	employee, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}

	employee, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListTasks devuelve las tareas asignadas a un empleado.
func (h *EmployeeHandler) ListTasks(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	tasks, err := h.taskSvc.ListByResponsable(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}
