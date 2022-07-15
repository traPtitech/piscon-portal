<template>
  <div class="row row-equal">
    <div class="flex xl12 xs12">
      <div class="row">
        <div
          class="flex xs12 sm6"
        >
          <va-card class="mb-6" :color="'success'">
            <va-card-content>
              <p style="color: white;">総ベンチマーク回数</p>
              <p class="display-2 mb-0" style="color: white;">{{ resultCount }}</p>
            </va-card-content>
          </va-card>
        </div>
        <div
          class="flex xs12 sm6"
        >
          <va-card class="mb-6" :color="'info'">
            <va-card-content>
              <p style="color: white;">参加チーム数</p>
              <p class="display-2 mb-0" style="color: white;">{{ length }}</p>
            </va-card-content>
          </va-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from '@vue/runtime-core'
import store from '@/store'
export default {
  name: 'DashboardInfoBlock',
  setup(){
    return{
      resultCount : computed(() => store.getters.resultCount),
      length : computed(() => store.state.AllResults ? store.state.AllResults.length : 0 ) //TODO:型修正
    }
  }
}
</script>

<style lang="scss" scoped>
  .row-separated {
    .flex + .flex {
      border-left: 1px solid var(--va-background);
    }

    // @include media-breakpoint-down(xs) {
    //   p:not(.display-2) {
    //     font-size: 0.875rem;
    //   }
    // }
  }

  .rich-theme-card-text {
    line-height: 24px;
  }

</style>
