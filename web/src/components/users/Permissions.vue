<template>
  <div>
    <h3>{{ $t('mangerUser.permissions') }}</h3>
    <p class="small">{{ $t('mangerUser.permissionsHelp') }}</p>

    <p><input type="checkbox" v-model="admin"> {{ $t('mangerUser.administrator') }}</p>

    <p><input type="checkbox" :disabled="admin" v-model="perm.create"> {{ $t('mangerUser.perm.create') }}</p>
    <p><input type="checkbox" :disabled="admin" v-model="perm.delete"> {{ $t('mangerUser.perm.delete') }}</p>
    <p><input type="checkbox" :disabled="admin" v-model="perm.download"> {{ $t('mangerUser.perm.download') }}</p>
    <p><input type="checkbox" :disabled="admin" v-model="perm.modify"> {{ $t('mangerUser.perm.modify') }}</p>
    <p><input type="checkbox" :disabled="admin" v-model="perm.rename"> {{ $t('mangerUser.perm.rename') }}</p>
  </div>
</template>

<script>
export default {
  name: 'permissions',
  props: ['perm'],
  computed: {
    admin: {
      get () {
        return this.perm.admin
      },
      set (value) {
        if (value) {
          for (const key in this.perm) {
            this.perm[key] = true
          }
        }

        this.perm.admin = value
      }
    }
  }
}
</script>
