<template>
  <div class="row row-equal">
    <div class="flex md12">
      <va-card color="background" style="padding: 0.75rem" class="item">
        <va-card-title>
          <div class="display-4">順位表</div>
        </va-card-title>
        <va-card-content>
          <div class="flex markup-tables">
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
                  <tr v-for="(r, i) in ranking" :key="r.name">
                    <td>{{ i + 1 }}</td>
                    <td>{{ r.name }}</td>
                    <td>{{ r.results.score }}</td>
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
          </div>
        </va-card-content>
      </va-card>
    </div>
    <div class="flex md12">
      <!-- TODO:間隔開ける%&分離 -->
      <va-card color="background" style="padding: 0.75rem" class="item">
        <va-card-title>
          <div class="display-4">最近のベンチマーク</div>
        </va-card-title>
        <va-card-content>
          <div class="flex markup-tables">
            <div class="va-table-responsive">
              <table class="va-table va-table--hoverable">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>TEAM NAME</th>
                    <th>PASS</th>
                    <th>SCORE</th>
                    <th>TIME</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="r in results" :key="r.result.id">
                    <td>{{ r.result.id }}</td>
                    <td>{{ r.name }}</td>
                    <td>{{ r.result.pass }}</td>
                    <td>{{ r.result.score }}</td>
                    <td>{{ r.result.created_at.slice(5, 16) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </va-card-content>
      </va-card>
    </div>
  </div>
</template>
<script lang="ts">
import { computed } from '@vue/runtime-core'
import store from '../../../../store'
export default {
  setup() {
    return {
      results: computed(() => store.getters.resentResults),
      ranking: computed(() => store.getters.rankingData)
    }
  }
}
</script>
<style lang="scss">
.markup-tables {
  .table-wrapper {
    overflow: auto;
  }

  .va-table {
    width: 100%;
  }
}
</style>