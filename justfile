set windows-shell := ["powershell.exe", "-NoLogo", "-NoProfile", "-Command"]

frontend_pnpm := if os_family() == "windows" { "pnpm.cmd" } else { "pnpm" }

default:
    @just --list

status:
    git --no-pager status --short --branch

worktrees:
    git worktree list

plans:
    git --no-pager ls-files docs/plans

pending:
    git --no-pager ls-files docs/plans/v1-mvp/subagent-work-items/pending

approved:
    git --no-pager ls-files docs/plans/v1-mvp/subagent-work-items/approved

completed:
    git --no-pager ls-files docs/plans/v1-mvp/subagent-work-items/completed

frontend-dir:
    node -e "console.log(require('path').resolve('frontend'))"

backend-dir:
    node -e "console.log(require('path').resolve('backend'))"

backend-start:
    cd backend; go run ./cmd/api

list-users *args:
    node scripts/list-users.mjs {{args}}

promote-manager *args:
    node scripts/promote-manager.mjs {{args}}

dev:
    node scripts/dev-launcher.mjs

backend-test:
    cd backend; go test ./...

backend-migrate:
    cd backend; go run ./cmd/migrate

backend-seed:
    cd backend; go run ./cmd/seed

frontend-install:
    cd frontend; {{ frontend_pnpm }} install

frontend-start:
    cd frontend; {{ frontend_pnpm }} exec ng serve --proxy-config src/proxy.conf.json

frontend-build:
    cd frontend; {{ frontend_pnpm }} exec ng build

frontend-test:
    cd frontend; {{ frontend_pnpm }} exec ng test --watch=false --browsers=ChromeHeadless
