# Generated app layout (`../src`)

Generation writes relative to the caller’s working tree. By default, `project.DistPath` is set to `../src` in `main.go` `readData`.

## top-level output

- `../src/types/`
  - generated types + config loader (from `templates/project/types/`)
- `../src/webServer/`
  - Gin server, auth, API wiring (templates + copied `sourceFiles/src/webServer/`)
- `../src/sql/model/`
  - table model TOML files (one folder per table, e.g. `01_User/`)
- `../src/sql/template/function/`
  - SQL functions/triggers (generic + per-doc under `_<DocName>/`)
- `../src/webClient/`
  - Quasar 2 SPA
  - part comes from templates (`templates/project/webClient/`)
  - part is copied from the static skeleton (`webClient/`)
- `../src/pg/`, `../src/jobs/`, `../src/utils/`, `../src/sse/`
  - copied from `sourceFiles/`

## generation notes

- `../src/sql/model` is deleted at the start of each run to avoid numbering collisions (`removeOldFiles`).
- copied Go sources are import-rewritten to your module path (`Config.LocalProjectPath`), so the generated app builds under your project’s module.

