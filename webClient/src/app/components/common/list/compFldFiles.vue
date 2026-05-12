<!--
  files doc field: labeled bar with optional list, upload (inline or dialog), drag-drop,
  image thumbnails (blob), download, and soft-delete wired to ext.tableName/tableId.
-->
<template>
  <div style="max-width: 450px;">
    <!-- bar is the drag target; color hints drop-ready vs readonly -->
    <q-bar
      @dragenter="dropWaiting = true"
      @dragleave="onDragLeave"
      @dragover.prevent="onDragOver"
      @drop.prevent="onFileDrop"
      class="text-white shadow-2"
      :class="qBarColorClass"
      :style="qBarRadiusStyle"
    >
      <div>{{ label }}</div>
      <q-space/>
      <q-btn
        dense
        flat
        :icon="expandListIcon"
        v-if="list.length && !dropWaiting"
        @click="isShowList = !isShowList"
      />
      <div v-if="!readonly">
        <!-- dialog mode: only + opens modal uploader; inline mode: hidden q-uploader + trigger btn -->
        <q-btn
          dense
          flat
          round
          icon="add"
          @click="isShowDialog = true"
          v-if="showUploaderDialog === true"
        />
        <div class="q-mr-xs q-ml-md">
          <q-uploader
            v-if="showUploaderDialog !== true"
            ref="uploaderInline"
            style="width: 0; height: 0; margin: -13px 15px -13px -13px; z-index: 99"
            autoUpload
            multiple
            :url="uploadUrl"
            :headers="headers"
            :accept="(ext && ext.accept) ? ext.accept : ''"
            :maxFileSize="effectiveMaxFileSize"
            :formFields="formFields"
            @rejected="rejected"
            @uploaded="uploaded"
            @failed="failed"
            @finish="finish"
          >
            <template v-slot:header="scope">
              <q-btn
                round
                dense
                flat
                type="a"
                @click="scope.pickFiles"
                :icon="inlineUploaderIcon"
                size="sm"
              >
                <q-uploader-add-trigger/>
              </q-btn>
            </template>
          </q-uploader>
        </div>
      </div>
    </q-bar>

    <!-- non-deleted rows only; thumbnails use blobUrls filled by fetchImagePreview -->
    <q-slide-transition :duration="100">
      <q-list
        bordered
        separator
        v-if="filteredList.length && isShowList"
        style="border-radius: 0 0 10px 10px"
      >
        <q-item v-for="item in filteredList" :key="item.filename">
          <q-item-section
            avatar
            @click="isImage(item.filename) ? openInDialog(item) : downloadFile(item)"
            class="cursor-pointer"
          >
            <q-avatar size="md">
              <img v-if="!isImage(item.filename)" src="image/file.svg">
              <comp-stat-img-src
                v-else-if="blobUrls[item.filename]"
                :src="blobUrls[item.filename]"
                style="width: 50px; height: 50px; min-height: 50px"
              />
              <img v-else src="image/file.svg">
            </q-avatar>
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ item.filename }}</q-item-label>
          </q-item-section>
          <q-item-section side v-if="!readonly">
            <q-btn flat round icon="delete" size="sm" @click="showDeleteDialog(item.filename)">
              <q-tooltip>{{ $t('message.delete') }}</q-tooltip>
            </q-btn>
          </q-item-section>
        </q-item>
      </q-list>
    </q-slide-transition>
  </div>

  <!-- full uploader UI when showUploaderDialog is true (add btn or drop) -->
  <q-dialog v-model="isShowDialog">
    <q-uploader
      ref="uploaderDialog"
      color="secondary"
      style="border-radius: 20px"
      :label="$t('message.select_file_for_upload')"
      autoUpload
      multiple
      :url="uploadUrl"
      :headers="headers"
      :accept="ext && ext.accept ? ext.accept : ''"
      :maxFileSize="effectiveMaxFileSize"
      :formFields="formFields"
      @rejected="rejected"
      @uploaded="uploaded"
      @failed="failed"
      @finish="finish"
    />
  </q-dialog>

  <!-- marks row deleted + remove_file API; persistence same as finish -->
  <q-dialog v-model="isShowDeleteDialog" persistent>
    <q-card>
      <q-card-section class="row items-center">
        <q-avatar rounded icon="warning" color="warning" text-color="white"/>
        <span class="q-ml-sm">{{ $t('message.delete') }}?</span>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('message.cancel')" v-close-popup/>
        <q-btn flat :label="$t('message.delete')" v-close-popup @click="remove"/>
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- large preview for images; src is blob or signed path via dialogImgSrc -->
  <q-dialog v-model="isShowImageDialog" @hide="onImageDialogHide">
    <q-card style="border-radius: 20px; width: 600px; max-width: 100%">
      <q-btn
        round
        push
        dense
        color="primary"
        icon="download"
        style="position: absolute; top: 7px; right: 7px; z-index: 99"
        @click="downloadFile(dialogImage)"
      >
        <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]">
          {{ $t('message.download') }}
        </q-tooltip>
      </q-btn>
      <comp-stat-img-src :src="dialogImgSrc"/>
    </q-card>
  </q-dialog>
