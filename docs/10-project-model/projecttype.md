# `ProjectType` reference

`ProjectType` is the root input to the generator (`Start(p, ...)`).

Source: `types/typeProject.go`.

## core fields

- `Name`
  - project identifier; must not contain spaces (validated in `main.go` `readData`)
- `Docs []DocType`
  - list of documents/entities to generate
- `DistPath`
  - output root; set to `../src` by default during generation (`main.go` `readData`)

## `Config ProjectConfig`

High level:

- `LocalProjectPath`
  - module import path prefix used when rewriting copied `.go` files
  - default: inferred from `go.mod` by `FillLocalPath()` (module path + `/src`)
- `Auth`
  - `ByPhone` toggles phone auth codegen; when false, generator defaults `ByEmail=true`
- `Postgres`
  - includes `TimeZone` (used to replace `[[Config.Postgres.TimeZone]]` in copied `.sql`)
  - includes `Version` (defaults to `"18"` when empty)
- `WebServer`
  - port/url; also injected into client `config.js`
- `Email`
  - SMTP sender/service settings used by features that send email (including email-based auth when enabled)
- integrations config blocks:
  - `Telegram`, `Odata`, `Backup`, `Graylog`

See the guided map: [configuration map](../00-getting-started/configuration-map.md).

## `Vue ProjectVue`

Controls project-level Vue shell generation:

- `Routes`: additional routes injected into `routes.js` during copy
- `Menu`: side menu items injected into `sidemenu/index.vue` during copy
- `Theme`, `QuasarBoot`, `IndexHtmlHead`, and other UI toggles

## `Sql ProjectSql`

Project-level SQL methods/hooks:

- `Methods`: additional SQL methods; also used for task-related injection slots
- `InitialData`: initial SQL data template content

## `Go ProjectGo`

Project-level Go wiring:

- jobs list and extra route blocks for generated server
- extra main imports and custom flags

## template overrides

- `OverridePathForTemplates map[string]string`
  - **project-level template overrides** keyed by destination (e.g. `"/webClient/index.html"`)
  - applied in `templates/project.go` `ReadTmplAndPrint`

