#!/usr/bin/env bash
# Полный прогон тестов desktop: ОС, ФС, инсталлятор
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
OUT="${1:-}"

run() {
  echo "########################################"
  bash "$1"
}

if [[ -n "$OUT" ]]; then
  exec > >(tee "$OUT") 2>&1
fi

run "$DIR/check-os-compat.sh"
run "$DIR/check-filesystem.sh"
run "$DIR/check-install.sh"

echo ""
echo "=== Все desktop-проверки выполнены ==="
