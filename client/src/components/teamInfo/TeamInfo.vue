<template>
<div class="modal-page">
<div class="row">
  <div class="col-md-12" v-if="$store.state.Me">
    <div v-if="$store.state.Team.instance && $store.state.Team.instance.ip_address">
      <vuestic-widget class="col-md-12">
        <div class="widget-header">サーバー情報</div>
        <div class="widget-body">
          <h6><span class="col-md-6">インスタンス名 : {{$store.state.Team.name}}</span></h6>
          <h6><span class="col-md-6">IP アドレス : {{$store.state.Team.instance.ip_address}}</span></h6>
          <h6><span class="col-md-6">ユーザー名 : isucon</span></h6>
          <h6><span class="col-md-6">パスワード : {{$store.state.Team.instance.password}}</span></h6>
          <h6><span class="col-md-6">ベンチマーク回数 : {{$store.state.Team.results.length}}</span></h6>
          <h6><span class="col-md-6">最高スコア : {{$store.getters.maxScore.score}}</span></h6>
          <h6><span class="col-md-6">作成時間 : {{$store.state.Team.instance.CreatedAt}}</span></h6>
          <div class="col-md-12"></div>
          <div class="form-group">
            <div class="input-group">
              <textarea type="text" id="simple-textarea" required v-model="betterize"></textarea>
              <label class="control-label" for="simple-textarea">改善点を入力してください(記入しないとベンチマークを行えません)</label><i class="bar"></i>
            </div>
          </div>
          <button class="btn btn-micro btn-info" @click="benchmark" :disabled="benchmarkButton || betterize === ''">ベンチマークを行う</button>
          <div v-if="error" class="type-articles">
            {{ error }}
          </div>
        </div>
      </vuestic-widget>
      <vuestic-widget class="col-md-12" headerText="最新の結果">
        <pre>{{$store.getters.lastResult}}</pre>
        <a :href="tweetURL" target="_blank"><button class="btn btn-info btn-with-icon rounded-icon"><div class="btn-with-icon-content"><i style="color:white; top:0.5rem; left:0.5rem;" class="brandico brandico-twitter-bird"></i></div></button></a>
      </vuestic-widget>
      <vuestic-widget class="col-md-12" headerText="これまでの結果">
        <div class="table-responsible">
          <table class="table table-striped table-sm">
            <thead>
              <tr>
                <td>ID</td>
                <td>PASS</td>
                <td>FAIL</td>
                <td>SUCCESS</td>
                <td>SCORE</td>
                <td>TIME</td>
                <td>INFO</td>
              </tr>
            </thead>
            <tbody>
              <tr :class="{'table-danger': !result.pass}" v-for="result in $store.state.Team.results" :key="result.id">
                <td>{{result.id}}</td>
                <td>{{result.pass}}</td>
                <td>{{result.fail}}</td>
                <td>{{result.suceess}}</td>
                <td>{{result.score}}</td>
                <td>{{result.created_at.slice(5, 19)}}</td>
                <td><button class="btn btn-nano btn-info" @click="showModal(result)">詳細</button></td>
              </tr>
            </tbody>
          </table>
        </div>
      </vuestic-widget>
    </div>
    <vuestic-widget v-else>
      <div class="typo-headers">
        <h2>インスタンスが作成されていません</h2>
      </div>
      <div class="type-articles">
        <p>
          インスタンスを作成する場合は下のインスタンスを作成するボタンを押してください
        </p>
      </div>
      <button class="btn btn-primary" @click="makeInstance" :disabled="makeInstanceButton">インスタンスを作成する</button>
      <div v-if="error" class="type-articles">
        {{ error }}
      </div>
    </vuestic-widget>
  </div>
  <div v-else class="col-md-12">
    <vuestic-widget  headerText="traP外の方へ">
      <div class="widget-body">
        このページはtraP部員専用です！<br>
      </div>
    </vuestic-widget>
  </div>
</div>
<vuestic-modal :large="true" :show.sync="show" :okText="'閉じる'" :cancelClass="'none'" ref="largeModal">
  <div slot="title">結果詳細</div>
  <pre>{{modalText}}</pre>
</vuestic-modal>
</div>
</template>

<script>
import axios from '../../services/axios'

export default {
  name: 'team-info',
  data () {
    return {
      makeInstanceButton: false,
      modalText: '',
      betterize: '',
      show: true,
      error: '',
    }
  },
  methods: {
    makeInstance () {
      if (this.makeInstanceButton) return
      this.makeInstanceButton = true
      axios.post('/api/team', {name: this.$store.state.Me.name})
        .then(_ => {
          this.$store.dispatch('getData')
        })
        .catch(err => {
          this.error = err.response.data.message
        })
    },
    benchmark () {
      if (this.benchmarkButton) return
      axios.post(`/api/benchmark/${this.$store.state.Me.name}`, {betterize: this.betterize})
        .then(_ => {
          this.betterize = ''
          this.$store.dispatch('getData')
        }).catch(err => {
          this.error = err.response.data.message
        })
    },
    showModal (data) {
      this.modalText = JSON.stringify(data, null, '  ')
      this.$refs.largeModal.open()
    }
  },
  computed: {
    benchmarkButton () {
      return this.$store.state.Que.find(que => que.team.name === this.$store.state.Me.name)
    },
    tweetURL () {
      try {
        const result = JSON.parse(this.$store.getters.lastResult)
        return `https://twitter.com/intent/tweet?text=PISCONで${result.score}点を取りました！%0d%0a&hashtags=traPiscon`
      } catch (e) {
        return `https://twitter.com/intent/tweet?text=PISCONはじめました！%0d%0a&hashtags=traPiscon`
      }
    }
  }
}
</script>

<style>
.btn-nano {
  font-size: 11px;
  padding: 6px;
  border-radius: 9px;
}

</style>

