# Template system

## Directories

| Path | Role |
|------|------|
| [`templates/project/`](../templates/project/) | Generated app: Go main, Docker, SQL models (`01_User`, …), Quasar 2 app templates |
| [`templates/webClient/doc/`](../templates/webClient/doc/) | Doc-level Vue: `index.vue`, `item.vue`, tabs, `tabTasks.vue`, state machine under `doc/comp/stateMachine/` |
| [`templates/sql/`](../templates/sql/) | Generic SQL function templates (`list`, `update`, triggers, …) |
| [`templates/integrations/odata/`](../templates/integrations/odata/) | OData sync code templates |

Delimiters: most doc/vue/sql templates use **`[[` `]]`**; some `project_*` templates use `{{` `}}`.

## Template helpers

The shared `FuncMap` lives in [`templates/main.go`](../templates/main.go). It exposes common string helpers (`ToUpper`, `ToLower`, `ToCamel`, `ToLowerCamel`, `UpperCaseFirst`), Vue field rendering (`PrintVueFldTemplate`), select option rendering, and small formatting helpers used by `.vue`, `.sql`, and project templates.

`UpperCaseFirst` is implemented in [`utils/main.go`](../utils/main.go) with `golang.org/x/text/cases` instead of the deprecated `strings.Title`, so new helpers should avoid reintroducing deprecated stdlib APIs.

## Doc template naming → output

Resolved by **`utils.ParseDocTemplateFilename`** ([`utils/main.go`](../utils/main.go)): e.g. `webClient_*` → `webClient/src/app/components/...`, `sql_main.toml` → numbered `sql/model/`, `sql_function_*` → `sql/template/function/_DocName/`.

## Post-copy hooks ([`main.go`](../main.go) `copyFiles`)

- Go files: replace `github.com/the-suleiman/nla_framework` with `Config.LocalProjectPath`.
- **`config.js`**: `configJsModify` — app name, URLs, Telegram snippet, task tables, etc.
- **`routes.js`**, **`sidemenu/index.vue`**: codegen slots.
- **`_Task/main.toml`**, **`trigger_task_update_table_name.sql`**: optional task hooks from `Project.Sql.Methods`.
- SQL files: timezone placeholder `[[Config.Postgres.TimeZone]]`.

## Secondary pass

[`templates/tmplGenerateStep2/`](../templates/tmplGenerateStep2/):

- **`TasksTmpl`** — builds `currentUser/tasks/list.vue` from `[[PrintImports]]` / `[[PrintComps]]` using the list template in [`webClient/.../tasks/list.vue`](../webClient/src/app/components/currentUser/tasks/list.vue) (framework root; copied to generated `../src/webClient/`).
- **`PluginUtilsJs`**, **`BootI18nJs`** — overlay `utils.js` and `i18n.js` after copy.

Template execution and secondary generators use `os.ReadFile`, `os.WriteFile`, `os.ReadDir`, and `io.ReadAll`; `io/ioutil` should not be used in new generator code.
