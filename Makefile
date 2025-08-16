
all: fmt vet lint test

fmt:
	@gofmt -s -w .
	@echo "✓ Форматирование завершено"

vet:
	@go vet ./...
	@echo "✓ go vet прошёл"

lint:
	@golangci-lint run
	@echo "✓ Линтинг прошёл"

test:
	@go test -race ./...
	@echo "✓ Тесты с детектором гонок прошли"

# Запуск тестов с покрытием
coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@echo "📊 Покрытие кода рассчитано. Файл: coverage.out"
	@echo "💡 Используйте 'make cover-html' для просмотра в браузере"

# Покрытие в HTML
cover-html: coverage
	@go tool cover -html=coverage.out
	@echo "🌐 Открываю отчёт о покрытии в браузере..."

# Краткий отчёт о покрытии
cover-func: coverage
	@go tool cover -func=coverage.out

# Очистка временных файлов
clean:
	@rm -f coverage.out
	@rm -f *.log
	@echo "🧹 Временные файлы удалены"

help:
	@echo "Доступные команды:"
	@echo "  make fmt        — форматировать код"
	@echo "  make vet        — проверить ошибки"
	@echo "  make lint       — запустить линтер"
	@echo "  make test       — запустить тесты с -race"
	@echo "  make coverage   — запустить тесты с покрытием"
	@echo "  make cover-html — открыть отчёт о покрытии в браузере"
	@echo "  make cover-func — показать отчёт о покрытии в терминале"
	@echo "  make clean      — удалить временные файлы"
	@echo "  make all        — всё подряд"
	@echo "  make help       — эта подсказка"