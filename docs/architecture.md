# Architecture

## Public API

- **`Start(p types.ProjectType, modifyFunc)`** — [`main.go`](../main.go): validates defaults, calls `readData`, parses templates, writes project + doc files, copies static trees from [`sourceFiles/`](../sourceFiles/) and [`webClient/quasar_2/`](../webClient/quasar_2/), runs secondary generators (`OtherTemplatesGenerate`).

## Pipeline (order)

1. **`readData`** — [`main.go`](../main.go): sets `DistPath` to `../src`, fills doc templates, grids, Vue fields, i18n; validates fields and email config.
2. **`templates.ParseTemplates`** — [`templates/main.go`](../templates/main.go): loads `project_*`, `webClient_*`, `sql_*` templates; attaches per-doc templates and integration templates.
3. **`removeOldFiles`** — deletes `../src/sql/model` so model numbering can change cleanly.
4. **`templates.WriteProjectFiles`** — [`templates/project.go`](../templates/project.go): renders project templates and explicit paths (types, webServer, SQL seeds, Quasar shell).
5. **Per-doc `ExecuteToFile`** — emits each doc’s template map.
6. **`copyFiles(sourceFiles)`** — copies backend/sql utilities into the generated tree with import rewrites and slot injection (`routes.js`, `sidemenu`, `config.js`, etc.).
7. **`copyFiles(webClient/quasar_2)`** — copies the Quasar 2 SPA skeleton into `../src/`.
8. **`OtherTemplatesGenerate`** — [`templates/project.go`](../templates/project.go): task list regeneration, `utils.js` / `i18n.js` overlays ([`templates/tmplGenerateStep2/`](../templates/tmplGenerateStep2/)).

## Frontend bundle

- Only one Quasar 2 SPA is supported. Generator paths reference `quasar_2` literally; there is no version selector. Static skeleton lives in [`webClient/quasar_2/`](../webClient/quasar_2/), project-shell templates in [`templates/project/webClient/quasar_2/`](../templates/project/webClient/quasar_2/), and per-doc templates in [`templates/webClient/quasar_2/`](../templates/webClient/quasar_2/).
