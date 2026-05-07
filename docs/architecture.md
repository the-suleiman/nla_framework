# Architecture

## Public API

- **`Start(p types.ProjectType, modifyFunc)`** — [`main.go`](../main.go): validates defaults, calls `readData`, parses templates, writes project + doc files, copies static trees from [`sourceFiles/`](../sourceFiles/) and [`webClient/`](../webClient/) (into `../src/webClient/`), runs secondary generators (`OtherTemplatesGenerate`).

## Pipeline (order)

1. **`readData`** — [`main.go`](../main.go): sets `DistPath` to `../src`, fills doc templates, grids, Vue fields, i18n; validates fields and email config.
2. **`templates.ParseTemplates`** — [`templates/main.go`](../templates/main.go): loads `project_*`, `webClient_*`, `sql_*` templates; attaches per-doc templates and integration templates.
3. **`removeOldFiles`** — deletes `../src/sql/model` so model numbering can change cleanly.
4. **`templates.WriteProjectFiles`** — [`templates/project.go`](../templates/project.go): renders project templates and explicit paths (types, webServer, SQL seeds, Quasar shell).
5. **Per-doc `ExecuteToFile`** — emits each doc’s template map and skips unchanged `webClient` output so Quasar does not restart unnecessarily.
6. **`copyFiles(sourceFiles)`** — copies backend/sql utilities into the generated tree with import rewrites and slot injection (`routes.js`, `sidemenu`, `config.js`, etc.). The copy path uses current `os.ReadFile` / `os.WriteFile` APIs; no `io/ioutil` remains.
7. **`copyFiles(webClient → ../src/webClient/)`** — copies the Quasar 2 SPA skeleton from the framework’s [`webClient/`](../webClient/) tree, also preserving the unchanged-file optimization for `webClient`.
8. **`OtherTemplatesGenerate`** — [`templates/project.go`](../templates/project.go): task list regeneration, `utils.js` / `i18n.js` overlays ([`templates/tmplGenerateStep2/`](../templates/tmplGenerateStep2/)).

## Frontend bundle

- Only one Quasar 2 SPA is supported; there is no alternate frontend or version switch. The static skeleton lives in [`webClient/`](../webClient/); project shell templates in [`templates/project/webClient/`](../templates/project/webClient/); per-doc Vue under [`templates/webClient/doc/`](../templates/webClient/doc/).

## Go module hygiene

- The framework root has a [`go.mod`](../go.mod). It includes direct dependencies needed by framework-side helpers, including `golang.org/x/text` for Unicode-aware title casing in template-exposed string helpers.
- Generated-app source trees under [`sourceFiles/`](../sourceFiles/) are still copied into the target project and import-rewritten during generation; the root module mirrors only the package surface needed for local framework checks.
