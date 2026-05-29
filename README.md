# org-structure-api
[![CI](https://github.com/boldlogic/org-structure-api/actions/workflows/go.yml/badge.svg)](https://github.com/boldlogic/org-structure-api/actions/workflows/go.yml)
[![Go](https://img.shields.io/badge/Go-1.26-blue?logo=go&logoColor=white)](https://go.dev/)

REST-API сервис организационной структуры компании. Позволяет управлять иерархией подразделений и создавать сотрудников внутри подразделений.

## 1. Стек
**Стек**: Go, PostgreSQL, Docker.
- Слоистая архитектура через интерфейсы
- HTTP через `net/http`
- Тесты через `testify`, репозиторный слой замокан
- Работа с БД через raw sql поверх `GORM`
- Миграции через отдельный сервис миграции на базе `goose`
- Повторяемый код - с моего библиотечного репозитория `packages`

## 2. Сверх требований ТЗ
Базовое задание: 5 эндпоинтов и стек из [docs/0. ТЗ Go - API организационной структуры.md](docs/0.%20%D0%A2%D0%97%20Go%20-%20API%20%D0%BE%D1%80%D0%B3%D0%B0%D0%BD%D0%B8%D0%B7%D0%B0%D1%86%D0%B8%D0%BE%D0%BD%D0%BD%D0%BE%D0%B9%20%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D1%8B.md). Ниже: расширения.

### 2.1. Эксплуатация
- GET /health, HEALTHCHECK в Docker, graceful shutdown
- CI (GitHub Actions), отдельный migrate-сервис в docker-compose
- YAML-конфиг с pool/timeouts, middleware logging и panic recover
- Схема БД org, общая библиотека packages (HTTP, валидация, конфиг)

### 2.2. API и ошибки
- RFC 7807 Problem Details; коды 422, 415, 413 (помимо 400/404/409 из ТЗ)
- Строгий PATCH: различие null/пусто, неизвестные поля
- Валидация id в диапазоне 1..2147483647

### 2.3. Бизнес-логика
- DELETE mode=reassign: перенос сотрудников и прямых дочерних подразделений; запрет reassign на себя и на прямого ребёнка; транзакция
- POST department: INSERT … WHERE NOT EXISTS и повтор при гонке уникальности
- POST employee: FK 23503 → 404 при гонке
- PATCH: проверка цикла через recursive CTE; raw SQL поверх GORM

### 2.4. Документация
- Алгоритмы эндпоинтов, SVG-диаграммы, Postman-коллекция в docs/
- Unit-тесты service с mock repository и decode DTO через testify

## 3. Производительность
Помимо базового CRUD проработаны сценарии высокой нагрузки: все эндпоинты проверены k6-тестами.

Метрики ниже сняты локально на железе: Intel Core i5-11400H (2.7 GHz), 16 GB RAM, Windows 10 x64. Сервис, БД и генератор нагрузки на одной машине.

Основные изменения на критическом пути:
- проверка уникальности названия у родителя: индексное сравнение через COALESCE(parent_id, 0) вместо полного прохода по таблице;
- чтение дерева: индекс (parent_id, id), для сотрудников: (department_id);
- после проверки существования подразделения независимые чтения children и employees выполняются параллельно (errgroup).

| Endpoint                         |     Факт |   P95 |
| -------------------------------- | -------: | ----: |
| POST /departments                | 4859 RPS | 38 ms |
| GET /departments/{id}            | 4970 RPS | 63 ms |
| POST /departments/{id}/employees | 4570 RPS |  8 ms |
| PATCH /departments/{id}          | 4583 RPS |  1 ms |
| DELETE /departments/{id}         | 4474 RPS | 11 ms |

Для POST /departments и GET /departments/{id} до оптимизации: 313 RPS (P95 837 ms) и 46 RPS (P95 281 ms) соответственно.

## 4. Запуск
```
docker compose up -d --build
```

## 5. API
Общие JSON-структуры, схема БД: [docs/1-data-structures.md](docs/1-data-structures.md).

Контракт эндпоинта, схема алгоритма и сценарии - по ссылкам.

| Endpoint | Документация |
| -------- | ------------ |
| POST /departments | [docs/2-post-department.md](docs/2-post-department.md) |
| POST /departments/{id}/employees | [docs/3-post-employee.md](docs/3-post-employee.md) |
| GET /departments/{id} | [docs/4-get-department.md](docs/4-get-department.md) |
| PATCH /departments/{id} | [docs/5-patch-department.md](docs/5-patch-department.md) |
| DELETE /departments/{id} | [docs/6-delete-department.md](docs/6-delete-department.md) |
