# Server-Sent Events (SSE)

The generated backend includes an SSE endpoint for pushing live updates to the SPA.

## activation

SSE is treated as always present in the template/runtime surface (not gated by a single config flag).

## what it generates

Typical generated output includes:

- `../src/sse/` (copied from `sourceFiles/src/sse/`)
- a server route (commonly `/api/sse`) wired in generated `webServer/main.go` (template-driven)

## generator touchpoints (source of truth)

- `sourceFiles/src/sse/`
- `templates/project/webServer/main.go` (route wiring)

## notes / caveats

- keep the client subscription behavior aligned with server wiring (see `docs/refactor-backlog.md` note on “SSE vs Quasar 2 client”).

