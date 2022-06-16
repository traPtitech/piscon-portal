<template>
  <div>
    <va-content>
      <div class="row row-equal">
        <div class="flex md12" v-if="user">
          <div v-if="team">
            <va-card
              class="flex md12 item mb-3"
              style="padding: 1rem"
              v-if="team.name"
            >
              <va-content>
                <div class="md12" style="display: flex">
                  <img
                    :src="`https://q.trap.jp/api/v3/public/icon/BOT_DevOps`"
                    class="item"
                    style="width: 55px"
                  />
                  <h3 class="item ml-3 h-fix">{{ team.name }}</h3>
                </div>
              </va-content>
            </va-card>
            <va-card class="flex md12 item mb-3" style="padding: 1.3rem">
              <va-content>
                <div v-for="m in teamMembers" :key="m.name">
                  <div class="md12 pb-1 pt-1" style="display: flex">
                    <img
                      :src="`https://q.trap.jp/api/v3/public/icon/${m.name}`"
                      class="item profile-image"
                    />
                    <h3 class="item ml-3 h-fix">
                      {{ m.screen_name }}(@{{ m.name }})
                    </h3>
                  </div>
                </div>
              </va-content>
            </va-card>
            <va-card class="flex md12 item mb-3">
              <instance-info
                :teamResults="teamResults"
                :sortedInstance="sortedInstance"
                :team="team"
              />
              <benchmark
                :teamResults="teamResults"
                :sortedInstance="sortedInstance"
              />
            </va-card>
            <va-card class="flex md12 item mb-3">
              <va-card-title>最新の結果</va-card-title>
              <va-card-content
                ><h6>{{ lastResult }}</h6></va-card-content
              >
            </va-card>
            <va-card class="flex md12 item mb-3">
              <va-card-title>これまでの結果</va-card-title>
              <va-card-content>
                <div class="flex markup-tables">
                  <div class="va-table-responsive">
                    <table class="va-table va-table--hoverable">
                      <thead>
                        <tr>
                          <th>ID</th>
                          <th>PASS</th>
                          <th>SCORE</th>
                          <th>TIME</th>
                          <th>INFO</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(r, i) in teamResults" :key="r.id">
                          <td>{{ r.id }}</td>
                          <td>{{ r.pass }}</td>
                          <td>{{ r.score }}</td>
                          <td>{{ r.created_at.slice(5, 16) }}</td>
                          <td>
                            <va-button
                              color="info"
                              size="small"
                              @click="showInfo(i)"
                              >Info</va-button
                            >
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </va-card-content>
            </va-card>
          </div>
          <div v-else>
            <va-card>
              <va-card-title> チーム登録 </va-card-title>

              <va-card-content>
                <div class="flex md12 item mb-3">
                  <div class="row">
                    <!-- <p>チーム名を入力してください。</p>
              <p>
                チーム名はユニークなものでお願いします(チーム名で区別しているため)
              </p>
              TODO:スタイルあてる -->
                    <va-input
                      class="mb-4"
                      v-model="teamName"
                      placeholder="Team Name"
                    />
                    <va-button
                      :rounded="false"
                      class="ml-2"
                      @click="registerTeam"
                      >登録</va-button
                    >
                  </div>
                </div>
              </va-card-content>
            </va-card>
          </div>
        </div>
        <div v-else class="flex md12 item mb-3">
          <va-card>
            <va-card-title> 参加者専用ページ </va-card-title>
            <va-card-content>
              <p>このページは参加者専用です！</p>
            </va-card-content>
          </va-card>
        </div>
      </div>
    </va-content>
    <va-modal v-model="showInfoModal">
      <va-content>
        <p>{{ infoModalMessage.better }}</p>
        <p>{{ infoModalMessage.message }}</p>
      </va-content>
    </va-modal>
  </div>
</template>
<script lang="ts">
//TODO: ファイル分割する
/* eslint-disable @typescript-eslint/camelcase */
import apis, { PostTeamRequest, Result, User, Message } from '../../../lib/apis'
import store from '../../../store'
import { computed, ref, watchEffect } from 'vue'
import InstanceInfo from './InstanceInfo.vue'
import Benchmark from './Benchmark.vue'
export default {
  components: {
    InstanceInfo,
    Benchmark
  },
  setup() {
    const teamMembers = ref<User[]>([])
    const teamName = ref('')
    const makeInstanceButton = ref(false)
    const modalText = ref('')
    const error = ref('')
    const largeModal = ref(false)
    const team = computed(() => store.state.Team)
    const user = computed(() => store.state.User)
    const showInfoModal = ref(false)
    const infoModalMessage = ref<{ better: string; message: string }>({
      better: '',
      message: ''
    })
    const lastResult = computed(() => store.getters.lastResult)
    const teamResults = computed(() =>
      store.state.Team && store.state.Team.results
        ? store.state.Team.results.sort((a, b) => {
            const datea = new Date(a.created_at)
            const dateb = new Date(b.created_at)
            return datea < dateb ? 1 : -1
          })
        : []
    )
    const sortedInstance = computed(() => {
      if (!store.state.Team || !store.state.Team.instance) {
        return []
      }
      const res = store.state.Team.instance
        .map(v => v)
        .sort((a, b) => {
          if (a.instance_number > b.instance_number) {
            return 1
          } else {
            return -1
          }
        })
      return res
    })

    const makeInstance = async () => {
      if (makeInstanceButton.value) {
        return
      }
      makeInstanceButton.value = true
    }
    const registerTeam = async () => {
      if (!store.state.User) {
        return
      } //TODO
      const group = await apis.meGroupGet().then(res => res.data)
      const req: PostTeamRequest = {
        name: teamName.value,
        group: group
      }
      const team = await apis.teamPost(req).then(res => res.data)
      const userReq: User = {
        name: store.state.User.name,
        screen_name: store.state.User.screen_name,
        team_id: team.ID
      }

      const user = await apis.userPost(userReq).then(res => res.data)
      store.commit.setUser(user)
      store.commit.setTeam(team)
      await store.dispatch.fetchData().catch(e => console.warn(e))
    }
    const showModal = (data: Result) => {
      modalText.value = JSON.stringify(data, null, '  ')
      largeModal.value = true
    }
    const showInfo = (i: number) => {
      showInfoModal.value = true
      const betterize =
        '改善点：' +
        (teamResults.value && teamResults.value[i].betterize
          ? teamResults.value[i].betterize
          : '')

      infoModalMessage.value = {
        better: betterize,
        message: (teamResults.value && teamResults.value[i].messages
          ? teamResults.value[i].messages.map((a: Message) =>
              a.text ? a.text : ''
            )
          : []
        ).join('\n')
      }
    }

    watchEffect(async () => {
      if (!store.state.Team) {
        return [] as User[]
      }
      teamMembers.value = (await apis.teamIdMemberGet(store.state.Team.ID)).data
    })

    return {
      makeInstance,
      showModal,
      registerTeam,
      showInfo,
      teamMembers,
      error,
      team,
      user,
      lastResult,
      teamName,
      sortedInstance,
      teamResults,
      showInfoModal,
      infoModalMessage
    }
  }
}
</script>

<style scoped>
.btn-nano {
  font-size: 11px;
  padding: 6px;
  border-radius: 9px;
}
.h-fix {
  margin-bottom: auto;
  margin-top: auto;
}

.server-block:first-of-type > .server-title {
  margin-top: auto;
}

.profile-image {
  width: 55px;
  border-radius: 50%;
}
</style>
