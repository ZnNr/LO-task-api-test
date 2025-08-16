.PHONY: fmt vet lint test all

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

help:
	@echo "Доступные команды:"
	@echo "  make fmt    — форматировать код"
	@echo "  make vet    — проверить ошибки"
	@echo "  make lint   — запустить линтер"
	@echo "  make test   — запустить тесты с -race"
	@echo "  make all    — всё подряд"
	@echo "  make help   — эта подсказка"
