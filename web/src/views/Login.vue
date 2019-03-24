<template>
  <div id="login" >
    <form @submit="submit">
      <!--<img :src="" alt="File man">-->
      <h1>{{ name }}</h1>
      <div v-if="error !== ''" class="wrong">{{ error }}</div>

      <input class="input input--block" type="text" v-model="username" :placeholder="$t('login.username')">
      <input class="input input--block" type="password" v-model="password" :placeholder="$t('login.password')">
      <input class="button button--block" type="submit" :value=" $t('login.submit')">
    </form>
  </div>
</template>

<script>
import * as auth from '@/utils/auth'
import { name,  } from '@/utils/constants'

export default {
  name: 'login',
  computed: {
    name: () => name,
    logoURL: () => logoURL
  },
  data: function () {
    return {
      error: '',
      username: '',
      password: '',
    }
  },
  methods: {
    toggleMode () {
      this.createMode = !this.createMode
    },
    async submit (event) {
      event.preventDefault();
      event.stopPropagation();

      let redirect = this.$route.query.redirect
      if (redirect === '' || redirect === undefined || redirect === null) {
        redirect = '/files/'
      }
      try {
        await auth.login(this.username, this.password)
        this.$router.push({ path: redirect })
      } catch (e) {
          console.log(e)
        if (e.message == 409) {
          this.error = this.$t('login.usernameTaken')
        } else {
          this.error = this.$t('login.wrongCredentials')
        }
      }
    }
  }
}
</script>
