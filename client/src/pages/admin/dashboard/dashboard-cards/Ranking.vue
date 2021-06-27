<template>
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
        <tr v-for="(r,i) in results" :key="r.id">
          <td>{{user.fullName}}</td>
          <td>{{user.email}}</td>
          <td>{{user.country}}</td>
          <td>
            <!-- <va-badge
              :text="user.status"
              :color="user.status"
              TODO:failedとかのステータスを表示する
            /> -->

          </td>
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
      results : computed(() => store.state.AllResults)
    }
  },
}
</script>