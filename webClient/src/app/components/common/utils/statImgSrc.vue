<!-- компонента для отображения ссылок на фото, загруженных на сервер. В зависимости от режима разработки подставляет адрес сервера-->
<template>
  <q-img :src="imgUrl" @error="v => $emit('error', v)" style="border-radius: 10px">
    <slot></slot>
  </q-img>
</template>

<script>
export default {
    props: ['src'],
    computed: {
        imgUrl() {
            if (!this.src) {
                return 'https://www.foot.com/wp-content/uploads/2017/03/placeholder.gif'
            }
            const s = this.src
            if (s.startsWith('blob:') || s.startsWith('data:')) {
                return s
            }
            if (/^https?:\/\//i.test(s)) {
                return s
            }
            if (s.startsWith('image')) {
                return s
            }
            return `${this.$config.apiUrl()}${s}`
        },
    },
}
</script>
