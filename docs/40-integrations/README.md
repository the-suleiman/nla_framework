# Integrations

This folder documents the optional integrations that still exist in the framework codebase and can be activated via `ProjectType.Config` / `DocType.Integrations`.

## index

- [OData / 1C](odata.md)
- [Telegram](telegram.md)
- [Graylog](graylog.md)
- [Phone auth](phone-auth.md)
- [Yandex Disk backup](yandex-backup.md)
- [SSE](sse.md)

## inventory (activation + source of truth)

| Integration | Activation | Source of truth (generator touchpoints) |
|-------------|------------|------------------------------------------|
| **OData / 1C** | `len(Config.Odata.Url) > 0` → `IsOdataIntegration()` | `types/typeProject.go`, `types/typeDoc.go`, `templates/project.go`, `templates/docIsIntegration.go`, `templates/integrations/odata/` |
| **Telegram** | `len(Config.Telegram.Token) > 0` | `types/typeProject.go`, `templates/project.go`, `main.go` (`configJsModify`), `templates/integrations/telegram/` |
| **Yandex Disk backup** | `len(Config.Backup.ToYandexDisk.Token) > 0` | `types/typeProject.go`, `templates/project/yandexDiskBackup/` |
| **Graylog** | `Config.Graylog.Host` set | `types/typeProject.go`, `templates/project/`, `sourceFiles/src/graylog/` |
| **Phone auth** | `Config.Auth.ByPhone` | `types/typeProject.go`, `main.go` (defaults), `templates/project.go` |
| **SSE** | always wired (server-side) | `sourceFiles/src/sse/`, `templates/project/webServer/main.go` |

