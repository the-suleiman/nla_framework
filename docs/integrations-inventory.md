# Optional integrations (inventory)

These remain in the codebase; none were removed except **Bitrix**. Use this list when planning the next refactor slice.

| Integration | Activation | Main generator touchpoints | Notes |
|-------------|------------|----------------------------|-------|
| **OData / 1C** | `len(Config.Odata.Url) > 0` → `IsOdataIntegration()` | [`templates/project.go`](../templates/project.go), [`templates/docIsIntegration.go`](../templates/docIsIntegration.go), [`templates/integrations/odata/`](../templates/integrations/odata/), [`types/docSqlFunc.go`](../types/docSqlFunc.go) (`uuid`-based update paths) | Per-doc `Integrations.Odata`; SQL uses `uuid` for upsert-style updates. |
| **Telegram** | `len(Config.Telegram.Token) > 0` | [`templates/project.go`](../templates/project.go), [`main.go`](../main.go) `configJsModify`, `sourceFiles` SQL notify helpers | Auth route, SQL functions, optional `tgBot` binary. |
| **Yandex Disk backup** | `len(Config.Backup.ToYandexDisk.Token) > 0` | [`templates/project/yandexDiskBackup/`](../templates/project/yandexDiskBackup/), [`templates/main.go`](../templates/main.go) optional `deployYandexBackup.ps1` | Standalone backup tooling; separate from Yandex Metrika in Vue. |
| **Graylog** | `Config.Graylog.Host` set | [`templates/project/main.go`](../templates/project/main.go), [`sourceFiles/src/graylog/`](../sourceFiles/src/graylog/), docker-compose logging | `/api/log`, GELF in compose. |
| **Phone auth** | `Config.Auth.ByPhone` | [`templates/project.go`](../templates/project.go) conditional SQL + Vue + `phone.go` | Email auth stays default when phone is off (`Start` in [`main.go`](../main.go)). |
| **SSE** | Always wired in templates | [`sourceFiles/src/sse/`](../sourceFiles/src/sse/), [`templates/project/webServer/main.go`](../templates/project/webServer/main.go) `/api/sse` | Live updates; Quasar `App.vue` template may or may not open `EventSource` depending on template version. |

**Removed in this pass:** Bitrix (project config, doc integration, SQL `btx_id` branches, `/bitrix` routes, templates under `templates/integrations/bitrix/`).
