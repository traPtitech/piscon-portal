<template>
  <va-card color="background" style="padding: 0.75rem;">
    <va-card-title>
      <div class="display-4">
        最近のベンチマーク
      </div>
    </va-card-title>
    <div class="va-table-responsive">
      <table class="va-table va-table--hoverable">
        <thead>
          <tr>
            <th>ID</th>
            <th>TEAM ID</th>
            <th>PASS</th>
            <th>SCORE</th>
            <th>TIME</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in results" :key="r.id">
            <td>{{r.id}}</td>
            <td>{{r.team_id}}</td>
            <td>{{r.pass}}</td>
            <td>{{r.score}}</td>
            <td>{{r.created_at.slice(5,16)}}</td>
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
      results : computed(() => store.getters.resentResults)
    }
  },
}
</script>