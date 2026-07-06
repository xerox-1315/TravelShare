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
- **Docker** — контейнеризация

## Требования

- Go 1.21+
- Docker & Docker Compose
- golang-migrate

## Установка и запуск

### 1. Клонировать репозиторий
\```bash
git clone https://github.com/xerox-1315/TravelShare.git

cd TravelShare
\```

### 2. Создать .env файл
\```
DB_HOST=localhost
DB_PORT=5433
DB_USER=user
DB_PASSWORD=your_password
DB_NAME=travelshare
JWT_SECRET=your_secret_key
\```

### 3. Запустить базу данных
\```bash
docker compose up -d
\```

### 4. Применить миграции
\```bash
migrate -path migrations -database "postgres://user:password@localhost:5433/travelshare?sslmode=disable" up
\```

### 5. Запустить сервер
\```bash
go run cmd/main.go
\```

## API Endpoints

### Авторизация
| Метод | URL | Описание | Авторизация |
|-------|-----|----------|-------------|
| POST | /auth/register | Регистрация | ❌ |
| POST | /auth/login | Вход | ❌ |

### Метки
| Метод | URL | Описание | Авторизация |
|-------|-----|----------|-------------|
| GET | /points | Все метки | ✅ |
| GET | /point/:id | Метка по ID | ✅ |
| GET | /point/nearby?lat=A&lng=B&radius=R | Все метки в радиусе определенной точки | ✅ |
| POST | /point | Создать метку | ✅ |

## Структура проекта

TravelShare/
-> cmd/
    -> main.go
-> internal/
    -> db/          # работа с БД
    -> handler/     # HTTP обработчики
    -> models/      # модели данных
    -> service/     # бизнес-логика
    -> dto/         # объекты передачи данных
    -> migrations/  # SQL миграции
-> errs/            # ошибки приложения
-> docker-compose.yml

## Автор

**Жиликов Данил** — [GitHub](https://github.com/xerox-1315)
