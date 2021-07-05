<template>
  <div>
    <div class="row">
      <div v-if="me" class="flex md12">
        <va-card>
          <va-card-title> レギュレーション </va-card-title>
          <va-card-content>
            <div class="box">
              <p>
                指定された競技用サーバー上のアプリケーションのチューニングを行い、それに対するベンチマーク走行のスコアで競技を行います。
                与えられた競技用サーバーのみでアプリケーションの動作が可能であれば、どのような変更を加えても構いません。
                ベンチマーカーとブラウザの挙動に差異がある場合、ベンチマーカーの挙動を正とします。
                また、初期実装は言語毎に若干の挙動の違いはありますが、ベンチマーカーに影響のない挙動に関しては仕様とします。
              </p>
              <h3>ベンチマーク走行</h3>
              <p>ベンチマーク走行は以下のように実施されます。</p>
              <ul>
                <li>初期化処理の実行 POST /initialize (20秒以内)</li>
                <li>
                  アプリケーション互換性チェックの走行 (適宜: 数秒〜数十秒)
                </li>
                <li>負荷走行 (60秒)</li>
                <li>負荷走行後の確認 (適宜: 数秒〜数十秒)</li>
              </ul>
              <p>
                各ステップで失敗が見付かった場合にはその時点で停止します。
                ただし、負荷走行中のエラーについては、タイムアウトや500エラーを含む幾つかのエラーについては無視され、ベンチマーク走行が継続します。
              </p>
              <p>
                また負荷走行が60秒行われた後、レスポンスが返ってきていないリクエストはすべて強制的に切断されます。
                その際にnginxのアクセスログにステータスコード499が記録されることがありますが、これらのリクエストについては減点の対象外です。
              </p>
            </div>
            <div class="box">
              <h3>スコア計算</h3>
              <p>
                スコアは<strong
                >取引が完了した商品（椅子）の価格の合計（ｲｽｺｲﾝ）</strong
                >
                をベースに以下の計算式で計算されます。
              </p>
              <pre><code>取引が完了した商品（椅子）の価格の合計（ｲｽｺｲﾝ） - 減点 = スコア（ｲｽｺｲﾝ）</code></pre>
              <p>
                以下の条件のエラーが発生すると、失格・減点の対象となります。
              </p>
              <ul>
                <li>
                  致命的なエラー
                  <ul>
                    <li>1回以上で失格</li>
                    <li>
                      メッセージの最後に
                      <code>(critical error)</code> が付与されます
                    </li>
                  </ul>
                </li>
                <li>
                  HTTPステータスコードやレスポンスの内容などに誤りがある
                  <ul>
                    <li>1回で500ｲｽｺｲﾝ減点、10回以上で失格</li>
                  </ul>
                </li>
                <li>
                  一定時間内にレスポンスが返却されない・タイムアウト
                  <ul>
                    <li>200回を超えたら100回毎に5000ｲｽｺｲﾝ減点、失格はなし</li>
                    <li>
                      メッセージの最後に
                      <code>（タイムアウトしました）</code> か
                      <code>（一時的なエラー）</code> が付与されます
                    </li>
                  </ul>
                </li>
              </ul>
              <p>
                HTTPステータスコードは、基本的に参照実装と同一のものを想定しています。またベンチマーカーのメッセージは同一のメッセージを1つにまとめます。表示されているメッセージの数とエラー数は一致しないことがあります。
              </p>
              <p>また減点により0ｲｽｺｲﾝ以下になった場合は失格となります。</p>
            </div>
            <div class="box">
              <h3>制約事項</h3>
              <p>
                以下の事項に抵触すると失格(fail)となり、点数が0点になります。
              </p>
              <ul>
                <li>POST /initialize へのレスポンスが20秒以内に戻らない場合</li>
                <li>アプリケーション互換性チェックに失敗した場合</li>
                <li>アプリケーションがレスポンスを10秒以内に返さない場合</li>
                <li>その他、ベンチマーカーのチェッカが失敗を検出したケース</li>
              </ul>
              <p>
                最初に呼ばれる初期化処理
                <code>POST /initialize</code>
                は用意された環境内で、チェッカツールが要求する範囲の整合性を担保します。
                サーバーサイドで処理の変更・データ構造の変更などを行う場合、この処理が行っている内容を漏れなく提供してください。
              </p>
              <p>
                予選終了後に行われる主催者による確認作業（追試）において下記の点が確認できなかった場合は失格となります。
              </p>
              <ul>
                <li>
                  アプリケーションは全て保存データを永続化する必要があります
                  <ul>
                    <li>
                      処理実施後に再起動が行われた場合、再起動前に行われた処理内容が再起動後に保存されている必要があります
                    </li>
                  </ul>
                </li>
                <li>
                  アプリケーションはブラウザ上での表示を初期状態と同様に保つ必要があります
                </li>
              </ul>
              また、以下に示す改変を行ってはいけません。
              <ul>
                <li>パスワードを平文で保存する</li>
              </ul>
            </div>
            <div class="box">
              <h3>禁止事項</h3>
              <p>以下の事項は特別に禁止します。</p>
              <ul>
                <li>他のチームへの妨害と主催者がみなす全ての行為</li>
              </ul>
            </div>
            <div class="box">
              <h3>その他</h3>
              <h4>キャンペーン機能</h4>
              <p>
                <code>POST /initialize</code>
                のレスポンスにて、ｲｽｺｲﾝ還元キャンペーンの「還元率の設定」を返すことができます。この還元率によりユーザが増減します。
              </p>
              <p><code>POST /initialize</code> のレスポンスは JSON 形式で</p>
              <pre><code>{
  "campaign": 0,
  "language": "実装言語"
}</code></pre>
              <p>
                としてください。campaignが還元率の設定となります。有効な値は 0
                以上 4 以下の整数で 0 の場合はキャンペーン機能が無効になります。
              </p>
              <p>
                なお、ｲｽｺｲﾝ還元の費用が下で説明するスコアから引かれることはありません。
              </p>
              <p>
                また、languageの値として、本競技で利用した言語を出力してください(ベンチマーカーの仕様によるもので特に意味はありません。なお、参考実装では最初から入っているので気にする必要はありません)。languageが空の場合はベンチマークは失敗と見なされます。
              </p>
              <h4>参照実装の切り替え方</h4>
              <p>初期状態ではGoによる実装が起動している状態になります。</p>
              <p>
                各言語実装は <code>systemd</code> で管理されています。
                例えば、参照実装をGoからPerlに切り替えるには次のようにします。
              </p>
              <pre>
