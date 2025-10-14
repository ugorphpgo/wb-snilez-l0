FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
WORKDIR /app/scripts/kafka_producer
RUN go build -o /kafka-producer .

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /kafka-producer .

EXPOSE 8082
CMD ["./kafka-producer"]