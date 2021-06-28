/* eslint-disable @typescript-eslint/member-delimiter-style*/

export interface LineChartDataSets {
  label: string
  fill: boolean
  lineTension: number
  pointBackgroundColor: string
  borderColor: string
  pointBorderColor: string
  pointMoverBackgroundColor: string
  data: {
    x: number
    y: number
    time: string
  }[]
}
