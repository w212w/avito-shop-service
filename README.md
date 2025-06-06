# Инструкция по запуску приложения

### 1. Клонирование репозитория

Клонирование репозитория:

```bash
git clone https://github.com/w212w/avito-shop-service.git
cd avito-shop-service
```

### 2. Настройка конфигурации
- Корневая директория - avito-shop-service.
- docker-compose.yml и Dockerfile расположен в avito-shop-service/deployments
- База данных развертывается с помощью Docker Compose, создавать отдельный .env файл в корневой директории не обязательно. Будут использоваться параметры из config/config.go либо docker-compose.yml. Пример используемых параметров приведен ниже (локальное развертывание приложения и БД через docker-compose)


```bash
# .env файл (локальное развертывание)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shop
JWT_SECRET=supersecretkey
```
```bash
# (Развертывание через docker, параметры из docker-compose)
APP_ENV=docker
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shop
JWT_SECRET=supersecretkey
```
### 3. Запуск базы данных с Docker Compose
Для развертывания базы данных используйте Docker Compose:
```bash
docker-compose up -d db
```
### 4. Применение миграций
При запуске приложения миграции для создания таблиц в базе данных будут применены автоматически. 
- Миграции расположены в avito-shop-service/internal/db/migrations
- В БД создаются следующие таблицы:
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
### 5. Запуск тестов
```bash
go test ./...
go test -cover  ./...
```
**Результаты тестов:**
- ok      avito-shop-service/internal/handlers    1.387s  coverage: 46.5% of statements
- ok      avito-shop-service/internal/service     (cached) coverage: 61.9% of statements<br>

Тесты расположены в следующих директориях:
- avito-shop-service/internal/service
- avito-shop-service/internal/handlers


### 6. Запуск линтера
```bash
golangci-lint run
```
- Описание конфигурации линтера (.golangci.yaml) расположено в корневой директории (avito-shop-service).
