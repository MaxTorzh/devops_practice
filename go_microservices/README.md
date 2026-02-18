# Задание 6: Многоконтейнерное приложение (Go + PostgreSQL + Redis)

Микросервисное приложение на Go с PostgreSQL для хранения данных и Redis для кэширования. Все компоненты запускаются в отдельных Docker-контейнерах и взаимодействуют через общую сеть.

## Архитектура

┌─────────────────┐ ┌──────────────┐ ┌─────────────────┐
│ Go App │────▶ │ PostgreSQL │ │ Users Table │
│ (порт 8080) │ │ (порт 5432) │ │ Products Table │
└────────┬────────┘ └──────────────┘ └─────────────────┘
│ ▲
│ │
▼ │
┌─────────────────┐ ┌────────┴────────┐
│ Redis │────▶│ Кэширование │
│ (порт 6379) │ │ GET /users/1 │
└─────────────────┘ │ TTL: 5 минут │

## Быстрый старт

```bash
# Запустить все сервисы
docker-compose up -d

# Проверить статус
docker-compose ps

# Посмотреть логи
docker-compose logs -f
После запуска API будет доступно на http://localhost:8080
```

Взаимодействие компонентов

PostgreSQL — основное хранилище данных

При первом запуске автоматически создаёт таблицы через миграции

Данные сохраняются в volume и не теряются после перезапуска

Тестовые данные добавляются автоматически

Redis — кэш для ускорения ответов

Хранит часто запрашиваемые данные в оперативной памяти

TTL (время жизни) кэша: 5 минут

Автоматически удаляет устаревшие записи

Go приложение — бизнес-логика

Обрабатывает HTTP запросы

При чтении: сначала Redis, потом PostgreSQL

При записи: PostgreSQL + инвалидация кэша

Схема взаимодействия при GET /users/1

```go
// 1. Проверка кэша
cacheKey := "user:1"
data, found := redis.Get(cacheKey)

if found {
    // 2. HIT - данные из Redis (быстро)
    return data
}

// 3. MISS - идём в PostgreSQL
user := postgres.Query("SELECT * FROM users WHERE id = 1")

// 4. Сохраняем в Redis на 5 минут
redis.Set(cacheKey, user, 5*time.Minute)

// 5. Возвращаем пользователю
return user
```

API Endpoints
Пользователи

```bash
# Создать пользователя
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Иван Петров","email":"ivan@example.com"}'

# Получить всех пользователей
curl http://localhost:8080/users

# Получить пользователя по ID (с кэшированием)
curl http://localhost:8080/users/1

# Обновить пользователя
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Новое имя"}'

# Удалить пользователя
curl -X DELETE http://localhost:8080/users/1
```

Товары

```bash
# Создать товар
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Ноутбук","price":999.99,"stock":10}'

# Получить все товары
curl http://localhost:8080/products

# Получить товар по ID
curl http://localhost:8080/products/1

# Обновить количество на складе
curl -X PATCH http://localhost:8080/products/1/stock \
  -H "Content-Type: application/json" \
  -d '{"quantity":2}'
```

Тестирование работы
Проверка PostgreSQL

```bash
# Подключиться к БД
docker exec -it microservices-postgres psql -U appuser -d appdb

# В psql выполнить:
\dt                    # список таблиц
SELECT * FROM users;   # данные пользователей
SELECT * FROM products;# данные товаров
\q                     # выход
```

Проверка Redis

```bash
# Подключиться к Redis
docker exec -it microservices-redis redis-cli
```

# В redis-cli выполнить:

KEYS \* # все ключи в кэше
GET user:1 # получить данные пользователя
TTL user:1 # время жизни ключа
exit
Проверка кэширования

```bash
# Первый запрос - MISS (из PostgreSQL)
curl -v http://localhost:8080/users/1 2>&1 | grep -i "x-cache" || echo "No cache header"

# Второй запрос - HIT (из Redis)
curl -v http://localhost:8080/users/1 2>&1 | grep -i "x-cache" || echo "No cache header"

# После обновления - кэш сбрасывается
curl -X PUT http://localhost:8080/users/1 -d '{"name":"Новое имя"}'
curl -v http://localhost:8080/users/1 2>&1 | grep -i "x-cache"
```

Структура базы данных
Таблица users

id SERIAL Уникальный идентификатор
name VARCHAR(100) Имя пользователя
email VARCHAR(100) Email (уникальный)
created_at TIMESTAMP Дата создания
updated_at TIMESTAMP Дата обновления

Таблица products

id SERIAL Уникальный идентификатор
name VARCHAR(200) Название товара
description TEXT Описание
price DECIMAL(10,2) Цена
stock INTEGER Остаток на складе
created_at TIMESTAMP Дата создания
updated_at TIMESTAMP Дата обновления

Управление сервисами

```bash
# Остановить все контейнеры
docker-compose down

# Остановить и удалить все данные (volumes)
docker-compose down -v

# Перезапустить конкретный сервис
docker-compose restart app

# Посмотреть логи конкретного сервиса
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f redis
```

Конфигурация
Все настройки задаются через переменные окружения в файле .env:

env

# PostgreSQL

DB_USER=appuser
DB_PASSWORD=apppass123
DB_NAME=appdb

# Redis

REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL=5m

# Приложение

APP_ENV=development
PORT=8080

## Ключевые концепции

PostgreSQL: надёжное хранение, транзакции, целостность данных

Redis: скорость, снижение нагрузки на БД, идеально для часто запрашиваемых данных

## Кэширование

При чтении: сначала Redis → если нет → PostgreSQL → запись в Redis

При записи: PostgreSQL → удаление из Redis (инвалидация кэша)

TTL 5 минут автоматически удаляет устаревшие данные
