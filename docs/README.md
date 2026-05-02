# nla_framework documentation

First-pass docs for maintainers after the legacy cleanup (Quasar 2 only, Bitrix removed).

## Contents

| Doc | Purpose |
|-----|---------|
| [architecture.md](architecture.md) | Generator entrypoint and pipeline |
| [project-model.md](project-model.md) | `ProjectType`, `DocType`, `FldType` |
| [template-system.md](template-system.md) | Template dirs, naming, post-processing |
| [generated-app.md](generated-app.md) | Shape of output under `../src` |
| [integrations-inventory.md](integrations-inventory.md) | Optional integrations still in tree |
| [refactor-backlog.md](refactor-backlog.md) | Follow-up refactors and removed surface |

External references (historical): [framework.nl-a.ru](https://framework.nl-a.ru), [old docs mirror](https://pepelazz.github.io/nla-framework-docs).

## Module layout

This repository is consumed as a Go library (`github.com/the-suleiman/nla_framework`). There is **no `go.mod` at the repo root** in many setups; consumers wire the module via their own `go.mod` / `replace`. Generated apps copy `sourceFiles` and rewrite imports to `Config.LocalProjectPath`.
