package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"entrevista/core/domain/dto"
	"entrevista/core/port"
)

type TaskHandler struct {
	svc port.TaskService
}

// NewTaskHandler recibe el caso de uso de tareas por inyeccion de dependencias.
func NewTaskHandler(svc port.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// RegisterRoutes monta las rutas de tareas sobre el grupo recibido.
func (h *TaskHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tasks := rg.Group("/tasks")
	{
		tasks.POST("", h.Create)
		tasks.GET("", h.List)
		tasks.GET("/:id", h.GetByID)
		tasks.PUT("/:id", h.Update)
		tasks.PATCH("/:id/responsable", h.Assign)
		tasks.PATCH("/:id/estado", h.UpdateEstado)
		tasks.DELETE("/:id", h.Delete)
	}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}

	task, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, task)
}

// List devuelve todas las tareas, o solo las de un responsable si se pasa
// el query param ?responsable_id=.
func (h *TaskHandler) List(c *gin.Context) {
	if raw := c.Query("responsable_id"); raw != "" {
		id, ok := parseUint(c, raw, "responsable_id")
		if !ok {
			return
		}
		tasks, err := h.svc.ListByResponsable(c.Request.Context(), id)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, tasks)
		return
	}

	tasks, err := h.svc.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	task, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}

	task, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// Assign reasigna la tarea a otro empleado.
func (h *TaskHandler) Assign(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}

	task, err := h.svc.Assign(c.Request.Context(), id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateEstado(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}

	task, err := h.svc.UpdateEstado(c.Request.Context(), id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
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
