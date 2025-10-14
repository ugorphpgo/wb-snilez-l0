#Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/wb-snilez-l0-main

RUN go mod download
RUN go build -o /wbl0 ./cmd/server


# Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /wbl0 .
COPY web ./web

EXPOSE 8081
CMD ["/app/wbl0"]
