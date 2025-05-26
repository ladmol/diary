# Diary App

Приложение дневника с бэкендом на Go, базой данных PostgreSQL и аутентификацией SuperTokens.

## Запуск с Docker Compose

### Предварительные требования
- Docker
- Docker Compose

### Шаги запуска

1. Клонировать репозиторий
```
git clone <your-repo-url>
cd diary
```

2. Создать файл .env или использовать значения по умолчанию
```
# Пример .env файла
DATABASE_URL=postgres://postgres:postgres@postgres:5432/diary?sslmode=disable
PORT=8080
SUPERTOKENS_CONNECTION_URI=http://supertokens:3567
API_DOMAIN=http://localhost:8080
WEBSITE_DOMAIN=http://localhost:3000
```

3. Запустить приложение с помощью Docker Compose
```
docker-compose up -d
```

4. Проверить работу приложения
```
curl http://localhost:8080/ping
```

5. Остановить приложение
```
docker-compose down
```

## Разработка

### Структура проекта
- `cmd/main.go` - точка входа приложения
- `internal/models` - модели данных
- `internal/repos` - репозитории для работы с базой данных
- `internal/services` - бизнес-логика
- `internal/handlers` - обработчики HTTP-запросов
- `internal/router` - конфигурация маршрутов

### Работа с базой данных
GORM автоматически создает схему базы данных при запуске приложения.

### API Endpoints
- `GET /ping` - проверка работы сервера
- `POST /api/entries` - создание новой записи
- `GET /api/entries/{id}` - получение записи по ID
- `PUT /api/entries/{id}` - обновление записи
- `DELETE /api/entries/{id}` - удаление записи
- `GET /api/entries` - получение всех записей пользователя
