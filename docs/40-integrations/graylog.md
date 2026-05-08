# Graylog integration

This integration wires Graylog logging into the generated app when configured.

## activation

Graylog is controlled via the `ProjectType.Config.Graylog` block (host/port/app name/attrs).

## what it generates

Depending on templates and copied runtime helpers, generation typically includes:

- Graylog-related runtime package under `../src/graylog/` (copied from `sourceFiles/src/graylog/`)
- backend wiring (templates under `templates/project/`)
- docker compose logging configuration (templates under `templates/project/docker-compose*.yml`)

## generator touchpoints (source of truth)

- `types/typeProject.go`: `GraylogConfig`
- `sourceFiles/src/graylog/`: runtime helpers
- `templates/project/`: server + compose template wiring

## notes

- the exact on/off condition is template-driven; treat `Config.Graylog` as the configuration source and verify generated output in your fixture project.

