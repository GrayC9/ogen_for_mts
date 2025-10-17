# ogen_for_mts

Тестовое задание: генерация Go кода с помощью ogen, добавление комментариев к методам из OpenAPI и реализация сервера/клиента с in-memory хранилищем.

Ключевой ответ: чтобы комментарий попал в сгенерированный ogen код для handle-метода, его нужно указать в поле `description` объекта операции OpenAPI (paths -> <path> -> <method> -> description). Именно `operation.description` ogen переносит в комментарии над соответствующим методом интерфейса/хендлера.

Проект содержит:
- openapi.yaml — OpenAPI спецификация с одной сущностью User и CRUD методами. У каждой операции заполнены `summary` и `description`.
- Makefile — команда для генерации кода `make generate`.
- tools.go — фиксация инструмента `ogen` через build tools, можно `go generate`.
- internal/server/ — реализация in-memory handlers, удовлетворяющих интерфейсу, который сгенерирует ogen. Файлы помечены build-тегом `//go:build ogen`, чтобы проект можно было держать в репозитории до генерации кода.
- cmd/server — запуск HTTP сервера на основе сгенерированного роутера/серверного кода.
- cmd/client — пример клиента, использующего сгенерированный клиент для CRUD операций.

Важное замечание:
- Мы не расширяем/не модифицируем сам ogen. Комментарии появляются благодаря полям `description` в OpenAPI. После генерации вы увидите эти комментарии в файлах в каталоге `internal/api` над методами интерфейса/хендлера.

Требования:
- Go 1.22+
- Установленный ogen

Установка ogen:
```
go install github.com/ogen-go/ogen/cmd/ogen@latest
```
Убедитесь, что `$GOPATH/bin` присутствует в PATH.

Генерация кода:
```
make generate
# или напрямую
# ogen --target internal/api --package api --clean openapi.yaml
```
После генерации в каталоге `internal/api` появятся Go файлы. В них найдите интерфейс хендлера/сервера (обычно `type Handler interface { ... }` или похожую сущность) — над каждым методом появится комментарий из `description` соответствующей операции.

Сборка и запуск сервера:
Файлы сервера и клиента помечены build-тегом `ogen`, чтобы они компилировались только после генерации кода. Поэтому при запуске нужно указать тег сборки.
```
# генерация кода
make generate

# модульные зависимости (опционально)
go mod tidy

# запустить сервер на :8080
go run -tags=ogen ./cmd/server
```
Логи покажут, что сервер стартовал на 8080.

Как проверить, что всё работает (автоматически):
- Один шаг, который делает всё за вас: сгенерирует код, поднимет сервер, выполнит CRUD клиентом и остановит сервер:
```
make verify
```
Что вы увидите:
- [verify] Generating code...
- [verify] Starting server...
- В логе клиента последовательность: Created, List, Get, Updated, Deleted
- Если порт :8080 занят, задайте другой адрес для сервера: `ADDR=":8081" go run -tags=ogen ./cmd/server` или временно остановите процесс, держащий порт.

Проверка запросами (curl):
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

Запуск клиента:
```
# в другом терминале, когда сервер уже запущен
go run -tags=ogen ./cmd/client
```
Клиент выполнит последовательные CRUD операции и выведет результаты в лог.

Как убедиться, что комментарии появились в сгенерированном коде:
- Откройте любой из файлов в `internal/api`, где объявлен интерфейс обработчика (обычно рядом с именованными tipами операций).
- Найдите методы, соответствующие `CreateUser`, `ListUsers`, `GetUser`, `UpdateUser`, `DeleteUser`.
- Над каждым методом должен быть комментарий. Текст комментария будет взят из `description` в `openapi.yaml`.

Структура OpenAPI:
- Сущность `User` с полями `id:int64`, `name:string`, `description:*string`.
- CRUD маршруты:
  - POST /users — CreateUser — description: "Creates a new user and returns it with an assigned ID."
  - GET /users — ListUsers — description: "Returns a list of all users."
  - GET /users/{id} — GetUser — description: "Retrieves a user by its unique identifier."
  - PUT /users/{id} — UpdateUser — description: "Updates an existing user and returns the updated object."
  - DELETE /users/{id} — DeleteUser — description: "Deletes a user and returns no content."

Примечания по реализованной серверной логике:
- Используется in-memory map[id]*User и атомарный инкремент id.
- Ожидаются типы, которые сгенерирует ogen. Реализация размещена в `internal/server/handler.go` и будет собираться при теге `ogen`.

FAQ:
- В какое поле добавить комментарий? — В operation.description (а не summary). Summary часто используется кратко, но для комментариев хендлеров ogen использует description.
- Можно ли поменять формат комментариев без форка ogen? — В рамках задачи — нет, но можно влиять на текст через описание в OpenAPI.