# Go Final Project - Планировщик задач

## Описание проекта

Это веб-приложение для управления задачами с поддержкой планирования повторяющихся событий. Проект представляет собой HTTP-сервер на Go с использованием SQLite для хранения данных и современным веб-интерфейсом для взаимодействия с пользователем.

### Основные возможности:

- **Управление задачами**: создание, редактирование, удаление и просмотр задач
- **Планирование повторяющихся задач**: поддержка правил повторения (ежедневно, еженедельно, ежемесячно, ежегодно)
- **Расчёт дат**: автоматический расчёт следующей даты выполнения для повторяющихся задач
- **Поиск задач**: полнотекстовый поиск по заголовкам и комментариям
- **Аутентификация**: защита паролем с использованием JWT-токенов
- **Docker поддержка**: возможность запуска приложения в контейнере

## Технологический стек

- **Backend**: Go 1.24
- **База данных**: SQLite (modernc.org/sqlite)
- **Аутентификация**: JWT-токены (github.com/golang-jwt/jwt/v5)
- **Frontend**: Vanilla JavaScript, HTML, CSS
- **Контейнеризация**: Docker

## Выполненные задания со звёздочкой

✅ **Задание 1**: Поиск задач по подстроке в заголовке и комментарии  
✅ **Задание 2**: Аутентификация с использованием JWT-токенов  
✅ **Задание 3**: Создание Docker-образа для деплоя  

## Структура проекта

```
go_final_project/
├── cmd/
│   └── scheduler/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── auth/                    # Пакет аутентификации
│   │   ├── jwt.go               # Генерация и валидация JWT-токенов
│   │   └── middleware.go        # HTTP middleware для проверки аутентификации
│   ├── config/
│   │   └── config.go            # Конфигурация приложения
│   ├── database/
│   │   ├── db.go                # Инициализация БД
│   │   └── task.go              # Операции с задачами
│   ├── handlers/
│   │   ├── addtask.go           # POST/PUT/DELETE /api/task
│   │   ├── done.go              # POST /api/task/done
│   │   ├── nextdate.go          # GET /api/nextdate
│   │   ├── signin.go            # POST /api/signin (аутентификация)
│   │   ├── task.go              # GET /api/task
│   │   └── tasks.go             # GET /api/tasks
│   ├── models/
│   │   └── task.go              # Модель задачи
│   ├── services/
│   │   └── nextdate.go          # Логика расчёта следующей даты
│   └── utils/
│       ├── response.go          # Утилиты для HTTP-ответов
│       └── validation.go        # Валидация данных
├── web/                         # Статические файлы фронтенда
│   ├── index.html               # Главная страница
│   ├── login.html               # Страница входа
│   ├── css/
│   ├── js/
│   └── favicon.ico
├── tests/                       # Интеграционные тесты
├── Dockerfile                   # Конфигурация Docker
├── go.mod                       # Зависимости Go
└── README.md                    # Этот файл

```

## Инструкция по запуску локально

### Требования

- Go 1.24 или выше
- Git

### Шаги для запуска

1. **Клонируйте репозиторий** (если ещё не сделали):
   ```bash
   git clone <repository-url>
   cd go_final_project
   ```

2. **Установите зависимости**:
   ```bash
   go mod download
   ```

3. **Соберите приложение**:
   ```bash
   go build -o scheduler ./cmd/scheduler
   ```

4. **Запустите приложение**:
   ```bash
   ./scheduler
   ```

   Или запустите напрямую:
   ```bash
   go run ./cmd/scheduler/main.go
   ```

5. **Откройте браузер** и перейдите по адресу:
   ```
   http://localhost:7540
   ```

### Переменные окружения

Вы можете настроить приложение с помощью переменных окружения:

- `TODO_PORT` - порт веб-сервера (по умолчанию: 7540)
- `TODO_DBFILE` - путь к файлу базы данных SQLite (по умолчанию: scheduler.db)
- `TODO_PASSWORD` - пароль для аутентификации (если не задан, аутентификация отключена)

#### Пример запуска с переменными окружения:

```bash
export TODO_PORT=8080
export TODO_DBFILE=/path/to/database.db
export TODO_PASSWORD=mysecretpassword
./scheduler
```

Или в одной строке:

```bash
TODO_PORT=8080 TODO_DBFILE=./mydb.db TODO_PASSWORD=12345 ./scheduler
```

### Запуск с аутентификацией

Если задана переменная окружения `TODO_PASSWORD`, приложение потребует аутентификацию:

1. Откройте страницу логина: `http://localhost:7540/login.html`
2. Введите пароль, указанный в `TODO_PASSWORD`
3. После успешной аутентификации вы будете перенаправлены на главную страницу

JWT-токен сохраняется в куки `token` и действителен 8 часов.

## Инструкция по запуску тестов

### Настройка тестов

Перед запуском тестов необходимо настроить параметры в файле `tests/settings.go`:

```go
package tests

var Port = 7540              // Порт, на котором запущен сервер
var DBFile = "../scheduler.db"  // Путь к файлу БД
var FullNextDate = false         // Полная реализация nextdate (необязательно)
var Search = true                // Поддержка поиска (задание со звёздочкой)
var Token = ``                   // JWT-токен (если используется аутентификация)
```

### Запуск сервера для тестов

1. **Без аутентификации**:
   ```bash
   ./scheduler
   ```

2. **С аутентификацией**:
   ```bash
   TODO_PASSWORD=testpassword ./scheduler
   ```

   Затем получите токен:
   ```bash
   curl -X POST http://localhost:7540/api/signin \
     -H "Content-Type: application/json" \
     -d '{"password":"testpassword"}'
   ```

   Скопируйте значение поля `token` из ответа и вставьте в `tests/settings.go`:
   ```go
   var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
   ```

