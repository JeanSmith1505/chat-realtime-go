# Imagen base
FROM golang:1.20-alpine AS builder
WORKDIR /app


# Instalar git para go modules
RUN apk add --no-cache git


# Copiar módulos y descargar dependencias
COPY go.mod .
COPY go.sum .
RUN go mod download


# Copiar el resto del código
COPY . .


# Construir binario
RUN go build -o /chat-server ./cmd/server


# Imagen final ejecutable pequeña
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /chat-server /chat-server
COPY web/static /web/static
EXPOSE 8080
ENTRYPOINT ["/chat-server"]