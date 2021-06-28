<template>
  <div class="row row-equal">
    <div class="flex xl12 xs12">
      <div class="row">
        <div
          class="flex xs12 sm12"
        >
          <va-card color="background" style="padding: 0.75rem;">
            <div class="flex xs12">
              <span class="mr-2"> 現在のキュー </span>
              <va-chip square class="mr-2" v-for="que in queue" :key="que.team_id">
                {{que.team}}
              </va-chip>
            </div>
          </va-card>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { computed } from '@vue/runtime-core'
import { useStore } from '@/store'
import { Task } from '@/lib/apis'
export default{
  setup(){
    const store = useStore()
    return{
      queue : computed(() => !store.state.Queue ? [] :store.state.Queue.filter(
        (a: Task): boolean => a.state === 'benchmark' 
      ))
    }
  }
}
</script>
<style lang='scss' scoped>
</style>