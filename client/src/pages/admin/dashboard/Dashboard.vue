<template>
  <div class="dashboard">
    <dashboard-info-block />
    <score-tables />
  </div>
</template>

<script>
import DashboardInfoBlock from './DashboardInfoBlock'
import ScoreTables from './dashboard-cards/ScoreTables.vue'

export default {
  name: 'dashboard',
  components: {
    DashboardInfoBlock,
    ScoreTables
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
