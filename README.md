# Gestión de empleados y tareas

API REST para registrar empleados y asignarles tareas, con su propia interfaz web.

---

## Deploy con Docker

```bash
cp .env.example .env      # ajustar credenciales si hace falta
docker compose up --build
```

Eso levanta los tres servicios:

| Servicio | Imagen / build | Puerto | Notas |
|---|---|---|---|
| `web` | `web/Dockerfile.dev` (node:24-alpine + Vite) | 5173 | Hot reload: el código va bind-montado |
| `api` | `Dockerfile` (multi-stage → alpine) | 8000 | Healthcheck en `/health` |
| `db` | `postgis/postgis:18-master` | 5432 | Datos en el volumen `postgresql_data` |

El arranque es secuencial por healthchecks: `db` healthy → `api` healthy → `web`. Las tablas de la base de datos se crean automaticamente con la funcion AutoMigrate de Gorm en platform/migrate.go

Cuando termine se puede acceder a la web con http://localhost:5173

### Comandos útiles

```bash
docker compose ps                      # estado y salud de cada servicio
docker compose logs -f api             # seguir logs del backend
docker compose down                    # parar (conserva los datos)
docker compose down -v                 # parar y BORRAR la base de datos, redes de docker
docker compose up --build -V web       # si se cambia web/package.json
```

El `-V` del último comando recrea el volumen anónimo de `node_modules`. Sino se utiliza, el contenedor sigue usando las dependencias viejas.
---

## Stack

| Tecnología | Por qué |
|---|---|
| **Go + Gin + GORM** | Desarrollo rápido y **build portable**: compila a un único binario estático (`CGO_ENABLED=0`, sin dependencias de librerías C), así que se puede instalar y correr sin Docker en cualquier Linux, ademas de posibilidad de cross compilation para Windows, MacOS y distintas arquitecturas como ARM. |
| **Svelte** (SPA con Vite) | Compila a archivos estáticos, **compatible con Cloudflare Workers**: despliegue sencillo y gratuito dentro de la capa gratuita. |
| **PostgreSQL** | Su ecosistema de extensiones deja la puerta abierta a **vectores** para escalar la app a futuro, por ejemplo, un chatbot con RAG para consultas sobre empleados. |
| **Tailwind CSS v4** | Estilos consistentes. |

> Sobre el chatbot: la imagen actual es `postgis/postgis`, que trae PostGIS pero no pgvector. Para usarlo habría que cambiar a una imagen que lo incluya, por ejemplo `pgvector/pgvector` y ejecutar `CREATE EXTENSION vector`.

---

## Arquitectura

**Hexagonal (puertos y adaptadores).** Facilita el escalado y permite añadir nuevos casos de uso sobre las mismas interfaces sin tocar el resto del código.

```
core/
├── domain/          Entidades, DTOs y errores. No depende de nada.
│   ├── entity/        Employee, Task (structs de GORM) y los estados válidos
│   ├── dto/           Requests con validación de Gin + mappers a response
│   └── errors.go      ErrNotFound / ErrConflict / ErrInvalidInput
├── port/            Las interfaces. El contrato entre capas.
│   ├── service.go     Puertos de entrada  (EmployeeService, TaskService)
│   └── repository.go  Puertos de salida   (EmployeeRepository, TaskRepository)
├── service/         Lógica de negocio. Implementa las interfaces definidas en port.
└── adapter/
    ├── handler/       Adaptador de entrada: HTTP con Gin
    └── repository/    Adaptador de salida: persistencia con GORM

platform/           Infraestructura: conexión a Postgres y migraciones
main.go             Inyecta dependencias
```

Siguiendo arquitectura hexagonal, los adaptadores importan los servicios
```
handler → port.TaskService → service → port.TaskRepository → repository → GORM
```

`service/task.go` espera la interfaz `port.TaskRepository` , no `*gorm.DB` que es el tipo de dato especifico de la conexion utilizando Gorm. Como las dependencias se inyectan en `main.go`, un cambio de base de datos o libreria no causa cambios en el codigo, de esta manera la escalabilidad hacia otras librerias de Go como pgx para un mejor rendimiento es sencillo:

```go
taskRepo := repository.NewTaskRepository(db)          // adaptador de salida
taskSvc  := service.NewTaskService(taskRepo, employeeRepo)  // caso de uso
taskH    := handler.NewTaskHandler(taskSvc)           // adaptador de entrada
```

---

## API

Todo cuelga de `/api/v1`, excepto health para verificar el estado del backend.

| Método | Ruta | Qué hace |
|---|---|---|
| `GET` | `/health` | Healthcheck |
| `POST` | `/api/v1/employees` | Registrar empleado |
| `GET` | `/api/v1/employees` | Listar empleados |
| `GET` | `/api/v1/employees/:id` | Obtener un empleado |
| `PUT` | `/api/v1/employees/:id` | Actualizar |
| `DELETE` | `/api/v1/employees/:id` | Eliminar |
| `GET` | `/api/v1/employees/:id/tasks` | Tareas de ese empleado |
| `POST` | `/api/v1/tasks` | Crear tarea y asignarla |
| `GET` | `/api/v1/tasks` | Listar, acepta el query string `?responsable_id=N` para busqueda especifica|
| `GET` | `/api/v1/tasks/:id` | Obtener una tarea especifica |
| `PUT` | `/api/v1/tasks/:id` | Actualizar una tarea especifica|
| `PATCH` | `/api/v1/tasks/:id/responsable` | Reasignar la tarea a otro empleado |
| `PATCH` | `/api/v1/tasks/:id/estado` | Cambiar el estado de una tarea especifica |
| `DELETE` | `/api/v1/tasks/:id` | Eliminar una tarea especifica |

Dos detalles a tener en cuenta al probar con Postman o Curl:

- **fecha_limite va en RFC3339** (2026-09-30T00:00:00Z). El formato 2026-09-30 falla la validación de la libreria time de Go
- **Estados válidos**: `pendiente`, `en_progreso`, `completada`, `cancelada` (ver `core/domain/entity/estado.go`).

Los errores siempre llegan como `{"error": "..."}` con 400, 404, 409 o 500. Los `DELETE` responden 204 sin cuerpo.

```bash
# Registrar un empleado
curl -X POST localhost:8000/api/v1/employees \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Nestor Gavilanes","correo":"negavilaneson@outlook.com","cargo":"Desarrollador"}'

# Crear una tarea y asignársela
curl -X POST localhost:8000/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -d '{"titulo":"Desarrollo de una app para registros de empleados","description":"Desarrollar una app con documentacion para registrar empleados y asignar tareas","responsable_id":1,"fecha_limite":"2026-10-01T12:00:00Z"}'
```

---

## Sin Docker

Debido a Go, el backend se puede compilar a un binario estático que no necesita nada más que Postgres accesible en el servidor o en un servicio cloud como Supabase:

```bash
CGO_ENABLED=0 go build -o server .
DATABASE_URL='postgresql://postgres:postgres@localhost:5432/entrevista?sslmode=disable' ./server
```

Frontend:

```bash
cd web
npm install
npm run dev     # http://localhost:5173
```

El dev server de Vite hace de proxy de `/api` hacia `localhost:8000`, así que el navegador nunca cruza de origen y el backend no necesita CORS. Para desplegarlo, `npm run build` genera `web/dist` con archivos estáticos listos para Cloudflare Pages o cualquier CDN.

---

## Mejoras
- Agregar JWT para manejo de sesiones y proteger endpoints
- Mejorar el esquema de la base de datos para permitir un manejo basado en roles y llevar un registro en cambios de tareas