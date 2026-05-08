# Refactor backlog

Status of cleanups and outstanding work for the framework. Grouped by intent so the file doesn't drift back into a flat changelog.

## Resolved

### Frontend bundle (Quasar 2 only)

- **Flat `webClient/` layout** — removed the `webClient/quasar_2/webClient` nesting and the `quasar_2` segment under [`templates/project/webClient/`](../templates/project/webClient/) and [`templates/webClient/doc/`](../templates/webClient/doc/). The SPA skeleton is [`webClient/`](../webClient/) at the framework repo root; [`main.go`](../main.go) copies it to `../src/webClient/` (with a normalized `copyFiles` source prefix so root-level files map correctly). **`types.QuasarWebClientDir` removed** from [`types/typeProject.go`](../types/typeProject.go).
- **Single SPA bundle (Quasar 2 only)** — one frontend; no project config to choose another stack or Quasar version. See [getting started overview](00-getting-started/overview.md) and [generated app layout](20-generator-pipeline/dist-layout.md). Bumping the CLI stack is a framework change in [`templates/project/webClient/package.json`](../templates/project/webClient/package.json) and [`templates/project/webClient/quasar.config.js`](../templates/project/webClient/quasar.config.js).
- **Quasar 1 removed** — `webClient/quasar_1`, `templates/project/webClient/quasar_1`, `templates/webClient/quasar_1` deleted.
- **State machine Vue** — templates live under [`templates/webClient/doc/comp/stateMachine/`](../templates/webClient/doc/comp/stateMachine/) (copied from former v1-only path).
- **Tasks list** — `TasksTmpl` reads [`webClient/src/app/components/currentUser/tasks/list.vue`](../webClient/src/app/components/currentUser/tasks/list.vue); see [`templates/tmplGenerateStep2/vue.go`](../templates/tmplGenerateStep2/vue.go).
- **Quasar CLI `.quasar`** — removed from copied [`webClient/`](../webClient/) tree where present (regenerate locally with `quasar dev`).
- **Quasar / Vue bumped to current** — [`templates/project/webClient/package.json`](../templates/project/webClient/package.json) pins `quasar@^2.19.3`, `vue@^3.5.33`, `vue-router@^4.6.4`, `@quasar/extras@^1.18.0`, replaces legacy `@quasar/app` with `@quasar/app-webpack@^4.4.5`, and bumps the eslint stack (`eslint@^8.57.0`, `eslint-plugin-vue@^9.30.0`, `eslint-webpack-plugin@^4.0.1`) plus `engines.node >= 20.0.0`.
- **`core-js` minimum** — `core-js@^3.39.0` (was `^3.6.5`). `@quasar/babel-preset-app` emits imports like `core-js/modules/es.iterator.*.js` for Iterator helpers; those modules did not exist in 3.6.x, which caused hundreds of “Module not found” webpack errors.
- **ESLint + generated Vue** — [`webClient/.eslintrc.js`](../webClient/.eslintrc.js) disables `vue/multi-word-component-names` and `vue/no-v-text-v-html-on-component` so `eslint-plugin-vue` v9 does not fail `quasar dev` on generated `index.vue`/`item.vue` naming or `v-html` on Quasar components.
- **`quasar.conf.js` → `quasar.config.js`** — template under [`templates/project/webClient/`](../templates/project/webClient/); generator in [`templates/project.go`](../templates/project.go) and [`types/projectVueFunc.go`](../types/projectVueFunc.go); router shell in [`webClient/src/router/index.js`](../webClient/src/router/index.js). Config uses the v4 `eslint:` block in place of per-section `chainWebpack` ESLint plugin wiring, and `import { defineConfig } from '#q-app/wrappers'` with `export default` (ESM) because `@quasar/app-webpack` loads the file as ESM; `require('quasar/wrappers')` is not supported in that pipeline.
- **`src/index.template.html` → root `index.html`** — Quasar CLI v4 expects the HTML shell next to `package.json` (same folder as `quasar.config.js`), not under `src/`. Template: [`templates/project/webClient/index.html`](../templates/project/webClient/index.html); generator writes `/webClient/index.html`. [`templates/project/webClient/quasar.config.js`](../templates/project/webClient/quasar.config.js) sets `sourceFiles.indexHtmlTemplate: 'index.html'`. Copy-hook in [`main.go`](../main.go) accepts both the new path and legacy `src/index.template.html`. **`OverridePathForTemplates`** keys for this file are `/webClient/index.html` (was `/webClient/src/index.template.html`).
- **`<!-- quasar:entry-point -->` in `index.html`** — current `@quasar/app-webpack` / Quasar v2 CLI wants a single `<!-- quasar:entry-point -->` inside `<body>` and no longer a hand-written `<div id="q-app"></div>`; the build injects the app at that comment (see the [Quasar CLI Webpack upgrade guide](https://quasar.dev/quasar-cli-webpack/upgrade-guide/)). Yandex Metrika (when enabled) is kept after that marker, still inside `<body>`.

### Generator surface

- **Legacy SQL** — removed `sourceFilesSQL_legacy/` and copy branch in [`main.go`](../main.go).
- **GOPATH `FillLocalPath`** — only **`go.mod`**-based inference remains (or explicit `LocalProjectPath`).
- **Bitrix** — types, doc integration, SQL branches, generated routes/config, `templates/integrations/bitrix/` removed.
- **`DevModeConfig.IsDocker`** — confirmed unused and removed from [`types/typeProject.go`](../types/typeProject.go) (no read sites in framework, templates, or generated runtime config).
- **Go stdlib deprecations in template/generator code** — [`utils.UpperCaseFirst`](../utils/main.go) now uses `golang.org/x/text/cases` instead of deprecated `strings.Title`, and template/generator file I/O now uses `os.ReadFile`, `os.WriteFile`, `os.ReadDir`, and `io.ReadAll` instead of deprecated `io/ioutil`. The root [`go.mod`](../go.mod) now records these framework-side dependencies so `go mod tidy` can run at the repository root.

### Database / Postgres

- **Default Postgres bumped 12 -> 18** — [`main.go`](../main.go) now sets `Config.Postgres.Version = "18"` when unset (PG 12 EOL Nov 2024). Field comment in [`types/typeProject.go`](../types/typeProject.go) updated accordingly. Only affects projects that don't pin `Config.Postgres.Version`.
- **Compose volume mount switched to PG 18 layout** — [`templates/project/docker-compose.yml`](../templates/project/docker-compose.yml) and [`templates/project/docker-compose.dev.yml`](../templates/project/docker-compose.dev.yml) mount at `/var/lib/postgresql` (was `/var/lib/postgresql/data`). PG 18's official image moved `PGDATA` to `/var/lib/postgresql/MAJOR/docker` and declares `VOLUME /var/lib/postgresql` ([docker-library/postgres#1259](https://github.com/docker-library/postgres/pull/1259)). **Caveat:** these templates now assume PG **18+**. Pinning `Config.Postgres.Version` to `<= 17` will not work without restoring the old `/var/lib/postgresql/data` mount (or setting `PGDATA=/var/lib/postgresql/data` explicitly on the container).
- **No auto-migration of existing dev DB.** The dev compose's named volume already includes the version (`postgres_data_<name>_<version>`), so a fresh `_18` volume is created on bump; old `_12` (or any prior pinned) data is not auto-migrated. Migrate via `pg_dumpall` against the old container, then restore against the new one — pattern in [`templates/project/restoreDump.sh`](../templates/project/restoreDump.sh). Generated backup tooling ([`templates/project/restoreDump.sh`](../templates/project/restoreDump.sh), [`templates/project/yandexDiskBackup/dbBackup.go`](../templates/project/yandexDiskBackup/dbBackup.go)) uses logical `pg_dumpall`, so it is version-portable and unchanged.

## Current decisions

- **Generated runtime config is a subset of `ProjectConfig`.** The generated app's [`templates/project/types/config.go`](../templates/project/types/config.go) only deserializes Postgres / WebServer / Graylog / Email / optional Telegram / optional OData. Adding new fields to `types.ProjectConfig` does not automatically surface them at runtime; they must also be wired into the generated `Config` and `ReadConfigFile`.
- **Dev vs. Docker behavior.** Runtime dev/Docker switching in the generated app is driven by env vars (`IS_DEVELOPMENT`, `PG_HOST`, `PG_PORT`, etc.) and the `-dev` flag, not by framework-level `ProjectConfig` flags.

## Candidate next refactors (out of scope)

### Trim optional integrations if unused in your product line

- **OData** — removing it touches [`types/typeOdata*.go`](../types/), `Config.IsOdataIntegration`, generated routes, [`templates/integrations/odata/`](../templates/integrations/odata/), and the optional `Odata` block in generated `config.toml` / `ReadConfigFile`.
- **Telegram** — touches `Config.Telegram`, `IsTelegramIntegration`, generated `tgBot/`, [`templates/integrations/telegram/`](../templates/integrations/telegram/), and SQL functions for telegram-based auth.
- **Yandex Disk backup** — generated `yandexDiskBackup/` package, systemd unit, and `Config.Backup`.
- **Graylog** — `Config.Graylog` + generated `Graylog` block in `ReadConfigFile`.

### Server / client alignment

- **SSE vs. Quasar 2 client** — align server-generated SSE endpoints (`sourceFiles/src/sse/`, `templates/project/webServer/main.go` `/api/sse`) with what the SPA actually subscribes to in [`webClient/src/boot/`](../webClient/src/boot/).
- **`ProjectConfig` <-> generated `Config` parity** — consider a single source of truth (or a generator) for the runtime config struct + `ReadConfigFile` so framework-side schema changes can't silently leave generated apps unable to read new fields.

### Project hygiene

- **Tests** — no `*_test.go` today; add golden tests for template output or a minimal fixture project.
- **Root `go test ./...` cleanup** — root module metadata exists now, but full-package test runs still trip over template-placeholder directories and generated-app source packages. Add a test strategy that excludes raw templates or renders a fixture app before running full checks.
- **Doc template delimiters** — `templates/project/*` mixes `[[ ]]` and `{{ }}` delimiters depending on file. Document or unify per directory.
