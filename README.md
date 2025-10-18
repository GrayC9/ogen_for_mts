# ogen_for_mts

Небольшой демо‑проект про генерацию Go‑кода из OpenAPI с помощью ogen, добавление комментариев к методам из описания и запуск простого сервера/клиента на in‑memory хранилище.

Короткий ответ на главный вопрос: чтобы комментарии попали в сгенерированные методы интерфейса/хендлера, опишите их в поле `operation.description` OpenAPI (paths → <path> → <method> → description). Именно это поле ogen переносит в комментарии над методами.

Что внутри проекта:

- openapi.yaml — спецификация с сущностью User и CRUD‑операциями. У каждой операции есть `summary` и подробный `description`.
- Makefile — удобные цели: `make generate`, `make run-server`, `make run-client`, `make verify`.
- generate.go, tools.go — фиксация инструмента ogen и пример `go generate`.
- internal/server — реализация in‑memory обработчика, удовлетворяющего сгенерированному интерфейсу (файлы под build‑тегом `//go:build ogen`).
- cmd/server — запуск HTTP‑сервера на базе сгенерированного роутера.
- cmd/client — пример клиента, использующего сгенерированный клиент.

Важные нюансы структуры кода

- В репозитории уже лежит один снэпшот сгенерированного кода в каталоге `internal/api_1` — именно его используют примеры сервера и клиента.
- Цель `make generate` генерирует свежий код в `internal/api` (другая папка). Это сделано специально, чтобы наглядно увидеть, как появляются комментарии из `description` прямо после генерации, не трогая рабочие примеры.
- Мы не патчим и не форкаем ogen — все комментарии берутся из OpenAPI как есть.

Требования

- Go 1.24+
- Установленный ogen

Установка ogen

```
go install github.com/ogen-go/ogen/cmd/ogen@latest
```

Убедитесь, что `$GOPATH/bin` (или `$GOBIN`) находится в `$PATH`.

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

Файлы сервера/клиента собраны под build‑тегом `ogen`, поэтому тег нужно указывать при запуске.

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

Клиент последовательно выполняет CRUD и печатает результаты: Created → List → Get → Updated → Deleted.

Быстрая проверка end‑to‑end

```
make verify
```

Скрипт сгенерирует код, поднимет сервер, выполнит CRUD запросы через клиент и корректно остановит сервер. Если порт 8080 занят — задайте `ADDR=":8081"` при запуске.

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
- Хранилище — простая in‑memory map с мьютексом и автоинкрементом ID (см. `internal/storage`).
- Обработчик реализует интерфейс, который сгенерировал(а) ogen (см. `internal/server/handler.go`). При отсутствии записи операции `Get/Update` возвращают 404‑ответы соответствующих типов из сгенерированного пакета.

FAQ

- Где писать текст комментариев? — В `operation.description`. Именно его переносит ogen в комментарии методов интерфейса/хендлеров. `summary` остаётся кратким заголовком.
- Можно ли поменять формат комментариев без изменений в ogen? — Нет, в рамках задачи мы влияем только на текст в OpenAPI.