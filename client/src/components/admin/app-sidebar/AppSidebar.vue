<template>
  <vuestic-sidebar :hidden="this.isOpen"
                   :toggleWithoutAnimation="toggleWithoutAnimation">
    <template slot="menu">
      <sidebar-link
        :to="{ name: 'dashboard' }">
        <span slot="title">
          <span class="sidebar-menu-item-icon vuestic-icon vuestic-icon-dashboard"></span>
          <span>Dashboard</span>
        </span>
      </sidebar-link>
      <sidebar-link
        :to="{ name: 'statistics' }">
        <span slot="title">
          <span class="sidebar-menu-item-icon vuestic-icon vuestic-icon-graph"></span>
          <span>Statistics</span>
        </span>
      </sidebar-link>
      <sidebar-link
        :to="{ name: 'team-info' }">
        <span slot="title">
          <span class="sidebar-menu-item-icon vuestic-icon vuestic-icon-user"></span>
          <span v-if="$store.state.Me">TeamInfo</span>
          <span v-else>TeamInfo (traP Only)</span>
        </span>
      </sidebar-link>
      <sidebar-link
        :to="{ name: 'readme' }">
        <span slot="title">
          <span class="sidebar-menu-item-icon vuestic-icon vuestic-icon-forms"></span>
          <span>Readme</span>
        </span>
      </sidebar-link>
      <sidebar-link
        :to="{ name: 'manual' }">
        <span slot="title">
          <span class="sidebar-menu-item-icon vuestic-icon vuestic-icon-files"></span>
          <span v-if="$store.state.Me">Manual</span>
          <span v-else>Manual (traP Only)</span>
        </span>
      </sidebar-link>
      <sidebar-link
        :to="{ name: 'articles' }">
        <span slot="title">
          <span class="sidebar-menu-item-icon vuestic-icon vuestic-icon-extras"></span>
          <span>Articles</span>
        </span>
      </sidebar-link>
      <sidebar-link
        :to="{ name: 'q-and-a' }">
        <span slot="title">
          <span class="sidebar-menu-item-icon vuestic-icon vuestic-icon-comments"></span>
          <span>Q & A</span>
        </span>
      </sidebar-link>
    </template>
  </vuestic-sidebar>
</template>

<script>

  import VuesticSidebar from '../../../vuestic-theme/vuestic-components/vuestic-sidebar/VuesticSidebar'
  import {mapGetters} from 'vuex'
  import SidebarLink from './components/SidebarLink'
  import SidebarLinkGroup from './components/SidebarLinkGroup'

  export default {
    name: 'app-sidebar',
    components: {
      VuesticSidebar,
      SidebarLink,
      SidebarLinkGroup
    },
    props: {
      isOpen: {
        type: Boolean,
        required: true
      }
    },
    computed: {
      ...mapGetters([
        'toggleWithoutAnimation'
      ]),
    }
  }


</script>