</template>

<script setup>
// quasar field component: see docs/30-templates/vue-shared-controls.md for contract
import axios from 'axios'
import {i18n} from 'boot/i18n'
import {useQuasar} from 'quasar'
import {
  computed,
  nextTick,
  onBeforeUnmount,
  onMounted,
  reactive,
  ref,
  watch
} from 'vue'
import config from '../../../plugins/config'
import utils from '../../../plugins/utils'

// bytes when ext.maxFileSize missing or invalid
const DEFAULT_MAX_FILE_SIZE = 10000000

const props = defineProps({
  fld: {default: undefined},
  fldName: {type: String, default: 'files'},
  label: {default: undefined},
  readonly: {type: Boolean, default: false},
  icon: {default: undefined},
  vif: {type: Boolean, default: true},
  ext: {type: Object, default: null},
  // true = add opens dialog only; omit or null/false = compact bar uploader (set explicitly when needed)
  showUploaderDialog: {default: null}
})

const emit = defineEmits(['update']) // payload: current file list json

const $q = useQuasar()
// legacy vue-i18n: useI18n() is not available; use global translator
const t = (...args) => i18n.global.t(...args)

// ui state
const isShowDialog = ref(false)
const isShowList = ref(true)
const list = ref([])
const isShowDeleteDialog = ref(false)
const selectedForDeleteFilename = ref(null)
const dropWaiting = ref(false)
const isShowImageDialog = ref(false)
const dialogImage = ref({url: null})
// object-url cache for list + dialog previews (revoked on unmount / row removal)
const blobUrls = reactive({})

// separate refs so both uploaders can exist in the template without clobbering
const uploaderInline = ref(null)
const uploaderDialog = ref(null)

let axiosBlobErrorInterceptorId = null

// multipart target for q-uploader
const uploadUrl = computed(() => `${config.apiUrl()}/api/upload_file`)

const headers = computed(() => {
  const authToken = localStorage.getItem(config.appName)
  return [{name: 'Auth-token', value: authToken}]
})

const filteredList = computed(() => list.value.filter(v => !v.deleted))

const effectiveMaxFileSize = computed(() => {
  const raw = props.ext && props.ext.maxFileSize
  if (raw == null || raw === '') return DEFAULT_MAX_FILE_SIZE
  const n = parseInt(raw, 10)
  return Number.isNaN(n) ? DEFAULT_MAX_FILE_SIZE : n
})

const formFields = computed(() => [
  {name: 'tableName', value: props.ext.tableName},
  {name: 'tableId', value: props.ext.tableId}
])

// prefer in-memory blob for dialog when we already prefetched the row
const dialogImgSrc = computed(() => {
  const d = dialogImage.value
  if (!d || !d.filename) return null
  return blobUrls[d.filename] || d.url
})

