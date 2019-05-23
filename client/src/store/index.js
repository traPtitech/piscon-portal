import Vue from 'vue'
import Vuex from 'vuex'
import VuexI18n from 'vuex-i18n' // load vuex i18n module

import app from './modules/app'
import * as getters from './getters'
import { setAuthToken, getMe, getRsults, getNewer, getTeam, getQueue } from '../api'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

const store = new Vuex.Store({
  strict: true, // process.env.NODE_ENV !== 'production',
  getters,
  modules: {
    app
  },
  state: {
    Me: null,
    Team: {},
    AllResults: [],
    Que: [],
    Newer: [],
    authToken: null
  },
  mutations: {
    setMe (state, data) {
      state.Me = data
    },
    setTeam (state, data) {
      state.Team = data
    },
    setAllResults (state, data) {
      state.AllResults = data
    },
    setQue (state, data) {
      state.Que = data
    },
    setNewer (state, data) {
      state.Newer = data
    },
    setToken (state, data) {
      state.authToken = data
      setAuthToken(data)
    }
  },
  actions: {
    async getData ({commit}) {
      const me = await getMe().then(data => { commit('setMe', data.data); return data.data })
      console.log(me)
      getRsults().then(data => commit('setAllResults', data.data))
      getNewer().then(data => commit('setNewer', data.data))

      // if (me.user_id === '-') return
      // TODO: Fix
      getTeam().then(data => commit('setTeam', data.data))
      getQueue().then(data => commit('setQue', data.data))
    }
  },
  plugins: [createPersistedState({
    paths: ['Me', 'Team', 'AllResults', 'Que', 'Newer', 'authToken']
  })]
})

Vue.use(VuexI18n.plugin, store)

export default store
