# Этап 1: Сборка бинарного файла
FROM golang:1.22-alpine AS builder

# Устанавливаем необходимые инструменты
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/storage .

# Этап 2: Финальный образ
FROM alpine:latest

# Устанавливаем CA-сертификаты и часовые пояса
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/storage .

# Открываем порт gRPC (стандартный 50051)
EXPOSE 50051

# Запускаем сервис
CMD ["./storage"]