### Запуск тестов

```bash
go test ./tests -v
```

Или запустите отдельные тесты:

```bash
go test ./tests/app_1_test.go -v         # Базовые тесты
go test ./tests/db_2_test.go -v          # Тесты базы данных
go test ./tests/nextdate_3_test.go -v    # Тесты nextdate
go test ./tests/addtask_4_test.go -v     # Тесты добавления задач
go test ./tests/tasks_5_test.go -v       # Тесты списка задач
go test ./tests/task_6_test.go -v        # Тесты получения задачи
go test ./tests/task_7_test.go -v        # Тесты удаления задач
```

## Инструкция по сборке и запуску Docker-образа

### Сборка Docker-образа

```bash
docker build -t go-scheduler:latest .
```

### Запуск контейнера

#### Базовый запуск (без аутентификации):

```bash
docker run -d \
  --name scheduler \
  -p 7540:7540 \
  -v $(pwd)/data:/data \
  go-scheduler:latest
```

#### Запуск с аутентификацией:

```bash
docker run -d \
  --name scheduler \
  -p 7540:7540 \
  -v $(pwd)/data:/data \
  -e TODO_PASSWORD=mysecretpassword \
  go-scheduler:latest
```

#### Запуск с кастомными настройками:

```bash
docker run -d \
  --name scheduler \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  -e TODO_PORT=8080 \
  -e TODO_DBFILE=/data/scheduler.db \
  -e TODO_PASSWORD=12345 \
  go-scheduler:latest
```

### Подключение к существующей базе данных на хосте

Чтобы использовать существующую БД с хоста:

```bash
docker run -d \
  --name scheduler \
  -p 7540:7540 \
  -v /path/to/your/scheduler.db:/data/scheduler.db \
  -e TODO_DBFILE=/data/scheduler.db \
  go-scheduler:latest
```

Пример для текущей директории проекта:

```bash
docker run -d \
  --name scheduler \
  -p 7540:7540 \
  -v $(pwd)/scheduler.db:/data/scheduler.db \
  -e TODO_DBFILE=/data/scheduler.db \
  go-scheduler:latest
```

### Управление контейнером

**Просмотр логов**:
```bash
docker logs scheduler
```

**Остановка контейнера**:
```bash
docker stop scheduler
```

**Запуск остановленного контейнера**:
```bash
docker start scheduler
```

**Удаление контейнера**:
```bash
docker rm -f scheduler
```

### Проверка работы

После запуска контейнера откройте браузер:
```
http://localhost:7540
```

## API Endpoints

### Публичные (без аутентификации)

- `GET /api/nextdate` - расчёт следующей даты для правила повторения

### Требующие аутентификации (если задан `TODO_PASSWORD`)

- `POST /api/signin` - аутентификация пользователя
- `GET /api/task?id={id}` - получение задачи по ID
- `GET /api/tasks` - получение списка задач (с поддержкой поиска через `?search=...`)
- `POST /api/task` - создание новой задачи
- `PUT /api/task` - обновление существующей задачи
- `DELETE /api/task?id={id}` - удаление задачи
- `POST /api/task/done` - отметка задачи как выполненной

## Примеры использования API

### Аутентификация

```bash
curl -X POST http://localhost:7540/api/signin \
  -H "Content-Type: application/json" \
  -d '{"password":"mysecretpassword"}'
```

Ответ:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Создание задачи

```bash
curl -X POST http://localhost:7540/api/task \
  -H "Content-Type: application/json" \
  -H "Cookie: token=YOUR_TOKEN_HERE" \
  -d '{
    "title": "Купить продукты",
    "comment": "Молоко, хлеб, яйца",
    "date": "20251105",
    "repeat": ""
  }'
```

### Получение списка задач

```bash
curl -X GET http://localhost:7540/api/tasks \
  -H "Cookie: token=YOUR_TOKEN_HERE"
```

### Поиск задач

```bash
curl -X GET "http://localhost:7540/api/tasks?search=продукты" \
  -H "Cookie: token=YOUR_TOKEN_HERE"
```

### Расчёт следующей даты

```bash
curl -X GET "http://localhost:7540/api/nextdate?now=20251103&date=20251103&repeat=d%201"
```

Ответ (ежедневное повторение):
```
20251104
```

## Формат правил повторения

- `d N` - каждые N дней (например, `d 1` - каждый день, `d 7` - каждую неделю)
- `y` - каждый год (в ту же дату)
- `w 1,3,5` - по дням недели (1=Пн, 2=Вт, ..., 7=Вс)
- `m 1,15,-1` - по дням месяца (-1 = последний день, -2 = предпоследний)

## Устранение неполадок

### Проблема: Порт уже занят

```bash
# Найдите процесс, использующий порт
lsof -i :7540

# Остановите процесс или используйте другой порт
TODO_PORT=8080 ./scheduler
```

### Проблема: База данных заблокирована

Убедитесь, что не запущено несколько экземпляров приложения, использующих один файл БД.

### Проблема: Ошибка аутентификации в тестах

Убедитесь, что:
1. Переменная `TODO_PASSWORD` установлена при запуске сервера
2. Токен в `tests/settings.go` актуален (действителен 8 часов)
3. Токен получен для того же пароля, что используется сервером

## Лицензия

Этот проект создан в учебных целях в рамках финального задания курса Go от Яндекс Практикум.

## Автор

Проект выполнен студентом Яндекс Практикум в рамках итогового задания по курсу Go-разработки.
