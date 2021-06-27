<template>
  <div class="va-table-responsive" headerText="順位表">
    <table class="va-table va-table--hoverable">
      <thead>
        <tr>
          <td>ID</td>
          <td>TEAM ID</td>
          <td>PASS</td>
          <td>SCORE</td>
          <td>TIME</td>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(r,i) in results" :key="r.name" headerText="最近のベンチマーク">
          <td>{{i+1}}</td>
          <td>{{r.name}}</td>
          <td>{{r.score}}</td>
        </tr>
      </tbody>
    </table>
  </div>
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