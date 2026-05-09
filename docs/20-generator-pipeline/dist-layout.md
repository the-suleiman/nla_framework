# Generated app layout (`../src`)

Generation writes relative to the caller’s working tree. By default, `project.DistPath` is set to `../src` in `main.go` `readData`.

## top-level output

- `../src/types/`
  - generated types + config loader (from `templates/project/types/`)
- `../src/webServer/`
  - Gin server, auth, API wiring (templates + copied `sourceFiles/src/webServer/`)
- `../src/sql/model/`
  - table model TOML files, one folder per table
  - **reserved/base** folders are written by `templates/project.go` `WriteProjectFiles` and use low numbers: `01_User/`, `02_UserAuth/`, `03_UserTempEmailAuth/`, `04_File/`
  - **doc-generated** folders use `<docIndex+10>_<SnakeToCamel(docName)>` (see `utils.ParseDocTemplateFilename`), so the first user doc lands at `10_<DocName>/`, the second at `11_<DocName>/`, and so on
- `../src/sql/template/function/`
  - SQL functions/triggers; generic ones at the top, and **per-doc** under `_<SnakeToCamel(docName)>/`
  - examples: doc `user_auth` -> `_UserAuth/`, doc `work_time` -> `_WorkTime/` (see `utils.ParseDocTemplateFilename`)
- `../src/webClient/`
  - Quasar 2 SPA
  - part comes from templates (`templates/project/webClient/`)
  - part is copied from the static skeleton (`webClient/`)
- `../src/pg/`, `../src/jobs/`, `../src/utils/`, `../src/sse/`
  - copied from `sourceFiles/`

## generation notes

- `../src/sql/model` is deleted at the start of each run to avoid numbering collisions (`removeOldFiles`). this is also why doc-generated folders are offset by `+10`: numbers `01..` are reserved for base/system tables that don't get regenerated from `DocType`.
- copied Go sources are import-rewritten to your module path (`Config.LocalProjectPath`), so the generated app builds under your project’s module.

## worked example

For a project that declares two user docs (`Docs: []DocType{ first_doc, second_doc }`), output is:

- base/reserved models written by `WriteProjectFiles`:
  - `01_User/`, `02_UserAuth/`, `03_UserTempEmailAuth/`, `04_File/`
- doc-generated models, indexed in `Docs` declaration order:
  - `10_FirstDoc/`, `11_SecondDoc/`
- per-doc SQL function folders, camel-cased from each `DocType.Name`:
  - `_User/`, `_UserAuth/`, `_UserTempEmailAuth/`, `_File/`, `_FirstDoc/`, `_SecondDoc/`

