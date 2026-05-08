# Copy and rewrite rules (`copyFiles`)

After templates are executed, the generator copies static trees and applies a set of rewrite rules.

Source: `main.go` `copyFiles`.

## what gets copied

- `sourceFiles/` → `../` (into the generated app tree)
- `webClient/` → `../src/webClient/` (Quasar 2 SPA skeleton)

## import rewriting (Go)

If the copied filename ends with `.go`, it rewrites imports:

- from `github.com/the-suleiman/nla_framework`
- to `ProjectType.Config.LocalProjectPath`

This is how copied runtime helpers become part of the generated app module.

## slot injection / special file hooks

During copy, `copyFiles` applies targeted modifications for certain files:

- `app/plugins/config.js`
  - `configJsModify`: app name, UI name, ports/urls, logo, Dadata token, breadcrumb icons
  - injects “tables allowed for tasks”
  - injects Telegram config if enabled

- `src/router/routes.js`
  - injects generated routes into the `// for codeGenerate ##routes_slot1` slot

- `components/sidemenu/index.vue`
  - injects generated menu into the `// for codeGenerate ##sidemenu_slot1` slot

- `_Task/main.toml`
  - injects additional SQL methods into the `# for codeGenerate task_methods_slot` slot

- `trigger_task_update_table_name.sql`
  - injects additional trigger blocks into the `-- for codeGenerate #trigger_task_update_table_name_slot` slot

- `webClient/index.html` and legacy `webClient/src/index.template.html`
  - replaces `[[appName]]`

- `loginPage.vue`, `home.vue`
  - replaces `[[appLogoSrc]]`

- any `.sql` file
  - replaces `[[Config.Postgres.TimeZone]]` with `Config.Postgres.TimeZone`

## user-provided file modifier hook

If `modifyFunc` was passed to `Start`, it is called for every copied file after internal rewrites:

- `file = modifyFunc(dirPath+info.Name(), file)`

Use this to implement project-specific patching without forking the framework.

## write optimization for Quasar

To reduce unnecessary Quasar restarts:

- files under `webClient/` are compared to existing target content; if identical, they are not rewritten
- files under `webClient/.quasar/` are never overwritten

Note: template execution (`templates.ExecuteToFile`) contains a similar “skip unchanged” optimization for `webClient` outputs.

