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
        icon: 'vuestic-iconset-statistics'
      },
      disabled: true,
      children: [
        {
          name: 'charts',
          displayName: 'menu.charts'
        },
        {
          name: 'progress-bars',
          displayName: 'menu.progressBars'
        }
      ]
    }
  ]
}
