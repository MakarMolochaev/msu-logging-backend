# Используем многоэтапную сборку для уменьшения размера финального образа
# Этап сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git make
RUN apk add --no-cache curl
# Рабочая директория
WORKDIR /app

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем основной сервер
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/msu-logging-backend

# Собираем мигратор
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/migrator ./cmd/migrator

# Этап запуска
FROM alpine:3.21

# Устанавливаем зависимости для runtime (если нужны)
RUN apk add --no-cache ca-certificates

# Рабочая директория
WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrator /app/migrator
COPY migrations /app/migrations
COPY certs /app/certs
COPY config /app/config
COPY .env.prod .env.local
# Копируем конфигурационные файлы (если есть)

# Открываем порты, которые использует сервис
EXPOSE 50051 8081 8082

# Команда запуска сервиса
# CMD ["/app/server"]