FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o diary ./cmd/main.go

# Многоэтапная сборка для уменьшения размера образа
FROM alpine:latest

WORKDIR /app

# Установка необходимых пакетов
RUN apk --no-cache add ca-certificates tzdata postgresql-client

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/diary .
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

# Делаем скрипт исполняемым
RUN chmod +x /app/docker-entrypoint.sh

# Порт по умолчанию
EXPOSE 8080

# Использование скрипта запуска
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["./diary"]
