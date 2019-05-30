<template>
<div class="modal-page">
<div class="row">
  <div class="col-md-12" v-if="$store.state.Me">
    <vuestic-widget class="col-md-12">
      <div>
        <img :src="`https://q.trap.jp/api/1.0/files/${$store.state.Me.iconFileId}`" class="profile-image" />
        <h3 style="padding: 1rem 0 0 5rem;">{{ $store.state.Me.displayName }}(@{{ $store.state.Me.name }})</h3>
      </div>
    </vuestic-widget>
    <div v-if="$store.state.Team.instance && $store.state.Team.instance.ip_address">
      <vuestic-widget class="col-md-12">
        <div class="widget-header">サーバー情報</div>
        <div class="widget-body">
          <table>
            <tr>
              <td><h6><span class="col-md-6">インスタンス名 :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.name}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">IP アドレス :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.instance.ip_address}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">ユーザー名 :</span></h6></td>
              <td><h6><span class="col-md-6">isucon</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">初期パスワード :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.instance.password}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">ベンチマーク回数 :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.results.length}}</span></h6></td>
            </tr>
            <tr v-if="$store.state.Team.results.length > 0">
              <td><h6><span class="col-md-6">現在のスコア :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.results.slice(-1)[0].score}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">最高スコア :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.getters.maxScore.score}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">作成時間 :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.instance.CreatedAt}}</span></h6></td>
            </tr>
            <tr v-if="$store.state.Team.instance.instance_logs.length > 0">
              <td><h6><span class="col-md-6">最終ログイン時間 :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.instance.instance_logs.slice(-1)[0].CreatedAt}}</span></h6></td>
            </tr>
          </table>
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
      <vuestic-widget v-if="$store.state.Team.group !== '054409cd-97bb-452e-a5ee-a28fa55ea127'" class="col-md-12">
        <div class="widget-header">広告</div>
        <div class="widget-body">
          <p>
            今回の部内ISUCONの鯖代は運営のポケットマネーから捻出されています。<br>
            部内ISUCONの運営を支援していただけるという方は投げ銭をしていただけると@nagatechが泣いて喜びます。
          </p>
        </div>
      </vuestic-widget>
      <vuestic-widget class="col-md-12" headerText="最新の結果">
        <pre>{{$store.getters.lastResult}}</pre>
        <a :href="tweetURL" target="_blank">
          <button class="btn btn-info btn-with-icon rounded-icon">
            <img src="../../assets/twitter_logo.png" style="width: 55px; border-radius: 50%;" />
          </button>
        </a>
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
        <p>このページはtraP部員専用です！</p>
        <p>traP部員は右上の「Signin with traQ」よりログインすることができます。</p>
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
import { getMeGroup } from '../../api'

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
    async makeInstance () {
      if (this.makeInstanceButton) return
      this.makeInstanceButton = true
      const group = await getMeGroup().then(res => res.data[0])
      axios.post('/api/team', {
        name: this.$store.state.Me.name,
        screenName: this.$store.state.Me.displayName,
        iconFileId: this.$store.state.Me.iconFileId,
        group: group
      })
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
        return `https://twitter.com/intent/tweet?text=Pisconで${result.score}点を取りました！%0d&url=https://piscon.nagatech.work&hashtags=traPiscon`
      } catch (e) {
        return `https://twitter.com/intent/tweet?text=Pisconはじめました！%0d&url=https://piscon.nagatech.work&hashtags=traPiscon`
      }
    }
  }
}
</script>

<style scoped>
.btn-nano {
  font-size: 11px;
  padding: 6px;
  border-radius: 9px;
}

.profile-image {
  float: left;
  width: 55px;
  border-radius: 50%;
}

</style>