const qBarColorClass = computed(() => {
  if (props.readonly) return 'bg-secondary'
  return dropWaiting.value ? 'bg-positive' : 'bg-secondary'
})

const qBarRadiusStyle = computed(() => {
  const borderRadius = list.value.length && isShowList.value ? '10px 10px 0 0' : '10px'
  return {borderRadius}
})

const expandListIcon = computed(() => `expand_${isShowList.value ? 'less' : 'more'}`)

const inlineUploaderIcon = computed(() => (dropWaiting.value ? 'download' : 'add'))

// drop blob URLs for removed files; prefetch blobs for new image rows
watch(
  filteredList,
  (nextList) => {
    const names = new Set(nextList.map(i => i.filename))
    for (const fn of Object.keys(blobUrls)) {
      if (!names.has(fn)) {
        URL.revokeObjectURL(blobUrls[fn])
        delete blobUrls[fn]
      }
    }
    for (const item of nextList) {
      if (isImage(item.filename) && !blobUrls[item.filename]) fetchImagePreview(item)
    }
  },
  {deep: true, immediate: true}
)

// one file finished uploading into list (full persist happens in finish)
function uploaded({xhr: {response}}) {
  const res = JSON.parse(response)
  if (!res.ok) {
    $q.notify({color: 'negative', position: 'bottom', message: res.message})
  } else {
    if (!list.value.find(v => v.filename === res.result.filename && !v.deleted)) list.value.push(res.result)
  }
}

// queue drained: reset uploader, close dialog, persist whole list + notify parent
function finish() {
  if (isShowDialog.value && uploaderDialog.value) {
    uploaderDialog.value.reset()
  } else if (uploaderInline.value) {
    uploaderInline.value.reset()
  }
  isShowDialog.value = false
  utils.postCallPgMethod({
    method: `${props.ext.tableName}_update`,
    params: {id: props.ext.tableId, [props.fldName]: list.value}
  }).subscribe(() => emit('update', list.value))
}

function failed(msg) {
  let msgText = t('message.upload_error')
  if (msg.xhr && msg.xhr.responseText) {
    const res = JSON.parse(msg.xhr.responseText)
    if (res.message) msgText = res.message
  }
  $q.notify({color: 'negative', position: 'bottom', message: msgText})
}

function rejected(msg) {
  const niceBytes = (x) => {
    const units = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
    let l = 0
    let n = parseInt(x, 10) || 0
    while (n >= 1024 && ++l) n = n / 1024
    return (n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l])
  }
  let msgText = t('message.file_rejected_constraints')
  if (msg.length && msg[0].failedPropValidation === 'accept') {
    msgText = `${t('message.only_files_with_extension')}: ${props.ext.accept} `
  }
  if (msg.length && msg[0].failedPropValidation === 'max-file-size') {
    const size = niceBytes(effectiveMaxFileSize.value)
    msgText = `${t('message.only_files_no_larger_than')}: ${size}`
  }
  $q.notify({color: 'negative', position: 'bottom', message: msgText})
}

// get file as blob for q-img / dialog (auth same as download)
function fetchImagePreview(item) {
  axios({
    url: `${config.apiUrl()}${item.url}`,
    method: 'GET',
    headers: {'Auth-token': localStorage.getItem(config.appName)},
    responseType: 'blob'
  }).then((response) => {
    const blobUrl = window.URL.createObjectURL(new Blob([response.data]))
    blobUrls[item.filename] = blobUrl
  }).catch((err) => {
    let msg = err.response && err.response.data
    if (err.response && err.response.data) {
      msg = err.response.data.message
      if (msg === 'not found') msg = t('message.file_not_found_by_link')
    }
    $q.notify({message: msg, type: 'negative', position: 'top-right'})
  })
}

