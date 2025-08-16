# Task API

REST API для управления задачами с асинхронным логированием через канал.

## Функционал

- `GET /tasks` — получить список задач (фильтрация по статусу)
- `GET /tasks/{id}` — получить задачу по ID
- `POST /tasks` — создать задачу

## Требования

- Использование только стандартных библиотек Go
- Асинхронное логирование через канал
- Хранилище в памяти
- Поддержка graceful shutdown
- Чистая архитектура

## Запуск

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/znnr/LO-task-api-test.git
   cd LO-task-api-test 
   ```
2. Установите зависимости:
   ```bash
   go mod download
   ```
3. Запустите сервер:
   ```bash
   go run cmd/main.go
   ```
4. Сервер будет доступен по адресу: http://localhost:8080

##  Тестирование и совместимость
Проект разрабатывался и тестировался в следующих средах:
- Основная разработка: Windows 11 (Go 1.25.0)
- Тестирование: WSL Ubuntu 20.04/22.04
- Сборка: Все Makefile команды работают в Linux/WSL/MacOS

Особенности запуска в разных средах
- Windows (PowerShell/CMD)
```bash
# Запуск сервера
go run cmd\main.go
# Тестирование
go test .\...
```
- Linux/WSL/MacOS
```bash
# Запуск сервера
go run cmd/main.go
# Тестирование
go test ./...
```

## Примеры запросов

1. Создание задачи

```bash
   curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Задача 1", "status": "pending"}'
```

2. Получение всех задач

 ```bash
curl http://localhost:8080/tasks
```

3. Получение задач с фильтрацией по статусу

 ```bash
curl "http://localhost:8080/tasks?status=pending"
```

4. Получение задачи по ID

  ```bash
curl http://localhost:8080/tasks/1
  ```

## Тестирование

1. Запуск всех тестов

```bash
go test ./...
```

2. Запуск тестов с детектором гонок

```bash
go test -race ./...
```

3. Запуск тестов с покрытием кода

```bash
# Показать покрытие в терминале
go test -cover ./...

# Создать подробный отчёт
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Открыть отчёт в браузере
go tool cover -html=coverage.out
```

## Использование Makefile

   ```bash
# Форматирование кода
make fmt

# Проверка ошибок
make vet

# Линтинг
make lint

# Тесты с детектором гонок
make test

# Тесты с покрытием
make coverage

# Открыть покрытие в браузере
make cover-html

# Запустить всё
make all
```

## Структура проекта

LO-task-api-test/
├── cmd/
│ └── main.go # Точка входа
├── internal/
│ ├── task/
│ │ ├── task.go # Модели задач
│ │ ├── repository.go # Хранилище в памяти
│ │ ├── service.go # Бизнес-логика
│ │ └── *.go # Тесты
│ ├── handler/
│ │ ├── task_handler.go # HTTP хендлеры
│ │ └── *.go # Тесты
│ └── logger/
│ ├── logger.go # Асинхронный логгер
│ └── *.go # Тесты
├── integration/
│ └── *.go # Интеграционные тесты
├── go.mod # Модуль Go
├── go.sum # Зависимости
├── Makefile # Сборка и тестирование
├── .golangci.yml # Конфигурация линтеров
└── README.md # Документация

## Логирование

Все операции логируются в файл app.log в формате:
[2024-01-01T12:00:00Z] GET /tasks: 3 tasks returned
[2024-01-01T12:00:01Z] POST /tasks: created task with title 'Новая задача'

## Graceful Shutdown

Для корректного завершения работы нажмите Ctrl+C. Сервер:

- Завершит обработку текущих запросов
- Остановит логгер
- Корректно освободит ресурсы

## Тестирование производительности

Для нагрузочного тестирования можно использовать hey:

# Установка hey

go install github.com/rakyll/hey@latest

## Статический анализ

Проект использует следующие инструменты:

gofmt — форматирование
go vet — проверка ошибок
golangci-lint — комплексный линтинг
staticcheck — статический анализ
errcheck — проверка обработки ошибок
revive — стилистические проверки

Автор
znnr