# --- Etapa de compilacion ---
FROM golang:1.27.1-alpine AS builder

WORKDIR /src

# Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Se copian unicamente los paquetes Go, asi un cambio en web/ no invalida
# el build del backend.
COPY main.go ./
COPY core/ ./core/
COPY platform/ ./platform/

# Build del binario del backend, CGO_ENABLED=0 permite que el binario no dependa de librerias C del sistema operativo
# -s -w descarta la tabla de simbolos y la informacion de depuracion.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/server /app/server

EXPOSE 8000

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
  CMD wget -qO- http://127.0.0.1:8000/health || exit 1

ENTRYPOINT ["/app/server"]
