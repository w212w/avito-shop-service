# Инструкция по запуску приложения

### 1. Клонирование репозитория

Клонирование репозитория:

```bash
git clone https://github.com/w212w/avito-shop-service.git
cd avito-shop-service
```

### 2. Настройка конфигурации
Корневая директория - avito-shop-service. 
- Когда база данных и сервис развертывается с помощью Docker Compose, .env файл НУЖНО УДАЛИТЬ.
- Если база данных развертывается с помощью Docker Compose, а сервис запускается локально, ОБЯЗАТЕЛЬНО СОЗДАЙТЕ .env в корневой директории проекта и добавьте следующие переменные:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shop
JWT_SECRET=supersecretkey
```
### 3. Запуск базы данных с Docker Compose
Для развертывания базы данных используйте Docker Compose:

Убедитесь, что у вас установлен Docker и Docker Compose.
В корневой директории проекта найдите файл docker-compose.yml, который описывает сервис для базы данных.

```bash
docker-compose up -d db
```

### 4. Запуск сервиса локально
После того как база данных развернута, можно запустить сервис локально:

Убедитесь, что в корне проекта создан файл .env, как указано в разделе выше.
```bash
go run main.go
```

### 5. Применение миграций
При запуске приложения миграции для создания таблиц в базе данных будут применены автоматически. В БД создаются следующие таблицы:
```bash
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    coins INT NOT NULL DEFAULT 1000,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_user_id INT REFERENCES users(id) ON DELETE CASCADE,
    to_user_id INT REFERENCES users(id) ON DELETE CASCADE,
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    item TEXT NOT NULL,
    price INT NOT NULL CHECK (price > 0),
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS shop (
    item TEXT PRIMARY KEY,
    price INT NOT NULL CHECK (price > 0)
);

INSERT INTO shop (item, price) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500)
ON CONFLICT (item) DO NOTHING;

```
### 6. Запуск тестов
```bash
go test ./...
```
Тесты расположены в следующих директориях:

- avito-shop-service/internal/service
- avito-shop-service/internal/handlers

  ### 6. Запуск линтера
```bash
golangci-lint run
```
