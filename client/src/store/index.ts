import { createDirectStore } from 'direct-vuex'
import { OAuth2Token, MyUserDetail } from '@traptitech/traq'
import TraqApis from '@/lib/apis/traq'
import { setAuthToken } from '@/lib/apis/api'
import apis, { Task, Team } from '@/lib/apis'

const { store, rootActionContext } = createDirectStore({
  state: {
    me: null as MyUserDetail | null,
    authToken: null as OAuth2Token | null,
    Team: null as Team | null,
    AllResults: null as Team[] | null,
    Queue: null as Task[] | null,
    Newer: null as Team[] | null
  },
  mutations: {
    setMe(state, me: MyUserDetail) {
      state.me = me
    },
    setToken(state, data: OAuth2Token) {
      state.authToken = data
      setAuthToken(data)
    },
    setTeam(state, data: Team) {
      state.Team = data
    },
    setAllResults(state, data: Team[]) {
      state.AllResults = data
    },
    setQueue(state, data: Task[]) {
      state.Queue = data
    },
    setNewer(state, data: Team[]) {
      state.Newer = data
    }
  },
  actions: {
    async fetchMe(context) {
      const { commit } = rootActionContext(context)
      const { data } = await TraqApis.getMe()
      commit.setMe(data)
    },
    async getData(context) {
      const { commit } = rootActionContext(context)
      apis.resultsGet().then(data => commit.setAllResults(data.data))
      apis.newerGet().then(data => commit.setNewer(data.data))
      apis.benchmarkQueueGet().then(data => commit.setQueue(data.data))
    }
  }
})

export default store.original

export type Store = typeof store
export const useStore = (): Store => store
