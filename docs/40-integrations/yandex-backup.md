# Yandex Disk backup integration

This integration generates a standalone backup tool/service that uploads Postgres dumps to Yandex Disk.

## activation

- `ProjectType.IsBackupOnYandexDisk()` → `len(Config.Backup.ToYandexDisk.Token) > 0`

## what it generates

When enabled, `templates/project.go` generates:

- `../src/yandexDiskBackup/main.go`
- `../src/yandexDiskBackup/yandexApi.go`
- `../src/yandexDiskBackup/dbBackup.go`
- systemd service file: `../src/yandexDiskBackup/<dbName>_yandexBackup.service`
- helper script: `../src/yandexDiskBackup/startYandexBackupService.sh`

## generator touchpoints (source of truth)

- `types/typeProject.go`: `BackupConfig`, `BackupConfigYandexDisk`, `IsBackupOnYandexDisk`
- `templates/project.go`: conditional generation for `yandexDiskBackup/`
- template sources under `templates/project/yandexDiskBackup/`

## notes

- backup approach is logical dump based (portable across Postgres versions); see `docs/refactor-backlog.md` for details and compose/version caveats.

