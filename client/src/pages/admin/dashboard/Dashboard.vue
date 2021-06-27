<template>
  <div class="dashboard">
    <ranking />
    <resent-benchmark />
    <dashboard-info-block />
    <div class="row row-equal">
      <div class="flex xs12 lg6">
        <dashboard-tabs @submit="addAddressToMap"/>
      </div>
      <div class="flex xs12 lg6">
        <dashboard-map ref="dashboardMap"/>
      </div>
    </div>
  </div>
</template>

<script>
import DashboardInfoBlock from './DashboardInfoBlock'
import DashboardTabs from './DashboardTabs'
import DashboardMap from './DashboardMap'
import Ranking from './dashboard-cards/Ranking.vue'
import ResentBenchmark from './dashboard-cards/ResentBenchmark.vue'

export default {
  name: 'dashboard',
  components: {
    DashboardInfoBlock,
    DashboardTabs,
    DashboardMap,
    Ranking,
    ResentBenchmark,
  },
  methods: {
    addAddressToMap ({ city, country }) {
      this.$refs.dashboardMap.addAddress({ city: city.text, country })
    },
  },
}
</script>

<style lang="scss">
  .row-equal .flex {
    .va-card {
      height: 100%;
    }
  }

  .dashboard {
    .va-card {
      margin-bottom: 0 !important;
      &__title {
        display: flex;
        justify-content: space-between;
      }
    }
  }
</style>
