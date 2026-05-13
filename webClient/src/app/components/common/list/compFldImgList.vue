<!--
  multi-image doc field: draggable thumbnails, q-uploader (inline or dialog), optional URL add,
  client sort-by-file (persisted via same PG update as reorder) / clear / bulk-delete when ext allows,
  preview dialog, rotate via /api/rotate_image.
  script setup + direct config/utils/i18n (same pattern as comp-fld-files).
  server: upload_image optional preserveOriginalFileName; see docs/90-internals/plan-image-list-upload-unify.md
-->
<template>
  <div>
    <p class="text-caption">{{ label }}</p>
    <div class="q-gutter-sm row items-start q-mb-md">
      <draggable
        :model-value="list"
        item-key="file"
        ghostClass="ghost"
        :disabled="selectMode || readonly"
        @start="onDragStart"
        @update:model-value="onListReorder"
      >
        <template #item="{element}">
          <comp-stat-img-src
            :src="actualImg(element.file)"
            @error="v => imgSrcError(element.file, v)"
            style="border-radius: 10px; cursor: pointer; width: 150px; margin: 0 8px 8px 0"
            :style="selectMode ? `border: 2px solid ${isImageSelected(element.file) ? 'rgba(25, 118, 210, .6)' : 'transparent'}` : ''"
            @click="selectMode ? toggleImageSelection(element.file) : openInDialog(element.file)"
            class="image-hover-container"
          >
            <!-- selection checkbox when in select mode -->
            <q-checkbox
              dense
              v-if="selectMode"
              v-model="selectedImages"
              :val="element.file"
              class="absolute-top-right select-checkbox"
              @click.stop
              style="margin: 4px"
            />

            <!-- regular buttons when not in select mode -->
            <template v-if="!selectMode">
              <a :href="imgUrl(element.file)" target="_blank" @click.stop>
                <q-btn round push dense size="sm" icon="open_in_new" color="grey-4" textColor="black"
                       class="absolute-top-right all-pointer-events hover-btn">
                  <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]" style="border-radius: 10px" align="center">
                    открыть в<br>новом окне
                  </q-tooltip>
                </q-btn>
              </a>
              <q-btn v-if="!readonly" outline round dense size="sm" icon="delete" color="red-14"
                     @click.stop="showDeleteDialog(element.file)" class="absolute-bottom-right all-pointer-events hover-btn">
                <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]" style="border-radius: 10px">
                  {{$t('message.delete')}} фото
                </q-tooltip>
              </q-btn>
              <q-btn v-if="!readonly" round push dense size="sm" icon="redo" color="grey-4" textColor="black"
                     @click.stop="rotateImage(element.file)" class="absolute-top-left all-pointer-events hover-btn">
                <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]" style="border-radius: 10px">повернуть</q-tooltip>
              </q-btn>
            </template>
          </comp-stat-img-src>
        </template>
      </draggable>

      <div class="row q-gutter-sm items-center" v-if="!readonly">
        <!-- upload: legacy dialog or inline hidden uploader -->
        <div v-if="useDialogUploader">
          <q-btn flat round icon="add" size="sm" @click="isShowUploadDialog = true">
            <q-tooltip>{{ $t('message.add_photo') }}</q-tooltip>
          </q-btn>
        </div>
        <div v-else>
          <q-uploader multiple
                      autoUpload
                      color="secondary"
                      style="height: 0; width: 0; margin-right: 30px"
                      ref="uploader"
                      :url="uploadUrl"
                      :headers="headers"
                      :accept="(ext && ext.accept) ? ext.accept : ''"
                      :maxFileSize="(ext && ext.maxFileSize) ? ext.maxFileSize : 10000000"
                      @rejected="rejected"
                      @uploaded="uploaded"
                      @failed="failed"
                      @finish="finish"
                      :formFields="formField"
          >
            <template v-slot:header="scope">
              <q-btn round dense flat type="a" icon="add" color="secondary" @click="scope.pickFiles" style="margin: -22px 0 0 -5px">
                <q-uploader-add-trigger/>
                <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]" style="border-radius: 10px">
                  {{ $t('message.add_photo') }}
                </q-tooltip>
              </q-btn>
            </template>
          </q-uploader>
        </div>

        <div v-if="list.length && canSortByFile">
          <q-btn flat round icon="sort" size="sm" :disable="selectMode" @click="isShowSortDialog = true">
            <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]" style="border-radius: 10px">сортировка</q-tooltip>
          </q-btn>
        </div>

        <div v-if="list.length && canClearAll">
          <q-btn flat round icon="delete" size="sm" :disable="selectMode" @click="isShowClearDialog = true">
            <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]" style="border-radius: 10px">очистить</q-tooltip>
          </q-btn>
        </div>

        <div v-if="list.length">
          <q-btn flat round :icon="selectMode ? 'check_box' : 'check_box_outline_blank'" size="sm"
                 @click="toggleSelectMode" :color="selectMode ? 'warning' : 'grey'">
            <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]" style="border-radius: 10px">
              {{ selectMode ? 'убрать выбор' : 'выбрать' }}
            </q-tooltip>
          </q-btn>
        </div>

        <div v-if="selectMode && selectedImages.length">
          <q-btn round color="negative" size="sm" icon="delete" @click="isShowBulkDeleteDialog = true">
            <q-badge floating rounded color="red" :label="selectedImages.length" style="margin: -3px -3px 0 0"/>
            <q-tooltip anchor="top middle" self="bottom middle" :offset="[10, 10]" style="border-radius: 10px">
              удалить выбранные фото ({{ selectedImages.length }})
            </q-tooltip>
          </q-btn>
        </div>

        <div v-if="ext && ext.canAddUrls">
          <q-btn size="sm" flat round icon="add" @click="isShowDialogAddUrl = true">
            <q-tooltip>{{ $t('message.add_link') }}</q-tooltip>
          </q-btn>
        </div>
      </div>
    </div>

    <!-- legacy upload dialog -->
    <q-dialog v-model="isShowUploadDialog">
      <q-uploader
        ref="uploaderDialog"
        :label="$t('message.select_file_for_upload')"
        multiple
        :url="uploadUrl"
        :headers="headers"
        :accept="(ext && ext.accept) ? ext.accept : ''"
        :max-file-size="(ext && ext.maxFileSize) ? ext.maxFileSize : 10000000"
        @rejected="rejected"
        @uploaded="uploaded"
        @failed="failed"
        @finish="finish"
        :form-fields="formField"
      />
    </q-dialog>

    <!-- диалог добавления url ссылку на фото  -->
    <q-dialog v-model="isShowDialogAddUrl">
      <q-card style="width: 300px" class="q-px-sm q-pb-md">
        <q-img :src="newImgUrl ? newImgUrl : 'https://www.cowgirlcontractcleaning.com/wp-content/uploads/sites/360/2018/05/placeholder-img-4.jpg'" />
        <q-card-section>
          <q-input v-model="newImgUrl" label="ссылка на фото"/>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('message.cancel')" v-close-popup/>
          <q-btn flat label="Ок" v-close-popup @click="addImgUrl"/>
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- диалог подтверждения удаления -->
    <q-dialog v-model="isShowDeleteDialog">
      <q-card style="border-radius: 20px" align="center">
        <q-card-section>
          <div style="font-size: 1.4rem; margin-top: -10px" class="q-mx-lg q-mb-md">точно?</div>
          <q-btn rounded push color="negative" :label="$t('message.delete')" v-close-popup @click="remove"/>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- диалог подтверждения сортировки -->
    <q-dialog v-model="isShowSortDialog">
      <q-card style="border-radius: 20px" align="center">
        <q-card-section>
          <div style="font-size: 1.4rem; margin-top: -10px" class="q-mx-lg q-mb-md">сортировать фото по имени?</div>
          <q-btn rounded push color="primary" label="да" v-close-popup @click="confirmSort" class="q-px-lg"/>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- диалог подтверждения очистки -->
    <q-dialog v-model="isShowClearDialog">
      <q-card style="border-radius: 20px" align="center">
        <q-card-section>
          <div style="font-size: 1.4rem; margin-top: -10px" class="q-mx-lg q-mb-md">очистить все фото?</div>
          <q-btn rounded push color="negative" label="да" v-close-popup @click="confirmClear" class="q-px-lg"/>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- диалог подтверждения массового удаления -->
    <q-dialog v-model="isShowBulkDeleteDialog">
      <q-card style="border-radius: 20px" align="center">
        <q-card-section>
          <div style="font-size: 1.4rem; margin-top: -10px" class="q-mx-lg q-mb-md">удалить {{ selectedImages.length }} фото?</div>
          <q-btn rounded push color="negative" :label="$t('message.delete')" v-close-popup @click="confirmBulkDelete" class="q-px-lg"/>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- диалог изображения -->
    <q-dialog v-model="isShowImageDialog">
      <q-card style="border-radius: 20px; width: 600px; max-width: 100%" v-if="dialogImageUrl">
        <comp-stat-img-src :src="dialogImageUrl" @error="v => imgSrcError(dialogImageUrl, v)" v-close-popup/>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
