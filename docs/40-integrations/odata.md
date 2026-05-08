# OData / 1C integration

This integration generates OData sync code when **project-level** config is enabled and/or **doc-level** integration metadata is provided.

## activation

- project-level: `ProjectType.IsOdataIntegration()` → `len(Config.Odata.Url) > 0`
- doc-level: `DocType.IsOdataIntegration()` → `len(DocType.Integrations.Odata.Name) > 0`

## what it generates

When project-level OData is enabled, `templates/project.go` generates:

- `../src/odata/main.go`
- `../src/odata/odataQueryType.go`

When a doc has OData integration enabled, `templates/docIsIntegration.go` attaches a doc template that generates:

- `../src/odata/<docName>.go` (camel-lower file name derived from doc name)

## generator touchpoints (source of truth)

- `types/typeProject.go`: `IsOdataIntegration`, `OdataConfig`
- `types/typeDoc.go`: `DocIntegrations.Odata`, `IsOdataIntegration`
- `templates/project.go`: project-level generation of `odata/` files
- `templates/docIsIntegration.go`: per-doc generation based on OData field metadata
- `templates/integrations/odata/`: OData templates (including `odataDoc.go`)

## field mapping

Fields can store OData metadata in `FldType.IntegrationData["odata"]` as a `types.OdataFld`.

The generator uses that to:

- decide which fields are included in OData structs
- decide casting behavior (custom `CastToGoType` can override default conversions)

## notes / caveats

- OData codegen relies on both **doc-level integration config** and **field-level mapping** to be useful.
- debug mode can expose additional routes/behavior (see `DocIntegrationsOdata.IsDebugMode`).

