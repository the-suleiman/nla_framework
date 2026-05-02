# Refactor backlog

Status of cleanups and outstanding work for the framework. Grouped by intent so the file doesn't drift back into a flat changelog.

## Resolved

### Frontend bundle (Quasar 2 only)

- **Quasar 1 removed** — `webClient/quasar_1`, `templates/project/webClient/quasar_1`, `templates/webClient/quasar_1` deleted; only `quasar_2` remains.
- **`QuasarWebClientDir` constant removed** — the `types.QuasarWebClientDir` indirection in [`types/typeProject.go`](../types/typeProject.go) is gone; all generator paths now reference `quasar_2` literally (see [main.go](../main.go), [templates/project.go](../templates/project.go), [templates/main.go](../templates/main.go), [templates/docIsRecursion.go](../templates/docIsRecursion.go), [templates/fldJsonListProcess.go](../templates/fldJsonListProcess.go), [templates/fldTagProcess.go](../templates/fldTagProcess.go), [templates/tmplGenerateStep2/](../templates/tmplGenerateStep2/), [types/shortcuts.go](../types/shortcuts.go), [types/typeFldVueComps.go](../types/typeFldVueComps.go), [types/docVueFunc.go](../types/docVueFunc.go), [types/typeDocStateMachine.go](../types/typeDocStateMachine.go)).
- **State machine Vue** — templates live under [`templates/webClient/quasar_2/doc/comp/stateMachine/`](../templates/webClient/quasar_2/doc/comp/stateMachine/) (copied from former v1-only path).
- **Tasks list** — `TasksTmpl` runs for Quasar 2; template path uses [`webClient/quasar_2/.../tasks/list.vue`](../templates/tmplGenerateStep2/vue.go).
- **Quasar CLI `.quasar`** — removed from copied `webClient/quasar_2` tree where present (regenerate locally with `quasar dev`).
- **Quasar / Vue bumped to current** — [`templates/project/webClient/quasar_2/package.json`](../templates/project/webClient/quasar_2/package.json) now pins `quasar@^2.19.3`, `vue@^3.5.33`, `vue-router@^4.6.4`, `@quasar/extras@^1.18.0`, replaces legacy `@quasar/app` with `@quasar/app-webpack@^4.4.5`, and bumps the eslint stack (`eslint@^8.57.0`, `eslint-plugin-vue@^9.30.0`, `eslint-webpack-plugin@^4.0.1`) plus `engines.node >= 20.0.0`.
- **`core-js` minimum** — `core-js@^3.39.0` (was `^3.6.5`). `@quasar/babel-preset-app` emits imports like `core-js/modules/es.iterator.*.js` for Iterator helpers; those modules did not exist in 3.6.x, which caused hundreds of “Module not found” webpack errors.
- **ESLint + generated Vue** — [`webClient/quasar_2/webClient/.eslintrc.js`](../webClient/quasar_2/webClient/.eslintrc.js) disables `vue/multi-word-component-names` and `vue/no-v-text-v-html-on-component` so `eslint-plugin-vue` v9 does not fail `quasar dev` on generated `index.vue`/`item.vue` naming or `v-html` on Quasar components (copied into generated apps with `webClient/quasar_2`).
- **`quasar.conf.js` → `quasar.config.js`** — template file renamed under [`templates/project/webClient/quasar_2/`](../templates/project/webClient/quasar_2/), generator output path updated in [`templates/project.go`](../templates/project.go), and inline references in [`types/projectVueFunc.go`](../types/projectVueFunc.go) and [`webClient/quasar_2/webClient/src/router/index.js`](../webClient/quasar_2/webClient/src/router/index.js) updated. Config uses the v4 `eslint:` block in place of per-section `chainWebpack` ESLint plugin wiring, and `import { defineConfig } from '#q-app/wrappers'` with `export default` (ESM) because `@quasar/app-webpack` loads the file as ESM; `require('quasar/wrappers')` is not supported in that pipeline.
- **`src/index.template.html` → root `index.html`** — Quasar CLI v4 expects the HTML shell next to `package.json` (same folder as `quasar.config.js`), not under `src/`. Template renamed to [`templates/project/webClient/quasar_2/index.html`](../templates/project/webClient/quasar_2/index.html); generator writes `/webClient/index.html`. [`templates/project/webClient/quasar_2/quasar.config.js`](../templates/project/webClient/quasar_2/quasar.config.js) sets `sourceFiles.indexHtmlTemplate: 'index.html'`. Copy-hook in [`main.go`](../main.go) accepts both the new path and legacy `src/index.template.html`. **`OverridePathForTemplates`** keys for this file are now `/webClient/index.html` (was `/webClient/src/index.template.html`).
- **`<!-- quasar:entry-point -->` in `index.html`** — current `@quasar/app-webpack` / Quasar v2 CLI wants a single `<!-- quasar:entry-point -->` inside `<body>` and no longer a hand-written `<div id="q-app"></div>`; the build injects the app at that comment (see the [Quasar CLI Webpack upgrade guide](https://quasar.dev/quasar-cli-webpack/upgrade-guide/)). Yandex Metrika (when enabled) is kept after that marker, still inside `<body>`.

### Generator surface

- **Legacy SQL** — removed `sourceFilesSQL_legacy/` and copy branch in [`main.go`](../main.go).
- **GOPATH `FillLocalPath`** — only **`go.mod`**-based inference remains (or explicit `LocalProjectPath`).
- **Bitrix** — types, doc integration, SQL branches, generated routes/config, `templates/integrations/bitrix/` removed.
- **`DevModeConfig.IsDocker`** — confirmed unused and removed from [`types/typeProject.go`](../types/typeProject.go) (no read sites in framework, templates, or generated runtime config).

## Current decisions

- **Single web client.** There is one Quasar 2 SPA bundle under [`webClient/quasar_2/`](../webClient/quasar_2/), [`templates/project/webClient/quasar_2/`](../templates/project/webClient/quasar_2/), and [`templates/webClient/quasar_2/`](../templates/webClient/quasar_2/). There is no project config to choose another frontend, no Quasar version selector, and no parallel `quasar_*` tree. Bumping Quasar / `@quasar/app-webpack` / Vue is a framework-level change in [`templates/project/webClient/quasar_2/package.json`](../templates/project/webClient/quasar_2/package.json) plus [`templates/project/webClient/quasar_2/quasar.config.js`](../templates/project/webClient/quasar_2/quasar.config.js).
- **Generated runtime config is a subset of `ProjectConfig`.** The generated app's [`templates/project/types/config.go`](../templates/project/types/config.go) only deserializes Postgres / WebServer / Graylog / Email / optional Telegram / optional OData. Adding new fields to `types.ProjectConfig` does not automatically surface them at runtime; they must also be wired into the generated `Config` and `ReadConfigFile`.
- **Dev vs. Docker behavior.** Runtime dev/Docker switching in the generated app is driven by env vars (`IS_DEVELOPMENT`, `PG_HOST`, `PG_PORT`, etc.) and the `-dev` flag, not by framework-level `ProjectConfig` flags.

## Candidate next refactors (out of scope)

### Trim optional integrations if unused in your product line

- **OData** — removing it touches [`types/typeOdata*.go`](../types/), `Config.IsOdataIntegration`, generated routes, [`templates/integrations/odata/`](../templates/integrations/odata/), and the optional `Odata` block in generated `config.toml` / `ReadConfigFile`.
- **Telegram** — touches `Config.Telegram`, `IsTelegramIntegration`, generated `tgBot/`, [`templates/integrations/telegram/`](../templates/integrations/telegram/), and SQL functions for telegram-based auth.
- **Yandex Disk backup** — generated `yandexDiskBackup/` package, systemd unit, and `Config.Backup`.
- **Graylog** — `Config.Graylog` + generated `Graylog` block in `ReadConfigFile`.

### Server / client alignment

- **SSE vs. Quasar 2 client** — align server-generated SSE endpoints (`sourceFiles/src/sse/`, `templates/project/webServer/main.go` `/api/sse`) with what the SPA actually subscribes to in [`webClient/quasar_2/webClient/src/boot/`](../webClient/quasar_2/webClient/src/boot/).
- **`ProjectConfig` <-> generated `Config` parity** — consider a single source of truth (or a generator) for the runtime config struct + `ReadConfigFile` so framework-side schema changes can't silently leave generated apps unable to read new fields.

### Project hygiene

- **Tests** — no `*_test.go` today; add golden tests for template output or a minimal fixture project.
- **Root `go.mod`** — optional: add module file + `replace` for `sourceFiles` import paths if you want `go test ./...` at framework root.
- **Doc template delimiters** — `templates/project/*` mixes `[[ ]]` and `{{ }}` delimiters depending on file. Document or unify per directory.
