<script lang="ts">
  import { fly } from 'svelte/transition';
  import { store } from '../lib/store.svelte';

  // Los avisos de exito se ocultan solos; los errores esperan al usuario.
  $effect(() => {
    const aviso = store.aviso;
    if (aviso?.tipo !== 'exito') return;

    const id = setTimeout(() => store.limpiarAviso(), 2500);
    return () => clearTimeout(id);
  });
</script>

{#if store.aviso}
  {@const esError = store.aviso.tipo === 'error'}
  <div
    role="status"
    transition:fly={{ y: 12, duration: 200 }}
    class="fixed bottom-5 left-1/2 z-50 flex max-w-md -translate-x-1/2 items-start gap-3 rounded-xl
           px-4 py-3 text-sm shadow-lg ring-1
           {esError ? 'bg-red-50 text-red-800 ring-red-200' : 'bg-emerald-50 text-emerald-800 ring-emerald-200'}"
  >
    <span aria-hidden="true">{esError ? '⚠️' : '✓'}</span>
    <p class="flex-1">{store.aviso.texto}</p>
    <button
      onclick={() => store.limpiarAviso()}
      class="shrink-0 opacity-50 transition hover:opacity-100"
      aria-label="Cerrar aviso"
    >
      ✕
    </button>
  </div>
{/if}
