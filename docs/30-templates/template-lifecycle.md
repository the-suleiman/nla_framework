# Template lifecycle

This page explains how templates are discovered, parsed, and written, and how template names map to output paths.

## where templates live

- `templates/project/`: project-level templates (backend, docker, SQL seeds, Quasar shell)
- `templates/webClient/doc/`: shared doc-level Vue templates (`index.vue`, `item.vue`, tabs, etc.)
- `templates/sql/` and `templates/sql/function/`: shared doc-level SQL templates
- `templates/integrations/*`: integration-specific templates (e.g. `odata/`, `telegram/`)

## delimiters

The generator uses both delimiter styles:

- `[[ ... ]]` (common for doc/vue and many project templates)
- `{{ ... }}` (used by some project and SQL templates)

Delimiter choice is set at parse time inside `templates.ParseTemplates`.

## parse and attach rules

Source: `templates/main.go` `ParseTemplates`.

High level:

1. parse global template files into a map (prefix + filename)
2. for each `DocType`:
   - parse user-provided templates from `DocType.Templates[*].Source`
   - decide base template set based on `DocType.IsBaseTemplates`
   - attach additional templates for tabs, recursion, integrations, and certain field types

## template name → output path

For doc templates, output destinations are derived from:

- template key (e.g. `webClient_item.vue`, `sql_function_list.sql`)
- doc name
- doc index (for numbered model directories)

Mapping logic lives in:

- `utils.ParseDocTemplateFilename` (see `utils/main.go`)

## overrides

There are two override mechanisms:

- project-level: `ProjectType.OverridePathForTemplates` (applied in `templates/project.go` `ReadTmplAndPrint`)
- doc-level: `DocType.TemplatePathOverride` (used by some doc integrations to override a specific template source)

See: [template overrides](overrides.md).

