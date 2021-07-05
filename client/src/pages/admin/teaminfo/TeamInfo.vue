<template>
  <div class="modal-page">
    <Modal
      v-if="showOperationModal"
      @close="showOperationModal = false"
      @operation="operationInstance(operationInstanceNumber)"
    >
      <slot class="header">確認</slot>
      <slot class="body"
      >この操作は取り消せません。間違えて行わないように注意してください</slot
      >
    </Modal>
    <div class="row">
      <div class="flex md12" v-if="me && team &&team.name">
        <va-card class="flex md12" v-if="team.name">
          <div>
            <img
              :src="`https://q.trap.jp/api/1.0/files/8fc25ba6-a0c9-493d-81d2-11f0870d711f`"
              class="profile-image"
            />
            <h3 style="padding: 1rem 0 0 5rem">{{ team.name }}</h3>
          </div>
        </va-card>
        <va-card class="flex md12">
          <div>
            <img
              :src="`https://q.trap.jp/api/1.0/public/icon/${me.name}`"
              class="profile-image"
            />
            <h3 style="padding: 1rem 0 0 5rem">
              {{ me.displayname }}(@{{ me.name }})
            </h3>
          </div>
        </va-card>
        <div v-if="user && team.instance">
          <va-card class="flex md12">
            <div class="widget-header">サーバー情報</div>
            <div class="widget-body">
              <h6>チーム名 : {{ team.name }}</h6>
              <br />
              <table v-for="n in team.max_instance_number" :key="'global' + n">
                <tr>
                  <td>
                    <h6>
                      <span class="md6">サーバ{{ n }}</span>
                    </h6>
                  </td>
                </tr>
                <tr>
                  <td>
                    <h6><span class="md6">状態</span></h6>
                  </td>
                  <td>
                    <h6>
                      <span :class="`md6 ${instanceStatusClass(n)}`">{{
                        sortedInstance.value[n - 1].status
                      }}</span>
                    </h6>
                  </td>
                </tr>
                <tr>
                  <td>
                    <h6>
                      <span class="md6">グローバル IP アドレス :</span>
                    </h6>
                  </td>
                  <td>
                    <h6>
                      <span class="md6">{{
                        sortedInstance.value[n - 1].global_ip_address
                      }}</span>
                    </h6>
                  </td>
                </tr>
                <tr>
                  <td>
                    <h6>
                      <span class="md6">プライベート IP アドレス :</span>
                    </h6>
                  </td>
                  <td>
                    <h6>
                      <span class="md6">{{
                        sortedInstance.value[n - 1].private_ip_address
                      }}</span>
                    </h6>
                  </td>
                </tr>
                <tr>
                  <td>
                    <h6><span class="md6">ユーザー名 :</span></h6>
                  </td>
                  <td>
                    <h6><span class="md6">isucon</span></h6>
                  </td>
                </tr>
                <tr>
                  <td>
                    <h6><span class="md6">初期パスワード :</span></h6>
                  </td>
                  <td>
                    <h6>
                      <span class="md6">{{
                        sortedInstance.value[n - 1].password
                      }}</span>
                    </h6>
                  </td>
                </tr>
                <tr>
                  <td>
                    <h6><span class="md6">サーバー作成時間 :</span></h6>
                  </td>
                  <td>
                    <h6>
                      <span class="md6">{{
                        sortedInstance.value[n - 1].CreatedAt
                      }}</span>
                    </h6>
                  </td>
                </tr>
                <br />
              </table>

              <tr>
                <td>
                  <h6><span class="md6">ベンチマーク回数 :</span></h6>
                </td>
                <td>
                  <h6>
                    <span class="md6">{{ team.results.length }}</span>
                  </h6>
                </td>
              </tr>
              <table>
                <tr v-if="team.results.length > 0">
                  <td>
                    <h6><span class="md6">現在のスコア :</span></h6>
                  </td>
                  <td>
                    <h6>
                      <span class="md6">{{
                        team.results.slice(-1)[0].score
                      }}</span>
                    </h6>
                  </td>
                </tr>
                <tr>
                  <td>
                    <h6><span class="md6">最高スコア :</span></h6>
                  </td>
                  <td>
                    <h6>
                      <span class="md6">{{ maxScore.score }}</span>
                    </h6>
                  </td>
                </tr>
              </table>

              <div class="flex md12"></div>
              <div class="form-group">
                <div class="input-group">
                  <textarea
                    type="text"
                    id="simple-textarea"
                    required
                    v-model="betterize"
                  ></textarea>
                  <label class="control-label" for="simple-textarea"
                  >改善点を入力してください(記入しないとベンチマークを行えません)</label
                  ><i class="bar"></i>
                </div>
              </div>
              <div
                class="flex md12 my-2"
                v-for="i in team.max_instance_number"
                :key="i"
              >
                <button
                  class="btn btn-micro btn-info"
                  @click="benchmark(i)"
                  :disabled="benchmarkButton(i) || betterize === ''"
                >
                  サーバ{{ i }}にベンチマークを行う
                </button>
                <button
                  :class="instanceButtonClass(i)"
                  @click="setOperationModal(i)"
                  :disabled="instanceButton(i) || waiting"
                >
                  {{ instanceButtonMessage(i) }}
                </button>
              </div>
              <div v-if="error" class="type-articles">
                {{ error }}
              </div>
            </div>
          </va-card>
          <!-- <va-card v-if="$store.state.Team.group !== '054409cd-97bb-452e-a5ee-a28fa55ea127'" class="flex md12">
            <div class="widget-header">広告</div>
            <div class="widget-body">
              <p>
                今回の部内ISUCONの鯖代は運営のポケットマネーから捻出されています。<br>
                部内ISUCONの運営を支援していただけるという方は投げ銭をしていただけるとSysAd班が泣いて喜びます。
              </p>
            </div>
          </va-card> -->
          <va-card class="flex md12" headerText="最新の結果">
            <pre>{{ lastResult }}</pre>
          </va-card>
          <va-card class="md12" headerText="これまでの結果">
            <div class="table-responsible">
              <table class="table table-striped table-sm">
                <thead>
                  <tr>
                    <td>ID</td>
                    <td>PASS</td>
                    <td>SCORE</td>
                    <td>TIME</td>
                    <td>INFO</td>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    :class="{ 'table-danger': !result.pass }"
                    v-for="result in team.results"
                    :key="result.id"
                  >
                    <td>{{ result.id }}</td>
                    <td>{{ result.pass }}</td>
                    <td>{{ result.score }}</td>
                    <td>{{ result.created_at.slice(5, 19) }}</td>
                    <td>
                      <button
                        class="btn btn-nano btn-info"
                        @click="showModal(result)"
                      >
                        詳細
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </va-card>
        </div>
      </div>
      <div v-else class="flex md12">
        <va-card>
          <va-card-title> 参加者専用ページ </va-card-title>
          <va-card-content>
            <p>このページは参加者専用です！</p>
          </va-card-content>
        </va-card>
      </div>
    </div>
    <va-modal :okText="'閉じる'" ref="largeModal" title="結果詳細">
      <pre>{{ modalText }}</pre>
    </va-modal>
  </div>
