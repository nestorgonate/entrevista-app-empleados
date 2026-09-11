// Espejo de los DTOs de Go en core/domain/dto/.

export type Estado = 'pendiente' | 'en_progreso' | 'completada' | 'cancelada';

/** Debe coincidir con entity.EstadoValido() en core/domain/entity/estado.go */
export const ESTADOS: Estado[] = ['pendiente', 'en_progreso', 'completada', 'cancelada'];

export const ESTADO_LABEL: Record<Estado, string> = {
  pendiente: 'Pendiente',
  en_progreso: 'En progreso',
  completada: 'Completada',
  cancelada: 'Cancelada',
};

export const ESTADO_CLASS: Record<Estado, string> = {
  pendiente: 'bg-amber-100 text-amber-800 ring-amber-600/20',
  en_progreso: 'bg-sky-100 text-sky-800 ring-sky-600/20',
  completada: 'bg-emerald-100 text-emerald-800 ring-emerald-600/20',
  cancelada: 'bg-slate-200 text-slate-600 ring-slate-500/20',
};

export interface Employee {
  id: number;
  nombre: string;
  correo: string;
  cargo: string;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: number;
  titulo: string;
  /** El backend usa "description" en ingles mientras el resto va en espanol. */
  description: string;
  responsable_id: number;
  /** Solo viene si el repositorio hizo Preload("Responsable"). */
  responsable?: Employee;
  fecha_limite: string;
  estado: Estado;
  created_at: string;
  updated_at: string;
}

export interface CreateEmployeeRequest {
  nombre: string;
  correo: string;
  cargo: string;
}

export type UpdateEmployeeRequest = Partial<CreateEmployeeRequest>;

export interface CreateTaskRequest {
  titulo: string;
  description: string;
  responsable_id: number;
  /** RFC3339: el binding de Gin rechaza un "2026-09-30" pelado. */
  fecha_limite: string;
  estado?: Estado;
}
