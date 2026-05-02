# Refactor backlog

## Done (this pass)

- **Quasar 1** — removed `webClient/quasar_1`, `templates/project/webClient/quasar_1`, `templates/webClient/quasar_1`; single bundle [`types.QuasarWebClientDir`](../types/typeProject.go) = `quasar_2`.
- **Legacy SQL** — removed `sourceFilesSQL_legacy/` and copy branch in [`main.go`](../main.go).
- **GOPATH `FillLocalPath`** — only **`go.mod`**-based inference remains (or explicit `LocalProjectPath`).
- **Bitrix** — types, doc integration, SQL branches, generated routes/config, [`templates/integrations/bitrix/`](../templates/integrations/bitrix/) removed.
- **State machine Vue** — templates live under `templates/webClient/quasar_2/doc/comp/stateMachine/` (copied from former v1-only path).
- **Tasks list** — `TasksTmpl` runs for Quasar 2; template path uses [`webClient/quasar_2/.../tasks/list.vue`](../templates/tmplGenerateStep2/vue.go).
- **Quasar CLI `.quasar`** — removed from copied `webClient/quasar_2` tree where present (regenerate locally with `quasar dev`).

## Candidate next refactors (out of scope)

- Trim **OData / Telegram / Yandex backup / Graylog** if unused in your product line (each touches SQL, compose, and/or binaries).
- **SSE** vs Quasar 2 client: align server-generated SSE with what the SPA actually subscribes to.
- **`DevModeConfig.IsDocker`** — appears unused in generation; validate and delete or wire up.
- **Tests** — no `*_test.go` today; add golden tests for template output or a minimal fixture project.
- **Root `go.mod`** — optional: add module file + `replace` for `sourceFiles` import paths if you want `go test ./...` at framework root.
