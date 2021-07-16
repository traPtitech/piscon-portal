import { createDirectStore } from 'direct-vuex'
import apis, { Result, Task, Team, User } from '@/lib/apis'

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
      const l = state.Team.results ? state.Team.results.length : 0
      return l > 0
        ? JSON.stringify(state.Team.results[l - 1], null, '  ')
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
      await apis
        .meGet()
        .then(data => commit.setUser(data.data))
        .catch(e => console.error(e))
    },
    async getData(context) {
      const { commit } = rootActionContext(context)
      apis.resultsGet().then(data => commit.setAllResults(data.data))
      apis.newerGet().then(data => commit.setNewer(data.data))
      apis.benchmarkQueueGet().then(data => commit.setQueue(data.data))
      if (!store.state.User) {
        return
      }
      await apis
        .userNameGet(store.state.User.name)
        .then(data => {
          commit.setUser(data.data)
          return data.data
        })
        .catch(() => {
          return null
        })
      if (store.state.User) {
        await apis.teamIdGet(store.state.User.team_id).then(data => {
          commit.setTeam(data.data)
        })
      } else {
        console.log('user is empty')
      }
    }
  }
})

export default store

export type Store = typeof store
export const useStore = (): Store => store
