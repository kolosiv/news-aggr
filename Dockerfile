FROM golang:1.22.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache libc6-compat

RUN mkdir -p /app/logs

COPY --from=builder /app/myapp .
COPY --from=builder /app/sources.json .
COPY --from=builder /app/.env . 
COPY --from=builder /app/internal/database/migrations /app/internal/database/migrations
COPY --from=builder /app/static /app/static

EXPOSE 8080

CMD ["./myapp"]
