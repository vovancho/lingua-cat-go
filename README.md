# Lingua-Cat-Go

**Lingua-Cat-Go** — это проект для обучения иностранным словам.
Система состоит из API Gateway (Traefik), сервиса аутентификации (Keycloak), трёх основных микросервисов (Dictionary, Exercise, Analytics), а также Kafka для асинхронной передачи событий и ClickHouse для аналитики.

![Architecture Diagram](project/scheme.png)

---

## Структура репозитория

```

/backend
├── analytics                      # Сервис аналитики: потребляет события из Kafka и отдаёт статистику
│  ├── cmd
│  │  └── main.go                  # Точка входа: запуск HTTP-сервера и Kafka-консьюмера
│  ├── delivery                    # Внешний слой
│  │  ├── http
│  │  │  └── exercise_complete_handler.go  # HTTP-эндпоинты
│  │  └── kafka
│  │     └── exercise_complete_handler.go  # Kafka-обработчик
│  ├── docker-compose.yml          # Конфигурация Docker
│  ├── docs
│  │  └── swagger.json             # OpenAPI-спецификация для analytics
│  ├── domain                      # Бизнес-модель
│  │  ├── errors.go                # Определение бизнес-ошибок
│  │  ├── exercise_complete.go     # Интерфейсы и сущность выполненного упражнения
│  │  ├── user.go                  # Интерфейсы и сущность пользователя
│  │  └── validator.go             # Кастомные валидаторы доменных структур
│  ├── go.mod
│  ├── go.sum
│  ├── internal                    # Внутренние пакеты
│  │  ├── config                   # Чтение конфигурации
│  │  │  └── config.go             # Конфигурация сервиса (из переменных окружения)
│  │  ├── validator
│  │  │  └── validator.go          # Интеграция кастомных валидаторов
│  │  └── wire
│  │     ├── wire.go               # DI-конфигурация с Google Wire
│  │     └── wire_gen.go
│  ├── migrations                  # SQL-миграции для ClickHouse
│  │  ├── 20250506122921_init_schema.down.sql
│  │  └── 20250506122921_init_schema.up.sql
│  ├── repository                  # Реализация хранилищ
│  │  ├── clickhouse
│  │  │  └── exercise_complete.go  # Запись и чтение статистики по выполненным упражнениям
│  │  └── http
│  │     └── user.go               # Получение имени пользователя в сервисе аутентификации
│  └── usecase                     # Бизнес-логика
│     ├── exercise_complete_usecase.go  # Получение и сохранение статистики
│     ├── exercise_complete_usecase_test.go  # Тесты для exercise_complete-usecase
│     ├── user_usecase.go          # Получение имени пользователя
│     └── user_usecase_test.go     # Тесты для user-usecase
├── dictionary                # Сервис словаря: CRUD-операции со словами и переводами
│  ├── cmd                  # Точка входа
│  │  └── main.go           # Запуск HTTP и gRPC серверов
│  ├── delivery             # Handlers для запросов
│  │  ├── grpc              # gRPC-сервер
│  │  │  ├── dictionary_handler.go  # Обработка gRPC-запросов
│  │  │  └── gen            # Сгенерированные файлы протобуфа
│  │  │     ├── dictionary.pb.go
│  │  │     ├── dictionary.pb.gw.go
│  │  │     └── dictionary_grpc.pb.go
│  │  └── http              # HTTP-сервер и валидация
│  │     ├── dictionary_handler.go  # Обработка HTTP-эндпоинтов
│  │     └── validator.go    # HTTP-валидация входных данных
│  ├── docker-compose.yml   # Конфиг Docker для dictionary
│  ├── docs                 # OpenAPI-и gRPC-документация
│  │  ├── grpc-gw-swagger.json
│  │  └── swagger.json
│  ├── domain               # Доменные модели и ошибки
│  │  ├── dictionary.go
│  │  ├── errors.go
│  │  ├── sentence.go
│  │  ├── translation.go
│  │  └── validator.go
│  ├── go.mod
│  ├── go.sum
│  ├── internal             # Внутренние реализации
│  │  ├── config
│  │  │  └── config.go
│  │  ├── validator
│  │  │  └── validator.go
│  │  └── wire
│  │     ├── wire.go
│  │     └── wire_gen.go
│  ├── migrations
│  │  ├── 20250424135436_init_schema.down.sql
│  │  ├── 20250424135436_init_schema.up.sql
│  │  ├── 20250429112627_init_data.down.sql
│  │  └── 20250429112627_init_data.up.sql
│  ├── repository
│  │  └── postgres
│  │     └── dictionary.go  # Реализация репозитория на PostgreSQL
│  └── usecase
│     ├── dictionary_usecase.go   # Бизнес-правила работы со словарём
│     └── dictionary_usecase_test.go  # Тесты для dictionary usecase
├── exercise                # Сервис упражнений: создание и прохождение
│  ├── cmd
│  │  └── main.go           # Запуск HTTP-сервера и продюсера Kafka (Outbox)
│  ├── delivery
│  │  └── http
│  │     ├── exercise_handler.go  # HTTP-эндпоинты для упражнений
│  │     └── task_handler.go      # HTTP-эндпоинты для задач внутри упражнения
│  ├── docker-compose.yml   # Docker-окружение для exercise
│  ├── docs
│  │  └── swagger.json      # OpenAPI-спецификация для exercise
│  ├── domain
│  │  ├── dictionary.go
│  │  ├── errors.go
│  │  ├── exercise.go       # Сущность упражнения
│  │  ├── sentence.go
│  │  ├── task.go           # Сущность задачи внутри упражнения
│  │  ├── translation.go
│  │  └── validator.go      # Валидация доменных структур
│  ├── go.mod
│  ├── go.sum
│  ├── internal
│  │  ├── config
│  │  │  └── config.go
│  │  ├── validator
│  │  │  └── validator.go
│  │  └── wire
│  │     ├── wire.go
│  │     └── wire_gen.go
│  ├── migrations
│  │  ├── 20250501080432_init_schema.down.sql
│  │  ├── 20250501080432_init_schema.up.sql
│  │  ├── 20250506034401_outbox_init.down.sql
│  │  └── 20250506034401_outbox_init.up.sql
│  ├── repository
│  │  ├── grpc
│  │  │  ├── dictionary.go   # К gRPC-репозиторию словаря
│  │  │  └── gen
│  │  │     ├── dictionary.pb.go
│  │  │     └── dictionary_grpc.pb.go
│  │  └── postgres
│  │     ├── exercise.go     # Реализация репозитория упражнений
│  │     └── task.go         # Репозиторий задач (tasks)
│  └── usecase
│     ├── dictionary_usecase.go  # Использование словаря в exercise
│     ├── dictionary_usecase_test.go
│     ├── exercise_usecase.go    # Логика создания и завершения упражнений
│     ├── exercise_usecase_test.go
│     ├── task_usecase.go        # Логика задач внутри упражнения
│     └── task_usecase_test.go
├── pkg                     # Общие пакеты, разделяемые между сервисами
│  ├── auth                 # JWT-аутентификация и middleware
│  │  ├── auth.go
│  │  ├── go.mod
│  │  ├── go.sum
│  │  ├── interceptor.go
│  │  └── middleware.go
│  ├── db                   # Инициализация и доступ к базе данных
│  │  ├── db.go
│  │  ├── go.mod
│  │  └── go.sum
│  ├── error                # Утилиты для работы с ошибками
│  │  ├── error.go
│  │  └── go.mod
│  ├── eventpublisher       # Абстракция публикации событий (Outbox, Kafka)
│  │  ├── eventpublisher.go
│  │  ├── go.mod
│  │  └── go.sum
│  ├── request              # Утилиты парсинга HTTP-запросов
│  │  ├── go.mod
│  │  ├── go.sum
│  │  └── request.go
│  ├── response             # Формирование HTTP-ответов и middleware
│  │  ├── go.mod
│  │  ├── go.sum
│  │  ├── middleware.go
│  │  └── response.go
│  ├── tracing              # Инструменты для распределённого трейсинга
│  │  ├── go.mod
│  │  ├── go.sum
│  │  └── tracing.go
│  ├── translator          # Утилиты для перевода текста (например, внешние API)
│  │  ├── go.mod
│  │  ├── go.sum
│  │  └── translator.go
│  ├── txmanager           # Управление транзакциями (DB и Outbox)
│  │  ├── go.mod
│  │  ├── go.sum
│  │  └── txmanager.go
│  └── validator           # Общие правила валидации и обёртка validator.v10
│     ├── go.mod
│     ├── go.sum
│     └── validator.go
└── proto                   # Папка с Proto-файлами для gRPC-генерации
   └── dictionary.proto    # Описание сервисов и сообщений для словаря

```

