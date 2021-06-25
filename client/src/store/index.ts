import { createDirectStore } from 'direct-vuex'
import { User, OAuth2Token } from '@traptitech/traq'
import { getMe } from '@/apis/api'

const { store, rootActionContext } = createDirectStore({
  state: {
    me: null as User | null,
    authToken: null as OAuth2Token | null
  },
  mutations: {
    setMe(state, me: User) {
      state.me = me
    }
  },
  actions: {
    async fetchMe(context) {
      const { commit } = rootActionContext(context)
      const { data: me } = await getMe()
      commit.setMe(me)
    }
  }
})

export default store.original

export type Store = typeof store
export const useStore = (): Store => store
