import { createDirectStore } from 'direct-vuex'
import { User } from '@traptitech/traq'
import { api } from '@/apis/api'

const { store, rootActionContext } = createDirectStore({
  state: {
    me: null as User | null
  },
  mutations: {
    setMe(state, me: User) {
      state.me = me
    }
  },
  actions: {
    async fetchMe(context) {
      const { commit } = rootActionContext(context)
      const { data: me } = await api.getMe()
      commit.setMe(me)
    }
  }
})

export default store.original

export type Store = typeof store
export const useStore = (): Store => store