// quasar field component: contract summarized in docs/30-templates/vue-shared-controls.md
import axios from 'axios'
import {i18n} from 'boot/i18n'
import {useQuasar} from 'quasar'
import draggable from 'vuedraggable'
import {
  computed,
  nextTick,
  onMounted,
  reactive,
  ref,
  watch
} from 'vue'
import config from '../../../plugins/config'
import utils from '../../../plugins/utils'

// same default as legacy template fallbacks
const DEFAULT_MAX_FILE_SIZE = 10000000

const props = defineProps({
  fld: {default: undefined},
  label: {default: undefined},
  readonly: {default: null},
  icon: {default: null},
  vif: {type: Boolean, default: true},
  ext: {type: Object, default: undefined},
  formFieldParams: {type: Array, default: () => []},
  // флаг, чтобы отключить обновления записи в БД после загрузки фото
  isUpdateFldsInPostgres: {type: Boolean, default: true}
})

const emit = defineEmits(['update'])

const $q = useQuasar()
const t = (...args) => i18n.global.t(...args)

const list = ref([])
const selectedForDeleteFilename = ref(null)
const newImgUrl = ref(null)
const imgKeys = reactive({})
const isShowUploadDialog = ref(false)
const isShowDeleteDialog = ref(false)
const isShowSortDialog = ref(false)
const isShowClearDialog = ref(false)
const isShowBulkDeleteDialog = ref(false)
const isShowImageDialog = ref(false)
const dialogImageUrl = ref(null)
const isShowDialogAddUrl = ref(false)
const selectMode = ref(false)
const selectedImages = ref([])
const dragReorderInProgress = ref(false)
// snapshot before drag; vuedraggable can briefly clear v-model so we never POST [] by mistake
const listSnapshotBeforeDrag = ref(null)

