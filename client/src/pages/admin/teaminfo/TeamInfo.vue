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
              <va-card-title><h4 class="h-fix">サーバー情報</h4></va-card-title>
              <va-card-content>
                <div class="flex markup-tables">
                  <div
                    class="va-table-responsive server-block"
                    v-for="n in sortedInstance.length"
                    :key="'global' + n"
                  >
                    <h5 class="server-title">サーバー {{ n }}</h5>
                    <table
                      class="va-table va-table--hoverable va-table--striped"
                    >
                      <tr>
                        <td>
                          <strong><span class="md6">状態</span></strong>
                        </td>
                        <td>
                          <span :class="`md6 ${instanceStatusClass(n)}`">{{
                            sortedInstance[n - 1].status
                          }}</span>
                        </td>
                      </tr>
                      <tr>
                        <td>
                          <strong>
                            <span class="md6">グローバル IP アドレス :</span>
                          </strong>
                        </td>
                        <td>
                          <span class="md6">{{
                            sortedInstance[n - 1].global_ip_address
                          }}</span>
                        </td>
                      </tr>
                      <tr>
                        <td>
                          <strong>
                            <span class="md6">プライベート IP アドレス :</span>
                          </strong>
                        </td>
                        <td>
                          <span class="md6">{{
                            sortedInstance[n - 1].private_ip_address
                          }}</span>
                        </td>
                      </tr>
                      <tr>
                        <td>
                          <strong><span class="md6">ユーザー名 :</span></strong>
                        </td>
                        <td>
                          <span class="md6">isucon</span>
                        </td>
                      </tr>
                      <tr>
                        <td>
                          <strong
                            ><span class="md6">初期パスワード :</span></strong
                          >
                        </td>
                        <td>
                          <span class="md6">{{
                            sortedInstance[n - 1].password
                          }}</span>
                        </td>
                      </tr>
                      <tr>
                        <td>
                          <strong
                            ><span class="md6">サーバー作成時間 :</span></strong
                          >
                        </td>
                        <td>
                          <span class="md6">{{
                            sortedInstance[n - 1].CreatedAt
                          }}</span>
                        </td>
                      </tr>
                    </table>
                  </div>
                </div>

                <table class="va-table va-table--hoverable">
                  <tr>
                    <td>
                      <strong
                        ><span class="md6">ベンチマーク回数 :</span></strong
                      >
                    </td>
                    <td>
                      <span class="md6">{{ teamResults.length }}</span>
                    </td>
                  </tr>
                  <tr v-if="teamResults.length > 0">
                    <td>
                      <strong><span class="md6">現在のスコア :</span></strong>
                    </td>
                    <td>
                      <span class="md6">{{
                        team.results.slice(-1)[0].score
                      }}</span>
                    </td>
                  </tr>
                  <tr>
                    <td>
                      <strong><span class="md6">最高スコア :</span></strong>
                    </td>
                    <td>
                      <span class="md6">{{ maxScore.score }}</span>
                    </td>
                  </tr>
                </table>

                <div class="flex md12"></div>
                <div class="form-group">
                  <div class="input-group">
                    <va-input
                      class="mb-4"
                      v-model="betterize"
                      type="textarea"
                      placeholder="改善点を入力してください(記入しないとベンチマークを行えません)"
                    />
                  </div>
                </div>
                <div
                  class="flex md12 my-2"
                  v-for="i in team.instance.length"
                  :key="i"
                >
                  <va-button
                    :rounded="false"
                    class="mr-4 item"
                    @click="benchmark(i)"
                    :disabled="benchmarkButton(i) || betterize === ''"
                  >
                    サーバ{{ i }}にベンチマークを行う
                  </va-button>
                  <va-button
                    :rounded="false"
                    class="mr-4 item"
                    :color="instanceButtonColor(i)"
                    @click="setOperationModal(i)"
                    :disabled="instanceButton(i) || waiting"
                  >
                    {{ instanceButtonMessage(i) }}
                  </va-button>
                </div>
                <div v-if="error" class="type-articles">
                  {{ error }}
                </div>
              </va-card-content>
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
    <!-- <va-modal :okText="'閉じる'" ref="largeModal" title="結果詳細">
    <pre>{{ modalText }}</pre>
  </va-modal> -->
    <va-modal hide-default-actions v-model="showOperationModal">
      <template #header>
        <h3>確認</h3>
      </template>
      <slot>
        <div>
          この操作は取り消せません。間違えて行わないように注意してください
        </div>
      </slot>
      <template #footer>
        <va-button @click="operationInstance(operationInstanceNumber)">
          実行
        </va-button>
        <va-button @click="showOperationModal = false"> キャンセル </va-button>
      </template>
    </va-modal>
    <va-modal v-model="showInfoModal" :message="infoModalMessage" />
  </div>
