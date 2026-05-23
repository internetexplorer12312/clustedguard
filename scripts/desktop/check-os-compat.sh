#!/usr/bin/env bash
# TC-OS: совместимость с версиями ОС (Linux desktop)
set -euo pipefail

section() { echo ""; echo "=== $1 ==="; }
ok() { echo "[OK] $*"; }
warn() { echo "[WARN] $*"; }
info() { echo "     $*"; }

section "TC-OS-01 Идентификация ОС"
uname -a
if [[ -f /etc/os-release ]]; then
  # shellcheck source=/dev/null
  . /etc/os-release
  info "ID=$ID VERSION_ID=${VERSION_ID:-?} PRETTY_NAME=$PRETTY_NAME"
  ok "Дистрибутив определён"
else
  warn "/etc/os-release отсутствует"
fi

section "TC-OS-02 Матрица поддержки (документированная)"
cat <<'MATRIX'
| ОС              | Версия      | Статус        |
|-----------------|-------------|---------------|
| Ubuntu          | 22.04 LTS   | Поддерживается|
| Ubuntu          | 24.04 LTS   | Поддерживается|
| Debian          | 12          | Поддерживается|
| Fedora          | 38+         | Тестируется   |
| Arch Linux      | rolling     | Тестируется   |
MATRIX

section "TC-OS-03 Зависимости сборки/рантайма"
pkgs=(libgtk-3-0 libwebkit2gtk-4.1-0)
for p in "${pkgs[@]}"; do
  if dpkg -l "$p" &>/dev/null 2>&1; then
    ver=$(dpkg -l "$p" | awk '/^ii/{print $3}')
    ok "$p $ver"
  elif rpm -q gtk3 webkit2gtk4.1 &>/dev/null 2>&1; then
    ok "RPM: gtk3 / webkit2gtk4.1 установлены"
    break
  else
    warn "Пакет $p не найден через dpkg (проверьте вручную для вашего дистрибутива)"
  fi
done

section "TC-OS-04 Go и Wails"
command -v go >/dev/null && ok "go $(go version | awk '{print $3}')" || warn "go не в PATH"
command -v wails >/dev/null && ok "wails $(wails version 2>/dev/null | head -1)" || warn "wails CLI не найден"

section "TC-OS-05 Архитектура процессора"
arch=$(uname -m)
case "$arch" in
  x86_64|amd64) ok "amd64 — целевая платформа сборки" ;;
  aarch64|arm64) ok "arm64 — поддерживается Go, проверьте WebKit пакеты" ;;
  *) warn "Архитектура $arch — требуется ручная проверка" ;;
esac

section "Итог"
ok "Проверка совместимости ОС завершена"
