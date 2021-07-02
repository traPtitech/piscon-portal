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
          <!-- <vuestic-widget v-if="$store.state.Team.group !== '054409cd-97bb-452e-a5ee-a28fa55ea127'" class="col-md-12">
            <div class="widget-header">広告</div>
            <div class="widget-body">
              <p>
                今回の部内ISUCONの鯖代は運営のポケットマネーから捻出されています。<br>
                部内ISUCONの運営を支援していただけるという方は投げ銭をしていただけるとSysAd班が泣いて喜びます。
              </p>
            </div>
          </vuestic-widget> -->
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

<script lang="ts">
import {AxiosError} from 'axios'
import traqApis from '@/lib/apis/traq'
import Modal from './Modal.vue'
import apis, { PostBenchmarkRequest, PostTeamRequest, Responce } from '@/lib/apis'
import store from '@/store'
import { computed, ref } from 'vue'
export default { 
  components: {
    Modal
  },
  setup() {
      const makeInstanceButton= ref(false)
      const modalText= ref('')
      const betterize= ref('')
      const show= ref(true)
      const error= ref('')
      const showOperationModal= ref(false)
      const operationInstanceNumber= ref(0)
      const waiting= ref(false)
      const sortedInstance = computed(() => store.state.Team?.instance.map(v => v).sort((a,b) =>{
        if(a.instance_number > b.instance_number){
          return 1
        }
        else{
          return -1
        }
      })) 
      const benchmarkButton = (i: number) =>  computed(()=> store.state.Queue?.find(que => que.team.name === store.state.Team?.name) || (!sortedInstance.value ? false : sortedInstance.value[i-1].status !== 'ACTIVE'))

      const makeInstance = async () => {
        if (makeInstanceButton.value) {return}
        makeInstanceButton.value = true
        const group = await traqApis.getMe().then(res => res.data.groups[0])
        if (!store.state.me) {return} //TODO
        const req: PostTeamRequest = {
          name: store.state.me.name,
          group: group
        } 
        apis.teamPost(req).then( () => {store.dispatch.getData()}).catch(err =>{error.value = err})
        }
      const benchmark = (id: number) => {
        if (benchmarkButton(id).value || !store.state.Team) {return}
        const req: PostBenchmarkRequest = {
          betterize: betterize.value
        }
        apis.benchmarkNameInstanceNumberPost(store.state.Team?.name,id,req).then(()=> {
          betterize.value = '' 
          store.dispatch.getData()
        }).catch((err: AxiosError<Responce>)=> {
          error.value = !err.response?.data.message ? "" : err.response?.data.message
        })

      }
      return{
        makeInstance,
        error,
      }
  }, 
  data () {
    return {
      
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

