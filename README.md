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
│  ├── docker-compose.yml          # Конфиг Docker для analytics
│  ├── docs
│  │  └── swagger.json             # OpenAPI-спецификация для HTTP analytics
│  ├── domain                      # Бизнес-модель
│  │  ├── errors.go                # Определение бизнес-ошибок
│  │  ├── exercise_complete.go     # Интерфейсы и сущность выполненного упражнения
│  │  ├── user.go                  # Интерфейсы и сущность пользователя
│  │  └── validator.go             # Кастомные валидаторы доменных структур
│  ├── internal                    # Внутренние пакеты
│  │  ├── config
│  │  │  └── config.go             # Конфигурация сервиса (из переменных окружения)
│  │  ├── validator
│  │  │  └── validator.go          # Интеграция кастомных валидаторов
│  │  └── wire
│  │     ├── wire.go               # DI-конфигурация с Google Wire
│  ├── migrations                  # SQL-миграции для ClickHouse
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
│
├── dictionary                     # Сервис словаря: Получение и сохранение слов и переводов
│  ├── cmd
│  │  └── main.go                  # Точка входа: запуск HTTP и gRPC серверов
│  ├── delivery                    # Внешний слой
│  │  ├── grpc
│  │  │  ├── dictionary_handler.go  # Обработка gRPC-запросов
│  │  │  └── gen                    # Сгенерированные proto-файлы
│  │  └── http
│  │     ├── dictionary_handler.go  # HTTP-эндпоинты
│  │     └── validator.go          # Кастомные валидаторы входных данных
│  ├── docker-compose.yml          # Конфиг Docker для dictionary
│  ├── docs
│  │  ├── grpc-gw-swagger.json     # OpenAPI-спецификация для gRPC-gateway analytics
│  │  └── swagger.json             # OpenAPI-спецификация для HTTP analytics
│  ├── domain                      # Бизнес-модель
│  │  ├── dictionary.go            # Интерфейсы и сущность слова
│  │  ├── errors.go                # Определение бизнес-ошибок
│  │  ├── sentence.go              # Сущность примеров предложений с переводимым словом
│  │  ├── translation.go           # Сущность перевода слова
│  │  └── validator.go             # Кастомные валидаторы доменных структур
│  ├── internal                    # Внутренние пакеты
│  │  ├── config
│  │  │  └── config.go             # Конфигурация сервиса (из переменных окружения)
│  │  ├── validator
│  │  │  └── validator.go          # Интеграция кастомных валидаторов
│  │  └── wire
│  │     ├── wire.go               # DI-конфигурация с Google Wire
│  ├── migrations                  # SQL-миграции для PostgreSQL
│  ├── repository                  # Реализация хранилищ
│  │  └── postgres
│  │     └── dictionary.go         # Запись и чтение слов с их переводами и предложениями
│  └── usecase                     # Бизнес-логика
│     ├── dictionary_usecase.go    # Получение и сохранение слов с их переводами и предложениями
│     └── dictionary_usecase_test.go  # Тесты для dictionary-usecase
│
├── exercise                       # Сервис упражнений: создание упражнений и прохождение заданий
│  ├── cmd
│  │  └── main.go                  # Точка входа: запуск HTTP-сервера и продюсера Kafka (Outbox)
│  ├── delivery                    # Внешний слой
│  │  └── http
│  │     ├── exercise_handler.go   # HTTP-эндпоинты для упражнений
│  │     └── task_handler.go       # HTTP-эндпоинты для заданий упражнения
│  ├── docker-compose.yml          # Конфиг Docker для exercise
│  ├── docs
│  │  └── swagger.json             # OpenAPI-спецификация для HTTP exercise
│  ├── domain                      # Бизнес-модель
│  │  ├── dictionary.go            # Интерфейсы и сущность слова
│  │  ├── errors.go                # Определение бизнес-ошибок
│  │  ├── exercise.go              # Интерфейсы и сущность упражнения
│  │  ├── sentence.go              # Сущность примеров предложений с переводимым словом
│  │  ├── task.go                  # Интерфейсы и сущность задания упражнения
│  │  ├── translation.go           # Сущность перевода слова
│  │  └── validator.go             # Валидация доменных структур
│  ├── internal                    # Внутренние пакеты
│  │  ├── config
│  │  │  └── config.go             # Конфигурация сервиса (из переменных окружения)
│  │  ├── validator
│  │  │  └── validator.go          # Интеграция кастомных валидаторов
│  │  └── wire
│  │     ├── wire.go               # DI-конфигурация с Google Wire
│  ├── migrations                  # SQL-миграции для PostgreSQL
│  ├── repository                  # Реализация хранилищ
│  │  ├── grpc
│  │  │  ├── dictionary.go         # Чтение случайного набора слов или слов по идентификаторам
│  │  │  └── gen                   # Сгенерированные proto-файлы
│  │  └── postgres
│  │     ├── exercise.go           # Запись и чтение упражнений, проверка владельца
│  │     └── task.go               # Запись и чтение заданий упражнения, проверка принадлежности
│  └── usecase                     # Бизнес-логика
│     ├── dictionary_usecase.go    # Получение случайного набора слов или слов по идентификаторам
│     ├── dictionary_usecase_test.go  # Тесты для dictionary-usecase
│     ├── exercise_usecase.go      # Получение и сохранение упражнений, проверка владельца
│     ├── exercise_usecase_test.go  # Тесты для exercise-usecase
│     ├── task_usecase.go          # Получение и сохранение заданий упражнения, проверка принадлежности
│     └── task_usecase_test.go     # Тесты для task-usecase
│
├── pkg                            # Общие пакеты, разделяемые между сервисами
│  ├── auth                        # JWT-аутентификация и middleware
│  │  ├── auth.go                  # Сервис аутентификации
│  │  ├── interceptor.go           # interceptor для gRPC (проверка JWT-токена, сохранение в контекст)
│  │  └── middleware.go            # middleware для HTTP (проверка JWT-токена, сохранение в контекст)
│  ├── db                          # Инициализация и доступ к базе данных
│  ├── error                       # Утилиты для работы с ошибками
│  ├── eventpublisher              # Публикация событий (Outbox, Kafka)
│  ├── request                     # Утилиты парсинга HTTP-запросов
│  ├── response                    # Формирование HTTP-ответов и middleware
│  │  ├── middleware.go            # middleware для обработки паник
│  │  └── response.go              # Утилиты формирования HTTP-ответов (в том числе ошибок)
│  ├── tracing                     # Инструменты для распределённого трейсинга
│  ├── translator                  # Инициализация переводчика
│  ├── txmanager                   # Управление транзакциями базы данных
│  └── validator                   # Инициализация валидатора
│
└── proto                          # Папка с Proto-файлами для gRPC-генерации
   └── dictionary.proto            # Описание сервисов и сообщений для dictionary

