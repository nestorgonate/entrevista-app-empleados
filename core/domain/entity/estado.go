package entity

// Estados validos de una Task.
const (
	EstadoPendiente  = "pendiente"
	EstadoEnProgreso = "en_progreso"
	EstadoCompletada = "completada"
	EstadoCancelada  = "cancelada"
)

// EstadoValido indica si el estado recibido es uno de los soportados.
func EstadoValido(estado string) bool {
	switch estado {
	case EstadoPendiente, EstadoEnProgreso, EstadoCompletada, EstadoCancelada:
		return true
	default:
		return false
	}
}
