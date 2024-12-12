# Используем официальный образ Go
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Сборка исполняемого файла
RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd

# Финальный образ
FROM alpine:3.16

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем исполняемый файл
COPY --from=builder /app/main .

# Копируем файл конфигурации
COPY config.env .

# Указываем команду запуска
CMD ["./main", "-c", "config.env"]