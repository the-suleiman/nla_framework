# Phone auth integration

Phone auth changes both backend and frontend generation to support registration/login via SMS.

## activation

- `ProjectType.Config.Auth.ByPhone == true`

Default behavior when phone auth is off:

- `Start` sets `Auth.ByEmail = true`, meaning email auth is the default.

## what it generates

When enabled, `templates/project.go` generates additional files:

- backend:
  - `../src/webServer/auth/phone.go`
- SQL helpers:
  - `user_get_by_phone_with_password.sql`
  - phone-auth create/check functions under `_UserTempEmailAuth`
- frontend components under `../src/webClient/src/app/components/auth/phone/`

## generator touchpoints (source of truth)

- `types/typeProject.go`: `AuthConfig.ByPhone`, SMS service config
- `main.go`: default `ByEmail` behavior when phone auth is disabled
- `templates/project.go`: conditional generation block under `if p.Config.Auth.ByPhone { ... }`
- phone template sources:
  - `templates/project/webServer/auth/phone.go`
  - `templates/project/webClient/app/components/auth/phone/...`

