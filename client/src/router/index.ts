import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import AppLayout from '@/layout/app-layout.vue'
import Page404Layout from '@/layout/page-404-layout.vue'

import UIRoute from '@/pages/admin/ui/route'
import store from '@/store'
import { redirectAuthorizationEndpoint } from '@/lib/apis/api'
import apis from '@/lib/apis'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/:catchAll(.*)',
    redirect: { name: 'dashboard' }
  },
  {
    name: 'admin',
    path: '/',
    component: AppLayout,
    children: [
      {
        name: 'dashboard',
        path: 'dashboard',
        component: () => import('@/pages/admin/dashboard/Dashboard.vue')
      },
      {
        name: 'statistics',
        path: 'statistics',
        component: () => import('@/pages/admin/statistics/Statistics.vue')
      },
      {
        name: 'team',
        path: 'team',
        component: () => import('@/pages/admin/teaminfo/TeamInfo.vue')
      },
      {
        name: 'readme',
        path: 'readme',
        component: () => import('@/pages/admin/readme/Readme.vue')
      },
      {
        name: 'manual',
        path: 'manual',
        component: () => import('@/pages/admin/manual/Manual.vue')
      },
      {
        name: 'faq',
        path: 'faq',
        component: () => import('@/pages/admin/pages/FaqPage.vue')
      },
      UIRoute
    ],
    beforeEnter: async (to, from, next) => {
      try {
        await store.dispatch.getData()
        if (to.path === '/') {
          next('/dashboard')
        }
        next()
      } catch (e) {
        console.error(e)
      }
    }
  },
  {
    path: '/auth/login',
    name: 'login',
    component: () => import('@/pages/auth/Callback.vue'),
    beforeEnter: async (to, from, next) => {
      try {
        await store.dispatch.getData()
        if (!store.state.User) {
          await redirectAuthorizationEndpoint()
        }
        next()
      } catch (e) {
        console.error(e)
      }
    }
  },
  {
    path: '/auth/callback',
    name: 'callback',
    component: () => import('@/pages/auth/Callback.vue'),
    beforeEnter: async (to, from, next) => {
      const code = String(to.query.code)
      await apis.authCallbackGet(code)
      const destination = sessionStorage.getItem('destination')
      if (destination) next(destination)
      else next('/')

      try {
        store.dispatch.getData()
        next('/team-info')
      } catch (e) {
        console.error(e)
      }
    }
  },
  // {
  //   path: '/auth/logout',
  //   name: 'logout',
  //   component: () => import('@/pages/auth/Callback.vue')
  //   // beforeEnter: async (to, from, next) => {
  //   //   if (!store.state.authToken) {
  //   //     return
  //   //   }
  //   //   try {
  //   //     await revokeAuthToken(store.state.authToken)
  //   //     await store.commit.destroySession()
  //   //     await store.dispatch.getData()
  //   //     next('/dashboard')
  //   //   } catch (e) {
  //   //     console.error(e)
  //   //   }
  //   // }
  // },
  {
    path: '/404',
    component: Page404Layout,
    children: [
      {
        name: 'not-found-advanced',
        path: 'not-found-advanced',
        component: () => import('@/pages/404-pages/VaPageNotFoundSearch.vue')
      },
      {
        name: 'not-found-simple',
        path: 'not-found-simple',
        component: () => import('@/pages/404-pages/VaPageNotFoundSimple.vue')
      },
      {
        name: 'not-found-custom',
        path: 'not-found-custom',
        component: () => import('@/pages/404-pages/VaPageNotFoundCustom.vue')
      },
      {
        name: 'not-found-large-text',
        path: '/pages/not-found-large-text',
        component: () => import('@/pages/404-pages/VaPageNotFoundLargeText.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  //  mode: process.env.VUE_APP_ROUTER_MODE_HISTORY === 'true' ? 'history' : 'hash',
  routes
})

router.beforeEach(async (to, from, next) => {
  if (to.path !== '/auth/callback') {
    // ログイン済みかどうか調べる
    if (!store.state.User) {
      try {
        await store.dispatch.fetchMe()
      } catch (e) {
        console.error(e)
      }
    }
    if (!store.state.User) {
      sessionStorage.setItem('destination', to.fullPath)
      redirectAuthorizationEndpoint()
    }
  }
  next()
})

export default router
