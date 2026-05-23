# Тестирование ClusterGuard

## Модульные тесты (Go)

Пакет `testing` из стандартной библиотеки, без внешних фреймворков.

### Запуск

```bash
# Все тесты проекта
go test ./...

# С покрытием
go test ./... -cover

# Только сервисы и API агента
go test ./internal/service/... ./pkg/agentapi/... -cover
```

### Что покрыто

| Пакет | Файл | Проверка |
|-------|------|----------|
| `internal/service` | `alert_service_test.go` | Нет алерта ниже порога; создание алерта при CPU > порога |
| `pkg/agentapi` | `validate_test.go` | Структура JSON-ответа `/metrics` |

Репозитории в unit-тестах подменяются in-memory заглушками (`memAlertRepo`).

## Интеграция с агентом (ручная)

```bash
cd agent
go build -o clusterguard-agent .
CLUSTERGUARD_TOKEN=my-secret ./clusterguard-agent -addr :9100
```

В другом терминале:

```bash
curl -s http://127.0.0.1:9100/health
curl -s -H "X-ClusterGuard-Token: my-secret" http://127.0.0.1:9100/metrics
# Неверный токен — ожидается HTTP 401
curl -s -o /dev/null -w "%{http_code}\n" -H "X-ClusterGuard-Token: wrong" http://127.0.0.1:9100/metrics
```

## Тестовый кластер в Docker

См. [DOCKER.md](DOCKER.md): три агента на портах 9101–9103.

После `docker compose up -d --build` добавьте в desktop-приложение серверы `127.0.0.1:9101`, `9102`, `9103` с токеном `dev-token-change-me`.

## Системное тестирование (UI)

1. Запустите агент(ы) или Docker Compose.
2. `wails dev -tags=webkit2_41`
3. Проверьте сценарии: Обзор → Серверы → добавление → метрики → алерты → кластеры → смена темы.

Интервал фонового опроса — **30 секунд** (`internal/app/container.go`).