---

## 🚀 Обзор компонентов

| Компонент                 | Описание                                                                                                                            |
|---------------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| **Traefik (API Gateway)** | Принимает все HTTP/HTTPS запросы, обрабатывает маршрутизацию, SSL-терминацию, аутентификацию JWT и пробрасывает запрос дальше.      |
| **Keycloak**              | Сервис Identity & Access Management. Выдаёт JWT-токены по протоколам OpenID Connect / OAuth2 и хранит учётные записи пользователей. |
| **Dictionary**            | Микросервис для CRUD-операций со словами и переводами. Предоставляет HTTP + gRPC API.                                               |
| **Exercise**              | Микросервис для создания и прохождения упражнений. После завершения шлёт статистику в Kafka.                                        |
| **Kafka**                 | Брокер сообщений — принимает события `lcg_exercise_completed` и передаёт их дальше.                                                 |
| **Analytics**             | Консумер Kafka + HTTP API. Пишет статистику в ClickHouse и выдаёт агрегированные данные.                                            |
| **Swagger & Jaeger**      | Самодокументируемые схемы API (Swagger UI) и трейсинг (Jaeger) для отладки и мониторинга.                                           |

---

## 🔗 Конечные точки (HTTP API)

### Gateway (Traefik)
- Все запросы начинаются с `/v1/...`
- JWT-авторизация: передаётся в заголовке `Authorization: Bearer <token>`

