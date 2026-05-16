#!/usr/bin/env bash
# Пересборка истории ClusterGuard по Git Flow (чистая симуляция с нуля)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SRC="/tmp/clusterguard-src-backup"
BUILD="/tmp/clusterguard-gitflow-build"

export GIT_AUTHOR_NAME="dinozavrik9911"
export GIT_COMMITTER_NAME="dinozavrik9911"
export GIT_AUTHOR_EMAIL="dinozavrik9911@mail.ru"
export GIT_COMMITTER_EMAIL="dinozavrik9911@mail.ru"
export GIT_EDITOR=true
export EDITOR=true

copy_src() {
  rm -rf "$SRC"
  mkdir -p "$SRC"
  rsync -a \
    --exclude='.git' \
    --exclude='build' \
    --exclude='node_modules' \
    --exclude='frontend/dist' \
    --exclude='frontend/package.json.md5' \
    --exclude='ClusterGuard' \
    --exclude='scripts/' \
    "$ROOT/" "$SRC/"
}

commit() {
  git add -A
  git diff --cached --quiet && return 0
  git commit -m "$1"
}

merge_feature() {
  local target
  target="$(git branch --show-current)"
  git merge --no-ff "$1" -m "Merge branch '$1' into ${target}"
}

# --- подготовка ---
[ -d "$SRC" ] || copy_src
rm -rf "$BUILD"
mkdir -p "$BUILD"
cd "$BUILD"
git init -b main
git config user.name "$GIT_AUTHOR_NAME"
git config user.email "$GIT_AUTHOR_EMAIL"

# ========== main: только старт ==========
cp "$SRC/.gitignore" .
cat > README.md <<'EOF'
# ClusterGuard

Desktop-приложение для мониторинга кластера серверов (в разработке).
EOF
commit "chore(main): инициализация репозитория"

git branch develop
git checkout develop

# ========== features ==========
git checkout -b feature/project-setup
cp "$SRC/go.mod" "$SRC/go.sum" "$SRC/wails.json" .
commit "chore: каркас Wails v2, go.mod и конфигурация сборки"
git checkout develop && merge_feature feature/project-setup

git checkout -b feature/domain-models
mkdir -p internal/domain
cp -r "$SRC/internal/domain/"* internal/domain/
commit "feat(domain): сущности сервера и кластера, интерфейсы портов"
git checkout develop && merge_feature feature/domain-models

git checkout -b feature/leveldb-repositories
mkdir -p internal/repository/leveldb
cp -r "$SRC/internal/repository/leveldb/"* internal/repository/leveldb/
commit "feat(storage): персистентность LevelDB"
git checkout develop && merge_feature feature/leveldb-repositories

git checkout -b feature/health-checkers
mkdir -p internal/health
cp -r "$SRC/internal/health/"* internal/health/
commit "feat(health): стратегии проверки TCP и HTTP"
git checkout develop && merge_feature feature/health-checkers

git checkout -b feature/business-services
mkdir -p internal/service
cp "$SRC/internal/service/server_service.go" \
   "$SRC/internal/service/cluster_service.go" \
   "$SRC/internal/service/monitor_service.go" \
   internal/service/
commit "feat(services): CRUD серверов и кластеров, сервис мониторинга"
git checkout develop && merge_feature feature/business-services

git checkout -b feature/clusterguard-agent
cp -r "$SRC/agent" "$SRC/pkg" .
mkdir -p internal/agentclient
cp "$SRC/internal/agentclient/http_fetcher.go" internal/agentclient/
commit "feat(agent): демон метрик ClusterGuard для установки на ноды"
git checkout develop && merge_feature feature/clusterguard-agent

git checkout -b feature/metrics-and-alerts
cp "$SRC/internal/service/metrics_service.go" \
   "$SRC/internal/service/alert_service.go" \
   "$SRC/internal/service/collector_service.go" \
   internal/service/
commit "feat(metrics): опрос агента, история метрик и алерты по порогам"
git checkout develop && merge_feature feature/metrics-and-alerts

git checkout -b feature/wails-application
mkdir -p internal/app
cp "$SRC/internal/app/container.go" internal/app/
cp "$SRC/app.go" "$SRC/main.go" .
commit "feat(app): DI-контейнер и фасад Wails API"
git checkout develop && merge_feature feature/wails-application