```

---

## Обзор компонентов

| Компонент                 | Описание                                                                                                              |
|---------------------------|-----------------------------------------------------------------------------------------------------------------------|
| **Traefik (API Gateway)** | Принимает все HTTP запросы, обрабатывает маршрутизацию и пробрасывает запрос дальше.                                  |
| **Keycloak**              | Сервис аутентификации. Выдаёт JWT-токены по протоколам OpenID Connect / OAuth2 и хранит учётные записи пользователей. |
| **Dictionary**            | Микросервис для CRUD-операций со словами и переводами. Предоставляет HTTP + gRPC API.                                 |
| **Exercise**              | Микросервис для создания упражнений и прохождения заданий. После завершения упражнения шлёт статистику в Kafka.       |
| **Kafka**                 | Брокер сообщений — принимает события о выполненных упражнениях и передаёт их дальше.                                  |
| **Analytics**             | Консьюмер Kafka + HTTP API. Пишет статистику в ClickHouse и выдаёт агрегированные данные по пользователю.             |
| **Swagger & Jaeger**      | Самодокументируемые схемы API (Swagger UI) и трейсинг (Jaeger) для отладки и мониторинга.                             |

---

## Конечные точки (HTTP API)

### Gateway (Traefik)
- JWT-авторизация: передаётся в заголовке `Authorization: Bearer <token>`

### Dictionary Service  
```

POST   /v1/dictionary                        # создать слово с переводами и примерами
POST   /v1/dictionary/{id}/name              # переименовать слово
GET    /v1/dictionary/{id}                   # получить слово по id
DELETE /v1/dictionary/{id}                   # удалить слово по id (заполнить deleted_at)
GET    /v1/dictionary/random                 # получить случайный набор слов по языку

GET    /grpc-gateway/v1/dictionary/random    # получить случайный набор слов по языку
GET    /grpc-gateway/v1/dictionary           # получить набор слов по идентификаторам

```

### Exercise Service  
```

POST   /v1/exercise                              # создать новое упражнение с заданным количеством заданий
GET    /v1/exercise/{exer_id}                    # получить упражнение по id
POST   /v1/exercise/{exer_id}/task               # создать задание для упражнения со случайным набором слов
GET    /v1/exercise/{exer_id}/task/{task_id}     # получить задание по id
POST   /v1/exercise/{exer_id}/task/{task_id}/word-selected    # выбрать слово в задании

```

- По завершении упражнения микросервис публикует в topic Kafka `lcg_exercise_completed` сообщение со статистикой.

### Analytics Service  
```

GET /v1/analytics/user/{user_id}    # получить статистику по пользователю

```
- Подписывается на topic Kafka `lcg_exercise_completed`, записывает данные в ClickHouse.

---

## Быстрый старт

1. Склонировать репозиторий:
```bash
   git clone https://github.com/your-org/lingua-cat-go.git
   cd lingua-cat-go
```

2. Инициализировать и запустить проект:

```bash
make init
```

- Эта команда подымет контейнеры, импортирует конфиг keycloak, выполнит миграции.

3. Открыть в браузере:

  * **Swagger UI**: `http://swagger.localhost/`
  * **Jaeger UI**: `http://jaeger.localhost/`
  * **Keycloak**: `http://keycloak.localhost/` (логин/пароль по умолчанию: `admin/admin`)
  * **Dictionary API**: `http://api.lingua-cat-go.localhost/dictionary/`
  * **Exercise API**: `http://api.lingua-cat-go.localhost/exercise/`
  * **Analytics API**: `http://api.lingua-cat-go.localhost/analytics/`

4. Получить JWT и начать отправлять запросы через Traefik.

Получение JWT-токена для встроенного dummy-user:
```bash
curl -X POST --location "http://keycloak.localhost/realms/lingua-cat-go/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H "Accept: application/json" \
    -d 'grant_type=password&scope=openid&client_id=lingua-cat-go-dev&client_secret=GatPbS9gsEfplvCpiNitwBdmIRc0QqyQ&username=dummy-user&password=password'
```

Получить из ответа `access_token` и использовать его в заголовке `Authorization: Bearer <token>`.

---

## Запуск тестов

```bash
make test-all
```

## Запуск A/B теста

```bash
make ab-test
```

Это запустит 10 потоков, которые будут отправлять 100 сценариев прохождения упражнений к API. Результаты теста будут сохранены в `project/ab/ab_test_results.txt`.

Сценарий:
1. Получение access_token
2. Запуск 100 сценариев прохождения упражнений
2.1. Создание упражнения с одним заданием
2.2. Выбор правильного слова в задании
