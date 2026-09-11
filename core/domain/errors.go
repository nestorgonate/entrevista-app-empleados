package domain

import "errors"

// Errores de dominio. Los adapters los traducen al protocolo que corresponda
// (por ejemplo, el handler HTTP los mapea a codigos de estado).
var (
	ErrNotFound     = errors.New("recurso no encontrado")
	ErrConflict     = errors.New("el recurso ya existe")
	ErrInvalidInput = errors.New("datos invalidos")
)
