# Telegram integration

This integration adds Telegram auth wiring and optional `tgBot` tooling to the generated app.

## activation

- `ProjectType.IsTelegramIntegration()` → `len(Config.Telegram.Token) > 0`

## what it generates

When enabled, `templates/project.go` generates:

- `../src/webServer/telegramAuth.go`
- SQL functions under `../src/sql/template/function/_User/`:
  - `user_telegram_auth.sql`
  - `user_get_by_telegram_id.sql`
- `../src/tgBot/main.go`

It also injects telegram config into the client `config.js` via `main.go` `configJsModify`.

## generator touchpoints (source of truth)

- `types/typeProject.go`: `TelegramConfig`, `IsTelegramIntegration`
- `templates/project.go`: generation of Telegram wiring and SQL helpers
- `templates/integrations/telegram/`: template sources
- `main.go`: `configJsModify` injects Telegram block into `app/plugins/config.js`

## runtime expectations

- the generated backend must have access to the token/config at runtime (via generated config)
- the client `config.js` contains telegram settings only when enabled at generation time

