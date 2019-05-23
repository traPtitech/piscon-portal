<template>
  <div class="dashboard">

    <dashboard-info-widgets></dashboard-info-widgets>

    <vuestic-widget headerText="直近24時間で得点を伸ばした人"> 
      <li v-for="team in $store.state.Newer">
        {{team.name}}
      </li>
    </vuestic-widget>

    <vuestic-widget class="no-padding no-v-padding" headerText="順位表">
      <table class="table table-striped table-sm">
        <thead>
          <tr>
            <td>RANK</td>
            <td>TEAM</td>
            <td>SUCCESS</td>
            <td>FAIL</td>
            <td>SCORE</td>
          </tr>
        </thead>
        <tbody v-for="(team, i) in $store.getters.rankingData">
          <tr>
            <td>{{i+1}}</td>
            <td>{{team.name}}</td>
            <td>{{team.result.score?team.result.suceess : '-'}}</td>
            <td>{{team.result.score?team.result.fail : '-'}}</td>
            <td>{{team.result.score?team.result.score : '-'}}</td>
          </tr>
        </tbody>
      </table>
    </vuestic-widget>

    <vuestic-widget class="no-padding no-v-padding" headerText="最近のベンチマーク">
      <table class="table table-striped table-sm">
        <thead>
          <tr>
            <td>ID</td>
            <td>TEAM ID</td>
            <td>PASS</td>
            <td>SCORE</td>
            <td>TIME</td>
          </tr>
        </thead>
        <tbody v-for="r in $store.getters.recentResults">
          <tr>
            <td>{{r.id}}</td>
            <td>{{r.team_id}}</td>
            <td>{{r.pass}}</td>
            <td>{{r.score}}</td>
            <td>{{r.created_at.slice(5,16)}}</td>
          </tr>
        </tbody>
      </table>
    </vuestic-widget>
  </div>
</template>

<script>
  import DashboardInfoWidgets from './DashboardInfoWidgets'
  import UsersMembersTab from './users-and-members-tab/UsersMembersTab.vue'
  import SetupProfileTab from './setup-profile-tab/SetupProfileTab.vue'
  import FeaturesTab from './features-tab/FeaturesTab.vue'
  import DataVisualisationTab from './data-visualisation-tab/DataVisualisation.vue'
  import DashboardBottomWidgets from './DashboardBottomWidgets.vue'

  export default {
    name: 'dashboard',
    components: {
      DataVisualisationTab,
      DashboardInfoWidgets,
      UsersMembersTab,
      SetupProfileTab,
      FeaturesTab,
      DashboardBottomWidgets
    },

    data () {
      return {
        fieleds: [
          {name: 'team', sortField: 'name'},
          {name: 'success', sortField: 'success'},
          {name: 'fail', sortField: 'fail'},
          {name: 'score', sortField: 'score'}
        ]
      }
    }
  }

</script>
<style lang="scss" scoped>
</style>
