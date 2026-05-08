# Template FuncMap reference

Templates use a shared function map defined in `templates/main.go` as `funcMap`.

This page is a quick reference for the “framework-provided” template helpers you can rely on across Vue/SQL/project templates.

## string helpers

- `ToUpper`, `ToLower`
- `ToCamel`, `ToLowerCamel`
- `UpperCaseFirst` (implemented in `utils/UpperCaseFirst`)

## Vue helpers

- `PrintVueFldTemplate(fld)`
  - renders the correct Vue control for a field (`q-input`, `q-select`, ref-search, files/images, jsonList, custom compositions)
- `PrintFldSelectOptions(doc, fldName)`
  - prints select options for a doc field by name

## misc helpers

- `ArrayStringJoin([]string)`
  - joins a slice into a JS-like quoted list
- `GetPgTimeZone()`
  - reads `project.Config.Postgres.TimeZone` when project context is set
- `StringContainsQuote(str)`
  - helper for template quoting decisions

## extension points

Doc templates can provide their own `FuncMap` on `DocTemplate.FuncMap`. The generator merges it with the global map so you can add doc-specific helpers without losing base helpers.