### Dictionary Service  
```

POST   /v1/dictionary
POST   /v1/dictionary/{id}/name
GET    /v1/dictionary/{id}
DELETE /v1/dictionary/{id}
GET    /v1/dictionary/random

```

### Exercise Service  
```

POST   /v1/exercise
GET    /v1/exercise/{id}
POST   /v1/exercise/{id}/task
GET    /v1/exercise/{id}/task/{task\_id}
POST   /v1/exercise/{id}/task/{task\_id}/word-selected

```

- По завершении упражнения микросервис публикует в Kafka тему `lcg_exercise_completed` сообщение со статистикой.

### Analytics Service  
```

GET /v1/analytics/user/{user\_id}

```
- Подписывается на тему `lcg_exercise_completed`, агрегирует данные в ClickHouse.

---

## 🔧 Быстрый старт

1. Склонировать репозиторий:
```bash
   git clone https://github.com/your-org/lingua-cat-go.git
   cd lingua-cat-go/backend
```

2. Запустить в Docker Compose:

```bash
docker-compose up
```

3. Открыть в браузере:

  * **Swagger UI**: `http://swagger.localhost/`
  * **Jaeger UI**: `http://jaeger.localhost/`
  * **Keycloak**: `http://keycloak.localhost/` (логин/пароль по умолчанию: `admin/admin`)
  * **Dictionary API**: `http://api.lingua-cat-go.localhost/dictionary/swagger.json`
  * **Exercise API**: `http://api.lingua-cat-go.localhost/exercise/swagger.json`
  * **Analytics API**: `http://api.lingua-cat-go.localhost/analytics/swagger.json`

4. Создать Realm и клиента в Keycloak, получить JWT и начать отправлять запросы через Traefik.

---

## 📖 Документация и трассировка

* **Swagger**: Самодокументируемые спецификации для каждого сервиса доступны по URL из раздела “Быстрый старт”.
* **Jaeger**: Собирает трейсы HTTP и gRPC, поможет отследить цепочки запросов между сервисами.

---

## 🤝 Вклад

1. Fork репозиторий
2. Создать ветку `feature/ваша-фича`
3. Сделать Commit & Push
4. Открыть Pull Request

---

## ⚖️ Лицензия

MIT © Your Organization
