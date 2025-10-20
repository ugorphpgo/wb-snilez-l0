FROM golang:1.25 as build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/wbservice ./cmd/wbservice

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /bin/wbservice /app/wbservice
COPY configs/ /app/configs/
COPY web/ /app/web/
COPY migrations/ /app/migrations/
EXPOSE 8081
USER 65532:65532
ENTRYPOINT ["/app/wbservice"]