const uploader = ref(null)
const uploaderDialog = ref(null)

const uploadUrl = computed(() =>
  `${config.apiUrl()}/api/${(props.ext && props.ext.uploadUrl) || 'upload_image'}`
)

const headers = computed(() => {
  const authToken = localStorage.getItem(config.appName)
  return [{name: 'Auth-token', value: authToken}]
})

const useDialogUploader = computed(() => !!(props.ext && props.ext.useDialogUploader))

const canClearAll = computed(() => !props.readonly && !(props.ext && props.ext.disableClearAll))

// sort-by-file only when uploads are not forced to random names (same signal as preserveOriginalFileName)
const canSortByFile = computed(() => {
  const e = props.ext
  if (!e || e.disableSortByFile === true) return false
  if (e.randomizeUploadFileNames === true) return false
  return true
})

const formField = computed(() => {
  const res = [
    {name: 'tableName', value: props.ext && props.ext.tableName},
    {name: 'tableId', value: props.ext && props.ext.tableId},
    ...props.formFieldParams
  ]
  if (props.ext && props.ext.width) res.push({name: 'width', value: props.ext.width})
  if (props.ext && props.ext.crop) res.push({name: 'crop', value: props.ext.crop})
  // default: preserve original basename; opt-out keeps server default (ULID)
  if (!(props.ext && props.ext.randomizeUploadFileNames === true)) {
    res.push({name: 'preserveOriginalFileName', value: '1'})
  }
  return res
})

const imgUrl = computed(() => (src) =>
  (src && src.includes('http') ? src : `${config.apiUrl()}${src}`)
)

const actualImg = computed(() => {
  const timestamp = Date.now()
  return (link) => {
    const fileKey = imgKeys[link] || 0
    return `${link}?key=${fileKey + timestamp}`
  }
})

function pgUpdateMethod() {
  const e = props.ext
  return `${e.methodUpdate || e.tableName + '_update'}`
}

function emitListUpdate() {
  // emit plain objects to avoid circular reference issues with Vue reactivity
  emit('update', list.value.map(v => ({...v})))
}

function ensureExtValidOrThrow() {
  if (props.readonly || !props.isUpdateFldsInPostgres) return
  if (!props.ext) {
    throw new Error('compFldFiles missed param: "ext"')
  }
  if (!props.ext.fldName) {
    throw new Error('compFldFiles missed param: "ext.fldName"')
  }
  if (!props.ext.methodUpdate && !(props.ext.tableId && props.ext.tableName)) {
    throw new Error('compFldFiles missed param: "ext.methodUpdate" OR "ext.tableId" AND "ext.tableName"')
  }
}

