# `FldType` reference

`FldType` describes one field: both how it is stored in SQL and how it is rendered/edited in Vue.

Source: `types/typeFld.go`.

## core fields

- `Name`, `NameRu`
- `Type`
  - generator-level type string (e.g. `string`, `int`, `date`, `jsonb`, `uuid`, `text[]`, etc.)
- `Sql FldSql`
  - column and constraint behavior (`IsRequired`, `IsUniq`, `Size`, `Default`, `Ref`, etc.)
- `Vue FldVue`
  - UI behavior: component type, options for selects, validation flags, layout classes, etc.

## SQL mapping helpers

`FldType` includes methods that affect how SQL templates treat the field:

- `GoType()`: Go type for generated structs (e.g. `double` → `float64`)
- `PgInsertType()`: SQL type for inserts (e.g. `date`/`datetime` → `timestamp`)
- `PgUpdateType()`: UI/JSON coercion type used by generated update functions

`PrintPgModel()` renders the TOML model line used in `../src/sql/model/**/main.toml`.

## Vue rendering

The framework renders Vue controls primarily via template helpers in `templates/main.go`:

- `PrintVueFldTemplate(fld FldType)` produces the `<q-input>`, `<q-select>`, or custom component markup based on:
  - `FldVue.Type` (or fallback to `FldType.Type`)
  - `FldSql.Ref` (turns `int` into a “ref” selector)
  - `FldVue.Composition` (custom renderer)
  - `FldVue.Ext` (extra params for some widgets)

For **`FldVue.Type == "files"`** (`FldVueTypeFiles`), the template emits `<comp-fld-files … :ext="{ tableName, tableId, … }"/>`. Runtime behavior (upload URL, drag-and-drop, optional dialog-only mode, image previews) lives in the shared Vue implementation. See [shared Vue controls](../30-templates/vue-shared-controls.md). Image fields use `comp-stat-img-src` for resolved URLs; that helper is documented on the same page.

## reserved/validated rules

During `readData` (`main.go`), the generator rejects any field named:

- `user_id`

because it is treated as a framework internal/service field.

## integrations metadata

`IntegrationData map[string]interface{}` can store integration-specific metadata.

Example:

- OData uses `IntegrationData["odata"]` with a `types.OdataFld` value (set via `FldType.SetOdataInfo` in `types/typeFld.go`).

