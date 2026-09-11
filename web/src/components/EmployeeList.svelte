<script lang="ts">
  import { store } from '../lib/store.svelte';
  import type { Employee } from '../lib/types';

  interface Props {
    onEditar: (employee: Employee) => void;
  }

  let { onEditar }: Props = $props();

  // Cuantas tareas tiene cada empleado, recalculado solo si cambian las tareas.
  const conteoTareas = $derived.by(() => {
    const conteo: Record<number, number> = {};
    for (const task of store.tasks) {
      conteo[task.responsable_id] = (conteo[task.responsable_id] ?? 0) + 1;
    }
    return conteo;
  });

  async function confirmarBorrado(employee: Employee) {
    const pendientes = conteoTareas[employee.id] ?? 0;
    const aviso = pendientes
      ? `${employee.nombre} tiene ${pendientes} tarea(s) asignada(s). ¿Eliminar de todos modos?`
      : `¿Eliminar a ${employee.nombre}?`;
    if (confirm(aviso)) await store.eliminarEmpleado(employee.id);
  }
</script>

<div class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
  {#if store.employees.length === 0}
    <p class="p-8 text-center text-sm text-slate-500">
      Todavía no hay empleados registrados.
    </p>
  {:else}
    <table class="w-full text-left text-sm">
      <thead class="border-b border-slate-200 bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
        <tr>
          <th class="px-4 py-3 font-medium">Nombre</th>
          <th class="px-4 py-3 font-medium">Correo</th>
          <th class="px-4 py-3 font-medium">Cargo</th>
          <th class="px-4 py-3 font-medium">Tareas</th>
          <th class="px-4 py-3"><span class="sr-only">Acciones</span></th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-100">
        {#each store.employees as employee (employee.id)}
          <tr class="transition hover:bg-slate-50/70">
            <td class="px-4 py-3 font-medium text-slate-900">{employee.nombre}</td>
            <td class="px-4 py-3 text-slate-600">{employee.correo}</td>
            <td class="px-4 py-3 text-slate-600">{employee.cargo}</td>
            <td class="px-4 py-3 text-slate-600 tabular-nums">
              {conteoTareas[employee.id] ?? 0}
            </td>
            <td class="px-4 py-3 text-right whitespace-nowrap">
              <button class="boton-secundario" onclick={() => onEditar(employee)}>Editar</button>
              <button
                class="boton-peligro"
                onclick={() => confirmarBorrado(employee)}
                disabled={store.cargando}
              >
                Eliminar
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
