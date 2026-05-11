# Shared Vue controls (`webClient/`)

The framework repo ships a Quasar SPA skeleton under [`webClient/`](../../webClient/). It is copied into generated apps as `../src/webClient/`. This page documents selected **shared** controls that back generated doc fields (not every component under `webClient/src/app/components/`).

## files field: `comp-fld-files`

- **source**: [`webClient/src/app/components/common/list/compFldFiles.vue`](../../webClient/src/app/components/common/list/compFldFiles.vue)
- **registration**: global component `comp-fld-files` in [`webClient/src/app/components/common/index.js`](../../webClient/src/app/components/common/index.js)

### generated markup

`PrintVueFldTemplate` emits `<comp-fld-files …/>` for `FldVue.Type == "files"` (see [`templates/main.go`](../../templates/main.go)). The `ext` object always includes `tableName` and `tableId` for upload/delete wiring.

### upload and download URLs

Uploads use:

- `` `${this.$config.apiUrl()}/api/upload_file` ``

Downloads and image previews use:

- `` `${this.$config.apiUrl()}${item.url}` ``

So file storage is assumed to live on the **same API origin** as the app (not a separate hard-coded files host).

### default size cap

If `ext.maxFileSize` is absent, the uploader defaults to **10_000_000** bytes (10 MB). `GetFldFiles` can set `Vue.Ext["maxFileSize"]` from Go (`FldVueFilesParams.MaxFileSize`); the component parses it as an integer string.

### props (high level)

| prop | role |
|------|------|
| `fld` | current file list (json array) |
| `fldName` | field name for `_update` payload (default `files`) |
| `label` | bar title |
| `readonly` | hides add/delete/drop |
| `ext` | must include `tableName`, `tableId`; optional `accept`, `maxFileSize` |
| `vif` | kept for API compatibility with other field components |
| `showUploaderDialog` | set **`true` explicitly** to use “add opens dialog” only; omit or use falsy for **compact bar** uploader with drag-and-drop on the bar |

There is **no default** for `showUploaderDialog`: generated templates omit it, so new docs get the **inline bar** uploader. Pass `:show-uploader-dialog="true"` in a custom template when you want dialog-only upload.

### behavior notes

- **Multiple files** per batch; the list is persisted with `postCallPgMethod` on the uploader **`finish`** event (after the queue completes), not on each individual `uploaded` response.
- **Drag-and-drop** on the bar adds files to the inline uploader, or opens the dialog and adds files there when `showUploaderDialog === true`.
- **Images** (common extensions): row shows a prefetched thumbnail; click opens a preview dialog with download. Thumbnails use blob URLs that are **revoked** when rows disappear or the component is destroyed.
- **Two refs**: `uploaderInline` and `uploaderDialog` so Vue does not overwrite a single `ref` when both exist in the tree.

### axios blob errors

The component registers a **global** axios response interceptor in `mounted` to turn JSON error bodies hidden inside `blob` responses into parsed objects (for failed downloads). If several file fields mount, interceptors stack; prefer moving this to app bootstrap if you hit odd behavior.

## static image URLs: `comp-stat-img-src`

- **source**: [`webClient/src/app/components/common/utils/statImgSrc.vue`](../../webClient/src/app/components/common/utils/statImgSrc.vue)
- **registration**: `comp-stat-img-src` in the same `index.js`

### URL resolution (`src` prop)

- `blob:` and `data:` URLs are used as-is.
- URLs matching **`/^https?:\/\//i`** are used as-is.
- Paths starting with **`image`** (e.g. bundled `image/...` assets) are used as-is.
- Any other relative path is prefixed with **`this.$config.apiUrl()`**.

If `src` is empty, a placeholder image URL is used (same stock URL as before this behavior split).

## i18n strings used by `comp-fld-files`

Defined in [`types/i18n.go`](../../types/i18n.go) under `message` (ru/en):

- `select_file_for_upload`, `upload_error`, `delete`, `cancel`, `download`, `upload_file_hint` (not used by the stock component after tooltip simplification; safe for custom templates)
- `only_files_with_extension`, `only_files_no_larger_than` (Russian copy uses **расширением**)
- `file_not_found_by_link`, `file_rejected_constraints`

Regenerate or merge i18n in consuming apps after changing `FillI18n()` output.
