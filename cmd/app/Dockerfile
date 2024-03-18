# Используем образ Ubuntu
FROM golang:latest

# Установка необходимых зависимостей
RUN apt-get update && apt-get install -y \
    build-essential \
    golang-go

# Установка рабочей директории
WORKDIR /app

# Копирование исходного кода приложения
COPY . .

# Сборка приложения
RUN GOOS=linux go build -o app ./cmd/app

# Копирование конфигурационного файла
COPY config/config.yaml /app

# Определение команды для запуска приложения
CMD ["./app"]

