<template>
  <div class="row row-equal">
    <div class="flex xs12 xl12">
      <va-card>
        <va-card-title>
          <div class="display-4">
            スコア推移(全体)
          </div>
        </va-card-title>
        <va-card-content>
          <va-chart class="chart" ref="lineChart" :data="scoreAllData" type="line"/>
        </va-card-content>
      </va-card>
    </div>
    <div class="flex xs12 xl12">
      <va-card>
        <va-card-title>
          <div class="display-4">
            スコア推移(21B)
          </div>
        </va-card-title>
        <va-card-content>
          <va-chart class="chart" ref="lineChart" :data="score21BData" type="line"/>
        </va-card-content>
      </va-card>
    </div>
  </div>
</template>
<script lang="ts">
import { useStore } from '@/store'
import chroma from 'chroma-js'
import {LineChartDataSets} from '@/lib/apis/types'
export default {
  setup(){
    const store = useStore()
    const scoreAllData= () => {
        const data = {
          datasets: [] as LineChartDataSets[]
        }
        data.datasets = !store.state.AllResults ? [] : store.state.AllResults
        .filter(a => a.results.filter(r => r.pass).length > 0)
        .map((team, i, c) => {
          const color = chroma(360 / c.length * i, 0.6, 0.4, 'hsl')
          const td: LineChartDataSets = {
            label: team.name,
            fill: false,
            lineTension: 0,
            pointBackgroundColor: color.alpha(0.8).css(),
            borderColor: color.alpha(0.6).css(),
            pointBorderColor: color.darken(0.4).alpha(0.8).css(),
            pointMoverBackgroundColor: color.darken(2).alpha(0.8).css(),
            data: [] as {x: number;y: number; time: string}[]
          }
          td.data = team.results.filter(r => r.pass && r.score !== 0).map(r => {
            return {x: r.id, y: r.score, time: r.created_at}
          })
          return td
        }).filter(a => a.data.length > 0).sort((a, b) => {
          const po = a.data.reduce((c, d) => { return c < d.y ? d.y : c }, 0)
          const pi = b.data.reduce((c, d) => { return c < d.y ? d.y : c }, 0)
          return pi - po
        })

        return data
      }
    const score21BData = () => {
        const data = {
          datasets: [] as  LineChartDataSets[]
        }
        data.datasets = !store.state.AllResults ? [] : store.state.AllResults
        .filter(a => a.group === '21B')
        .filter(a => a.results.filter(r => r.pass).length > 0)
        .map((team, i, c) => {
          const color = chroma(360 / c.length * i, 0.6, 0.4, 'hsl')
          const td: LineChartDataSets = {
            label: team.name,
            fill: false,
            lineTension: 0,
            pointBackgroundColor: color.alpha(0.8).css(),
            borderColor: color.alpha(0.6).css(),
            pointBorderColor: color.darken(0.4).alpha(0.8).css(),
            pointMoverBackgroundColor: color.darken(2).alpha(0.8).css(),
            data: []
          }
          td.data = team.results.filter(r => r.pass && r.score !== 0).map(r => {
            return {x: r.id, y: r.score, time: r.created_at}
          })
          return td
        }).filter(a => a.data.length > 0).sort((a, b) => {
          const po = a.data.reduce((c, d) => { return c < d.y ? d.y : c }, 0)
          const pi = b.data.reduce((c, d) => { return c < d.y ? d.y : c }, 0)
          return pi - po
        })

        return data
      }
    return {
      scoreAllData,score21BData
    }
  }
}
</script>