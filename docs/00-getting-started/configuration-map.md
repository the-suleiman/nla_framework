# Configuration map (ProjectConfig)

This page is a guided map of `ProjectType.Config` sections: what each block controls in generated output.

## `Config.LocalProjectPath`

- **what it is**: module import path prefix used in generated/copied Go sources
- **defaulting**: inferred from `go.mod` (module path + `/src`) via `ProjectType.FillLocalPath()`
- **where used**: `main.go` `copyFiles` rewrites imports from `github.com/the-suleiman/nla_framework` → `LocalProjectPath`

## `Config.Auth`

- **email vs phone**:
  - if `Auth.ByPhone == false`, the generator sets `Auth.ByEmail = true` by default (see `Start` in `main.go`)
  - `Auth.ByEmail` / `Auth.ByPhone` toggle which auth templates and flows are generated

## `Config.Postgres`

- affects:
  - generated `config.toml`
  - docker compose templates (version, volume layout, etc.)
  - timezone placeholder replacement in `.sql` files during copy (`[[Config.Postgres.TimeZone]]`)

## `Config.WebServer`

- affects:
  - generated backend server config and client `config.js` (port/url)

## `Config.Email`

- what it is:
  - SMTP sender/service settings (sender address, host, port, password) used by features that send email (including email-based auth when enabled)

## optional integrations (high level)

- `Config.Telegram`: enables Telegram auth + bot tooling
- `Config.Odata`: enables OData sync tooling
- `Config.Graylog`: enables Graylog wiring
- `Config.Backup.ToYandexDisk`: enables generated Yandex Disk backup tooling

Details: see [integrations](../40-integrations/README.md).

