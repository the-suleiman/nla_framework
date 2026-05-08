# Template overrides

There are two override layers for template sources.

## 1) project-level overrides (`ProjectType.OverridePathForTemplates`)

Source of truth: `templates/project.go` `ReadTmplAndPrint`.

- key: destination path + filename, formatted as `"<distPath>/<filename>"`
  - example destination key: `"/webClient/index.html"`
- value: filesystem path to the template source file to parse instead of the default

This affects project-level generation performed via `ReadTmplAndPrint`.

## 2) doc-level overrides (`DocType.TemplatePathOverride`)

Some doc features/integrations parse templates directly and allow overriding the template source path by name.

Example: OData uses a template named `odataDoc.go` and checks `DocType.TemplatePathOverride["odataDoc.go"]`.

## when to use which

- use **project-level** overrides when you want to replace a framework template for a generated file like `webServer/main.go` or `webClient/index.html`
- use **doc-level** overrides when a doc integration or feature provides an explicit override hook for a specific integration template

