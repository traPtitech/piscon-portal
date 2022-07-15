<template>
  <div class="row row-equal">
    <div class="flex xs12 xl12">
      <va-card>
        <va-card-title>
          <div class="display-4">スコア推移</div>
        </va-card-title>
        <va-card-content>
          <!-- <va-chart :data="scoreAllData" type="line" /> -->
          <LineChart :chartData="scoreAllData" height="200" />
        </va-card-content>
      </va-card>
    </div>
    <!-- <div class="flex xs12 xl12">
      <va-card>
        <va-card-title>
          <div class="display-4">スコア推移(21B)</div>
        </va-card-title>
        <va-card-content>
          <va-chart :data="score21BData" type="line" />
        </va-card-content>
      </va-card>
    </div> -->
  </div>
</template>
<script lang="ts">
import store from '../../../store'
import chroma from 'chroma-js'
import { Datasets, LineChartDataSets } from '../../../lib/apis/types'
import { LineChart } from 'vue-chart-3'
export default {
  components: { LineChart },
  setup() {
    const fetchAllData = () => {
      const datasets = !store.state.AllResults
        ? []
        : store.state.AllResults.filter(
            a => a.results.filter(r => r.pass).length > 0 //少なくとも一回は成功しているチームの
          )
            .map((team, i, c) => {
              const color = chroma((360 / c.length) * i, 0.6, 0.4, 'hsl')
              const td: Datasets = {
                label: team.name,
                tension: 0,
                fill: false,
                pointBackgroundColor: color.alpha(0.8).css(),
                borderColor: color.alpha(0.6).css(),
                pointBorderColor: color.darken(0.4).alpha(0.8).css(),
                pointMoverBackgroundColor: color.darken(2).alpha(0.8).css(),
                spanGaps: true,
                data: []
              }
              td.data = team.results
                .filter(r => r.pass && r.score !== 0)
                .map(r => {
                  return { x: r.id.toString(), y: r.score }
                })
              td.data.push({
                x: store.getters.resultCount.toString(),
                y: td.data[td.data.length - 1]
              })
              return td
            })
            .filter(a => a.data.length > 0)
      const labels = new Array(store.getters.resultCount)
        .fill(0)
        .map((_, i) => (i + 1).toString())
      const res: LineChartDataSets = {
        labels: labels,
        datasets: datasets
      }
      return res
    }

    // const fetch21BData = () => {
    //   const data = {
    //     datasets: [] as LineChartDataSets[]
    //   }
    //   data.datasets = !store.state.AllResults
    //     ? []
    //     : store.state.AllResults.filter(a => a.group === '21B')
    //         .filter(a => a.results.filter(r => r.pass).length > 0)
    //         .map((team, i, c) => {
    //           const color = chroma((360 / c.length) * i, 0.6, 0.4, 'hsl')
    //           const td: LineChartDataSets = {
    //             label: team.name,
    //             fill: false,
    //             lineTension: 0,
    //             pointBackgroundColor: color.alpha(0.8).css(),
    //             borderColor: color.alpha(0.6).css(),
    //             pointBorderColor: color.darken(0.4).alpha(0.8).css(),
    //             pointMoverBackgroundColor: color.darken(2).alpha(0.8).css(),
    //             data: []
    //           }
    //           td.data = team.results
    //             .filter(r => r.pass && r.score !== 0)
    //             .map(r => {
    //               return { x: r.id, y: r.score, time: r.created_at }
    //             })
    //           return td
    //         })
    //         .filter(a => a.data.length > 0)
    //         .sort((a, b) => {
    //           const po = a.data.reduce((c, d) => {
    //             return c < d.y ? d.y : c
    //           }, 0)
    //           const pi = b.data.reduce((c, d) => {
    //             return c < d.y ? d.y : c
    //           }, 0)
    //           return pi - po
    //         })

    //   return data

    const scoreAllData = fetchAllData()
    return {
      scoreAllData
    }
  }
}
</script>