# Builders and shortcuts (`types/shortcuts.go`)

The `types/shortcuts.go` file provides “sugar” functions to construct common `FldType` and `DocType` structures quickly.

Source: `types/shortcuts.go`.

## common field builders

- `GetFldTitle(...)`
  - creates `title` with `IsRequired`, `IsUniq`, `IsSearch`, and default size
- `GetFldString(name, nameRu, size, rowCol, ...)`
  - creates a basic string field with Vue layout classes and optional readonly
- `GetFldInt`, `GetFldInt64`, `GetFldDouble`, `GetFldDate`, `GetFldDateTime`
- `GetFldRef(name, nameRu, refTable, ...)`
  - creates an `int` ref field with additional UI behaviors controlled via params:
    - `isShowLink` (pathUrl/avatar filled later by `ProjectType.FillVueFlds`)
    - `isAddNew` (addNewUrl filled later)
    - `isClearable`
    - `ext:{...}` (raw JSON ext payload)
- `GetFldSelectString`, `GetFldSelectMultiple`
- `GetFldFiles`, `GetFldImg`, `GetFldImgList`
  - creates JSON-backed file/image widgets with upload constraints stored in `Vue.Ext`
- `GetFldJsonList`
  - jsonb list editor with nested fields

## doc helpers

Some helpers modify the document itself:

- `DocType.AddVueTaskAndTabs()`
  - enables tasks and adds default tabs
- `DocType.SetIsRecursion(title)`
  - enables recursion and injects recursion-specific fields (`parent_id`, `is_folder`)
- `VueTab.AddCounter(...)`
  - adds a counter mixin and generates a `tabCounter*.js` template for the doc

## when to use shortcuts vs manual structs

- use shortcuts for common fields to get consistent Vue+SQL defaults
- use manual `FldType` / `DocType` when you need precise control over hooks/templates or unusual Vue compositions

