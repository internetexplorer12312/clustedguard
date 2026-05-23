# Docker — ClusterGuard Agent

Контейнеризация обеспечивает воспроизводимое тестовое окружение (§3.7.4.3 отчёта).

## Быстрый старт

```bash
cd /home/lesha/clusterguard
docker compose up -d --build
```

Три агента слушают порты **9101**, **9102**, **9103**. Токен по умолчанию: `dev-token-change-me`.

## Проверка

```bash
curl -s http://localhost:9101/health
curl -s -H "X-ClusterGuard-Token: dev-token-change-me" http://localhost:9101/metrics | head
```

## Подключение из приложения

| Сервер | Хост | Порт | Токен |
|--------|------|------|-------|
| node-1 | 127.0.0.1 | 9101 | dev-token-change-me |
| node-2 | 127.0.0.1 | 9102 | dev-token-change-me |
| node-3 | 127.0.0.1 | 9103 | dev-token-change-me |

## Файлы

- `Dockerfile.agent` — образ агента
- `docker-compose.yml` — три ноды для тестового кластера
- `Dockerfile.dev` — вспомогательная сборка (frontend + agent binary)

Десктопное приложение (Wails) запускается на хосте: `wails dev -tags=webkit2_41`.
