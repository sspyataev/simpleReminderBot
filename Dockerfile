# Stage 1: Сборка бинарного файла
FROM golang:1.25-alpine AS builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник без CGO (статическая сборка)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mybot .

# Stage 2: Минимальный образ для запуска
FROM alpine:latest

# Установим ca-certificates (для HTTPS-запросов к Telegram API)
RUN apk --no-cache add ca-certificates

# Рабочая директория
WORKDIR /app/

# Копируем бинарник из первого этапа
COPY --from=builder /app/mybot .

# Копируем .env (если нужно, но лучше передавать через переменные при запуске)
# Рекомендуется: передавать BOT_TOKEN через -e, а не в образе

# Запуск бота
CMD ["./mybot"]
