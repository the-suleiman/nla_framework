# Project model

## `ProjectType` ([`types/typeProject.go`](../types/typeProject.go))

- **`Name`** — project identifier (no spaces).
- **`Docs`** — slice of `DocType` (entities / tables / screens).
- **`Config`** — `ProjectConfig`: Postgres, web server, email, auth, Vue (e.g. Dadata token), Telegram, OData, Yandex (Metrika id), backup, Docker, Graylog, etc.
- **`Vue`** — menu, routes, hooks, theme, `QuasarBoot`, message templates.
- **`Sql`** — project-level SQL methods and initial data hooks.
- **`Go`** — jobs, extra routes/imports/flags.
- **`Roles`**, **`I18n`**, **`OverridePathForTemplates`** — localization and template overrides.

## `DocType` ([`types/typeDoc.go`](../types/typeDoc.go))

Describes one document: fields, Vue (grid, tabs, routes, filters), SQL (triggers, methods), templates map, optional **state machine**, recursion, **Odata** integration metadata, i18n entries.

## `FldType` ([`types/typeFld.go`](../types/typeFld.go))

Field SQL + Vue typing; builder helpers in [`types/shortcuts.go`](../types/shortcuts.go). Integration metadata for OData lives in `IntegrationData["odata"]` (`types.OdataFld`).

## Local module path

**`FillLocalPath`** finds **`go.mod`** by walking **up from the current working directory** (so running from e.g. `projectTemplate/` still picks the app’s module). It sets `Config.LocalProjectPath` to `modulePath` + `/src` (with optional `/projectTemplate` suffix trimmed). Set `LocalProjectPath` explicitly if you work outside a normal module tree.
