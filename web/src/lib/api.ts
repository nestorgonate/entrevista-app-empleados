import type {
  CreateEmployeeRequest,
  CreateTaskRequest,
  Employee,
  Estado,
  Task,
  UpdateEmployeeRequest,
} from './types';

const BASE = '/api/v1';

/** Error con el mensaje real que devuelve el backend en {"error": "..."}. */
export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  let res: Response;
  try {
    res = await fetch(`${BASE}${path}`, {
      ...init,
      headers: init?.body ? { 'Content-Type': 'application/json' } : undefined,
    });
  } catch {
    throw new ApiError(0, 'No se pudo contactar con el servidor. ¿Está corriendo la API en :8000?');
  }

  if (!res.ok) {
    // El handler siempre responde {"error": "..."}; si no, caemos al statusText.
    const body = await res.json().catch(() => null);
    const message =
      body && typeof body === 'object' && typeof body.error === 'string'
        ? body.error
        : `Error ${res.status}: ${res.statusText}`;
    throw new ApiError(res.status, message);
  }

  // Los DELETE responden 204 sin cuerpo: no intentar parsear JSON.
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}

const json = (body: unknown): RequestInit['body'] => JSON.stringify(body);

// --- Empleados ---

export const listEmployees = () => request<Employee[]>('/employees');

export const createEmployee = (data: CreateEmployeeRequest) =>
  request<Employee>('/employees', { method: 'POST', body: json(data) });

export const updateEmployee = (id: number, data: UpdateEmployeeRequest) =>
  request<Employee>(`/employees/${id}`, { method: 'PUT', body: json(data) });

export const deleteEmployee = (id: number) =>
  request<void>(`/employees/${id}`, { method: 'DELETE' });

// --- Tareas ---

export const listTasks = (responsableId?: number) =>
  request<Task[]>(responsableId ? `/tasks?responsable_id=${responsableId}` : '/tasks');

export const createTask = (data: CreateTaskRequest) =>
  request<Task>('/tasks', { method: 'POST', body: json(data) });

/** Reasigna la tarea a otro empleado. */
export const assignTask = (id: number, responsableId: number) =>
  request<Task>(`/tasks/${id}/responsable`, {
    method: 'PATCH',
    body: json({ responsable_id: responsableId }),
  });

export const updateEstado = (id: number, estado: Estado) =>
  request<Task>(`/tasks/${id}/estado`, { method: 'PATCH', body: json({ estado }) });

export const deleteTask = (id: number) => request<void>(`/tasks/${id}`, { method: 'DELETE' });
