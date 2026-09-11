<script lang="ts">
  import { store } from '../lib/store.svelte';
  import { ESTADO_CLASS, ESTADO_LABEL, ESTADOS, type Estado, type Task } from '../lib/types';

  const formatoFecha = new Intl.DateTimeFormat('es', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  });

  const hoy = new Date().setHours(0, 0, 0, 0);

  /** Una tarea esta vencida si paso su fecha limite y sigue abierta. */
  function vencida(task: Task) {
    return (
      new Date(task.fecha_limite).getTime() < hoy &&
      task.estado !== 'completada' &&
      task.estado !== 'cancelada'
    );
  }

  function nombreResponsable(task: Task) {
    return (
      task.responsable?.nombre ??
      store.employees.find((e) => e.id === task.responsable_id)?.nombre ??
      `Empleado #${task.responsable_id}`
    );
  }

  async function borrar(task: Task) {
    if (confirm(`¿Eliminar la tarea "${task.titulo}"?`)) await store.eliminarTarea(task.id);
  }
</script>

<div class="mb-4 flex items-center gap-3">
  <label class="text-sm text-slate-600" for="filtro">Filtrar por responsable</label>
  <select
    id="filtro"
    class="campo max-w-xs"
    value={store.filtroResponsable}
    onchange={(e) => store.aplicarFiltro(Number(e.currentTarget.value))}
  >
    <option value={0}>Todos</option>
    {#each store.employees as employee (employee.id)}
      <option value={employee.id}>{employee.nombre}</option>
    {/each}
  </select>
  <span class="ml-auto text-sm text-slate-500 tabular-nums">
    {store.tasks.length} tarea{store.tasks.length === 1 ? '' : 's'}
  </span>
</div>

{#if store.tasks.length === 0}
  <p class="rounded-xl border border-slate-200 bg-white p-8 text-center text-sm text-slate-500">
    No hay tareas que mostrar.
  </p>
{:else}
  <ul class="grid gap-3">
    {#each store.tasks as task (task.id)}
      <li class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:shadow-md">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <h4 class="truncate font-medium text-slate-900">{task.titulo}</h4>
              <span
                class="shrink-0 rounded-full px-2 py-0.5 text-xs font-medium ring-1 ring-inset {ESTADO_CLASS[task.estado]}"
              >
                {ESTADO_LABEL[task.estado] ?? task.estado}
              </span>
            </div>
            <p class="mt-1 text-sm text-slate-600">{task.description}</p>
            <p class="mt-2 text-xs text-slate-500">
              <span class="font-medium text-slate-700">{nombreResponsable(task)}</span>
              <span class="mx-1.5">·</span>
              <span class:text-red-600={vencida(task)} class:font-medium={vencida(task)}>
                Vence el {formatoFecha.format(new Date(task.fecha_limite))}
                {#if vencida(task)}(vencida){/if}
              </span>
            </p>
          </div>

          <button class="boton-peligro shrink-0" onclick={() => borrar(task)} disabled={store.cargando}>
            Eliminar
          </button>
        </div>

        <div class="mt-3 flex flex-wrap gap-4 border-t border-slate-100 pt-3">
          <label class="flex items-center gap-2 text-xs text-slate-500">
            Reasignar a
            <select
              class="campo w-auto py-1 text-xs"
              value={task.responsable_id}
              disabled={store.cargando}
              onchange={(e) => store.reasignarTarea(task.id, Number(e.currentTarget.value))}
            >
              {#each store.employees as employee (employee.id)}
                <option value={employee.id}>{employee.nombre}</option>
              {/each}
            </select>
          </label>

          <label class="flex items-center gap-2 text-xs text-slate-500">
            Estado
            <select
              class="campo w-auto py-1 text-xs"
              value={task.estado}
              disabled={store.cargando}
              onchange={(e) => store.cambiarEstado(task.id, e.currentTarget.value as Estado)}
            >
              {#each ESTADOS as opcion (opcion)}
                <option value={opcion}>{ESTADO_LABEL[opcion]}</option>
              {/each}
            </select>
          </label>
        </div>
      </li>
    {/each}
  </ul>
{/if}
