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

