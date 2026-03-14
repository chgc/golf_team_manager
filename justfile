set windows-shell := ["powershell.exe", "-NoLogo", "-NoProfile", "-Command"]

default:
    @just --list

status:
    git --no-pager status --short --branch

worktrees:
    git worktree list

plans:
    Get-ChildItem docs\plans -Recurse -File | Select-Object -ExpandProperty FullName

pending:
    Get-ChildItem docs\plans\v1-mvp\subagent-work-items\pending -File | Select-Object -ExpandProperty Name

approved:
    Get-ChildItem docs\plans\v1-mvp\subagent-work-items\approved -File | Select-Object -ExpandProperty Name

completed:
    Get-ChildItem docs\plans\v1-mvp\subagent-work-items\completed -Recurse -File | Select-Object -ExpandProperty FullName

frontend-dir:
    Get-Item frontend | Select-Object FullName

backend-dir:
    Get-Item backend | Select-Object FullName

backend-start:
    Set-Location backend; go run ./cmd/api

backend-test:
    Set-Location backend; go test ./...

backend-migrate:
    Set-Location backend; go run ./cmd/migrate

backend-seed:
    Set-Location backend; go run ./cmd/seed

frontend-install:
    Set-Location frontend; pnpm.cmd install

frontend-start:
    Set-Location frontend; pnpm.cmd exec ng serve

frontend-build:
    Set-Location frontend; pnpm.cmd exec ng build

frontend-test:
    Set-Location frontend; pnpm.cmd exec ng test --watch=false --browsers=ChromeHeadless
