import { createDirectStore } from 'direct-vuex'
import apis, { Instance, Result, Task, Team, User } from '@/lib/apis'

const { store, rootActionContext } = createDirectStore({
  state: {
    User: null as User | null,
    Team: null as Team | null,
    AllResults: null as Team[] | null,
    Queue: null as Task[] | null,
    Newer: null as Team[] | null,
    isSidebarMinimized: false
  },
  getters: {
    rankingData(state) {
      if (!state.AllResults) {
        return
      }
      return state.AllResults.map(team => {
        const res = {
          name: team.name,
          results: (team.results || ([] as Result[]))
            .filter(result => result.pass)
            .reduce(
              (a, b) => {
                return a.score < b.score ? b : a
              },
              { score: 0 } as Result
            )
        }
        return res
      }).sort((a, b) => b.results.score - a.results.score)
    },
    resentResults(state) {
      if (!state.AllResults) {
        return
      }
      const results = state.AllResults.map(a =>
        a.results.map(r => {
          const res = {
            name: a.name,
            result: r
          }
          return res
        })
      )
        .reduce((a, b) => a.concat(b || []), [])
        .sort((a, b) => b.result.id - a.result.id)
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
      const l = state.Team.results ? state.Team.results.length : 0
      return l > 0
        ? JSON.stringify(state.Team.results[0], null, '  ')
        : 'まだベンチマークは行われていません'
    },
    maxScore(state) {
      if (!state.Team) {
        return
      } else {
        return state.Team.results
          ? state.Team.results.reduce(
            (a, b) => {
              return a.score < b.score ? b : a
            },
            { score: 0 }
          )
          : []
      }
    }
  },
  mutations: {
    setUser(state, data: User) {
      state.User = data
    },
    setInstances(state, data: Instance[]) {
      if (!state.Team) {
        return
      }
      state.Team.instance = data
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
    },
    destroySession(state) {
      state.AllResults = null
      state.User = null
      state.Newer = null
      state.Queue = null
      state.Team = null
    },
    updateSidebarCollapsedState(state, isSidebarMinimized) {
      state.isSidebarMinimized = isSidebarMinimized
    }
  },
  actions: {
    async fetchMe(context) {
      const { commit } = rootActionContext(context)
      const res = await apis.meGet()
      commit.setUser(res.data)
    },
    async fetchInstances(context) {
      const { commit } = rootActionContext(context)
      if (!store.state.Team) {
        return
      }
      apis.teamIdInstancesPut(store.state.Team?.ID).then(data => {
        if (!store.state.Team) {
          return
        }
        commit.setInstances(data.data)
      })
    },
    async fetchUser(context) {
      if (!store.state.User) {
        throw new Error('no user information')
      }

      const { commit } = rootActionContext(context)
      const res = await apis.userNameGet(store.state.User.name)
      commit.setUser(res.data)
    },
    async fetchTeam(context) {
      if (!store.state.User) {
        throw new Error('no user information')
      }

      const { commit } = rootActionContext(context)
      const res = await apis.teamIdGet(store.state.User.team_id)
      commit.setTeam(res.data)
    },
    async getData(context) {
      const { commit } = rootActionContext(context)
      apis.resultsGet().then(data => commit.setAllResults(data.data))
      apis.newerGet().then(data => commit.setNewer(data.data))
      apis.benchmarkQueueGet().then(data => commit.setQueue(data.data))
    }
  }
})

export default store

export type Store = typeof store
export const useStore = (): Store => store
