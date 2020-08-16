<template>
<div class="modal-page">
  <Modal v-if="showOperationModal" @close="showOperationModal = false" @operation="operationInstance(operationInstanceNumber)">
    <div slot="header">確認</div>
    <div slot="body">この操作は取り消せません。間違えて行わないように注意してください</div>
  </Modal>
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
          <h6>チーム名 : {{$store.state.Team.name}}</h6>
          <br>
          <table v-for="n in $store.state.Team.max_instance_number" :key="'global'+n">
            <tr>
              <td><h6><span class="col-md-6">サーバ{{n}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">状態</span></h6></td>
              <td><h6><span :class="`col-md-6 ${instanceStatusClass(n)}`">{{sortedInstance[n-1].status}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">グローバル IP アドレス :</span></h6></td>
              <td><h6><span class="col-md-6">{{sortedInstance[n-1].global_ip_address}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">プライベート IP アドレス :</span></h6></td>
              <td><h6><span class="col-md-6">{{sortedInstance[n-1].private_ip_address}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">ユーザー名 :</span></h6></td>
              <td><h6><span class="col-md-6">isucon</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">初期パスワード :</span></h6></td>
              <td><h6><span class="col-md-6">{{sortedInstance[n-1].password}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">サーバー作成時間 :</span></h6></td>
              <td><h6><span class="col-md-6">{{sortedInstance[n-1].CreatedAt}}</span></h6></td>
            </tr>
            <br>
          </table>
            
            <tr>
              <td><h6><span class="col-md-6">ベンチマーク回数 :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.results.length}}</span></h6></td>
            </tr>
          <table>
            <tr v-if="$store.state.Team.results.length > 0">
              <td><h6><span class="col-md-6">現在のスコア :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.state.Team.results.slice(-1)[0].score}}</span></h6></td>
            </tr>
            <tr>
              <td><h6><span class="col-md-6">最高スコア :</span></h6></td>
              <td><h6><span class="col-md-6">{{$store.getters.maxScore.score}}</span></h6></td>
            </tr>
          </table>

          <div class="col-md-12"></div>
          <div class="form-group">
            <div class="input-group">
              <textarea type="text" id="simple-textarea" required v-model="betterize"></textarea>
              <label class="control-label" for="simple-textarea">改善点を入力してください(記入しないとベンチマークを行えません)</label><i class="bar"></i>
            </div>
          </div>
          <div class="col-md-12 my-2" v-for="i in $store.state.Team.max_instance_number" :key="i">
            <button class="btn btn-micro btn-info" @click="benchmark(i)" :disabled="benchmarkButton(i) || betterize === ''">サーバ{{i}}にベンチマークを行う</button>
            <button :class="instanceButtonClass(i)" @click="setOperationModal(i)" :disabled="instanceButton(i)||waiting">{{instanceButtonMessage(i)}}</button>
            
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
            部内ISUCONの運営を支援していただけるという方は投げ銭をしていただけるとSysAd班が泣いて喜びます。
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
import Modal from './Modal'

export default {
  name: 'team-info',
  data () {
    return {
      makeInstanceButton: false,
      modalText: '',
      betterize: '',
      show: true,
      error: '',
      showOperationModal: false,
      operationInstanceNumber: 0,
      waiting: false,
    }
  },
  components: {
    Modal
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
      if (this.benchmarkButton(id)) return
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
    },
    setOperationModal (id) {
      this.showOperationModal = true
      console.log(id)
      this.operationInstanceNumber = id
    },
    operationInstance (id) {
      this.showOperationModal = false
      this.waiting = true
      if (this.instanceButton(id)) return
      switch (this.sortedInstance[id - 1].status) {
        case 'ACTIVE':
          console.log('delete instance')
          axios.delete(`/api/instance/${this.$store.state.Team.ID}/${id}`)
          .then(_ => {
            this.$store.dispatch('getData')
            this.waiting = false
          })
          .catch(err => {
            this.err = err.response.data.message
            this.waiting = false
          })
          break
        case 'NOT_EXIST':
          console.log('make instance')
          axios.post(`/api/instance/${this.$store.state.Team.ID}/${id}`)
          .then(_ => {
            this.$store.dispatch('getData')
            this.waiting = false
          })
          .catch(err => {
            this.err = err.response.data.message
            this.waiting = false
          })
          break
      }
    },
  },
  computed: {
    benchmarkButton (i) {
      return function (i) { return this.$store.state.Que.find(que => que.team.name === this.$store.state.Team.name) || this.sortedInstance[i - 1].status !== 'ACTIVE' }
    },
    tweetURL () {
      try {
        const result = JSON.parse(this.$store.getters.lastResult)
        return `https://twitter.com/intent/tweet?text=Pisconで${result.score}点を取りました！%0d&url=https%3A%2F%2Fpiscon.nagatech.work&hashtags=traPiscon`
      } catch (e) {
        return `https://twitter.com/intent/tweet?text=Pisconはじめました！%0d&url=https%3A%2F%2Fpiscon.nagatech.work&hashtags=traPiscon`
      }
    },
    sortedInstance () {
      // console.log(this.$store.state.Team.instance)
      return this.$store.state.Team.instance.map(v => v).sort((a, b) => {
        if (a.instance_number > b.instance_number) {
          return 1
        } else {
          return -1
        }
      })
    },
    instanceButtonMessage () {
      return function (i) {
        switch (this.sortedInstance[i - 1].status) {
          case 'ACTIVE':
            return `インスタンス${this.sortedInstance[i - 1].instance_number}を削除する`

          case 'NOT_EXIST':
            return `インスタンス${this.sortedInstance[i - 1].instance_number}を作成する`

          case 'BUILDING':
            return `作成中`

          default:
            return this.sortedInstance[i - 1].status
        }
      }
    },
    instanceButton () {
      return function (i) {
        return this.sortedInstance[i - 1].status !== 'ACTIVE' && this.sortedInstance[i - 1].status !== 'NOT_EXIST'
      }
    },
    instanceButtonClass (i) {
      return function (i) {
        switch (this.sortedInstance[i - 1].status) {
          case 'ACTIVE':
            return `btn btn-micro btn-danger`

          case 'NOT_EXIST':
            return `btn btn-micro btn-info`

          default:
            return `btn btn-micro btn-info`
        }
      }
    },
    instanceStatusClass (i) {
      return function (i) {
        switch (this.sortedInstance[i - 1].status) {
          case 'ACTIVE':
            return 'text-primary'

          case 'NOT_EXIST':
            return 'text-muted'

          case 'BUILDING':
            return 'text-info'

          default:
            return 'text-primary'
        }
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

