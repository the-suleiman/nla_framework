# Testing and verification (maintainers)

The generator is easiest to validate via **golden output**: run it against a known `ProjectType` input, then inspect diffs in the generated `../src` tree.

## what to verify after changes

- **generator flow**: `Start` still runs end-to-end and writes expected files
- **template parsing**: no missing template errors; expected delimiters per directory
- **write optimizations**: unchanged `webClient/` files are not rewritten (to avoid Quasar restarts)
- **copy rewrites**:
  - Go import rewrite from `github.com/the-suleiman/nla_framework` → `Config.LocalProjectPath`
  - slot injection still lands in `config.js`, `routes.js`, `sidemenu/index.vue`
- **integrations**: if you touched integration code, validate both “off” and “on” paths

## recommended approach

- keep one small “fixture project” (minimal docs + 1–2 fields) and generate into a scratch tree
- keep one “kitchen sink” fixture (tabs, tasks, files/images, optional integration toggles) for broader coverage

This repository currently does not include `*_test.go` golden tests; see `docs/refactor-backlog.md` for candidates to add them.

