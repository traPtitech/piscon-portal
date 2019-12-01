<template>
<div class="modal-page">
<div class="row">
  <div class="col-md-12" v-if="$store.state.Me && $store.state.Team.name">
    <vuestic-widget class="col-md-12" v-if="$store.state.Team.name">
      <div>
        <img :src="`https://q.trap.jp/api/1.0/files/8fc25ba6-a0c9-493d-81d2-11f0870d711f`" class="profile-image" />
        <h3 style="padding: 1rem 0 0 5rem;">{{ $store.state.Team.name }}</h3>
      </div>
    </vuestic-widget>
    <vuestic-widget class="col-md-12">
      <div>
        <img :src="`https://q.trap.jp/api/1.0/public/icon/${$store.state.Me.name}`" class="profile-image" />
        <h3 style="padding: 1rem 0 0 5rem;">{{ $store.state.Me.displayName }}(@{{ $store.state.Me.name }})</h3>
      </div>
    </vuestic-widget>
    <div v-if="$store.state.User && $store.state.Team.instance && $store.state.Team.instance">
      <vuestic-widget class="col-md-12">
        <div class="widget-header">サーバー情報</div>
        <div class="widget-body">
          <table>
            <tr>
              <td><h6><span class="col-md-6">チーム名 :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.name}}</span></h6></td>
            </tr>
            <tr v-for="(instance, index) in $store.state.Team.instance" :key="'global'+index">
              <td><h6><span class="col-md-6">サーバ{{index+1}} グローバル IP アドレス :</span></h6></td>
              <td><h6><span class="col-md-6">{{instance.global_ip_address}}</span></h6></td>
            </tr>
            <tr v-for="(instance, index) in $store.state.Team.instance" :key="'private'+index">
              <td><h6><span class="col-md-6">サーバ{{index+1}} プライベート IP アドレス :</span></h6></td>
              <td><h6><span class="col-md-6">{{instance.private_ip_address}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">ユーザー名 :</span></h6></td>
              <td><h6><span class="col-md-6">isucon</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">初期パスワード(共通) :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.instance[0].password}}</span></h6></td>
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
            <tr v-for="(instance, index) in $store.state.Team.instance" :key="'created'+index">
              <td><h6><span class="col-md-6">サーバ{{index+1}} 作成時間 :</span></h6></td>
              <td><h6><span class="col-md-6">{{instance.CreatedAt}}</span></h6></td>
            </tr>
          </table>
          <div class="col-md-12"></div>
          <div class="form-group">
            <div class="input-group">
              <textarea type="text" id="simple-textarea" required v-model="betterize"></textarea>
              <label class="control-label" for="simple-textarea">改善点を入力してください(記入しないとベンチマークを行えません)</label><i class="bar"></i>
            </div>
          </div>
          <div class="col-md-12 my-2" v-for="i in $store.state.Team.instance.length" :key="i">
            <button class="btn btn-micro btn-info" @click="benchmark(i)" :disabled="benchmarkButton || betterize === ''">サーバ{{i}}にベンチマークを行う</button>
          </div>
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
      </vuestic-widget>
      <vuestic-widget class="col-md-12" headerText="これまでの結果">
        <div class="table-responsible">
          <table class="table table-striped table-sm">
            <thead>
              <tr>
                <td>ID</td>
                <td>PASS</td>
                <td>SCORE</td>
                <td>TIME</td>
                <td>INFO</td>
              </tr>
            </thead>
            <tbody>
              <tr :class="{'table-danger': !result.pass}" v-for="result in $store.state.Team.results" :key="result.id">
                <td>{{result.id}}</td>
                <td>{{result.pass}}</td>
                <td>{{result.score}}</td>
                <td>{{result.created_at.slice(5, 19)}}</td>
                <td><button class="btn btn-nano btn-info" @click="showModal(result)">詳細</button></td>
              </tr>
            </tbody>
          </table>
        </div>
      </vuestic-widget>
    </div>
  </div>
  <div v-else class="col-md-12">
    <vuestic-widget  headerText="参加者専用ページ">
      <div class="widget-body">
        <p>このページは参加者専用です！</p>
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
    benchmark (id) {
      if (this.benchmarkButton) return
      axios.post(`/api/benchmark/${this.$store.state.Team.name}/${id}`, {betterize: this.betterize})
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
      return this.$store.state.Que.find(que => que.team.name === this.$store.state.Team.name)
    },
    tweetURL () {
      try {
        const result = JSON.parse(this.$store.getters.lastResult)
        return `https://twitter.com/intent/tweet?text=Pisconで${result.score}点を取りました！%0d&url=https%3A%2F%2Fpiscon.nagatech.work&hashtags=traPiscon`
      } catch (e) {
        return `https://twitter.com/intent/tweet?text=Pisconはじめました！%0d&url=https%3A%2F%2Fpiscon.nagatech.work&hashtags=traPiscon`
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