function resyncListFromFld() {
  list.value = Array.isArray(props.fld)
    ? props.fld.filter(x => x && !x.deleted).map(x => ({...x}))
    : []
}

function listSnapshotForSave() {
  try {
    return JSON.parse(JSON.stringify(list.value))
  } catch (e) {
    return list.value.map(v => ({...v}))
  }
}

function onDragStart() {
  listSnapshotBeforeDrag.value = listSnapshotForSave()
}

function onListReorder(newList) {
  // apply order from draggable explicitly (avoids v-model/end timing where list was still [])
  const next = Array.isArray(newList)
    ? newList.filter(x => x && !x.deleted).map(x => ({...x}))
    : []
  list.value = next
  dragReorderInProgress.value = true
  nextTick(() => {
    nextTick(() => {
      saveListOrder()
    })
  })
}

function saveListOrder() {
  if (!props.isUpdateFldsInPostgres) {
    dragReorderInProgress.value = false
    emitListUpdate()
    return
  }
  let payload = listSnapshotForSave()
  const fldActive = Array.isArray(props.fld) ? props.fld.filter(x => x && !x.deleted) : []
  const snap = listSnapshotBeforeDrag.value
  // never persist an empty jsonb from a reorder when we still had images before drag
  if (payload.length === 0 && snap && snap.length > 0) {
    list.value = snap.map(x => ({...x}))
    payload = listSnapshotForSave()
  }
  if (payload.length === 0 && fldActive.length > 0) {
    resyncListFromFld()
    listSnapshotBeforeDrag.value = null
    dragReorderInProgress.value = false
    return
  }
  // reorder path must not wipe the field
  if (payload.length === 0) {
    listSnapshotBeforeDrag.value = null
    dragReorderInProgress.value = false
    return
  }
  utils.postCallPgMethod({
    method: pgUpdateMethod(),
    params: {id: props.ext.tableId, [props.ext.fldName]: payload}
  }).subscribe(
    () => {
      listSnapshotBeforeDrag.value = null
      dragReorderInProgress.value = false
      emitListUpdate()
    },
    () => {
      listSnapshotBeforeDrag.value = null
      dragReorderInProgress.value = false
    }
  )
}

function uploaded({xhr: {response}}) {
  const res = JSON.parse(response)
  if (!res.ok) {
    $q.notify({color: 'negative', position: 'bottom', message: res.message})
  } else {
    list.value.push(res.result)
  }
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

  let msgText = 'данный файл не соответствует ограничениям'
  if (msg.length && msg[0].failedPropValidation === 'accept') {
    msgText = `${t('message.only_files_with_extension')}: ${props.ext.accept} `
  }
  if (msg.length && msg[0].failedPropValidation === 'max-file-size') {
    const size = niceBytes((props.ext && props.ext.maxFileSize) || DEFAULT_MAX_FILE_SIZE)
    msgText = `${t('message.only_files_no_larger_than')}: ${size}`
  }
  $q.notify({color: 'negative', position: 'bottom', message: msgText})
}

function showDeleteDialog(filename) {
  selectedForDeleteFilename.value = filename
  isShowDeleteDialog.value = true
}

function imgSrcError(filename) {
  $q.notify({color: 'negative', position: 'bottom', message: `ошибка: url "${filename}" не является ссылкой на фото `})
  selectedForDeleteFilename.value = filename
  remove()
}

function remove() {
  const i = list.value.findIndex(v => v.file === selectedForDeleteFilename.value && !v.deleted)
  if (i === -1) return
  list.value[i].deleted = true

  if (!props.isUpdateFldsInPostgres) {
    list.value = list.value.filter(v => !v.deleted)
    emitListUpdate()
    return
  }

  utils.postCallPgMethod({
    method: pgUpdateMethod(),
    params: {id: props.ext.tableId, [props.ext.fldName]: list.value.map(v => ({...v}))}
  }).subscribe(() => {
    list.value = list.value.filter(v => !v.deleted)
    emitListUpdate()
  })
}

function toggleSelectMode() {
  selectMode.value = !selectMode.value
  if (!selectMode.value) selectedImages.value = []
}

function toggleImageSelection(filename) {
  const index = selectedImages.value.indexOf(filename)
  if (index > -1) selectedImages.value.splice(index, 1)
  else selectedImages.value.push(filename)
}

function isImageSelected(filename) {
  return selectedImages.value.includes(filename)
}

