FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/bin/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/bin/migrate ./cmd/migrate

FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/bin/api ./meus_gastos
COPY --from=builder /app/bin/migrate ./migrate

COPY --from=builder /app/migrations/*.sql ./migrations/

RUN chmod +x /root/meus_gastos /root/migrate \
    && apk --no-cache add tzdata

ENV SERVER_PORT=4000 \
    MIGRATE_PATH=/root/migrations

ENTRYPOINT ["/root/meus_gastos"]