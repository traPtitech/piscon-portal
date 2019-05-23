import Vue from 'vue'
import Router from 'vue-router'
import AppLayout from '../components/admin/AppLayout'
import lazyLoading from './lazyLoading'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '*',
      redirect: { name: 'dashboard' }
    },
    {
      name: 'Admin',
      path: '/',
      component: AppLayout,
      children: [
        {
          name: 'dashboard',
          path: 'dashboard',
          component: lazyLoading('dashboard/Dashboard'),
          default: true
        },
        {
          name: 'team-info',
          path: 'team-info',
          component: lazyLoading('teamInfo/TeamInfo')
        },
        {
          name: 'readme',
          path: 'readme',
          component: lazyLoading('readme/Readme')
        },
        {
          name: 'q-and-a',
          path: 'qa',
          component: lazyLoading('qanda/QandA')
        },
        {
          name: 'statistics',
          path: 'statistics',
          component: lazyLoading('statistics/charts/Charts')
        }
      ]
    }
  ]
})
