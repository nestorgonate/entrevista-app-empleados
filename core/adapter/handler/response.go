package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"entrevista/core/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

// respondError traduce los errores de dominio al codigo HTTP correspondiente.
// Es el unico punto del adaptador que conoce esa equivalencia.
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "error interno del servidor"})
	}
}

func respondBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
}

// parseIDParam lee un parametro de ruta numerico y responde 400 si es invalido.
func parseIDParam(c *gin.Context, name string) (uint, bool) {
	return parseUint(c, c.Param(name), name)
}

// parseUint convierte un valor a uint positivo y responde 400 si es invalido.
func parseUint(c *gin.Context, raw, name string) (uint, bool) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "el parametro " + name + " debe ser un entero positivo"})
		return 0, false
	}
	return uint(id), true
}
