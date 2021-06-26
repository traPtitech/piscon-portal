import { createDirectStore } from 'direct-vuex'
import { User, OAuth2Token, MyUserDetail } from '@traptitech/traq'
import TraqApis from '@/lib/apis/traq'
import { setAuthToken } from '@/lib/apis/api'

const { store, rootActionContext } = createDirectStore({
  state: {
    me: null as MyUserDetail | null,
    authToken: null as OAuth2Token | null
  },
  mutations: {
    setMe(state, me: MyUserDetail) {
      state.me = me
    },
    setToken(state, data: OAuth2Token) {
      state.authToken = data
      setAuthToken(data)
    }
  },
  actions: {
    async fetchMe(context) {
      const { commit } = rootActionContext(context)
      const { data } = await TraqApis.getMe()
      commit.setMe(data)
    }
  }
})

export default store.original

export type Store = typeof store
export const useStore = (): Store => store
