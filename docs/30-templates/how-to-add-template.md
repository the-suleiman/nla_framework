# How to add a template

This is a practical guide for extending generation by adding templates.

## add a doc-level template

Doc-level templates live on `DocType.Templates` as `map[string]*types.DocTemplate`.

Typical workflow:

1. create a template file in your project (outside the framework), for example under `client/tmpl/`
2. register it in `DocType.Templates` (set `Source`, or let `FillDocTemplatesFields()` infer it)
3. set `DistPath`/`DistFilename` or let filename mapping derive them from the key
4. ensure the template uses the correct delimiters (`[[` `]]` by default for doc templates)

The generator will parse and execute these in `templates.ParseTemplates` / `templates.ExecuteToFile`.

### inferred vs explicit destinations

When a `DocTemplate` entry leaves fields empty, `ProjectType.FillDocTemplatesFields` and `utils.ParseDocTemplateFilename` derive defaults from the template key:

- `Source` -> `<DocType.PathPrefix>/<snakeToCamelLower(docName)>/tmpl/<templateKey>`
- `DistPath` / `DistFilename` -> standard locations for `webClient_*` (under `webClient/src/app/components/<doc>/...`) and `sql_*` (under `sql/model/...` or `sql/template/function/_<DocCamel>/`)

Default mapping covers the common cases (`webClient_index.vue`, `webClient_item.vue`, `sql_function_list.sql`, `sql_function_update.sql`, etc.). For nested component folders or non-standard filenames, set `Source`, `DistPath`, and `DistFilename` explicitly.

### example: mixing inferred and explicit

A typical doc combines both styles:

```go
PathPrefix: "docs",
Templates: map[string]*t.DocTemplate{
    // inferred source via PathPrefix + key, default destination
    "sql_function_update.sql": {},
    "webClient_item.vue":      {},
    "webClient_index.vue":     {},

    // explicit source + nested destination
    "webClient_customNested.vue": {
        Source:       "docs/<docName>/tmpl/nested/webClient_customNested.vue",
        DistPath:     "../src/webClient/src/app/components/<docName>/comp/nested",
        DistFilename: "customNested.vue",
    },
    "webClient_customComp.vue": {
        Source:       "docs/<docName>/tmpl/webClient_customComp.vue",
        DistPath:     "../src/webClient/src/app/components/<docName>/comp",
        DistFilename: "customComp.vue",
    },
},
```

Notes:

- short-form entries rely on `PathPrefix` + `snakeToCamelLower(DocType.Name)` to resolve sources under `<PathPrefix>/<docCamelLower>/tmpl/`.
- explicit entries are required when the destination is a deeper folder (e.g. `.../comp/<sub>`) or when the template key wouldn't naturally map to that path via `ParseDocTemplateFilename`.
- the `webClient_*` / `sql_*` key prefixes still matter for explicit entries: the `[[` `]]` vs `{{` `}}` delimiter set is selected by key prefix in `templates.ParseTemplates`.

## add a project-level template

Project-level templates are written by `templates.WriteProjectFiles` via `ReadTmplAndPrint`.

If you want to replace an existing project template without forking the framework, use:

- `ProjectType.OverridePathForTemplates["/some/dist/path/filename.ext"] = "/path/to/your/template"`

See: [template overrides](overrides.md).

## adding new built-in framework templates (maintainers)

If you’re changing the framework itself:

- add the template file under `templates/` (appropriate subtree)
- wire it into `templates.ParseTemplates` (for parsing) or `templates.WriteProjectFiles` / `ReadTmplAndPrint` (for writing)
- update docs pages:
  - `docs/30-templates/template-lifecycle.md`
  - `docs/90-internals/feature-inventory.md`