</template>
<script lang="ts">
//TODO: ファイル分割する
/* eslint-disable @typescript-eslint/camelcase */
import { AxiosError } from 'axios'
import apis, {
  PostBenchmarkRequest,
  PostTeamRequest,
  Response,
  Result,
  User
} from '../../../lib/apis'
import store from '../../../store'
import { computed, ref, watchEffect } from 'vue'
export default {
  setup() {
    const teamMembers = ref<User[]>([])
    const teamName = ref('')
    const makeInstanceButton = ref(false)
    const modalText = ref('')
    const betterize = ref('')
    const error = ref('')
    const showOperationModal = ref(false)
    const showInfoModal = ref(false)
    const infoModalMessage = ref('')
    const operationInstanceNumber = ref(0)
    const waiting = ref(false)
    const largeModal = ref(false)
    const team = computed(() => store.state.Team)
    const user = computed(() => store.state.User)
    const lastResult = computed(() => store.getters.lastResult)
    const maxScore = computed(() => store.getters.maxScore)
    const teamResults = computed(() =>
      store.state.Team && store.state.Team.results
        ? store.state.Team.results
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
    const tweetURL = computed(() => {
      try {
        if (!store.getters.lastResult) {
          return
        }
        const result = JSON.parse(store.getters.lastResult)
        return `https://twitter.com/intent/tweet?text=Pisconで${result.score}点を取りました！%0d&url=https%3A%2F%2Fpiscon.nagatech.work&hashtags=traPiscon`
      } catch (e) {
        return `https://twitter.com/intent/tweet?text=Pisconはじめました！%0d&url=https%3A%2F%2Fpiscon.nagatech.work&hashtags=traPiscon`
      }
    })
    const showInfo = (i: number) => {
      showInfoModal.value = true
      console.log(teamResults.value)
      const betterize: Array<string> = [
        '改善点：' +
          (teamResults.value[i].bettrize ? teamResults.value[i].bettrize : '')
      ]
      infoModalMessage.value = betterize
        .concat(
          teamResults.value[i].messages
            ? teamResults.value[i].messages.map(a => (a.text ? a.text : ''))
            : []
        )
        .join('\n')
    }
    const instanceButtonMessage = (i: number) =>
      computed(() => {
        if (!sortedInstance.value) {
          return
        }
        switch (sortedInstance.value[i - 1].status) {
          case 'ACTIVE':
            return `インスタンス${
              sortedInstance.value[i - 1].instance_number
            }を削除する`

          case 'NOT_EXIST':
            return `インスタンス${
              sortedInstance.value[i - 1].instance_number
            }を作成する`

          case 'BUILDING':
            return `作成中`

          default:
            return sortedInstance.value[i - 1].status
        }
      }).value
    const benchmarkButton = (i: number) =>
      computed(
        () =>
          store.state.Queue?.find(
            que => que.team.name === store.state.Team?.name
          ) ||
          (!sortedInstance.value
            ? false
            : sortedInstance.value[i - 1].status !== 'ACTIVE')
      ).value
    const instanceButton = (i: number) => {
      return computed(
        () =>
          (!sortedInstance.value || !sortedInstance.value[i - 1]
            ? false
            : sortedInstance.value[i - 1].status !== 'ACTIVE') &&
          (!sortedInstance.value || !sortedInstance.value[i - 1]
            ? false
            : sortedInstance.value[i - 1].status !== 'NOT_EXIST')
      ).value
    }
    const instanceButtonColor = (i: number) =>
      computed(() => {
        if (!sortedInstance.value) {
          return
        }
        switch (sortedInstance.value[i - 1].status) {
          case 'ACTIVE':
            return 'danger'

          case 'NOT_EXIST':
            return 'info'

          default:
            return 'info'
        }
      }).value
    const instanceStatusClass = (i: number) =>
      computed(() => {
        if (!sortedInstance.value) {
          return
        }
        switch (sortedInstance.value[i - 1].status) {
          case 'ACTIVE':
            return 'text-primary'

          case 'NOT_EXIST':
            return 'text-muted'

          case 'BUILDING':
            return 'text-info'

          default:
            return 'text-primary'
        }
      }).value
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
      const userReq = await apis
        .teamPost(req)
        .then(data => {
          if (!store.state.User) {
            return
          }
          const req: User = {
            name: store.state.User.name,
            screen_name: store.state.User.screen_name,
            team_id: data.data.ID
          }
          return req
        })
        .catch(err => {
          error.value = err
          return
        })
      if (!userReq) {
        return
      }
      const user = await apis.userPost(userReq).then(res => res.data)
      store.commit.setUser(user)
      await store.dispatch.getData().catch(e => console.warn(e))
    }
    const benchmark = (id: number) => {
      if (benchmarkButton(id) || !store.state.Team) {
        return
      }
      const req: PostBenchmarkRequest = {
        betterize: betterize.value
      }
      apis
        .benchmarkNameInstanceNumberPost(store.state.Team?.name, id, req)
        .then(() => {
          betterize.value = ''
          store.dispatch.getData()
        })
        .catch((err: AxiosError<Response>) => {
          error.value = !err.response?.data.message
            ? ''
            : err.response?.data.message
        })
    }
    const showModal = (data: Result) => {
      modalText.value = JSON.stringify(data, null, '  ')
      largeModal.value = true
    }
    const setOperationModal = (id: number) => {
      showOperationModal.value = true
      operationInstanceNumber.value = id
    }
    const operationInstance = (id: number) => {
      showOperationModal.value = false
      waiting.value = true
      if (
        instanceButton(id) ||
        !sortedInstance.value ||
        !sortedInstance.value[id - 1] ||
        !store.state.Team
      ) {
        return
      }
      switch (sortedInstance.value[id - 1].status) {
        case 'ACTIVE':
          apis
            .instanceTeamIdInstanceNumberDelete(store.state.Team?.ID, id)
            .then(() => {
              store.dispatch.getData()
              waiting.value = false
            })
            .catch((err: AxiosError<Response>) => {
              error.value = !err.response?.data.message
                ? ''
                : err.response?.data.message
              waiting.value = false
            })
          break
        case 'NOT_EXIST':
          apis
            .instanceTeamIdInstanceNumberPost(store.state.Team.ID, id)
            .then(() => {
              store.dispatch.getData()
              waiting.value = false
            })
            .catch((err: AxiosError<Response>) => {
              error.value = !err.response?.data.message
                ? ''
                : err.response?.data.message
              waiting.value = false
            })
          break
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
      benchmark,
      showModal,
      setOperationModal,
      operationInstance,
      instanceButtonMessage,
      instanceButtonColor,
      instanceStatusClass,
      registerTeam,
      benchmarkButton,
      instanceButton,
      showInfo,
      betterize,
      teamMembers,
      error,
      tweetURL,
      team,
      user,
      lastResult,
      maxScore,
      teamName,
      sortedInstance,
      teamResults,
      waiting,
      showOperationModal,
      operationInstanceNumber,
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
