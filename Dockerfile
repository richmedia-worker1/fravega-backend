FROM golang:1.21-alpine

# Устанавливаем необходимые зависимости для CGO
RUN apk add --no-cache gcc musl-dev

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем модули Go
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код приложения
COPY . .



# Собираем приложение с поддержкой CGO
ENV CGO_ENABLED=1
RUN go build -o main ./cmd/fravega-app/main.go


# Указываем команду для запуска приложения
CMD ["/app/main"]