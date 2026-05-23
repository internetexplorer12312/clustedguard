# Тестирование desktop-приложения ClusterGuard

Раздел соответствует требованиям ВКР: **инсталлятор**, **совместимость с ОС**, **файловая система** (LevelDB в `~/.config/clusterguard/data`).

## Быстрый прогон

```bash
chmod +x scripts/desktop/*.sh
./scripts/desktop/run-all-desktop-tests.sh
```

Сохранить отчёт в файл:

```bash
./scripts/desktop/run-all-desktop-tests.sh ~/clusterguard-desktop-test.log
```

## 1. Тестирование инсталлятора (Linux)

ClusterGuard распространяется как **один исполняемый файл** после `wails build` (аналог portable-инсталлятора). Локальная «установка» — копирование в `~/.local/bin`.

| ID | Сценарий | Как проверить | Ожидаемый результат |
|----|----------|---------------|---------------------|
| TC-INST-01 | Сборка дистрибутива | `wails build -tags=webkit2_41` | Файл `build/bin/clusterguard` существует |
| TC-INST-02 | Установка в PATH | `scripts/desktop/check-install.sh` | Бинарник в `~/.local/bin/clusterguard` |
| TC-INST-03 | Зависимости | `ldd build/bin/clusterguard` | Нет строк `not found` |
| TC-INST-04 | Первый запуск | Запуск из меню/терминала | Окно открывается, нет ошибки GTK |

```bash
# Сборка
wails build -tags=webkit2_41

# Проверка инсталляции
./scripts/desktop/check-install.sh
```

**Примечание:** для SSH/CI без `DISPLAY` скрипт пропускает полный UI smoke и проверяет только бинарник и библиотеки.

## 2. Совместимость с версиями ОС

| ID | Сценарий | Инструмент | Ожидаемый результат |
|----|----------|------------|---------------------|
| TC-OS-01 | Версия ОС | `check-os-compat.sh` | Зафиксированы ID и VERSION_ID |
| TC-OS-02 | Матрица поддержки | `docs/DESKTOP-TESTING.md` | Ubuntu 22.04/24.04, Debian 12 |
| TC-OS-03 | GTK / WebKit | `dpkg -l libwebkit2gtk-4.1-0` | Пакеты установлены |
| TC-OS-04 | Toolchain | `go version`, `wails version` | Go 1.22+, Wails v2 |
| TC-OS-05 | Архитектура | `uname -m` | amd64 (основная), arm64 — вручную |

```bash
./scripts/desktop/check-os-compat.sh
```

Рекомендуется дополнительно прогнать приложение на **второй** машине/ВМ с другой версией Ubuntu и отметить результат в таблице отчёта.

## 3. Тестирование файловой системы

Данные хранятся в:

```
~/.config/clusterguard/data/db/   # LevelDB (серверы, метрики, алерты)
```

| ID | Сценарий | Как проверить | Ожидаемый результат |
|----|----------|---------------|---------------------|
| TC-FS-01 | Запись/чтение LevelDB | `go test ./internal/repository/leveldb/...` | Все тесты PASS |
| TC-FS-02 | Создание каталога | Запуск app + добавление сервера | Появляется `~/.config/clusterguard/data` |
| TC-FS-03 | Персистентность | Закрыть app → открыть снова | Серверы и алерты на месте |
| TC-FS-04 | Резервная копия | `cp -a data/db /backup/` | Копия не пустая |
| TC-FS-05 | Права | `stat ~/.config/clusterguard` | Доступ только у пользователя |

```bash
./scripts/desktop/check-filesystem.sh
```

### Автотесты (Go)

```bash
go test ./internal/repository/leveldb/... -v
```

Покрыто: запись JSON в LevelDB, повторное открытие БД, `ServerRepository` после перезапуска процесса.

## 4. Ручной чек-лист (для отчёта)

1. Собрать и «установить» → скриншот `build/bin` и `which clusterguard`.
2. `check-os-compat.sh` → скриншот версии ОС и WebKit.
3. Добавить сервер в UI → скриншот `ls ~/.config/clusterguard/data/db`.
4. Перезапустить приложение → сервер остался в списке.
5. Скопировать папку `db` в бэкап → восстановить → данные на месте.

## См. также

- [TESTING.md](TESTING.md) — модульные и интеграционные тесты
- [DOCKER.md](DOCKER.md) — тестовый кластер агентов
