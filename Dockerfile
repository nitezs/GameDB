FROM golang:1.21-alpine AS builder
LABEL authors="Nite07"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gamedb .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gamedb /app/gamedb
ENTRYPOINT ["/app/gamedb", "server"]