<template>
  <vuestic-layout>
    <app-navbar :isOpen="opened" @toggle-menu="toggleSidebar"/>
    <app-sidebar :isOpen="opened"/>
    <main slot="content" id="content" class="content" role="main">
      <vuestic-breadcrumbs/>
      <vuestic-pre-loader v-show="isLoading" ref="preLoader" class="pre-loader"></vuestic-pre-loader>
      <router-view></router-view>
    </main>
    <span slot="footer">©2019. Made by&nbsp;<a href="https://to-hutohu.com" target="_blank">@to-hutohu</a>, <a href="https://nagatech.work" target="_blank">@nagatech</a></span>
  </vuestic-layout>
</template>

<script>

  import VuesticLayout from '../../vuestic-theme/vuestic-components/vuestic-layout/VuesticLayout'
  import AppNavbar from './app-navbar/AppNavbar'
  import AppSidebar from './app-sidebar/AppSidebar'
  import {mapGetters} from 'vuex'

  export default {
    name: 'app-layout',

    components: {
      VuesticLayout,
      AppNavbar,
      AppSidebar
    },
    data () {
      return {
        opened: true
      }
    },
    methods: {
      toggleSidebar (opened) {
        this.opened = opened
      }
    },
    computed: {
      ...mapGetters([
        'isLoading'
      ])
    }
  }
</script>

<style lang="sass" scoped>
  a
    color: #4ab1e2;
</style>
