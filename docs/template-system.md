# Template system

## Directories

| Path | Role |
|------|------|
| [`templates/project/`](../templates/project/) | Generated app: Go main, Docker, SQL models (`01_User`, …), Quasar 2 app templates |
| [`templates/webClient/quasar_2/doc/`](../templates/webClient/quasar_2/doc/) | Doc-level Vue: `index.vue`, `item.vue`, tabs, `tabTasks.vue`, state machine under `doc/comp/stateMachine/` |
| [`templates/sql/`](../templates/sql/) | Generic SQL function templates (`list`, `update`, triggers, …) |
| [`templates/integrations/odata/`](../templates/integrations/odata/) | OData sync code templates |

Delimiters: most doc/vue/sql templates use **`[[` `]]`**; some `project_*` templates use `{{` `}}`.

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

- **`TasksTmpl`** — builds `currentUser/tasks/list.vue` from `[[PrintImports]]` / `[[PrintComps]]` using templates under `webClient/quasar_2/.../tasks/list.vue`.
- **`PluginUtilsJs`**, **`BootI18nJs`** — overlay `utils.js` and `i18n.js` after copy.
