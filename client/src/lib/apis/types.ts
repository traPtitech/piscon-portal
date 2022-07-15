/* eslint-disable @typescript-eslint/member-delimiter-style*/

export interface LineChartDataSets {
  labels: string[] //横軸のラベル
  datasets: Datasets[]
}

export interface Datasets {
  label: string //線につくラベル
  fill: boolean
  borderColor: string
  pointBackgroundColor: string
  pointMoverBackgroundColor: string
  pointBorderColor: string
  tension: number
  spanGaps: boolean
  data: {}[]
}
