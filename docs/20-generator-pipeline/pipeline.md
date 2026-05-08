# Generator pipeline (end-to-end)

This page is the authoritative “what happens when you run the generator” guide, grounded in `main.go` and `templates/*`.

## entrypoint

The public API is:

- `Start(p types.ProjectType, modifyFunc)` in `main.go`

`modifyFunc` is optional and can mutate **copied** files (from `sourceFiles/` and `webClient/`) before they’re written.

## step 0: defaults and validation

`Start` sets several defaults before generation:

- if `Config.Auth.ByPhone == false` then `Config.Auth.ByEmail = true` (email auth is default)
- if `Config.Postgres.Version` is empty, it defaults to `"18"`

Then `readData(p)`:

- sets `project.Config.LocalProjectPath = project.FillLocalPath()` (module path inference from `go.mod` if empty)
- sets `project.DistPath = "../src"`
- fills doc template fields and Vue helpers:
  - `FillDocTemplatesFields()`
  - `GenerateGrid()`
  - `FillVueFlds()`
  - fills i18n defaults if missing (`FillI18n()`)
- validates:
  - project name has no spaces
  - no field is named `user_id` (reserved)
  - option values contain no spaces (Vue select options)
  - `Config.Email` has required fields (sender/host/port)
  - for “unique link docs”: `title` must not be unique

## step 1: template parsing and attachment

`templates.ParseTemplates(project)` loads:

- project templates (docker, root config, etc.)
- global doc-level templates (Vue, SQL)
- per-doc custom templates from `DocType.Templates`
- additional templates based on doc features:
  - tabs
  - recursion
  - integrations (e.g. OData)
  - special field types (tags, jsonList)
  - state machine template overrides (SQL update, Vue item template)

Important behaviors:

- directories use different delimiters (`[[ ]]` vs `{{ }}`) depending on template group
- base templates are attached only if `DocType.IsBaseTemplates` requests them

## step 2: clean up old generated files

`removeOldFiles(project.DistPath)` currently deletes:

- `../src/sql/model`

This is intentional: model numbering may change between runs.

## step 3: write project-level generated files

`templates.WriteProjectFiles(project, tmplMap)` writes:

- “project_” templates parsed in step 1
- additional project-level templates via `ReadTmplAndPrint` (types, webServer, SQL seed, Quasar shell, etc.)
- i18n index generation (`FillDocI18n`, `PrintI18nJs`, per-language i18n files)
- integration-specific project files when enabled (telegram, yandex backup, odata)

Template overrides:

- `ProjectType.OverridePathForTemplates` can override sources by *destination key* (see `templates/project.go` `ReadTmplAndPrint`)

## step 4: execute doc templates

For each doc (`for _, d := range p.Docs`), the generator executes every `d.Templates[*]` using:

- `templates.ExecuteToFile(dt.Tmpl, d, dt.DistPath, dt.DistFilename)`

Write optimization:

- for `webClient` output, `ExecuteToFile` compares to existing file and skips writing if identical (to avoid unnecessary Quasar restarts)

## step 5: copy static trees + apply rewrites

The generator copies two trees using `copyFiles` in `main.go`:

- `sourceFiles/` → `../` (into `../src/...` paths contained in the tree)
- `webClient/` → `../src/webClient/`

During copy, it applies:

- Go import rewrite (`github.com/the-suleiman/nla_framework` → `Config.LocalProjectPath`)
- special hooks (config, routes, sidemenu, timezone placeholder, etc.)
- optional external `modifyFunc`
- write optimizations for `webClient/` files
- a hard “do not overwrite” rule for `webClient/.quasar/`

Details: [copy and rewrite rules](copy-and-rewrite.md).

## step 6: secondary generators (“step2”)

After copy, `templates.OtherTemplatesGenerate(project)` runs:

- `TasksTmpl(p)` (regenerates tasks list Vue)
- `PluginUtilsJs(p)` (overlays generated utils)
- `BootI18nJs(p)` (overlays i18n boot wiring)

These live under `templates/tmplGenerateStep2/`.

