FROM golang:1.23

# Создаём рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum из родительской директории
COPY ../go.mod ../go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем все остальные файлы проекта
COPY . .

# Сборка бинарника
RUN go build -o avito-shop-service ./cmd/shop-service \
    && go clean -cache -modcache

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["/app/avito-shop-service"]
