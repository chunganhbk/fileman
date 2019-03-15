<template>
  <div>
    <router-view></router-view>
  </div>
</template>

<script>

export default {
  name: 'app',
  mounted () {
      const loading = document.getElementById('loading')
      loading.classList.add('done')

      setTimeout(function () {
          loading.parentNode.removeChild(loading)
      }, 200)
  }
}
</script>

<style>
  @import './assets/css/styles.css';
</style>

