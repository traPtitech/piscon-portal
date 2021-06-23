import Vue from 'vue'
import Vuex from 'vuex'
import VuexI18n from 'vuex-i18n' // load vuex i18n module

import app from './modules/app'
import * as getters from './getters'
import { setAuthToken, getMe, getRsults, getNewer, getQueue, getTeam, getUser } from '../api'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

const initState = {
  Me: null,
  User: {},
  Team: {},
  AllResults: [],
  Que: [],
  Newer: [],
  authToken: null
}

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
    setUser (state, data) {
      state.User = data
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
    },
    destroySession (state) {
      for (let key in initState) {
        state[key] = initState[key]
      }
    }
  },
  actions: {
    async getData ({commit}) {
      getRsults().then(data => commit('setAllResults', data.data))
      getNewer().then(data => commit('setNewer', data.data))
      getQueue().then(data => commit('setQue', data.data))
      const me = await getMe()
        .then(data => {
          commit('setMe', data.data)
          return data.data
        })
        .catch(() => {
          return null
        })

      if (!me) return
      const user = await getUser(me.name)
        .then(data => {
          commit('setUser', data.data)
          return data.data
        })
        .catch(() => {
          return null
        })
      if (!user) return
      getTeam(user.team_id).then(data => commit('setTeam', data.data))
    }
  },
  plugins: [createPersistedState({
    paths: ['authToken']
  })]
})

Vue.use(VuexI18n.plugin, store)

export default store
