<template>

  <vuestic-navbar>
    <header-selector slot="selector" :isOpen.sync="valueProxy"/>
    <span slot="logo">Piscon</span>
    <profile-dropdown>
      <img :src="$store.state.Me.user_id!=='-'?`https://q.trapti.tech/static/icon/${$store.state.Me.user_id}/128.png`:`http://i.imgur.com/nfa5itq.png`"/>
    </profile-dropdown>
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
  }
</script>
