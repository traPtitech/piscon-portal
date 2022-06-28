<template>
  <va-content class="typography content">
    <div class="row">
      <div v-if="user" class="flex md12">
        <va-card>
          <va-card-content>
            <div class="mb-4">
              <h3 class="h-fix">Links</h3>
              <ul>
                <li>
                  <a href="https://github.com/isucon/isucon11-qualify/blob/main/docs/manual.md">
                    ISUCON11 予選当日マニュアル
                  </a>
                </li>
                <li>
                  <a href="https://github.com/isucon/isucon11-qualify/blob/main/docs/isucondition.md">
                    ISUCONDITION アプリケーションマニュアル
                  </a>
                </li>
              </ul>
            </div>

            <div class="mb-4">
              <h3 class="h-fix">補足事項</h3>
              <div class="mb-4">
                <h5>競技環境について</h5>
                <p>
                  ISUCON11 予選当日とは異なり, PISCONでは各チームで競技環境の構築を行う必要はありません.
                  TeamInfo ページより, インスタンスの作成, 及びサーバー情報の確認を行ってください.
                </p>
                <p>
                  競技用インスタンスには <a href="https://aws.amazon.com/jp/ec2/instance-types/t2/">Amazon EC2 T2 インスタンス</a> を使用しています.
                  そのため, 短時間に多数回ベンチマークを行うと, CPU クレジットの不足により, サーバーのパフォーマンスが低下する場合があります.
                  急なパフォーマンスの低下が見られた場合, しばらく時間を置いて, 再度ベンチマークを行ってください.
                  なお, 競技の性質上, 基本的にこの現象が発生することはありません.
                </p>
                <p>
                  参考:
                  <a href="https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/burstable-credits-baseline-concepts.html">
                    バーストパフォーマンスインスタンスに関する主要な概念と定義
                  </a>
                </p>
              </div>

              <div class="mb-4">
                <h5>ブラウザでのアクセスにおける留意点</h5>
                <p>競技用インスタンスで動作している isucondition にブラウザからアクセスする際の留意点です.</p>
                <div class="mb-4">
                  <h6>ログイン</h6>
                  <p>
                    「JIAのアカウントでログイン」を押すと <code>http://localhost:5000</code> に遷移するようになっています.
                    このアクセスは競技用サーバ上で動作する <code>jiaapi-mock.service</code> が受ける想定です.
                  </p>
                  <p>
                    以下のコマンドより <code>localhost:5000</code> が競技用サーバ上の 5000 番ポートにローカルフォワードされるようにした上でブラウザ操作を行ってください.
                  </p>
                  <pre><code>ssh isucon@[競技用サーバのグローバルアドレス] -L 5000:localhost:5000</code></pre>
                </div>
                <div class="mb-4">
                  <h6>ISU の登録</h6>
                  <p>
                    ブラウザより ISU の登録を行う際にも JIA API Mock が必要です. こちらについては
                    <a href="https://github.com/isucon/isucon11-qualify/blob/main/docs/isucondition.md">
                      アプリケーションマニュアル
                    </a> をご確認ください.
                  </p>
                </div>
              </div>
            </div>
          </va-card-content>
        </va-card>
      </div>
      <div v-else class="flex md12">
        <va-card>
          <va-card-title> traP外の方へ </va-card-title>
          <va-card-content>
            <p>このページはtraP部員専用です！</p>
            <p>
              traP部員は右上の「Signin with
              traQ」よりログインすることができます。
            </p>
          </va-card-content>
        </va-card>
      </div>
    </div>
  </va-content>
</template>

<script lang="ts">
import store from '../../../store'
export default {
  setup() {
    const user = store.state.User
    return {
      user
    }
  }
}
</script>

<style lang="sass" scoped>
.mb-4
  margin-bottom: 48px

a
  color: #4ae287

.well
  margin-bottom: 48px
  padding: 1rem 1.5rem

pre
  padding: 16px
  margin-bottom: 1.5em
  line-height: 1.45
  border-radius: 6px
  background: #25292f
  color: #fff
  overflow: auto

code
  font-family: monospace
  line-height: inherit
  overflow: visible
  font-size: 95%

.h-fix
  margin-bottom: auto
  margin-top: auto
</style>