$ sudo systemctl stop    isucari.golang.service
$ sudo systemctl disable isucari.golang.service
$ sudo systemctl start   isucari.perl.service
$ sudo systemctl enable isucari.perl.service</pre
              >
              <p>
                ただし、PHPを使う場合のみ、
                <code>systemd</code> の設定変更の他に、次のように nginx
                の設定ファイルの変更が必要です。
              </p>
              <pre>
$ sudo unlink /etc/nginx/sites-enabled/isucari.conf
$ sudo ln -s /etc/nginx/sites-available/isucari.php.conf /etc/nginx/sites-enabled/isucari.conf
$ sudo systemctl restart nginx.service</pre
              >
              <p>
                なお、参考実装としてGo, Perl, PHP, Ruby, Python,
                Node.jsによるアプリケーションが用意されています。
              </p>
              <h4>DBのリカバリ方法</h4>
              <p>
                DB(isucari)を初期状態にもどすには、次のコマンドを実行します。
              </p>
              <pre>$ /home/isucon/isucari/webapp/sql/init.sql</pre>
              <hr />
              <p>この他に、なにかあったらhosshiiまで教えてください。</p>
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
  </div>
</template>

<script lang="ts">
import store from '@/store'
export default {
  setup() {
    const me = store.state.me
    console.log(me ? 'true' : 'false0')
    return {
      me
    }
  }
}
</script>

<style lang="sass" scoped>
.box
  margin-bottom: 48px

a
  color: #4ae287

.well
  margin-bottom: 48px
  padding: 1rem 1.5rem

pre
  padding: 1em
  margin-left: 20px
  background-color: #f6f8fa
  white-space: -moz-pre-wrap; /* Mozilla */
  white-space: -pre-wrap; /* Opera 4-6 */
  white-space: -o-pre-wrap; /* Opera 7 */
  white-space: pre-wrap; /* CSS3 */
  word-wrap: break-word; /* IE 5.5+ */

code
  font-family: monospace
  font-weight: normal
  line-height: 150%
  font-size: 110%
  text-align: left
  margin-bottom: 10px
</style>

