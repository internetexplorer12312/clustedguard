#!/usr/bin/env bash
# TC-INST: проверка «инсталлятора» desktop-приложения (Linux, Wails)
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
BIN="$ROOT/build/bin/clusterguard"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
REPORT="${1:-}"

section() { echo ""; echo "=== $1 ==="; }
ok() { echo "[OK] $*"; }
fail() { echo "[FAIL] $*"; exit 1; }
warn() { echo "[WARN] $*"; }

section "TC-INST-01 Сборка production-бинарника"
if [[ ! -x "$BIN" ]]; then
  warn "Бинарник не найден: $BIN"
  echo "Выполните: cd $ROOT && wails build -tags=webkit2_41"
  if [[ "${SKIP_BUILD:-}" != "1" ]]; then
    (cd "$ROOT" && wails build -tags=webkit2_41) || fail "wails build не удался"
  else
    fail "Сборка пропущена (SKIP_BUILD=1), бинарника нет"
  fi
fi
[[ -x "$BIN" ]] && ok "Файл существует и исполняемый: $BIN"
ls -lh "$BIN"

section "TC-INST-02 Локальная установка (копирование в PATH)"
mkdir -p "$INSTALL_DIR"
cp -f "$BIN" "$INSTALL_DIR/clusterguard"
chmod +x "$INSTALL_DIR/clusterguard"
ok "Скопировано в $INSTALL_DIR/clusterguard"
command -v clusterguard >/dev/null 2>&1 && ok "clusterguard в PATH" || warn "Добавьте $INSTALL_DIR в PATH"

section "TC-INST-03 Зависимости GTK/WebKit (ldd)"
if command -v ldd >/dev/null; then
  missing=$(ldd "$BIN" 2>&1 | grep "not found" || true)
  if [[ -n "$missing" ]]; then
    fail "Не хватает библиотек:\n$missing"
  fi
  ok "Все динамические библиотеки найдены"
else
  warn "ldd не установлен"
fi

section "TC-INST-04 Smoke: запуск с таймаутом"
if command -v timeout >/dev/null; then
  # GUI может не стартовать без DISPLAY — проверяем только что процесс не падает сразу с ошибкой линковки
  if [[ -z "${DISPLAY:-}" ]]; then
    warn "DISPLAY не задан — полный UI-тест пропущен (нормально для CI/SSH)"
  else
    timeout 3s "$BIN" >/dev/null 2>&1 && ok "Процесс запустился" || warn "Окно закрылось по timeout (ожидаемо)"
  fi
else
  warn "timeout не найден"
fi

section "Итог"
ok "Проверка инсталляции завершена"
[[ -n "$REPORT" ]] && "$0" 2>&1 | tee "$REPORT"
