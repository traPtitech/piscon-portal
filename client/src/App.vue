<template>
  <div id="app" class="app">
    <router-view/>
  </div>
</template>

<script>
  export default {
    name: 'app',
    created () {
      this.$store.dispatch('getData')
      setInterval(() => {
        this.$store.dispatch('getData')
      }, 20000)
    }
  }
</script>

<style lang="scss">
  @import "sass/main";

  body {
    height: 100%;
    #app {
      height: 100%;
    }
  }
</style>
