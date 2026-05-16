# ClusterGuard

Desktop-приложение для мониторинга и управления кластером серверов.

**Стек:** Go · Wails v2 · LevelDB · TypeScript

## Возможности

- CRUD серверов (только с **агентом ClusterGuard** на каждой ноде)
- Мониторинг через агент: `/health`, `/metrics` — CPU, RAM, диск (`agent/`)
- Графики загрузки и пороги с алертами
- Тёмная и светлая тема
- Группировка серверов в кластеры
- Локальное хранение в LevelDB (`~/.config/clusterguard/data/`)

## Архитектура (SOLID)

```
internal/
  domain/          # Сущности + интерфейсы (порты) — DIP
  repository/      # LevelDB — реализация репозиториев
  agentclient/     # HTTP-клиент агента (health + metrics)
  service/         # Бизнес-логика — SRP
  app/             # DI-контейнер
app.go             # Wails adapter (тонкий фасад)
```

| Принцип | Реализация |
|---------|------------|
| **S** | `ServerService`, `ClusterService`, `MonitorService` — отдельные обязанности |
| **O** | Новые checkers через `HealthChecker` без изменения сервисов |
| **L** | Любая реализация `ServerRepository` взаимозаменяема |
| **I** | Узкие интерфейсы: `ServerRepository`, `HealthChecker`, `HealthCheckerRegistry` |
| **D** | Сервисы зависят от `domain.*`, не от LevelDB |

## Git Flow

| Ветка | Назначение |
|-------|------------|
| `main` | Стабильные релизы (`v1.0.0`, `v1.0.1`) |
| `develop` | Интеграция фич |
| `feature/*` | Новая функциональность → merge в `develop` |
| `release/*` | Подготовка релиза → `main` + `develop` |
| `hotfix/*` | Срочные правки от `main` → `main` + `develop` |

Текущая рабочая ветка: **`develop`**. Пересборка истории: `scripts/rebuild-gitflow.sh`.

## Запуск

```bash
# Dev mode (hot reload)
wails dev

# Production build
wails build
```

## Требования

- Go 1.22+
- Node.js 18+
- Wails CLI v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- Linux: `libgtk-3-dev`, `libwebkit2gtk-4.1-dev`
