# Quickstart

This quickstart shows the **minimal** way to run the generator: define a `types.ProjectType` and call `nla_framework.Start(p, modifyFunc)`.

## prerequisites

- a Go module with a `go.mod` (used to infer `Config.LocalProjectPath` unless you set it explicitly)
- a place where the generator is allowed to write output under `../src`

## minimal example

Create a small `main.go` in your project (outside this framework repo) and run it.

```go
package main

import (
	nla "github.com/the-suleiman/nla_framework"
	t "github.com/the-suleiman/nla_framework/types"
)

func main() {
	p := t.ProjectType{
		Name: "demo",
		Config: t.ProjectConfig{
			Postgres: t.PostrgesConfig{
				DbName:    "demo",
				Password:  "demo",
				Port:      5432,
				Host:      "localhost",
				TimeZone:  "Europe/Moscow",
				Version:   "18",
				Command:   "",
			},
			WebServer: t.WebServerConfig{
				Port: 8080,
				Url:  "localhost",
			},
			Email: t.EmailConfig{
				Sender:   "noreply@example.com",
				Host:     "smtp.example.com",
				Port:     587,
				Password: "change-me",
			},
		},
		Docs: []t.DocType{
			{
				Name:           "client",
				NameRu:         "клиент",
				IsBaseTemplates: t.DocIsBaseTemplates{Vue: true, Sql: true},
				Flds: []t.FldType{
					t.GetFldTitle(),
					t.GetFldString("phone", "телефон", 30, [][]int{{2, 1}}, "col-4"),
				},
			},
		},
	}

	// modifyFunc is optional: you can mutate copied/static files before they’re written.
	nla.Start(p, nil)
}
```

## what you should see

Generation writes to `../src` by default. Key output folders:

- `../src/webServer/` (backend)
- `../src/sql/` (models + functions)
- `../src/webClient/` (Quasar 2 SPA)

See the detailed map: [generated app layout](../20-generator-pipeline/dist-layout.md).

## common pitfalls

- **project name must not contain spaces** (validated in `readData`)
- if phone auth is disabled, **email auth is enabled by default** (see `Start` in `main.go`)
- if you run the generator outside a module tree, set `Config.LocalProjectPath` explicitly (otherwise `FillLocalPath` will fail while searching for `go.mod`)

