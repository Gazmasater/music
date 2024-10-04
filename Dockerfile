FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN swag init

# Собираем приложение
RUN go build -o main .

# Используем Alpine для выполнения приложения
FROM alpine:latest

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем только исполняемый файл и необходимые директории
COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs

# Устанавливаем необходимые зависимости
RUN apk add --no-cache libc6-compat ca-certificates


CMD ["./main"]