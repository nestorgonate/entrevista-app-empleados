import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [svelte(), tailwindcss()],
  server: {
    port: 5173,
    // El proxy hace que el navegador nunca cruce de origen: la web pide
    // /api/v1/... a su propio host y Vite lo reenvia al backend Gin.
    // Por eso no hace falta middleware de CORS en Go durante el desarrollo.
    proxy: {
      '/api': {
        // En el host el backend esta en localhost:8000; dentro de Docker es el
        // servicio `api`, y docker-compose inyecta API_PROXY_TARGET.
        // Este archivo se ejecuta en Node, asi que no necesita el prefijo VITE_.
        target: process.env.API_PROXY_TARGET ?? 'http://localhost:8000',
        changeOrigin: true,
      },
    },
  },
});
