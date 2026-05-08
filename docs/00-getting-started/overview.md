# Getting started overview

`nla_framework` is a **code generator**: you describe your app as a `types.ProjectType` (documents/entities + config), then the generator writes a full application tree (backend + SQL + Quasar 2 SPA).

## what it generates

By default, generation writes to **`../src`** relative to where the generator is run (see `project.DistPath` in `main.go` `readData`).

At a high level, output includes:

- `../src/webServer/`: Gin web server + auth + API endpoints (templates + copied `sourceFiles/`)
- `../src/sql/model/`: TOML model files per table
- `../src/sql/template/function/`: SQL functions/triggers templates per doc
- `../src/webClient/`: Quasar 2 (Vue 3) SPA (templates + copied `webClient/` skeleton)
- other copied utilities: `../src/pg/`, `../src/jobs/`, `../src/utils/`, `../src/sse/`

See the dedicated layout page: [generated app layout](../20-generator-pipeline/dist-layout.md).

## how generation works (mental model)

The main entrypoint is `Start(p, modifyFunc)` in `main.go`.

Roughly:

1. Fill defaults and validate project/doc/field rules (`readData`)
2. Parse global + per-doc templates (`templates.ParseTemplates`)
3. Write project-level templates (`templates.WriteProjectFiles`)
4. Execute per-doc templates (`templates.ExecuteToFile`)
5. Copy static trees (`sourceFiles/`, `webClient/`) with rewrites and slot injection
6. Run a second-pass generator for some overlay files (`OtherTemplatesGenerate`)

The full step-by-step is in [pipeline](../20-generator-pipeline/pipeline.md).

## constraints and current decisions

- **frontend**: Quasar 2 SPA (Vue 3). No alternate frontend stack selection.
- **output tree**: generator assumes it owns `../src` (it removes `../src/sql/model` on each run so numbering can change safely).
- **module path rewriting**: copied Go files rewrite imports from `github.com/the-suleiman/nla_framework` to `ProjectType.Config.LocalProjectPath` (which defaults to “module path + `/src`” inferred from `go.mod`).

