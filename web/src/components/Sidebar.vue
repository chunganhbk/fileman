<template>
  <nav :class="{active}">
    <template v-if="isLogged">
      <router-link class="action" to="/files/" :aria-label="$t('sidebar.myFiles')" :title="$t('sidebar.myFiles')">
        <i class="material-icons">folder</i>
        <span>{{ $t('sidebar.myFiles') }}</span>
      </router-link>

      <div v-if="user.perm.create">
        <button @click="$store.commit('showHover', 'newDir')" class="action" :aria-label="$t('sidebar.newFolder')" :title="$t('sidebar.newFolder')">
          <i class="material-icons">create_new_folder</i>
          <span>{{ $t('sidebar.newFolder') }}</span>
        </button>

        <button @click="$store.commit('showHover', 'newFile')" class="action" :aria-label="$t('sidebar.newFile')" :title="$t('sidebar.newFile')">
          <i class="material-icons">note_add</i>
          <span>{{ $t('sidebar.newFile') }}</span>
        </button>
      </div>
      <div  v-if="user.perm.admin">
        <router-link class="action" to="/users" :aria-label="$t('sidebar.users')" :title="$t('sidebar.users')">
          <i class="material-icons">people</i>
          <span>{{ $t('sidebar.users') }}</span>
        </router-link>
      </div>

    </template>
    <template v-else>
      <router-link class="action" to="/login" :aria-label="$t('sidebar.login')" :title="$t('sidebar.login')">
        <i class="material-icons">exit_to_app</i>
        <span>{{ $t('sidebar.login') }}</span>
      </router-link>

    </template>

  </nav>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import * as auth from '@/utils/auth'

export default {
  name: 'sidebar',
  computed: {
    ...mapState([ 'user' ]),
    ...mapGetters([ 'isLogged' ]),
    active () {
      return this.$store.state.show === 'sidebar'
    },
  },
  methods: {
    help () {
      this.$store.commit('showHover', 'help')
    },
    logout: auth.logout
  }
}
</script>
