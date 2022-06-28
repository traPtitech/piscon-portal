export default {
  root: {
    name: '/',
    displayName: 'navigationRoutes.home'
  },
  routes: [
    {
      name: 'dashboard',
      displayName: 'menu.dashboard',
      meta: {
        icon: 'vuestic-iconset-dashboard'
      }
    },
    {
      name: 'statistics',
      displayName: 'menu.statistics',
      meta: {
        icon: 'vuestic-iconset-graph'
      }
    },
    {
      name: 'team',
      displayName: 'menu.team',
      meta: {
        icon: 'vuestic-iconset-user'
      }
    },
    {
      name: 'readme',
      displayName: 'menu.readme',
      meta: {
        icon: 'vuestic-iconset-forms'
      }
    },
    {
      name: 'manual',
      displayName: 'menu.manual',
      meta: {
        icon: 'vuestic-iconset-files'
      }
    }
  ]
}
