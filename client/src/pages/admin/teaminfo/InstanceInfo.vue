<template>
  <div>
    <va-card-title
      ><h4 class="h-fix">サーバー情報</h4>
      <va-button :rounded="false" class="ml-2 mt-3" @click="fetchInstanceInfo"
        >更新</va-button
      ></va-card-title
    >
    <va-card-content>
      <div class="flex markup-tables">
        <div
          class="va-table-responsive server-block"
          v-for="n in sortedInstance.length"
          :key="'global' + n"
        >
          <h5 class="server-title">サーバー {{ n }}</h5>
          <table class="va-table va-table--hoverable va-table--striped">
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
                <strong><span class="md6">初期パスワード :</span></strong>
              </td>
              <td>
                <span class="md6">{{ sortedInstance[n - 1].password }}</span>
              </td>
            </tr>
            <tr>
              <td>
                <strong><span class="md6">サーバー作成時間 :</span></strong>
              </td>
              <td>
                <span class="md6">{{ sortedInstance[n - 1].CreatedAt }}</span>
              </td>
            </tr>
          </table>
        </div>
      </div>

      <table class="va-table va-table--hoverable">
        <tr>
          <td>
            <strong><span class="md6">ベンチマーク回数 :</span></strong>
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
            <span class="md6">{{ team.results.slice(1)[0].score }}</span>
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
    </va-card-content>
  </div>
</template>
<script lang="ts">
import { defineComponent, computed, PropType } from 'vue'
import { Instance, Result, Team } from '../../../lib/apis'
import store from '../../../store'
export default defineComponent({
  props: {
    sortedInstance: {
      type: Array as PropType<Array<Instance>>,
      required: true
    },
    teamResults: { type: Array as PropType<Array<Result>>, required: true },
    team: { type: Object as PropType<Team>, required: true }
  },
  setup(props) {
    const maxScore = computed(() => store.getters.maxScore)
    const instanceStatusClass = (i: number) =>
      computed(() => {
        if (!props.sortedInstance) {
          return
        }
        switch (props.sortedInstance[i - 1].status) {
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
    const fetchInstanceInfo = () => store.dispatch.fetchInstances()
    return {
      instanceStatusClass,
      fetchInstanceInfo,
      maxScore
    }
  }
})
</script>