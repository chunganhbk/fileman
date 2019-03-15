<template>
  <div>
    <p v-if="!isDefault">
      <label for="username">{{ $t('mangerUser.username') }}</label>
      <input class="input input--block" type="text" v-model="user.username" id="username">
    </p>

    <p v-if="!isDefault">
      <label for="password">{{ $t('mangerUser.password') }}</label>
      <input class="input input--block" type="password" :placeholder="passwordPlaceholder" v-model="user.password" id="password">
    </p>

    <p>
      <label for="scope">{{ $t('mangerUser.scope') }}</label>
      <input class="input input--block" type="text" v-model="user.scope" id="scope">
    </p>

   <!-- <p>
      <label for="locale">{{ $t('settings.language') }}</label>
      <languages class="input input&#45;&#45;block" id="locale" :locale.sync="user.locale"></languages>
    </p>-->

    <p v-if="!isDefault">
      <input type="checkbox" :disabled="user.perm.admin" v-model="user.lockPassword"> {{ $t('mangerUser.lockPassword') }}
    </p>

    <permissions :perm.sync="user.perm" />

    <!--<div v-if="!isDefault">
      <h3>{{ $t('settings.rules') }}</h3>
      <p class="small">{{ $t('settings.rulesHelp') }}</p>
      <rules :rules.sync="user.rules" />
    </div>-->
  </div>
</template>

<script>
import Rules from './Rules'
import Permissions from './Permissions'


export default {
  name: 'user',
  components: {
    Permissions,
    Rules,
  },
  props: [ 'user', 'isNew', 'isDefault' ],
  computed: {
    passwordPlaceholder () {
      return this.isNew ? '' : this.$t('mangerUser.avoidChanges')
    }
  },
  watch: {
    'user.perm.admin': function () {
      if (!this.user.perm.admin) return
      this.user.lockPassword = false
    }
  }
}
</script>
