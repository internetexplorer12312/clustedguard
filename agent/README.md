# ClusterGuard Agent

Устанавливается на каждый мониторируемый сервер. Отдаёт метрики CPU, RAM и диска по HTTP.

## Сборка

```bash
cd agent
go build -o clusterguard-agent .
```

## Запуск

```bash
./clusterguard-agent -addr :9100 -token "your-secret-token"
```

Переменная окружения: `CLUSTERGUARD_TOKEN`.

## Endpoints

| Path | Описание |
|------|----------|
| `GET /health` | Проверка живости |
| `GET /metrics` | JSON с метриками (заголовок `X-ClusterGuard-Token` если задан token) |

## systemd

```bash
sudo cp clusterguard-agent /usr/local/bin/
sudo cp clusterguard-agent.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now clusterguard-agent
```

Откройте порт `9100` в firewall.

## В ClusterGuard

При добавлении сервера укажите:
- **Agent port**: `9100`
- **Agent token**: тот же секрет
- Пороги CPU / Memory / Disk (%)