// one-off blob download link (object url revoked immediately after click)
function downloadFile(item) {
  axios({
    url: `${config.apiUrl()}${item.url}`,
    method: 'GET',
    headers: {'Auth-token': localStorage.getItem(config.appName)},
    responseType: 'blob'
  }).then((response) => {
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', item.filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  }).catch((err) => {
    let msg = err.response && err.response.data
    if (err.response && err.response.data) {
      msg = err.response.data.message
      if (msg === 'not found') msg = t('message.file_not_found_by_link')
    }
    $q.notify({message: msg, type: 'negative', position: 'top-right'})
  })
}

function showDeleteDialog(filename) {
  selectedForDeleteFilename.value = filename
  isShowDeleteDialog.value = true
}

// marks deleted, moves row to end of array for PG payload, hits remove_file then _update
function remove() {
  const i = list.value.findIndex(v => v.filename === selectedForDeleteFilename.value && !v.deleted)
  list.value[i].deleted = true
  const fileToken = list.value[i].url.split('/', -1).slice(-1)
  const item = list.value.splice(i, 1)
  list.value.push(item[0])
  utils.postApiRequest({url: `/api/remove_file/${fileToken}`}).subscribe((res) => {
    if (res.ok) selectedForDeleteFilename.value = null
  })
  utils.postCallPgMethod({
    method: `${props.ext.tableName}_update`,
    params: {id: props.ext.tableId, [props.fldName]: list.value}
  }).subscribe(() => {
    emit('update', list.value)
  })
}

function onDragOver(event) {
  event.dataTransfer.dropEffect = 'copy'
}

function onDragLeave(event) {
  if (!event.relatedTarget || !event.currentTarget.contains(event.relatedTarget)) dropWaiting.value = false
}

// route dropped files to dialog uploader or hidden inline uploader
function onFileDrop(event) {
  dropWaiting.value = false
  if (props.readonly) return
  const files = event.dataTransfer.files
  if (!files.length) return
  if (props.showUploaderDialog === true) {
    isShowDialog.value = true
    nextTick(() => {
      const d = uploaderDialog.value
      if (d) d.addFiles(Array.from(files))
    })
  } else {
    const u = uploaderInline.value
    if (u) u.addFiles(Array.from(files))
  }
}

function openInDialog(item) {
  dialogImage.value = item
  isShowImageDialog.value = true
}

function onImageDialogHide() {
  dialogImage.value = {url: null}
}

// thumbnail + preview dialog use this set (must match fetchImagePreview / comp-stat-img-src)
function isImage(fileName) {
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.svg', '.webp']
  return imageExtensions.some(ext => fileName.toLowerCase().endsWith(ext))
}

onMounted(() => {
  if (!props.ext) {
    throw new Error('compFldFiles missed param: "ext"')
  }
  if (!props.ext.tableName) {
    throw new Error('compFldFiles missed param: "ext.tableName"')
  }
  if (!props.ext.tableId) {
    throw new Error('compFldFiles missed param: "ext.tableId"')
  }
  list.value = props.fld || []
  // parse JSON errors returned with responseType blob (shared with download/preview)
  axiosBlobErrorInterceptorId = axios.interceptors.response.use(
    (response) => response,
    (error) => {
      if (
        error.request &&
        error.request.responseType === 'blob' &&
        error.response &&
        error.response.data instanceof Blob &&
        error.response.data.type &&
        error.response.data.type.toLowerCase().indexOf('json') !== -1
      ) {
        return new Promise((resolve, reject) => {
          const reader = new FileReader()
          reader.onload = () => {
            error.response.data = JSON.parse(reader.result)
            resolve(Promise.reject(error))
          }
          reader.onerror = () => {
            reject(error)
          }
          reader.readAsText(error.response.data)
        })
      }
      return Promise.reject(error)
    }
  )
})

onBeforeUnmount(() => {
  if (axiosBlobErrorInterceptorId != null) axios.interceptors.response.eject(axiosBlobErrorInterceptorId)
  // release thumbnails
  for (const fn of Object.keys(blobUrls)) URL.revokeObjectURL(blobUrls[fn])
})
</script>
