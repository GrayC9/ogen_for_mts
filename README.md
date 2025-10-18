# ogen_for_mts

Чтобы комментарии попали в сгенерированные методы интерфейса/хендлера, опишите их в поле `operation.description` OpenAPI (paths - <path> - <method> - description). Это поле ogen переносит в комментарии над методами.

Структура проекта:

- openapi.yaml — спецификация с сущностью User и CRUD‑операциями. У каждой операции есть `summary` и подробный `description`.
- Makefile — запуск: `make generate`, `make run-server`, `make run-client`, `make verify`.
- generate.go, tools.go — фиксация инструмента ogen и пример `go generate`.
- internal/server — реализация in‑memory обработчика, удовлетворяющего сгенерированому интерфейсу (файлы под build‑тегом `//go:build ogen`).
- cmd/server — запуск HTTP‑сервера.
- cmd/client — пример клиента, использующего сгенерированный клиент.

Установка ogen

```
go install github.com/ogen-go/ogen/cmd/ogen@latest
```

Генерация кода из OpenAPI

```
make generate
# или напрямую:
# ogen --target internal/api --package api --clean openapi.yaml
```

После генерации в `internal/api` появятся файлы. Откройте интерфейс `Handler` — над методами будут комментарии из OpenAPI:
- CreateUser — «Creates a new user and returns it with an assigned ID.»
- ListUsers — «Returns a list of all users.»
- GetUser — «Retrieves a user by its unique identifier.»
- UpdateUser — «Updates an existing user and returns the updated object.»
- DeleteUser — «Deletes a user and returns no content.»

Запуск сервера

```
# при необходимости подтянуть зависимости
make tidy

# запустить сервер (по умолчанию на :8080)
make run-server
# или
ADDR=":8081" go run -tags=ogen ./cmd/server
```

Сервер логирует «Starting server on :8080». Адрес можно задать через переменную окружения `ADDR`.

Запуск клиента

```
# в отдельном терминале, когда сервер уже запущен
make run-client
# или
BASE_URL="http://localhost:8080" go run -tags=ogen ./cmd/client
```

Клиент последовательно выполняет CRUD и печатает результаты: Created -List - Get - Updated - Deleted.

Быстрая проверка end‑to‑end

```
make verify
```
ц
Скрипт сгенерирует код, поднимет сервер, выполнит CRUD запросы через клиент и корректно остановит сервер. 

Проверка через curl

```
# Создать пользователя
curl -sS -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice"}' | jq .

# Список пользователей
curl -sS http://localhost:8080/users | jq .

# Получить по id
curl -sS http://localhost:8080/users/1 | jq .

# Обновить
curl -sS -X PUT http://localhost:8080/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Updated","description":"Updated user"}' | jq .

# Удалить
curl -sS -X DELETE http://localhost:8080/users/1 -i
```

Как это работает внутри

- Сущность `User`: поля `id:int64`, `name:string`, `description:*string` (nullable).
- Хранилище — простая in‑memory map с мьютексом и автоинкрементом ID (`internal/storage`).
- Обработчик реализует интерфейс, который сгенерировал(а) ogen (`internal/server/handler.go`). При отсутствии записи операции `Get/Update` возвращают 404‑ответы соответствующих типов из сгенерированного пакета.
