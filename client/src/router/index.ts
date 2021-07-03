import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import AuthLayout from '@/layout/auth-layout.vue'
import AppLayout from '@/layout/app-layout.vue'
import Page404Layout from '@/layout/page-404-layout.vue'

import UIRoute from '@/pages/admin/ui/route'

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
    ]
  },
  {
    path: '/auth',
    component: AuthLayout,
    children: [
      {
        name: 'login',
        path: 'login',
        component: () => import('@/pages/auth/login/Login.vue')
      },
      {
        name: 'signup',
        path: 'signup',
        component: () => import('@/pages/auth/signup/Signup.vue')
      },
      {
        name: 'recover-password',
        path: 'recover-password',
        component: () =>
          import('@/pages/auth/recover-password/RecoverPassword.vue')
      },
      {
        path: '',
        redirect: { name: 'login' }
      }
    ]
  },
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

export default router
