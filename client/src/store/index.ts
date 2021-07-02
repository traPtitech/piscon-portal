import { createDirectStore } from 'direct-vuex'
import { OAuth2Token, MyUserDetail } from '@traptitech/traq'
import TraqApis from '@/lib/apis/traq'
import { setAuthToken } from '@/lib/apis/api'
import apis, { Result, Task, Team, User } from '@/lib/apis'

const { store, rootActionContext } = createDirectStore({
  state: {
    me: null as MyUserDetail | null,
    User: null as User | null,
    authToken: null as OAuth2Token | null,
    Team: null as Team | null,
    AllResults: null as Team[] | null,
    Queue: null as Task[] | null,
    Newer: null as Team[] | null
  },
  getters: {
    rankingData(state) {
      if (!state.AllResults) {
        return
      }
      return state.AllResults.map(team => {
        const res = {
          name: team.name,
          results: (team.results || [])
            .filter(result => result.pass)
            .reduce((a, b) => {
              return a.score < b.score ? b : a
            })
        }
        return res
      }).sort((a, b) => b.results.score - a.results.score)
    },
    resentResults(state) {
      if (!state.AllResults) {
        return
      }
      const results = state.AllResults.reduce(
        (a, b) => a.concat(b.results || []),
        [] as Result[]
      ).sort((a, b) => b.id - a.id)
      if (results.length > 20) {
        return results.splice(0, 20)
      }
      return results
    },
    resultCount(state) {
      if (!state.AllResults) {
        return 0
      }
      return state.AllResults.reduce((a, b) => a + (b.results || []).length, 0)
    },
    lastResult(state) {
      if (!state.Team) {
        return
      }
      const l = state.Team.results.length
      return l > 0
        ? JSON.stringify(state.Team.results[l - 1], null, '  ')
        : 'まだベンチマークは行われていません'
    }
  },
  mutations: {
    setMe(state, me: MyUserDetail) {
      state.me = me
    },
    setUser(state, data: User) {
      state.User = data
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
      const me = await TraqApis.getMe()
        .then(data => {
          commit.setMe(data.data)
          return data.data
        })
        .catch(() => {
          return null
        })

      if (!me) {
        return
      }
      const user = await apis
        .userNameGet(me.name)
        .then(data => {
          commit.setUser(data.data)
          return data.data
        })
        .catch(() => {
          return null
        })
      if (!user || !user.team_id) {
        //TODO:仕様をちゃんと整える
        return
      } else {
        apis.teamIdGet(user.team_id).then(data => commit.setTeam(data.data))
      }
    }
  }
})

export default store

export type Store = typeof store
