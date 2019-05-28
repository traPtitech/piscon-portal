<template>
<div>
  <div class="row">
    <div class="col-md-12" v-if="$store.state.Me">
      <vuestic-widget headerText="環境構築マニュアル">
        <div class="well">
          今回は環境構築も競技の一部なのでわざと雑にかいてあります！<br>
          分からないことばかりだとは思いますが頑張りましょう！
        </div>
        <div class="box">
        </div>
        <div class="box">
          <h3>1. インスタンスを作成する</h3>
          <p>
            下のリンクを開き、インスタンスを作成するボタンを押してください。<br>
            <router-link to="team-info">インスタンス作成</router-link><br>
            IPアドレスやパスワードが表示されます。<br>
            しばらく待ってリロードしても作成されない場合は、@nagatechまで連絡してください。
          </p>
        </div>
        <div class="box">
          <h3>2. 起動したインスタンスに isucon ユーザで SSH ログインする</h3>
          <p>
            例:
          </p>
          <div class="well">$ ssh isucon@xx.xx.xx.xx</div>
          <p>
            xx.xx.xx.xxに入るのは表示されているIPアドレスで、<br>
            パスワードは上のページで表示されている初期パスワードです。<br>
            分かる人はこの段階でパスワードを変更しても構いません。(忘れた場合の保証は致しかねます。)
          </p>
        </div>
        <div class="box">
          <h3>3. ソースコードのダウンロード</h3>
          <p>
            今回は温かみのある運営を目指しているので、皆さんのインスタンスにはすでにソースコードをダウンロードしておきました！<br>
            /home/isucon/webapp/以下にソースコードがあるはずなので確認してください。<br>
            見当たらない場合は@nagatechまで連絡してください。
          </p>
        </div>
        <div class="box">
          <h3>4. アプリケーションに必要なソフトのインストール</h3>
          <p>
            <ul>
              <li>各言語の実行環境</li>
              <li>Mysql(MariaDB)</li>
              <li>Memcached</li>
            </ul>
            上記の3つがアプリケーションの実行に最低限必要なソフトです。<br>
            Nginx等その他に必要なものがあれば、適宜インストールしてください。<br>
            ここが一番つらいところですが頑張りましょう。
          </p>
        </div>
        <div class="box">
          <h3>5. データのインポート</h3>
          <p>
            /home/isucon/webapp/dump.sql.bz2 にMySQLの初期データに必要なデータを置いておきました。
          </p>
          <div class="well"> $ dump.sql.bz2 | mysql -uroot</div>
          <p>
            上のコマンドを実行し、データのインポートを行います。<br>
            数分かかります。気長に待ちましょう。<br>

            (ユーザー名・パスワードを変更した場合は適宜変更してください。)
          </p>
        </div>
        <div class="box">
          <h3>6. アプリケーションの実行</h3>
          <p>
            アプリケーションを実行してみてください。<br>
            環境変数等の設定はコードを見て行ってください。
          </p>
        </div>
        <div class="box">
          <h3>7. アプリケーションの動作を確認</h3>
          <p>
            Team Infoに表示されているIPアドレスにブラウザでアクセスし、動作を確認してください。以下の画面が表示されるはずです。
          </p>
          <img src="../../assets/manual.png" style="width: 80%; border: 1px solid;"/>
          <p>
            例として、「アカウント名」は mary、 「パスワード」は marymary を入力することでログインが行えます。
          </p>
        </div>
        <div class="box">
          <h3>8. 負荷走行を実行</h3>
          <p>
            ポータルのTeamInfoでベンチマークを行うボタンを押して、ベンチマークを開始してください。<br>
            キューはページ上部に表示されています。<br>
            この操作後、ポータルにてスコアが反映されているか確認して下さい。
          </p>
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
</div>
</template>

<script>
export default {
  name: 'manual',
}
</script>

<style lang="sass" scoped>
.box 
  margin-bottom: 48px;

a
  color: #4ae287;

</style>

