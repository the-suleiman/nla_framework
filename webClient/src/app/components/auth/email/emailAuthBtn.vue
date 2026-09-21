<template>
  <div>
    <q-btn size='md' outline color='primary' @click='login' :disabled='disabled'>
      <q-icon left size="3em" name="perm_identity"/>
      <div>{{$t('auth.login')}}</div>
    </q-btn>
    <q-dialog v-model="modalIsOpened" persistent>
      <q-card style="min-width: 25%">
        <q-toolbar>
          <q-toolbar-title style="display: flex"><span class="text-weight-bold" style="text-transform: capitalize; margin-left: auto; padding-left: 45px; margin-right: auto;">{{dialogTitle}}</span></q-toolbar-title>
          <q-btn flat round dense icon="close" v-close-popup/>
        </q-toolbar>

        <q-card-section style="padding-top: 10px">
          <!-- форма логин + пароль -->
          <comp-login-form v-if='!emailRegister.isPasswordRecover'
                           @passwordRecover="showPasswordRecover"/>
          <!--форма восстановления пароля -->
          <comp-recover-password-form v-if='emailRegister.isPasswordRecover' @cancel="resetPasswordRecover"/>
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>

<style lang="sass" scoped>
  .bg-facebook
    background: #4267b2
</style>

<script>
  import compLoginForm from './components/compLoginForm'
  import compRecoverPasswordForm from './components/compRecoverPasswordForm'

  export default {
    props: ['disabled'],
    components: {compLoginForm, compRecoverPasswordForm},
    computed: {
      dialogTitle() {
        if (this.emailRegister.isPasswordRecover) return this.$t('auth.password_recovery')
        return this.$t('auth.authorization')
      },
    },
    data() {
      return {
        modalIsOpened: false,
        emailRegister: {
          isPasswordRecover: false,
          isPasswordRecoverSuccess: false,
        },
      }
    },
    methods: {
      login() {
        this.modalIsOpened = true
      },
      showPasswordRecover() {
        this.emailRegister.isPasswordRecover = true
      },
      resetPasswordRecover() {
        this.emailRegister.isPasswordRecover = false
      },
    },
  }
</script>
