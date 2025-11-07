# Scheduler - планировщик задач

Это простое приложение для управления задачами на Go. Можно добавлять задачи, ставить им даты, настраивать повторения и искать нужные записи.

## Что умеет

- Создавать, редактировать и удалять задачи
- Искать задачи по тексту
- Настраивать повторяющиеся задачи (каждый день, неделю, месяц и т.д.)
- Защита паролем (JWT токены)
- Работает в Docker

## Что нужно для запуска

- Go версии 1.24+
- Или Docker если хотите запустить в контейнере

## Как запустить

### Обычный запуск

```bash
# Скачать зависимости
go mod download

# Собрать
go build -o scheduler ./cmd/scheduler

# Запустить
./scheduler
```

Приложение откроется на http://localhost:7540

### С паролем

```bash
TODO_PASSWORD=mypassword ./scheduler
```

### Другой порт

```bash
TODO_PORT=8080 ./scheduler
```

## Запуск через Docker

Сначала собрать образ:

```bash
./scripts/build-docker
```

Потом запустить:

```bash
docker run -p 7540:7540 -v $(pwd)/data:/data scheduler:latest
```

С паролем:

```bash
docker run -p 7540:7540 \
  -e TODO_PASSWORD=mypassword \
  -e JWT_SECRET_KEY=your-secret-key \
  -v $(pwd)/data:/data \
  scheduler:latest
```

## API

### Без аутентификации

**GET /api/nextdate** - рассчитать следующую дату для повторяющейся задачи

Параметры:
- `now` - текущая дата (формат: YYYYMMDD)
- `date` - начальная дата задачи
- `repeat` - правило повторения

Пример:
```bash
curl "http://localhost:7540/api/nextdate?now=20251110&date=20251110&repeat=d%201"
```

**POST /api/signin** - войти и получить токен

```bash
curl -X POST http://localhost:7540/api/signin \
  -H "Content-Type: application/json" \
  -d '{"password": "mypassword"}'
```

### С аутентификацией (нужен токен)

**GET /api/tasks** - получить список задач

Можно добавить параметр `search` для поиска.

**GET /api/task?id=123** - получить одну задачу

**POST /api/task** - создать задачу

```json
{
  "title": "Купить продукты",
  "date": "20251120",
  "comment": "Молоко, хлеб, яйца",
  "repeat": "d 3"
}
```

**PUT /api/task** - обновить задачу

**DELETE /api/task?id=123** - удалить задачу

**POST /api/task/done** - отметить задачу выполненной

```json
{
  "id": "123"
}
```

## Правила повторения

Формат: `тип интервал [параметры]`

- `y` - каждый год
- `d 7` - каждые 7 дней
- `w 1,3,5` - каждую неделю в понедельник, среду, пятницу (1-7)
- `m 15 1,6,12` - 15 числа в январе, июне и декабре
- `m -1` - последний день месяца

## Тесты

```bash
go test ./tests/... -v
```

## Структура проекта

```
cmd/scheduler/main.go - точка входа
internal/
  auth/ - аутентификация (JWT)
  config/ - настройки
  database/ - работа с SQLite
  handlers/ - HTTP обработчики
  models/ - модели данных
  services/ - бизнес логика
  utils/ - вспомогательные функции
web/ - фронтенд (HTML, CSS, JS)
tests/ - тесты
```
