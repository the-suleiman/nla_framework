<template>
  <div style="max-width: 450px;">
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

<script>
import axios from 'axios'

const DEFAULT_MAX_FILE_SIZE = 10000000

export default {
    props: {
        fld: {},
        fldName: {
            default: 'files',
        },
        label: {},
        readonly: null,
        icon: null,
        vif: {
            default: true,
        },
        ext: null,
        // true = add opens dialog only; omit or null/false = compact bar uploader (set explicitly when needed)
        showUploaderDialog: null,
    },
    computed: {
        uploadUrl() {
            return `${this.$config.apiUrl()}/api/upload_file`
        },
        headers() {
            const authToken = localStorage.getItem(this.$config.appName)
            return [{name: 'Auth-token', value: authToken}]
        },
        filteredList() {
            return this.list.filter(v => !v.deleted)
        },
        effectiveMaxFileSize() {
            const raw = this.ext && this.ext.maxFileSize
            if (raw == null || raw === '') {
                return DEFAULT_MAX_FILE_SIZE
            }
            const n = parseInt(raw, 10)
            return Number.isNaN(n) ? DEFAULT_MAX_FILE_SIZE : n
        },
        formFields() {
            return [
                {name: 'tableName', value: this.ext.tableName},
                {name: 'tableId', value: this.ext.tableId},
            ]
        },
        dialogImgSrc() {
            if (!this.dialogImage || !this.dialogImage.filename) {
                return null
            }
            return this.blobUrls[this.dialogImage.filename] || this.dialogImage.url
        },
        qBarColorClass() {
            if (this.readonly) {
                return 'bg-secondary'
            }
            return this.dropWaiting ? 'bg-positive' : 'bg-secondary'
        },
        qBarRadiusStyle() {
            const borderRadius = this.list.length && this.isShowList ? '10px 10px 0 0' : '10px'
            return {borderRadius}
        },
        expandListIcon() {
            return `expand_${this.isShowList ? 'less' : 'more'}`
        },
        inlineUploaderIcon() {
            return this.dropWaiting ? 'download' : 'add'
        },
    },
    data() {
        return {
            isShowDialog: false,
            isShowList: true,
            list: [],
            isShowDeleteDialog: false,
            selectedForDeleteFilename: null,
            dropWaiting: false,
            isShowImageDialog: false,
            dialogImage: {url: null},
            blobUrls: {},
        }
    },
    watch: {
        filteredList: {
            deep: true,
            immediate: true,
            handler(list) {
                const names = new Set(list.map(i => i.filename))
                for (const fn of Object.keys(this.blobUrls)) {
                    if (!names.has(fn)) {
                        URL.revokeObjectURL(this.blobUrls[fn])
                        this.$delete(this.blobUrls, fn)
                    }
                }
                for (const item of list) {
                    if (this.isImage(item.filename) && !this.blobUrls[item.filename]) {
                        this.fetchImagePreview(item)
                    }
                }
            },
        },
    },
    methods: {
        uploaded({xhr: {response}}) {
            const res = JSON.parse(response)
            if (!res.ok) {
                this.$q.notify({
                    color: 'negative',
                    position: 'bottom',
                    message: res.message,
                })
            } else {
                if (!this.list.find(v => v.filename === res.result.filename && !v.deleted)) {
                    this.list.push(res.result)
                }
            }
        },
        finish() {
            if (this.isShowDialog && this.$refs.uploaderDialog) {
                this.$refs.uploaderDialog.reset()
            } else if (this.$refs.uploaderInline) {
                this.$refs.uploaderInline.reset()
            }
            this.isShowDialog = false
            this.$utils.postCallPgMethod({
                method: `${this.ext.tableName}_update`,
                params: {id: this.ext.tableId, [this.fldName]: this.list},
            }).subscribe(() => this.$emit('update', this.list))
        },
        failed(msg) {
            let msgText = this.$t('message.upload_error')
            if (msg.xhr && msg.xhr.responseText) {
                let res = JSON.parse(msg.xhr.responseText)
                if (res.message) msgText = res.message
            }
            this.$q.notify({
                color: 'negative',
                position: 'bottom',
                message: msgText,
            })
        },
        rejected(msg) {
            const niceBytes = (x) => {
                const units = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
                let l = 0; let n = parseInt(x, 10) || 0
                while (n >= 1024 && ++l) {
                    n = n / 1024
                }
                return (n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l])
            }
            let msgText = this.$t('message.file_rejected_constraints')
            if (msg.length > 0 && msg[0].failedPropValidation === 'accept') {
                msgText = `${this.$t('message.only_files_with_extension')}: ${this.ext.accept} `
            }
            if (msg.length > 0 && msg[0].failedPropValidation === 'max-file-size') {
                let size = niceBytes(this.effectiveMaxFileSize)
                msgText = `${this.$t('message.only_files_no_larger_than')}: ${size}`
            }
            this.$q.notify({
                color: 'negative',
                position: 'bottom',
                message: msgText,
            })
        },
        fetchImagePreview(item) {
            axios({
                url: `${this.$config.apiUrl()}${item.url}`,
                method: 'GET',
                headers: {'Auth-token': localStorage.getItem(this.$config.appName)},
                responseType: 'blob',
            }).then((response) => {
                const blobUrl = window.URL.createObjectURL(new Blob([response.data]))
                this.$set(this.blobUrls, item.filename, blobUrl)
            }).catch((err) => {
                let msg = err.response && err.response.data
                if (err.response && err.response.data) {
                    msg = err.response.data.message
                    if (msg === 'not found') msg = this.$t('message.file_not_found_by_link')
                }
                this.$q.notify({message: msg, type: 'negative', position: 'top-right'})
            })
        },
        downloadFile(item) {
            axios({
                url: `${this.$config.apiUrl()}${item.url}`,
                method: 'GET',
                headers: {'Auth-token': localStorage.getItem(this.$config.appName)},
                responseType: 'blob',
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
                    if (msg === 'not found') msg = this.$t('message.file_not_found_by_link')
                }
                this.$q.notify({message: msg, type: 'negative', position: 'top-right'})
            })
        },
        showDeleteDialog(filename) {
            this.selectedForDeleteFilename = filename
            this.isShowDeleteDialog = true
        },
        remove() {
            let i = this.list.findIndex(v => v.filename === this.selectedForDeleteFilename && !v.deleted)
            this.list[i].deleted = true
            const fileToken = this.list[i].url.split('/', -1).slice(-1)
            let item = this.list.splice(i, 1)
            this.list.push(item[0])
            this.$utils.postApiRequest({url: `/api/remove_file/${fileToken}`}).subscribe((res) => {
                if (res.ok) {
                    this.selectedForDeleteFilename = null
                }
            })
            this.$utils.postCallPgMethod({
                method: `${this.ext.tableName}_update`,
                params: {id: this.ext.tableId, [this.fldName]: this.list},
            }).subscribe(() => {
                this.$emit('update', this.list)
            })
        },
        onDragOver(event) {
            event.dataTransfer.dropEffect = 'copy'
        },
        onDragLeave(event) {
            if (!event.relatedTarget || !event.currentTarget.contains(event.relatedTarget)) {
                this.dropWaiting = false
            }
        },
        onFileDrop(event) {
            this.dropWaiting = false
            if (this.readonly) return
            const files = event.dataTransfer.files
            if (!files.length) return
            if (this.showUploaderDialog === true) {
                this.isShowDialog = true
                this.$nextTick(() => {
                    const d = this.$refs.uploaderDialog
                    if (d) d.addFiles(Array.from(files))
                })
            } else {
                const u = this.$refs.uploaderInline
                if (u) u.addFiles(Array.from(files))
            }
        },
        openInDialog(item) {
            this.dialogImage = item
            this.isShowImageDialog = true
        },
        onImageDialogHide() {
            this.dialogImage = {url: null}
        },
        isImage(fileName) {
            const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.svg', '.webp']
            return imageExtensions.some(ext => fileName.toLowerCase().endsWith(ext))
        },
    },
    mounted() {
        if (!this.ext) {
            throw new Error('compFldFiles missed param: "ext"')
        }
        if (!this.ext.tableName) {
            throw new Error('compFldFiles missed param: "ext.tableName"')
        }
        if (!this.ext.tableId) {
            throw new Error('compFldFiles missed param: "ext.tableId"')
        }
        this.list = this.fld || []
        // global axios interceptor: stacks if multiple comp-fld-files mount; prefer app-level registration long-term
        axios.interceptors.response.use(
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
                        let reader = new FileReader()
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
    },
    beforeDestroy() {
        for (const fn of Object.keys(this.blobUrls)) {
            URL.revokeObjectURL(this.blobUrls[fn])
        }
    },
}
</script>
