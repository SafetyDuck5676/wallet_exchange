# Используем официальный образ Go
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы
ADD ./gw-exchanger .
RUN go mod download

# Сборка исполняемого файла
RUN go build -o main ./cmd

# Указываем команду запуска
CMD ["./main", "-c", "config.env"]