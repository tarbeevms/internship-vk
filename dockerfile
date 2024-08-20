# Используем официальный образ Golang как базовый образ
FROM golang:latest

# Устанавливаем зависимости для работы с OpenSSL
RUN apt-get update && \
    apt-get install -y \
    libssl-dev \
    pkg-config

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /intern

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости перед копированием остальных файлов
RUN go mod download

# Копируем остальные файлы вашего приложения внутрь контейнера
COPY . .

# Изменяем права доступа к файлу конфигурации
RUN chmod 777 /intern/config/config.yml

# Переходим в директорию cmd/redditclone
WORKDIR /intern/cmd/internship

# Собираем ваше приложение
RUN go build -o myapp .

# Указываем порт, который будет использоваться приложением
EXPOSE 8080

# Команда для запуска вашего приложения при старте контейнера
CMD ["./myapp"]
