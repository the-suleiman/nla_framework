<!--
  wraps q-img for server paths and local assets: resolves `src` to a full URL
  when needed (see docs/30-templates/vue-shared-controls.md).
-->
<template>
  <q-img :src="imgUrl" @error="v => emit('error', v)" style="border-radius: 10px">
    <slot/>
  </q-img>
</template>

<script setup>
import {computed} from 'vue'
import config from '../../../plugins/config'

const props = defineProps({
    src: {default: null},
})

// re-emits q-img load failures for parent handling
const emit = defineEmits(['error'])

// same resolution rules as documented: blob/data/http(s)/image/* as-is; else api prefix
const imgUrl = computed(() => {
  if (!props.src) return 'https://www.foot.com/wp-content/uploads/2017/03/placeholder.gif'
  const s = props.src
  if (s.startsWith('blob:') || s.startsWith('data:')) return s
  if (/^https?:\/\//i.test(s)) return s
  if (s.startsWith('image')) return s
  return `${config.apiUrl()}${s}`
})
</script>