git checkout -b feature/frontend-foundation
mkdir -p frontend/src/assets/fonts frontend/src/assets/images
cp "$SRC/frontend/package.json" "$SRC/frontend/package-lock.json" \
   "$SRC/frontend/tsconfig.json" "$SRC/frontend/index.html" frontend/
cp "$SRC/frontend/src/vite-env.d.ts" "$SRC/frontend/src/style.css" frontend/src/
cp -r "$SRC/frontend/src/assets/"* frontend/src/assets/
commit "feat(ui): каркас фронтенда Vite, глобальные стили и ассеты"
git checkout develop && merge_feature feature/frontend-foundation

git checkout -b feature/frontend-i18n-theme
cp "$SRC/frontend/src/i18n.ts" "$SRC/frontend/src/theme.ts" frontend/src/
commit "feat(ui): русификация через i18n и светлая/тёмная тема"
git checkout develop && merge_feature feature/frontend-i18n-theme

git checkout -b feature/frontend-dashboard
cp "$SRC/frontend/src/charts.ts" frontend/src/
cp -r "$SRC/frontend/wailsjs" frontend/
cp "$SRC/frontend/src/main.ts" frontend/src/
cp "$SRC/frontend/src/app.css" frontend/src/
commit "feat(ui): дашборд, таблицы, графики и экран алертов"
git checkout develop && merge_feature feature/frontend-dashboard

git checkout -b feature/agent-monitoring
cp "$SRC/internal/agentclient/health.go" internal/agentclient/
cp "$SRC/internal/service/monitor_service.go" internal/service/
cp "$SRC/internal/service/server_service.go" internal/service/
cp "$SRC/internal/app/container.go" internal/app/
cp "$SRC/app.go" .
cp "$SRC/frontend/src/main.ts" "$SRC/frontend/src/i18n.ts" frontend/src/
cp "$SRC/frontend/src/app.css" frontend/src/
cp "$SRC/README.md" .
commit "feat(agent): обязательный агент ClusterGuard и проверка /health"
git checkout develop && merge_feature feature/agent-monitoring

# ========== release 1.0.0 ==========
git checkout -b release/1.0.0
echo "1.0.0" > VERSION
commit "chore(release): подготовка релиза 1.0.0"

git checkout main
merge_feature release/1.0.0
git tag -a v1.0.0 -m "ClusterGuard 1.0.0 — первый релиз"

git checkout develop
git merge --no-ff release/1.0.0 -m "Merge branch 'release/1.0.0' into develop"

# ========== hotfix от main ==========
git checkout main
git checkout -b hotfix/alert-notifications
cp "$SRC/frontend/src/style.css" frontend/src/
commit "fix: контраст полей форм и стили автозаполнения"
git checkout main
git merge --no-ff hotfix/alert-notifications -m "Merge branch 'hotfix/alert-notifications'"
git tag -a v1.0.1 -m "ClusterGuard 1.0.1 — исправление UI форм"

git checkout develop
git merge --no-ff hotfix/alert-notifications -m "Merge branch 'hotfix/alert-notifications' into develop"

# --- перенос .git в рабочий каталог ---
cd "$ROOT"
rm -rf .git
mv "$BUILD/.git" "$ROOT/.git"
rsync -a --delete \
  --exclude='.git' \
  --exclude='build' \
  --exclude='node_modules' \
  --exclude='frontend/dist' \
  --exclude='frontend/package.json.md5' \
  --exclude='ClusterGuard' \
  "$SRC/" "$ROOT/"
mkdir -p "$ROOT/scripts"
cp /tmp/rebuild-gitflow.sh.bak "$ROOT/scripts/rebuild-gitflow.sh" 2>/dev/null || true
chmod +x "$ROOT/scripts/rebuild-gitflow.sh" 2>/dev/null || true
git add scripts/rebuild-gitflow.sh 2>/dev/null && git commit -m "chore: скрипт пересборки истории Git Flow" 2>/dev/null || true

git checkout develop 2>/dev/null || git checkout main

echo ""
echo "=== Git Flow ==="
git log --oneline --graph --decorate --all -35
echo ""
git branch -a
echo "Теги: $(git tag -l | tr '\n' ' ')"
