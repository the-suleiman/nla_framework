# Feature inventory (docs coverage map)

This page is the **coverage checklist**: every major capability should have a doc page, and every doc page should point to the **source of truth** in the repository (Go code, templates, or copied static trees).

## generator entrypoint and pipeline

- **Start / readData / copyFiles**
  - **Docs**: `docs/20-generator-pipeline/pipeline.md`, `docs/20-generator-pipeline/copy-and-rewrite.md`
  - **Source of truth**:
    - `main.go` (`Start`, `readData`, `copyFiles`, `removeOldFiles`)
    - `templates/main.go` (`ParseTemplates`, `ExecuteToFile`)
    - `templates/project.go` (`WriteProjectFiles`, `OtherTemplatesGenerate`, `ReadTmplAndPrint`)

## project model (what you can describe)

- **ProjectType / ProjectConfig / ProjectVue / ProjectSql / ProjectGo**
  - **Docs**: `docs/10-project-model/projecttype.md`, `docs/00-getting-started/configuration-map.md`
  - **Source of truth**:
    - `types/typeProject.go`
    - `templates/project.go` (which config fields affect generated output)

- **DocType / DocVue / DocSql / templates per doc**
  - **Docs**: `docs/10-project-model/doctype.md`
  - **Source of truth**:
    - `types/typeDoc.go` (`DocType`, `DocVue`, `DocSql`, `DocTemplate`, `.Init()`)
    - `templates/main.go` (how base templates are attached; tabs; state machine overrides)

- **FldType / SQL mapping / Vue mapping / builder sugar**
  - **Docs**: `docs/10-project-model/fldtype.md`, `docs/10-project-model/builders-and-shortcuts.md`
  - **Source of truth**:
    - `types/typeFld.go` (`FldType`, `FldSql`, `FldVue`, conversions like `GoType`, `PgInsertType`)
    - `types/shortcuts.go` (field builders like `GetFldTitle`, `GetFldRef`, `GetFldFiles`, etc.)
    - `templates/main.go` (`PrintVueFldTemplate` and related Vue rendering helpers)

## templates system

- **Template discovery, delimiters, naming → output paths**
  - **Docs**: `docs/30-templates/template-lifecycle.md`
  - **Source of truth**:
    - `templates/main.go` (`ParseTemplates`, base template selection rules)
    - `utils/main.go` (`ParseDocTemplateFilename`)
    - `templates/project.go` (`ReadTmplAndPrint`, `OverridePathForTemplates` behavior)

- **Template helper functions (FuncMap)**
  - **Docs**: `docs/30-templates/funcmap-reference.md`
  - **Source of truth**:
    - `templates/main.go` (global `funcMap` and helper functions it points to)
    - `utils/main.go` (`UpperCaseFirst`, template-adjacent helpers)

## generated app shape (what gets written/copied)

- **Generated folders and responsibilities**
  - **Docs**: `docs/20-generator-pipeline/dist-layout.md`
  - **Source of truth**:
    - `templates/project/` (template-written files)
    - `sourceFiles/` (copied backend/runtime utilities)
    - `webClient/` (copied Quasar 2 SPA skeleton)

## copy / rewrite / slot injection

- **Import rewrite and special file hooks**
  - **Docs**: `docs/20-generator-pipeline/copy-and-rewrite.md`
  - **Source of truth**:
    - `main.go` (`copyFiles`, `configJsModify`, `routesJsModify`, `sidemenuJsModify`)

## optional integrations

- **OData**
  - **Docs**: `docs/40-integrations/odata.md`
  - **Source of truth**:
    - `types/typeProject.go` (`IsOdataIntegration`)
    - `types/typeDoc.go` (`DocIntegrations.Odata`, `IsOdataIntegration`)
    - `templates/project.go` (project-level `odata/` generation)
    - `templates/docIsIntegration.go` + `templates/integrations/odata/`

- **Telegram**
  - **Docs**: `docs/40-integrations/telegram.md`
  - **Source of truth**:
    - `types/typeProject.go` (`IsTelegramIntegration`)
    - `templates/project.go` (`telegramAuth.go`, SQL functions, `tgBot/`)
    - `main.go` (`configJsModify` injects telegram config into `config.js`)

- **Graylog**
  - **Docs**: `docs/40-integrations/graylog.md`
  - **Source of truth**:
    - `types/typeProject.go` (`ProjectConfig.Graylog`)
    - `templates/project/` (compose logging + main wiring)
    - `sourceFiles/src/graylog/`

- **Phone auth**
  - **Docs**: `docs/40-integrations/phone-auth.md`
  - **Source of truth**:
    - `types/typeProject.go` (`AuthConfig.ByPhone`, `ByEmail`)
    - `main.go` (default `ByEmail` when phone is off)
    - `templates/project.go` (extra SQL + Vue + `phone.go`)

- **Yandex Disk backup**
  - **Docs**: `docs/40-integrations/yandex-backup.md`
  - **Source of truth**:
    - `types/typeProject.go` (`BackupConfig`, `IsBackupOnYandexDisk`)
    - `templates/project.go` (generated `yandexDiskBackup/` files)

- **SSE**
  - **Docs**: `docs/40-integrations/sse.md`
  - **Source of truth**:
    - `sourceFiles/src/sse/`
    - `templates/project/webServer/main.go` (route wiring)

