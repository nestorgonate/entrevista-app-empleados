<script lang="ts">
  import { store } from '../lib/store.svelte';
  import type { Employee } from '../lib/types';

  interface Props {
    /** Si viene, el formulario edita en vez de crear. */
    editando?: Employee | null;
    onCerrar?: () => void;
  }

  let { editando = null, onCerrar }: Props = $props();

  // Valores iniciales tomados del empleado en edicion. Capturar solo el valor
  // inicial es deliberado: el padre envuelve este componente en un {#key}, asi
  // que al cambiar de empleado se vuelve a montar y estos $state se
  // reinicializan solos, sin necesidad de un $effect que los reasigne.
  /* svelte-ignore state_referenced_locally */
  let nombre = $state(editando?.nombre ?? '');
  /* svelte-ignore state_referenced_locally */
  let correo = $state(editando?.correo ?? '');
  /* svelte-ignore state_referenced_locally */
  let cargo = $state(editando?.cargo ?? '');

  async function enviar(event: SubmitEvent) {
    event.preventDefault();
    const datos = { nombre, correo, cargo };

    const ok = editando
      ? await store.actualizarEmpleado(editando.id, datos)
      : await store.crearEmpleado(datos);

    if (!ok) return;
    nombre = correo = cargo = '';
    onCerrar?.();
  }
</script>

<form onsubmit={enviar} class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
  <h3 class="mb-4 text-sm font-semibold text-slate-900">
    {editando ? `Editando a ${editando.nombre}` : 'Registrar empleado'}
  </h3>

  <div class="grid gap-4 sm:grid-cols-3">
    <div>
      <label class="etiqueta" for="emp-nombre">Nombre</label>
      <!-- maxlength 25 replica el binding "max=25" del DTO de Go -->
      <input
        id="emp-nombre"
        class="campo"
        bind:value={nombre}
        maxlength="25"
        required
        placeholder="Ana Pérez"
      />
    </div>
    <div>
      <label class="etiqueta" for="emp-correo">Correo</label>
      <input
        id="emp-correo"
        class="campo"
        type="email"
        bind:value={correo}
        maxlength="255"
        required
        placeholder="ana@empresa.com"
      />
    </div>
    <div>
      <label class="etiqueta" for="emp-cargo">Cargo</label>
      <input id="emp-cargo" class="campo" bind:value={cargo} required placeholder="Desarrolladora" />
    </div>
  </div>

  <div class="mt-4 flex gap-2">
    <button type="submit" class="boton-primario" disabled={store.cargando}>
      {editando ? 'Guardar cambios' : 'Registrar'}
    </button>
    {#if editando}
      <button type="button" class="boton-secundario" onclick={() => onCerrar?.()}>Cancelar</button>
    {/if}
  </div>
</form>
