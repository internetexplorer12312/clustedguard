# ClusterGuard

Desktop-приложение для мониторинга кластера Linux-серверов. Работает **только** с предустановленным [агентом ClusterGuard](agent/) на каждой ноде: проверка доступности, сбор метрик CPU/RAM/диск, графики и алерты по порогам.

**Стек:** Go · Wails v2 · LevelDB · TypeScript · Vite

**Версия:** см. файл `VERSION` (релизы помечены тегами `v1.0.x`).

---

## Содержание

- [Возможности](#возможности)
- [Требования](#требования)
- [Быстрый старт](#быстрый-старт)
- [Установка агента на сервер](#установка-агента-на-сервер)
- [Архитектура](#архитектура)
- [Хранение данных](#хранение-данных)
- [HTTP API агента](#http-api-агента)
- [API desktop-приложения (Wails)](#api-desktop-приложения-wails)
- [События UI](#события-ui)
- [Справочник полей и констант](#справочник-полей-и-констант)
- [Фоновый сбор и алерты](#фоновый-сбор-и-алерты)
- [Структура проекта](#структура-проекта)
- [Сборка и разработка](#сборка-и-разработка)
- [Модульное тестирование](#модульное-тестирование)
- [Docker (тестовый кластер)](#docker-тестовый-кластер)
- [Git Flow](#git-flow)

---

## Возможности

| Функция | Описание |
|---------|----------|
| Серверы | CRUD, привязка к кластеру, роли `master` / `worker` / `any` |
| Мониторинг | Проверка `GET /health` агента, задержка в мс |
| Метрики | CPU, ОЗУ, диск с агента; история для графиков |
| Алерты | Пороги по ЦП/ОЗУ/диску; push в UI; список в приложении |
| Кластеры | Группы серверов, сводка online/offline |
| UI | Обзор, серверы, кластеры, алерты, экран метрик с графиками |
| Тема | Светлая / тёмная (localStorage) |
| Данные | Локально в LevelDB, без внешней БД |

---

## Требования

### Для сборки desktop-приложения

- Go 1.22+
- Node.js 18+
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation): `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Linux: `libgtk-3-dev`, `libwebkit2gtk-4.1-dev`

### На каждом мониторируемом сервере

- Запущенный бинарник `clusterguard-agent` (порт по умолчанию **9100**)
- Сетевой доступ с машины, где работает ClusterGuard, до `host:agentPort`

---

## Быстрый старт

### 1. Сборка и запуск приложения

```bash
git clone <url-репозитория>
cd clusterguard

# Фронтенд (для production-сборки Wails)
cd frontend && npm install && npm run build && cd ..

# Режим разработки
wails dev -tags=webkit2_41

# Production-бинарник
wails build -tags=webkit2_41
# Результат: build/bin/clusterguard
```

### 2. Добавление сервера в UI

1. Вкладка **Серверы** → **Добавить сервер**
2. Укажите **имя**, **хост** (IP или DNS), **роль**, при необходимости **кластер**
3. В секции **Агент ClusterGuard**: **порт** (9100) и **токен** (как на сервере)
4. Задайте **пороги** ЦП / ОЗУ / диска в процентах (по умолчанию 90)
5. Сохраните

Приложение начнёт опрашивать агент каждые **30 секунд** (фоновый коллектор).

---

## Установка агента на сервер

```bash
cd agent
go build -o clusterguard-agent .

# Запуск
export CLUSTERGUARD_TOKEN="ваш-секретный-токен"
./clusterguard-agent -addr :9100
```

### systemd

```bash
sudo cp clusterguard-agent /usr/local/bin/
sudo cp clusterguard-agent.service /etc/systemd/system/
# В unit-файле или Environment= задайте CLUSTERGUARD_TOKEN
sudo systemctl daemon-reload
sudo systemctl enable --now clusterguard-agent
```

Откройте порт **9100** в firewall. Тот же токен укажите в форме сервера в ClusterGuard.

Подробнее: [agent/README.md](agent/README.md).

---

## Архитектура

Слои построены по **SOLID**: домен не зависит от LevelDB и Wails; UI вызывает тонкий фасад `app.go`.

```
┌─────────────────────────────────────────────────────────┐
│  frontend (TypeScript)  ←→  Wails bindings (app.go)      │
└────────────────────────────┬────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────┐
│  internal/service   Server, Cluster, Monitor, Metrics,   │
│                     Alert, Collector                     │
└────────────────────────────┬────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────┐
│  internal/repository/leveldb                               │
│  internal/agentclient  →  HTTP к агенту на нодах         │
└───────────────────────────────────────────────────────────┘

На сервере:  agent (HTTP :9100)  →  gopsutil
```

| Принцип | Где |
|---------|-----|
| **S** | Отдельные сервисы: серверы, кластеры, мониторинг, метрики, алерты |
| **O** | Новые типы проверок через `domain.HealthChecker` (пакет `health/`) |
| **L** | Взаимозаменяемые реализации репозиториев |
| **I** | Узкие порты: `ServerRepository`, `MetricsFetcher`, `AlertNotifier` |
| **D** | Сервисы зависят от интерфейсов в `domain`, не от LevelDB |

---

## Хранение данных

| Путь | Содержимое |
|------|------------|
| Linux | `~/.config/clusterguard/data/db/` |
| Fallback | `~/.clusterguard/data/db/` |

В LevelDB:

| Префикс ключа | Данные |
|---------------|--------|
| `server:` | Серверы и кэш последних метрик |
| `cluster:` | Кластеры |
| `metrics:` | Временные ряды для графиков |
| `alert:` | Записи алертов |

---

## HTTP API агента

Базовый URL: `http://<host>:<agentPort>` (по умолчанию порт **9100**).

### `GET /health`

Проверка доступности агента. Используется desktop-приложением для статуса сервера.

**Заголовки:** при настроенном токене на агенте — `X-ClusterGuard-Token: <token>` (опционально для `/health` в текущей реализации агента токен не проверяется).

**Ответ:** `200 OK`, тело `ok\n`

**Пример:**

```bash
curl -s http://192.168.1.10:9100/health
```

---

### `GET /metrics`

Текущие метрики хоста (gopsutil).

**Заголовки:**

| Заголовок | Обязательный | Описание |
|-----------|--------------|----------|
| `X-ClusterGuard-Token` | Да, если при запуске агента задан `-token` или `CLUSTERGUARD_TOKEN` | Секрет |

**Ответ:** `200 OK`, `Content-Type: application/json`

**Тело (JSON):**

```json
{
  "timestamp": 1715798400,
  "cpuPercent": 12.5,
  "memoryUsedPercent": 64.2,
  "memoryAvailableBytes": 2147483648,
  "memoryTotalBytes": 8589934592,
  "diskUsedPercent": 71.0,
  "diskFreeBytes": 50000000000,
  "diskTotalBytes": 200000000000
}
```

| Поле | Тип | Описание |
|------|-----|----------|
| `timestamp` | int64 | Unix time сбора |
| `cpuPercent` | float | Загрузка CPU, % |
| `memoryUsedPercent` | float | Использование ОЗУ, % |
| `memoryAvailableBytes` | uint64 | Свободная память, байт |
| `memoryTotalBytes` | uint64 | Всего памяти, байт |
| `diskUsedPercent` | float | Занято диска `/`, % |
| `diskFreeBytes` | uint64 | Свободно на `/`, байт |
| `diskTotalBytes` | uint64 | Размер `/`, байт |

**Ошибки:**

| Код | Причина |
|-----|---------|
| `401` | Неверный или отсутствующий токен |
| `500` | Ошибка сбора метрик на хосте |

**Пример:**

```bash
curl -s -H "X-ClusterGuard-Token: your-secret-token" \
  http://192.168.1.10:9100/metrics | jq .
```

---

## API desktop-приложения (Wails)

Методы Go-структуры `App` доступны из фронтенда через сгенерированные биндинги:

```typescript
import {
  ListServers, CreateServer, UpdateServer, DeleteServer,
  CheckServer, CheckAllServers,
  ListClusters, CreateCluster, UpdateCluster, DeleteCluster, CheckCluster,
  GetDashboardStats, GetMetricsHistory, CollectServerMetrics,
  ListAlerts, MarkAlertRead, DeleteAlert,
} from '../wailsjs/go/main/App';
import type { main } from '../wailsjs/go/models';
```

После изменения `app.go` перегенерируйте биндинги:

```bash
wails generate module
```

### Серверы

#### `ListServers() → ServerDTO[]`

Список всех серверов с последними метриками и статусом.

#### `CreateServer(input: ServerInputDTO) → ServerDTO`

Создание сервера. Поле `useAgent` на бэкенде всегда приводится к `true`. Рекомендуется передавать `checkType: "agent"`, `port` = порт агента.

**Пример input:**

```json
{
  "name": "web-01",
  "host": "192.168.1.10",
  "port": 9100,
  "role": "worker",
  "clusterId": "",
  "agentPort": 9100,
  "agentToken": "your-secret-token",
  "cpuThreshold": 90,
  "memThreshold": 90,
  "diskThreshold": 90,
  "notes": ""
}
```

#### `UpdateServer(input: ServerInputDTO) → ServerDTO`

Обновление. В `input.id` — UUID существующего сервера.

#### `DeleteServer(id: string) → void`

Удаление сервера и истории его метрик.

#### `CheckServer(id: string) → ServerDTO`

Ручная проверка: `GET /health` агента + немедленный сбор метрик.

#### `CheckAllServers() → ServerDTO[]`

Проверка всех серверов и сбор метрик с каждого.

#### `CollectServerMetrics(serverID: string) → ServerDTO`

Только сбор метрик с агента (без отдельного health, если не вызывался ранее).

#### `GetMetricsHistory(serverID: string, limit: number) → MetricsSampleDTO[]`

История для графиков. `limit <= 0` → **60** точек.

---

### Кластеры

#### `ListClusters() → ClusterSummaryDTO[]`

Список кластеров с полями `totalServers`, `onlineCount`, `offlineCount`.

#### `CreateCluster(input: ClusterInputDTO) → ClusterDTO`

#### `UpdateCluster(input: ClusterInputDTO) → ClusterDTO`

#### `DeleteCluster(id: string) → void`

#### `CheckCluster(clusterID: string) → ServerDTO[]`

Проверка и сбор метрик для всех серверов кластера.

**ClusterInputDTO:**

```json
{
  "id": "",
  "name": "production",
  "description": "Продакшен-кластер",
  "serverIds": []
}
```

Привязка сервера к кластеру — через поле `clusterId` в `ServerInputDTO`.

---

### Алерты и обзор

#### `ListAlerts(limit: number) → AlertDTO[]`

Последние алерты, сортировка по времени (новые первые). `limit <= 0` → **50**.

#### `MarkAlertRead(id: string) → void`

#### `DeleteAlert(id: string) → void`

#### `GetDashboardStats() → DashboardStatsDTO`

```json
{
  "totalServers": 5,
  "onlineServers": 4,
  "totalClusters": 2,
  "unreadAlerts": 1
}
```

---

### Типы данных (DTO)

#### `ServerDTO`

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | string | UUID |
| `name` | string | Отображаемое имя |
| `host` | string | Хост агента |
| `port` | number | Порт (синхронизируется с `agentPort`) |
| `role` | string | `master` \| `worker` \| `any` |
| `status` | string | `unknown` \| `online` \| `offline` \| `degraded` |
| `tags` | string[] | Теги |
| `checkType` | string | `agent` |
| `checkPath` | string | Не используется для агента |
| `lastCheck` | number | Unix time последней проверки |
| `latencyMs` | number | Задержка health-запроса, мс |
| `clusterId` | string | ID кластера или пусто |
| `notes` | string | Заметки |
| `useAgent` | boolean | Всегда `true` |
| `agentPort` | number | Порт HTTP агента |
| `agentToken` | string | Токен |
| `cpuThreshold` | number | Порог ЦП, % |
| `memThreshold` | number | Порог ОЗУ, % |
| `diskThreshold` | number | Порог диска, % |
| `cpuPercent` | number | Последнее значение ЦП |
| `memPercent` | number | Последнее значение ОЗУ |
| `diskPercent` | number | Последнее значение диска |
| `memAvailBytes` | number | Свободная память |
| `diskFreeBytes` | number | Свободное место на диске |

#### `AlertDTO`

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | string | UUID алерта |
| `serverId` | string | ID сервера |
| `serverName` | string | Имя сервера |
| `kind` | string | `cpu` \| `memory` \| `disk` |
| `value` | number | Фактическое значение, % |
| `threshold` | number | Порог, % |
| `message` | string | Текст на русском |
| `createdAt` | number | Unix time |
| `read` | boolean | Прочитан |

#### `MetricsSampleDTO`

| Поле | Тип |
|------|-----|
| `serverId` | string |
| `timestamp` | number |
| `cpuPercent` | number |
| `memPercent` | number |
| `diskPercent` | number |
| `memAvailBytes` | number |
| `diskFreeBytes` | number |

---

## События UI

Бэкенд отправляет события в WebView через Wails Runtime.

### `alert`

Срабатывает при новом уведомлении о превышении порога (с учётом интервала повторов).

**Подписка (TypeScript):**

```typescript
import { EventsOn } from '../wailsjs/runtime/runtime';
import type { main } from '../wailsjs/go/models';

EventsOn('alert', (alert: main.AlertDTO) => {
  console.log('Новый алерт:', alert.message);
});
```

**Payload:** объект `AlertDTO` (см. выше).

---

## Справочник полей и констант

### Статус сервера (`status`)

| Значение | Условие |
|----------|---------|
| `online` | Агент ответил, задержка ≤ 500 мс |
| `degraded` | Агент ответил, задержка > 500 мс |
| `offline` | Агент не ответил |
| `unknown` | Ещё не проверялся |

### Роль (`role`)

`master`, `worker`, `any`

### Вид алерта (`kind`)

`cpu`, `memory`, `disk`

---

## Фоновый сбор и алерты

| Параметр | Значение |
|----------|----------|
| Интервал коллектора | **30 с** |
| Health | `GET http://host:agentPort/health` |
| Метрики | `GET http://host:agentPort/metrics` + заголовок токена |
| Toast в UI | Повтор не чаще **30 с** на пару сервер+метрика, пока порог превышен |
| Запись в БД | Не чаще **5 мин** на ту же пару |
| Сброс | Когда метрика снова ниже порога |

Сообщение алерта (пример): `web-01: ЦП 92.3% (порог 90.0%)`

---

## Структура проекта

```
clusterguard/
├── main.go                 # Точка входа Wails
├── app.go                  # API для фронтенда, DTO, EventsEmit
├── wails.json
├── agent/                  # Демон на серверах
│   ├── main.go
│   └── clusterguard-agent.service
├── pkg/agentapi/           # Общие типы JSON метрик
├── internal/
│   ├── domain/             # Модели и интерфейсы
│   ├── repository/leveldb/
│   ├── agentclient/        # HTTP к агенту
│   ├── health/             # TCP/HTTP checkers (расширяемость)
│   ├── service/
│   └── app/                # DI-контейнер
├── frontend/
│   ├── src/
│   │   ├── main.ts         # UI
│   │   ├── i18n.ts         # Русские строки
│   │   ├── charts.ts
│   │   └── theme.ts
│   └── wailsjs/            # Сгенерированные биндинги
└── scripts/
    └── rebuild-gitflow.sh
```

---

## Сборка и разработка

```bash
# Компиляция Go
go build -o /dev/null .

# Модульные тесты
go test ./... -cover

# Только агент
cd agent && go build -o clusterguard-agent .

# Фронтенд
cd frontend && npm run build

# Dev (hot reload UI + Go)
wails dev -tags=webkit2_41

# Production
wails build -tags=webkit2_41
```

### Типичные проблемы

| Симптом | Решение |
|---------|---------|
| Сервер `offline` | Проверьте firewall, `curl host:9100/health`, токен |
| Нет метрик | Токен в UI = `CLUSTERGUARD_TOKEN` на сервере |
| Пустые графики | Подождите 1–2 цикла коллектора или **Обновить** на экране метрик |
| После правки `app.go` | `wails generate module` и пересборка фронтенда |

---

## Модульное тестирование

```bash
go test ./... -cover
```

Подробнее: [docs/TESTING.md](docs/TESTING.md) (интеграция с агентом, ручные сценарии UI).

Покрыты `AlertService` (пороги алертов) и типы `pkg/agentapi`.

---

## Docker (тестовый кластер)

Для воспроизводимого окружения без установки агента на VM:

```bash
docker compose up -d --build
curl -s http://localhost:9101/health
curl -s -H "X-ClusterGuard-Token: dev-token-change-me" http://localhost:9101/metrics
```

Поднимаются **три** контейнера-агента на портах **9101**, **9102**, **9103**. Токен по умолчанию: `dev-token-change-me`.

Полная инструкция: [docs/DOCKER.md](docs/DOCKER.md).

| Файл | Назначение |
|------|------------|
| `Dockerfile.agent` | Образ агента |
| `docker-compose.yml` | Три ноды для тестов |
| `Dockerfile.dev` | Вспомогательная сборка (agent + frontend dist) |

Десктопное приложение (Wails) по-прежнему запускается на хосте.

---

## Git Flow

| Ветка | Назначение |
|-------|------------|
| `main` | Стабильные релизы (`v1.0.0`, `v1.0.1`, …) |
| `develop` | Интеграция фич |
| `feature/*` | Новая функциональность → merge в `develop` |
| `release/*` | Подготовка релиза → `main` + `develop` |
| `hotfix/*` | Срочные правки от `main` → `main` + `develop` |

Рабочая ветка по умолчанию: **`develop`**.

```bash
# Новая фича
git checkout develop
git checkout -b feature/my-feature
# ... коммиты ...
git checkout develop
git merge --no-ff feature/my-feature

# Релиз
git checkout -b release/1.1.0 develop
# ... версия, правки ...
git checkout main && git merge --no-ff release/1.1.0
git tag -a v1.1.0 -m "ClusterGuard 1.1.0"
git checkout develop && git merge --no-ff release/1.1.0
```

Пересборка учебной истории коммитов: `scripts/rebuild-gitflow.sh`.

---

## Лицензия

Проект распространяется под лицензией [MIT](LICENSE).
