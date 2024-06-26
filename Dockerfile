FROM golang:1.21-alpine AS builder
LABEL authors="Nite07"
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod download
RUN apk add --no-cache nodejs npm
WORKDIR /app/internal/server/frontend
RUN npm install
RUN npm run build
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o gamedb .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gamedb /app/gamedb
ENTRYPOINT ["/app/gamedb", "server"]