</template>

<script lang="ts">
import { AxiosError } from 'axios'
import traqApis from '@/lib/apis/traq'
import Modal from './Modal.vue'
import apis, {
  PostBenchmarkRequest,
  PostTeamRequest,
  Response,
  Result
} from '@/lib/apis'
import store from '@/store'
import { computed, ref } from 'vue'
export default {
  components: {
    Modal
  },
  setup() {
    const makeInstanceButton = ref(false)
    const modalText = ref('')
    const betterize = ref('')
    const error = ref('')
    const showOperationModal = ref(false)
    const operationInstanceNumber = ref(0)
    const waiting = ref(false)
    const largeModal = ref(false)
    const team = computed(() => store.state.Team)
    const me = computed(() => store.state.me)
    const user = computed(() => store.state.User)
    const lastResult = computed(() => store.getters.lastResult)
    const maxScore = computed(() => store.getters.maxScore)
    const sortedInstance = computed(() =>
      store.state.Team?.instance
        .map(v => v)
        .sort((a, b) => {
          if (a.instance_number > b.instance_number) {
            return 1
          } else {
            return -1
          }
        })
    )
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
      })
    const benchmarkButton = (i: number) =>
      computed(
        () =>
          store.state.Queue?.find(
            que => que.team.name === store.state.Team?.name
          ) ||
          (!sortedInstance.value
            ? false
            : sortedInstance.value[i - 1].status !== 'ACTIVE')
      )
    const instanceButton = (i: number) =>
      computed(
        () =>
          (!sortedInstance.value
            ? false
            : sortedInstance.value[i - 1].status !== 'ACTIVE') &&
          (!sortedInstance.value
            ? false
            : sortedInstance.value[i - 1].status !== 'NOT_EXIST')
      )
    const instanceButtonClass = (i: number) =>
      computed(() => {
        if (!sortedInstance.value) {
          return
        }
        switch (sortedInstance.value[i - 1].status) {
          case 'ACTIVE':
            return `btn btn-micro btn-danger`

          case 'NOT_EXIST':
            return `btn btn-micro btn-info`

          default:
            return `btn btn-micro btn-info`
        }
      })
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
      })
    const makeInstance = async () => {
      if (makeInstanceButton.value) {
        return
      }
      makeInstanceButton.value = true
      const group = await traqApis.getMe().then(res => res.data.groups[0])
      if (!store.state.me) {
        return
      } //TODO
      const req: PostTeamRequest = {
        name: store.state.me.name,
        group: group
      }
      apis
        .teamPost(req)
        .then(() => {
          store.dispatch.getData()
        })
        .catch(err => {
          error.value = err
        })
    }
    const benchmark = (id: number) => {
      if (benchmarkButton(id).value || !store.state.Team) {
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
      if (instanceButton(id) || !sortedInstance.value || !store.state.Team) {
        return
      }
      switch (sortedInstance.value[id - 1].status) {
        case 'ACTIVE':
          apis
            .instanceTeamIdInstanceNumberDelete(store.state.Team?.id, id)
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
            .instanceTeamIdInstanceNumberPost(store.state.Team.id, id)
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
    return {
      makeInstance,
      benchmark,
      showModal,
      setOperationModal,
      operationInstance,
      instanceButtonMessage,
      instanceButtonClass,
      instanceStatusClass,
      error,
      tweetURL,
      team,
      me,
      user,
      lastResult,
      maxScore
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

.profile-image {
  float: left;
  width: 55px;
  border-radius: 50%;
}
</style>