function confirmBulkDelete() {
  selectedImages.value.forEach(filename => {
    const i = list.value.findIndex(v => v.file === filename && !v.deleted)
    if (i !== -1) list.value[i].deleted = true
  })

  if (!props.isUpdateFldsInPostgres) {
    list.value = list.value.filter(v => !v.deleted)
    selectedImages.value = []
    emitListUpdate()
    return
  }

  utils.postCallPgMethod({
    method: pgUpdateMethod(),
    params: {id: props.ext.tableId, [props.ext.fldName]: list.value.map(v => ({...v}))}
  }).subscribe(() => {
    list.value = list.value.filter(v => !v.deleted)
    selectedImages.value = []
    emitListUpdate()
  })
}

function confirmSort() {
  if (!canSortByFile.value) return
  const next = [...list.value]
    .filter((x) => x && !x.deleted)
    .sort((a, b) =>
      String(a.file || '').localeCompare(String(b.file || ''), undefined, {
        numeric: true,
        sensitivity: 'base'
      })
    )
    .map((x) => ({...x}))
  list.value = next

  if (!props.isUpdateFldsInPostgres) {
    emitListUpdate()
    return
  }

  utils.postCallPgMethod({
    method: pgUpdateMethod(),
    params: {id: props.ext.tableId, [props.ext.fldName]: list.value.map((v) => ({...v}))}
  }).subscribe(
    () => emitListUpdate(),
    () => {}
  )
}

function confirmClear() {
  const fieldName = (props.ext && props.ext.clearPhotoField) || (props.ext && props.ext.fldName)
  if (!fieldName) return

  if (!props.isUpdateFldsInPostgres) {
    list.value = []
    emitListUpdate()
    return
  }

  utils.postCallPgMethod({
    method: pgUpdateMethod(),
    params: {id: props.ext.tableId, [fieldName]: null}
  }).subscribe(() => {
    list.value = []
    emitListUpdate()
  })
}

function finish() {
  // reset whichever uploader is used
  if (uploader.value) uploader.value.reset()
  if (uploaderDialog.value) uploaderDialog.value.reset()
  emitListUpdate()

  if (!props.isUpdateFldsInPostgres) return
  utils.postCallPgMethod({
    method: pgUpdateMethod(),
    params: {id: props.ext.tableId, [props.ext.fldName]: list.value.map(v => ({...v}))}
  }).subscribe(() => {})
}

function rotateImage(filename) {
  const authToken = localStorage.getItem(config.appName)
  const hdrs = {'Content-Type': 'application/json', 'Auth-token': authToken}
  const urlStr = `${config.apiUrl()}/api/rotate_image`

  axios.post(urlStr, {params: {filename}}, {headers: hdrs})
    .then((response) => {
      if (response.data && response.data.ok) {
        imgKeys[filename] = (imgKeys[filename] || 0) + 1
      }
    })
    .catch((error) => {
      const msg = (error && error.response && error.response.data) ? error.response.data : t('message.upload_error')
      $q.notify({message: msg, color: 'negative', position: 'top-right'})
    })
}

function openInDialog(url) {
  dialogImageUrl.value = imgUrl.value(url)
  isShowImageDialog.value = true
}

function addImgUrl() {
  list.value.push({file: newImgUrl.value})

  if (!props.isUpdateFldsInPostgres) {
    newImgUrl.value = null
    emitListUpdate()
    return
  }

  utils.postCallPgMethod({
    method: pgUpdateMethod(),
    params: {id: props.ext.tableId, [props.ext.fldName]: list.value.map(v => ({...v}))}
  }).subscribe((res) => {
    if (res && res.ok) newImgUrl.value = null
  })
}

watch(
  () => props.fld,
  (v) => {
    if (dragReorderInProgress.value) return
    list.value = Array.isArray(v) ? v.filter(x => x && !x.deleted).map(x => ({...x})) : []
  },
  {immediate: true}
)

onMounted(() => {
  ensureExtValidOrThrow()
})
</script>

<style scoped>
.ghost {
  opacity: 0.1;
  border: 1px solid black;
}

/* image container styling */
.image-hover-container {
  transition: all 0.1s ease-in-out;
  border-radius: 4px;
}

/* hide buttons by default */
.image-hover-container .hover-btn {
  opacity: 0;
  transition: opacity 0.1s ease-in-out;
}

/* show buttons and add shadow/outline on hover */
.image-hover-container:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, .2);
}

.image-hover-container:hover .hover-btn {
  opacity: 1;
}

/* Selection mode styling */
.select-checkbox {
  padding: 2px;
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}
</style>
