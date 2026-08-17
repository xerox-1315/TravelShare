# TravelShare 🗺️

Бэкенд-сервис для обмена метками и советами для путешественников.
Пользователи могут оставлять метки на карте с описанием интересных мест,
лайфхаки и советы для других путешественников.

## Стек технологий

- **Go** — основной язык
- **Gin** — HTTP фреймворк
- **PostgreSQL + PostGIS** — база данных с поддержкой геоданных
- **GORM** — ORM для работы с БД
- **JWT** — аутентификация
- **Redis** — хранение данных регистрации до верификации email
- **gRPC** — связь с Email Service для верификации почты
- **Docker** — контейнеризация

## Требования

- Go 1.21+
- Docker & Docker Compose
- golang-migrate
- Запущенный [TravelShare-Email](https://github.com/xerox-1315/TravelShare-Email) сервис

## Установка и запуск

### 1. Клонировать репозиторий
```bash
git clone https://github.com/xerox-1315/TravelShare.git
cd TravelShare
```

### 2. Создать .env файл
```
DB_HOST=localhost
DB_PORT=5433
DB_USER=user
DB_PASSWORD=your_password
DB_NAME=travelshare
JWT_SECRET=your_secret_key

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

EMAIL_SERVICE_ADDR=localhost:50051
```

### 3. Запустить базу данных и Redis
```bash
docker compose up -d
```

### 4. Применить миграции
```bash
migrate -path migrations -database "postgres://user:password@localhost:5433/travelshare?sslmode=disable" up
```

### 5. Запустить Email Service
Перед запуском основного сервиса необходимо запустить [TravelShare-Email](https://github.com/xerox-1315/TravelShare-Email).

### 6. Запустить сервер
```bash
go run cmd/main.go
```

## API Endpoints

### Авторизация
| Метод | URL | Описание | Авторизация |
|-------|-----|----------|-------------|
| POST | /auth/register | Регистрация (отправка кода на email) | ❌ |
| POST | /auth/verify | Верификация email и создание аккаунта | ❌ |
| POST | /auth/login | Вход | ❌ |

### Метки
| Метод | URL | Описание | Авторизация |
|-------|-----|----------|-------------|
| GET | /points | Все метки | ✅ |
| GET | /point/:id | Метка по ID | ✅ |
| GET | /points/nearby?lat=A&lng=B&radius=R | Метки в радиусе от точки | ✅ |
| POST | /point | Создать метку | ✅ |
| POST | /vote | Поставить голос (like/dislike) метке | ✅ |
| PATCH | /vote | Изменить голос метке | ✅ |

### Пользователь
| Метод | URL | Описание | Авторизация |
|-------|-----|----------|-------------|
| GET | /profile | Получить профиль | ✅ |

## Архитектура

Сервис разделён на два независимых приложения:

```
TravelShare (основной сервис :8080)
        ↕ gRPC (:50051)
TravelShare-Email (Email микросервис)
        ↕
      Redis (коды верификации)
        +
      SMTP (отправка писем)
```

### Процесс регистрации
```
1. POST /auth/register — данные сохраняются в Redis, код отправляется на email
2. POST /auth/verify  — код проверяется, аккаунт создаётся в PostgreSQL
```

## Структура проекта

```
TravelShare/
├── cmd/
│   └── main.go
├── internal/
│   ├── db/          # работа с БД
│   ├── handler/     # HTTP обработчики
│   ├── models/      # модели данных
│   ├── service/     # бизнес-логика
│   ├── dto/         # объекты передачи данных
│   ├── cache/       # Redis клиент
│   └── grpc/        # gRPC клиент для Email Service
├── migrations/      # SQL миграции
├── errs/            # ошибки приложения
└── docker-compose.yml
```

## Автор

**Жиликов Данил** — [GitHub](https://github.com/xerox-1315) 