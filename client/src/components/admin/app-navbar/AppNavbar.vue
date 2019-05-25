<template>

  <vuestic-navbar>
    <header-selector slot="selector" :isOpen.sync="valueProxy"/>
    <span slot="logo">Piscon</span>
    <profile-dropdown v-if="$store.state.Me"
      :options="profileOptions">
      <img :src="`https://q.trap.jp/api/1.0/files/${$store.state.Me.iconFileId}`" />
    </profile-dropdown>
    <div v-else class="signin">
      <router-link :to="{name: 'signin'}">Signin with traQ</router-link>
    </div>
  </vuestic-navbar>

</template>

<script>
  import VuesticNavbar from '../../../vuestic-theme/vuestic-components/vuestic-navbar/VuesticNavbar'
  import HeaderSelector from './components/HeaderSelector'
  import LanguageDropdown from './components/dropdowns/LanguageDropdown'
  import ProfileDropdown from './components/dropdowns/ProfileDropdown'
  import NotificationDropdown from './components/dropdowns/NotificationDropdown'
  import MessageDropdown from './components/dropdowns/MessageDropdown'

  export default {
    name: 'app-navbar',

    components: {
      VuesticNavbar,
      HeaderSelector,
      MessageDropdown,
      NotificationDropdown,
      LanguageDropdown,
      ProfileDropdown
    },

    props: {
      isOpen: {
        type: Boolean,
        required: true
      }
    },
    data () {
      return {
        profileOptions: [
          {
            name: 'name',
            redirectTo: 'team-info'
          },
          {
            name: 'Logout',
            redirectTo: 'logout'
          }
        ]
      }
    },
    computed: {
      valueProxy: {
        get () {
          return this.isOpen
        },
        set (opened) {
          this.$emit('toggle-menu', opened)
        },
      }
    },
    mounted () {
      if (this.$store.state.Me) {
        this.profileOptions[0]['name'] = this.$store.state.Me.displayName
      }
    }
  }
</script>

<style scoped>
.signin {
  margin: 15px 0 0 auto;
}
</style>
