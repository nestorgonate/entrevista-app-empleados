<script lang="ts">
  import { onMount } from 'svelte';
  import EmployeeForm from './components/EmployeeForm.svelte';
  import EmployeeList from './components/EmployeeList.svelte';
  import TaskForm from './components/TaskForm.svelte';
  import TaskList from './components/TaskList.svelte';
  import Toast from './components/Toast.svelte';
  import { store } from './lib/store.svelte';
  import type { Employee } from './lib/types';

  type Pestana = 'empleados' | 'tareas';

  let pestana = $state<Pestana>('empleados');
  let editando = $state<Employee | null>(null);

  const pestanas: { id: Pestana; texto: string }[] = [
    { id: 'empleados', texto: 'Empleados' },
    { id: 'tareas', texto: 'Tareas' },
  ];

  // Carga inicial, una sola vez al montar.
  onMount(() => {
    store.cargarTodo();
  });
</script>

<div class="min-h-screen bg-slate-50">
  <header class="border-b border-slate-200 bg-white">
    <div class="mx-auto flex max-w-5xl items-center gap-4 px-6 py-4">
      <div class="flex-1">
        <h1 class="text-lg font-semibold text-slate-900">Gestión de empleados y tareas</h1>
        <p class="text-sm text-slate-500">Registra empleados y asígnales trabajo.</p>
      </div>
      {#if store.cargando}
        <span class="text-xs text-slate-400">Cargando…</span>
      {/if}
    </div>

    <nav class="mx-auto flex max-w-5xl gap-1 px-6">
      {#each pestanas as item (item.id)}
        <button
          onclick={() => (pestana = item.id)}
          class="-mb-px border-b-2 px-4 py-2.5 text-sm font-medium transition
                 {pestana === item.id
                   ? 'border-indigo-600 text-indigo-700'
                   : 'border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700'}"
          aria-current={pestana === item.id ? 'page' : undefined}
        >
          {item.texto}
          <span class="ml-1.5 text-xs text-slate-400 tabular-nums">
            {item.id === 'empleados' ? store.employees.length : store.tasks.length}
          </span>
        </button>
      {/each}
    </nav>
  </header>

  <main class="mx-auto max-w-5xl space-y-5 px-6 py-8">
    {#if pestana === 'empleados'}
      <!-- La key remonta el formulario al cambiar de empleado, reiniciando sus campos. -->
      {#key editando?.id}
        <EmployeeForm {editando} onCerrar={() => (editando = null)} />
      {/key}
      <EmployeeList onEditar={(employee) => (editando = employee)} />
    {:else}
      <TaskForm />
      <TaskList />
    {/if}
  </main>
</div>

<Toast />
