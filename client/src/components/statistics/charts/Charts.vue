<template>
  <div class="charts-page">
    <div class="row">
      <div class="col-md-12">
        <vuestic-widget class="chart-widget" headerText="スコア推移(全体)">
          <vuestic-chart :options="options" :data="scoreAllData" type="line"></vuestic-chart>
        </vuestic-widget>
      </div>
    </div>

    <div class="row">
      <div class="col-md-12">
        <vuestic-widget class="chart-widget" headerText="スコア推移(新入生)">
          <vuestic-chart :options="options" :data="scorefreshData" type="line"></vuestic-chart>
        </vuestic-widget>
      </div>
    </div>

    <div v-if="false">
    <div class="row">
      <div class="col-md-6">
        <vuestic-widget class="chart-widget" :headerText="'charts.verticalBarChart' | translate">
          <vuestic-chart :data="verticalBarChartData" type="vertical-bar"></vuestic-chart>
        </vuestic-widget>
      </div>
      <div class="col-md-6">
        <vuestic-widget class="chart-widget" :headerText="'charts.horizontalBarChart' | translate">
          <vuestic-chart :data="horizontalBarChartData" type="horizontal-bar"></vuestic-chart>
        </vuestic-widget>
      </div>
    </div>


    <div class="row">
      <div class="col-md-6">
        <vuestic-widget class="chart-widget" :headerText="'charts.pieChart' | translate">
          <vuestic-chart :data="pieChartData" type="pie"></vuestic-chart>
        </vuestic-widget>
      </div>
      <div class="col-md-6">
        <vuestic-widget class="chart-widget" :headerText="'charts.donutChart' | translate">
          <vuestic-chart :data="donutChartData" type="donut"></vuestic-chart>
        </vuestic-widget>
      </div>
    </div>

    <div class="row">
      <div class="col-md-12">
        <vuestic-widget class="chart-widget" :headerText="'charts.bubbleChart' | translate">
          <vuestic-chart :data="bubbleChartData" type="bubble"></vuestic-chart>
        </vuestic-widget>
      </div>
    </div>
    </div>
  </div>
</template>

<script>
  import LineChartData from 'data/charts/LineChartData'
  import BubbleChartData from 'data/charts/BubbleChartData'
  import PieChartData from 'data/charts/PieChartData'
  import DonutChartData from 'data/charts/DonutChartData'
  import VerticalBarChartData from 'data/charts/VerticalBarChartData'
  import HorizontalBarChartData from 'data/charts/HorizontalBarChartData'

  import chroma from 'chroma-js'

  export default {
    name: 'charts',
    data: function () {
      return {
        bubbleChartData: BubbleChartData,
        lineChartData: LineChartData,
        pieChartData: PieChartData,
        donutChartData: DonutChartData,
        verticalBarChartData: VerticalBarChartData,
        horizontalBarChartData: HorizontalBarChartData,
        options: {
          scales: {
            xAxes: [{
              type: 'time'
            }]
          },
          lineTension: 0,
        }
      }
    },
    computed: {
      scoreAllData () {
        const data = {
          datasets: []
        }
        data.datasets = this.$store.state.AllResults
        .filter(a => a.results.filter(r => r.pass).length > 0)
        .map((team, i, c) => {
          const color = chroma(360 / c.length * i, 0.6, 0.4, 'hsl')
          const td = {
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
            return {x: r.created_at, y: r.score}
          })
          return td
        }).filter(a => a.data.length > 0).sort((a, b) => {
          const po = a.data.reduce((c, d) => { return c < d.y ? d.y : c }, 0)
          const pi = b.data.reduce((c, d) => { return c < d.y ? d.y : c }, 0)
          return pi - po
        })

        return data
      },
      scorefreshData () {
        const data = {
          datasets: []
        }
        data.datasets = this.$store.state.AllResults
        .filter(a => a.group === '054409cd-97bb-452e-a5ee-a28fa55ea127')
        .filter(a => a.results.filter(r => r.pass).length > 0)
        .map((team, i, c) => {
          const color = chroma(360 / c.length * i, 0.6, 0.4, 'hsl')
          const td = {
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
            return {x: r.created_at, y: r.score}
          })
          return td
        }).filter(a => a.data.length > 0).sort((a, b) => {
          const po = a.data.reduce((c, d) => { return c < d.y ? d.y : c }, 0)
          const pi = b.data.reduce((c, d) => { return c < d.y ? d.y : c }, 0)
          return pi - po
        })

        return data
      }
    }
  }
</script>

<style lang="scss">
  .widget.chart-widget {
    .widget-body {
      height: 550px;
    }
  }
</style>
