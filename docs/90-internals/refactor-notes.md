# Refactor notes (maintainers)

This page is a curated entrypoint into the existing backlog and decisions log.

- primary backlog: [`docs/refactor-backlog.md`](../refactor-backlog.md)

## current decisions that affect docs

- generated runtime config is a subset of `ProjectConfig` (generated `types/config.go` only reads a subset)
- dev vs docker behavior is mostly env-var driven in the generated runtime

## candidate next refactors (high signal)

- align SSE server wiring with what the Quasar client actually subscribes to
- improve `ProjectConfig` ↔ generated runtime `Config` parity (single source of truth / generator)
- add a test strategy (golden fixture generation)

