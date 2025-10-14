# Build stage
FROM golang:1.25-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./
# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код (в пределах контекста сборки)
COPY . .

# Собираем приложение УКАЖИТЕ ПРАВИЛЬНЫЙ ПУТЬ К ПАКЕТУ
# Например, если main.go находится в ./cmd/producer:
RUN go build -mod=vendor -o /kafka-producer ./scripts/kafka_producer/producer.go



# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /kafka-producer .

EXPOSE 8082
CMD ["./kafka-producer"]