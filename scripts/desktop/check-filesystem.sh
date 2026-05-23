#!/usr/bin/env bash
# TC-FS: работа с файловой системой (каталог данных LevelDB)
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

section() { echo ""; echo "=== $1 ==="; }
ok() { echo "[OK] $*"; }
fail() { echo "[FAIL] $*"; exit 1; }

DATA_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/clusterguard/data"
DB_DIR="$DATA_DIR/db"

section "TC-FS-01 Автотесты репозитория (LevelDB на диске)"
(cd "$ROOT" && go test ./internal/repository/leveldb/... -v -count=1) || fail "go test leveldb"

section "TC-FS-02 Каталог данных приложения"
info_path() {
  if [[ -d "$DATA_DIR" ]]; then
    ok "Каталог существует: $DATA_DIR"
    du -sh "$DATA_DIR" 2>/dev/null || true
    ls -la "$DATA_DIR" 2>/dev/null || true
  else
    echo "[INFO] Каталог ещё не создан (запустите приложение и добавьте сервер): $DATA_DIR"
  fi
}
info_path

section "TC-FS-03 Файлы LevelDB после эксплуатации"
if [[ -d "$DB_DIR" ]]; then
  count=$(find "$DB_DIR" -type f 2>/dev/null | wc -l)
  if [[ "$count" -gt 0 ]]; then
    ok "В $DB_DIR найдено файлов: $count"
    find "$DB_DIR" -maxdepth 1 -type f | head -5
  else
    echo "[INFO] Папка db пуста"
  fi
else
  echo "[INFO] $DB_DIR отсутствует до первого запуска"
fi

section "TC-FS-04 Резервное копирование (учебный сценарий)"
BACKUP="/tmp/clusterguard-backup-test"
rm -rf "$BACKUP"
if [[ -d "$DB_DIR" ]]; then
  mkdir -p "$BACKUP"
  cp -a "$DB_DIR" "$BACKUP/db"
  ok "Копия БД: $BACKUP/db"
  rm -rf "$BACKUP"
else
  echo "[SKIP] Нет данных для копирования — выполните после работы с UI"
fi

section "TC-FS-05 Права доступа"
if [[ -d "$DATA_DIR" ]]; then
  perm=$(stat -c '%a' "$DATA_DIR" 2>/dev/null || stat -f '%OLp' "$DATA_DIR")
  ok "Права на $DATA_DIR: $perm (ожидается доступ только пользователю)"
fi

section "Итог"
ok "Проверка файловой системы завершена"
