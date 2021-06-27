<template>
  <va-card color="background" style="padding: 0.75rem;">
    <va-card-title>
      <div class="display-4">
        順位表
      </div>
    </va-card-title>
    <div class="va-table-responsive">
      <table class="va-table va-table--hoverable">
        <thead>
          <tr>
            <th>RANK</th>
            <th>TEAM</th>
            <th>SCORE</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(r,i) in results" :key="r.name">
            <td>{{i+1}}aaa</td>
            <td>{{r.name}}aaa</td>
            <td>{{r.score}}aaa</td>
            <!-- <td>
              <va-badge
                :text="user.status"
                :color="user.status"
                TODO:failedとかのステータスを表示する
              />
            </td> -->
          </tr>
        </tbody>
      </table>
    </div>
  </va-card>
</template>
<script lang="ts">
import { computed } from '@vue/runtime-core'
import { useStore } from '@/store'
export default {
  setup(){
    const store = useStore()
    return{
      results : computed(() => store.getters.rankingData)
    }
  },
}
</script>