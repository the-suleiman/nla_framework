# `DocType` reference

`DocType` describes one generated “document” (entity/table + UI + SQL).

Source: `types/typeDoc.go`.

## identity and structure

- `Name`, `NameRu`
  - used for SQL naming, i18n keys, route/component naming, file paths
- `Flds []FldType`
  - field list; `DocType.Init()` sets `FldType.Doc = d` for every field

## templates

- `Templates map[string]*DocTemplate`
  - both user-provided and auto-attached templates live here
  - base templates are attached in `templates.ParseTemplates` depending on `IsBaseTemplates`
- `IsBaseTemplates`
  - `Vue: true` attaches default Vue templates (`index.vue`, `item.vue`)
  - `Sql: true` attaches default SQL templates (`main.toml`, list/update/triggers)
- `TemplatePathOverride map[string]TmplPathOverride`
  - **doc-level template override** mechanism used by some integrations (e.g. OData) to replace a specific template source
- `PathPrefix string`
  - prefix prepended when the generator infers a template's `Source` in `ProjectType.FillDocTemplatesFields`
  - inferred source layout: `<PathPrefix>/<snakeToCamelLower(Name)>/tmpl/<templateKey>`
  - use it when your doc folders aren't siblings of the generator's `main.go` (e.g. you keep all docs under `docs/<docName>/`, in which case set `PathPrefix: "docs"`)
  - if a `DocTemplate.Source` is set explicitly, `PathPrefix` is ignored for that entry

## Vue (`DocVue`)

`DocVue` controls UI generation:

- `RouteName`, `MenuIcon`, `BreadcrumbIcon`, `Roles`
- `Tabs []VueTab`
  - if tabs are present, base template `item.vue` is replaced with `itemWithTabs.vue`
  - tabs are also attached as templates under `.../tabs/<tabTitle>/index.vue`
- `Readonly`
  - if set (not `"false"`), it is applied to fields that don’t explicitly set their own readonly condition
- list customization: filters, sorts, create modal, add buttons slot, etc.
- hooks:
  - `DocVueHooks` slices allow injecting code blocks into generated Vue templates

## SQL (`DocSql`)

`DocSql` controls SQL generation:

- triggers (`IsBeforeTrigger`, `IsAfterTrigger`, `CustomTriggers`, `IsNotifyEventTrigger`)
- constraints (`CheckConstrains`, `UniqConstrains` with custom error messages)
- hooks (`DocSqlHooks`) allow injecting SQL fragments into the generated templates
- `Methods map[string]*DocSqlMethod`
  - doc-specific SQL methods that become callable via generated server endpoints

## special features

- **state machine**
  - `StateMachine != nil` marks the doc as a state machine doc
  - affects:
    - extra SQL templates (`action.sql`, `create.sql`)
    - overrides `sql_function_update.sql` and `webClient_item.vue` templates with state-machine versions
- **recursion**
  - `IsRecursion` adds fields like `parent_id` and `is_folder` and changes some UI titles
- **tasks**
  - `IsTaskAllowed` enables attaching tasks to the doc; affects `config.js` task tables list injection

## integrations

- `Integrations.Odata` activates per-doc OData codegen and adds an `odataDoc.go`-based template.

## typical project layout

A common project organization is:

- a generator entrypoint (e.g. `projectTemplate/main.go`) aggregates docs into `ProjectType.Docs`
- each doc lives in its own package under `<prefix>/<docName>/main.go` (commonly `docs/<docName>/`)
- per-doc templates live next to that package under `tmpl/`, and the doc sets `PathPrefix: "<prefix>"` so the generator resolves inferred sources like `<prefix>/<docName>/tmpl/webClient_index.vue` automatically

