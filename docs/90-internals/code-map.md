# Code map (maintainers)

This is a package-by-package map of where behavior lives.

## entrypoints

- `main.go`: `Start` + `readData` + `copyFiles` (the top-level generator orchestration)

## project model

- `types/`: `ProjectType`, `DocType`, `FldType` and builder helpers
  - `types/typeProject.go`
  - `types/typeDoc.go`
  - `types/typeFld.go`
  - `types/shortcuts.go`

## template engine (generator-side)

- `templates/main.go`: parses templates and attaches doc templates; contains the global template `FuncMap`
- `templates/project.go`: writes project-level generated files and runs “step2” generators
- `templates/sql/`, `templates/webClient/`, `templates/project/`: template trees
- `templates/integrations/`: integration-specific templates (e.g. `odata/`, `telegram/`)

## static trees copied into generated apps

- `sourceFiles/`: backend/runtime utilities copied into generated apps (with import rewriting)
- `webClient/`: Quasar 2 SPA skeleton copied into generated apps

## shared utilities

- `utils/`: filename/path mapping, string helpers, byte compare helpers used by generator and template funcs

