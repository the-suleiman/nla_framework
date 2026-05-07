# Generated application layout

Generation writes relative to the caller’s tree; **`project.DistPath`** defaults to **`../src`** ([`main.go`](../main.go) `readData`).

Typical output:

- **`../src/types/`** — config loader, generated types ([`templates/project/types/`](../templates/project/types/)).
- **`../src/webServer/`** — Gin app, auth, API, SSE ([`sourceFiles/src/webServer/`](../sourceFiles/src/webServer/) + templates).
- **`../src/sql/model/`** — TOML models per table.
- **`../src/sql/template/function/`** — SQL functions and triggers.
- **`../src/webClient/`** — Quasar 2 app (from [`templates/project/webClient/`](../templates/project/webClient/) templates + copy of [`webClient/`](../webClient/)).
- **`../src/pg/`**, **`../src/jobs/`**, **`../src/utils/`**, **`../src/sse/`** — from `sourceFiles`.

Root-level **`config.toml`**, **`Dockerfile`**, **`docker-compose*.yml`** come from `project_*` templates ([`templates/project/`](../templates/project/)).

Optional packages when flags are set (see [integrations-inventory.md](integrations-inventory.md)): **`odata/`**, **`tgBot/`**, **`yandexDiskBackup/`**.

Generated Go sources are intended to build in the target app after import rewriting to `Config.LocalProjectPath`. Runtime helpers copied from [`sourceFiles/`](../sourceFiles/) use current `io` APIs (`io.ReadAll`) rather than deprecated `io/ioutil`.
