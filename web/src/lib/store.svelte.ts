import * as api from './api';
import { ApiError } from './api';
import type { CreateEmployeeRequest, CreateTaskRequest, Employee, Estado, Task } from './types';

export type Aviso = { tipo: 'error' | 'exito'; texto: string } | null;

/**
 * Estado compartido de la app. Se usa una clase con campos `$state` porque
 * exportar un `let $state` desde un modulo no propaga la reactividad al
 * importarlo: las propiedades de una instancia si lo hacen.
 */
class AppStore {
  employees = $state<Employee[]>([]);
  tasks = $state<Task[]>([]);
  cargando = $state(false);
  aviso = $state<Aviso>(null);

  /** Filtro de la vista de tareas; 0 = todas. */
  filtroResponsable = $state(0);

  private notificar(tipo: 'error' | 'exito', texto: string) {
    this.aviso = { tipo, texto };
  }

  limpiarAviso() {
    this.aviso = null;
  }

  /**
   * Envuelve cada llamada a la API: marca el estado de carga y convierte
   * cualquier ApiError en un aviso legible. Devuelve `true` si salio bien,
   * para que los formularios sepan si deben limpiarse.
   */
  private async ejecutar(accion: () => Promise<void>, exito?: string): Promise<boolean> {
    this.cargando = true;
    try {
      await accion();
      if (exito) this.notificar('exito', exito);
      return true;
    } catch (err) {
      const texto = err instanceof ApiError ? err.message : 'Ocurrió un error inesperado';
      this.notificar('error', texto);
      return false;
    } finally {
      this.cargando = false;
    }
  }

  // --- Carga ---

  async cargarTodo() {
    await this.ejecutar(async () => {
      const [employees, tasks] = await Promise.all([api.listEmployees(), api.listTasks()]);
      this.employees = employees;
      this.tasks = tasks;
    });
  }

  async recargarTareas() {
    const tasks = await api.listTasks(this.filtroResponsable || undefined);
    this.tasks = tasks;
  }

  async aplicarFiltro(responsableId: number) {
    this.filtroResponsable = responsableId;
    await this.ejecutar(() => this.recargarTareas());
  }

  // --- Empleados ---

  crearEmpleado(data: CreateEmployeeRequest) {
    return this.ejecutar(async () => {
      const employee = await api.createEmployee(data);
      this.employees = [...this.employees, employee];
    }, 'Empleado registrado');
  }

  actualizarEmpleado(id: number, data: CreateEmployeeRequest) {
    return this.ejecutar(async () => {
      const employee = await api.updateEmployee(id, data);
      this.employees = this.employees.map((e) => (e.id === id ? employee : e));
      // El nombre del responsable aparece dentro de cada tarea.
      await this.recargarTareas();
    }, 'Empleado actualizado');
  }

  eliminarEmpleado(id: number) {
    return this.ejecutar(async () => {
      await api.deleteEmployee(id);
      this.employees = this.employees.filter((e) => e.id !== id);
      if (this.filtroResponsable === id) this.filtroResponsable = 0;
      await this.recargarTareas();
    }, 'Empleado eliminado');
  }

  // --- Tareas ---

  crearTarea(data: CreateTaskRequest) {
    return this.ejecutar(async () => {
      await api.createTask(data);
      await this.recargarTareas();
    }, 'Tarea creada');
  }

  reasignarTarea(id: number, responsableId: number) {
    return this.ejecutar(async () => {
      await api.assignTask(id, responsableId);
      await this.recargarTareas();
    }, 'Tarea reasignada');
  }

  cambiarEstado(id: number, estado: Estado) {
    return this.ejecutar(async () => {
      const task = await api.updateEstado(id, estado);
      // El PATCH de estado no precarga el responsable: se conserva el que ya teniamos.
      this.tasks = this.tasks.map((t) => (t.id === id ? { ...task, responsable: t.responsable } : t));
    }, 'Estado actualizado');
  }

  eliminarTarea(id: number) {
    return this.ejecutar(async () => {
      await api.deleteTask(id);
      this.tasks = this.tasks.filter((t) => t.id !== id);
    }, 'Tarea eliminada');
  }
}

export const store = new AppStore();
