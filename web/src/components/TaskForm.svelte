<script lang="ts">
  import { store } from '../lib/store.svelte';
  import { ESTADO_LABEL, ESTADOS, type Estado } from '../lib/types';

  let titulo = $state('');
  let description = $state('');
  let responsableId = $state(0);
  let fecha = $state('');
  let estado = $state<Estado>('pendiente');

  const sinEmpleados = $derived(store.employees.length === 0);

  async function enviar(event: SubmitEvent) {
    event.preventDefault();

    const ok = await store.crearTarea({
      titulo,
      description,
      responsable_id: responsableId,
      // Gin espera RFC3339; <input type="date"> entrega "2026-09-30".
      fecha_limite: new Date(`${fecha}T00:00:00`).toISOString(),
      estado,
    });

    if (!ok) return;
    titulo = description = fecha = '';
    responsableId = 0;
    estado = 'pendiente';
  }
</script>

<form onsubmit={enviar} class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
  <h3 class="mb-4 text-sm font-semibold text-slate-900">Nueva tarea</h3>

  {#if sinEmpleados}
    <!-- El backend rechaza con 400 un responsable_id inexistente. -->
    <p class="mb-4 rounded-lg bg-amber-50 px-3 py-2 text-sm text-amber-800 ring-1 ring-amber-200">
      Registra al menos un empleado antes de crear tareas.
    </p>
  {/if}

  <fieldset disabled={sinEmpleados} class="grid gap-4">
    <div class="grid gap-4 sm:grid-cols-2">
      <div>
        <label class="etiqueta" for="task-titulo">Título</label>
        <input
          id="task-titulo"
          class="campo"
          bind:value={titulo}
          maxlength="255"
          required
          placeholder="Preparar informe trimestral"
        />
      </div>
      <div>
        <label class="etiqueta" for="task-responsable">Responsable</label>
        <select id="task-responsable" class="campo" bind:value={responsableId} required>
          <option value={0} disabled>Selecciona un empleado…</option>
          {#each store.employees as employee (employee.id)}
            <option value={employee.id}>{employee.nombre} — {employee.cargo}</option>
          {/each}
        </select>
      </div>
    </div>

    <div>
      <label class="etiqueta" for="task-desc">Descripción</label>
      <textarea
        id="task-desc"
        class="campo resize-y"
        rows="2"
        bind:value={description}
        required
        placeholder="Detalle de lo que hay que hacer…"
      ></textarea>
    </div>

    <div class="grid gap-4 sm:grid-cols-2">
      <div>
        <label class="etiqueta" for="task-fecha">Fecha límite</label>
        <input id="task-fecha" class="campo" type="date" bind:value={fecha} required />
      </div>
      <div>
        <label class="etiqueta" for="task-estado">Estado inicial</label>
        <select id="task-estado" class="campo" bind:value={estado}>
          {#each ESTADOS as opcion (opcion)}
            <option value={opcion}>{ESTADO_LABEL[opcion]}</option>
          {/each}
        </select>
      </div>
    </div>

    <div>
      <button type="submit" class="boton-primario" disabled={store.cargando || responsableId === 0}>
        Crear y asignar
      </button>
    </div>
  </fieldset>
</form>
