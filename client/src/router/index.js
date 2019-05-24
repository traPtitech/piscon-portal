import Vue from 'vue'
import Router from 'vue-router'
import AppLayout from '../components/admin/AppLayout'
import lazyLoading from './lazyLoading'
import store from '../store/index'
import { fetchAuthToken, redirectAuthorizationEndpoint, setAuthToken, revokeAuthToken } from '../api'

setAuthToken(store.state.authToken)

Vue.use(Router)

export default new Router({
  mode: 'history',
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
      ],
      beforeEnter: async (to, from, next) => {
        try {
          await store.dispatch('getData')
          next()
        } catch (e) {
          next()
        }
      }
    },
    {
      path: '/auth/login',
      name: 'login',
      component: () => import('../components/auth/login/Login'),
      beforeEnter: async (to, from, next) => {
        try {
          await store.dispatch('getData')
          next()
        } catch (e) {
          await redirectAuthorizationEndpoint()
        }
      }
    },
    {
      path: '/auth/callback',
      name: 'callback',
      component: () => import('../components/auth/Callback'),
      beforeEnter: async (to, from, next) => {
        const code = to.query.code
        const state = to.query.state
        const codeVerifier = sessionStorage.getItem(`login-code-verifier-${state}`)
        if (!code || !codeVerifier) {
          next('/')
        }

        try {
          const res = await fetchAuthToken(code, codeVerifier)
          store.commit('setToken', res.data.access_token)
          store.dispatch('getData')
          next('/team-info')
        } catch (e) {
          console.error(e)
        }
      }
    },
    {
      path: '/auth/logout',
      name: 'logout',
      component: () => import('../components/auth/Callback'),
      beforeEnter: async (to, from, next) => {
        try {
          await revokeAuthToken(store.state.authToken)
          await store.commit('destroySession')
          next('/dashboard')
        } catch (e) {
          console.error(e)
        }
      }
    },
  ]
})
