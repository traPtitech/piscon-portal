<template>
<div>
  <div class="row">
    <div class="col-md-12" v-if="$store.state.Me">
      <vuestic-widget headerText="環境構築マニュアル">
        <div class="well">
          わからないことはどんどん<router-link to="qa">質問</router-link>しましょう！
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
          <h4>パッケージのアップデート</h4>
          <p>新しくサーバーを立てたときには、きちんと利用しているパッケージのアップデートをしましょう。 古いパッケージを利用していると、パフォーマンスが悪かったり脆弱性があったりするので良くないです。</p>
          <div class="well">
            $ sudo apt update<br>
            $ sudo apt upgrade
          </div>
          <h4>利用する言語</h4>
          <p>サーバーアプリケーションを実行するための、実行環境・コンパイラを導入する必要があります。</p>
          <p>自分が好きな言語・得意な言語が下のリストにある人はそれを使ってもらって構いません。 そんなのないよ！って人はGolangを使うことをおすすめします。</p>
          <h4>利用できる言語のリスト</h4>
          <ul>
            <li>Go</li>
            <li>Node.js(JavaScript)</li>
            <li>PHP</li>
            <li>Python</li>
            <li>Ruby</li>
            <li>Rust</li>
            <li>Scala</li>
          </ul>
          <p>ここではGolangとして進めていきます。</p>
          <p>ubuntuのパッケージ管理ツールである<code>apt</code>を利用してインストールします。</p>
          <div class="well">
            $ sudo apt install -y golang-go
          </div>
          <h4>データベース</h4>
          <p>データベースを入れます。 MySQLとMariaDBがありますが、MariaDBのほうがISUCON的には早くて良いです。互換性があるのでどちらでもMySQLの構文は動きます。</p>
          <div class="well">$ sudo apt install -y mariadb-server</div>
          <p><a href="https://mariadb.org/">mariadb 公式ページ (英語)</a> <a href="https://mariadb.com/kb/ja/mariadb/">mariadb ドキュメント
              (日本語訳)</a></p>
          <h4>nginx</h4>
          <div class="well">$ sudo apt install -y nginx</div>
          <p>設定ファイルを書く必要があります(後述)。</p>
          <h4>memcached</h4>
          <p>今回のアプリケーションではセッションの管理にmemcachedを使っているのでインストールする必要があります。</p>
          <div class="well">$ sudo apt install -y memcached</div>
          <p><a href="https://memcached.org/">memcached 公式ページ (英語)</a></p>
          <h4>データベースの権限設定</h4>
          <div class="well">$ sudo mysql -uroot</div>
          <p>以下mysqlのコンソール</p>
          <div class="well">
            &gt; use mysql;<br>
            &gt; update user set plugin=&#39;&#39; where user=&#39;root&#39;;<br>
            &gt; flush privileges;<br>
            &gt; exit
          </div>
          <p><code>$ mysql -uroot</code> は、mysql に ユーザー root でアクセスする、という意味です。 <code>$ mysql -u root</code> や
            <code>$ mysql --user=root</code> も同じ意味になります。</p>
        </div>
        <div class="box">
          <h3>5. データのインポート</h3>
          <p>
            /home/isucon/dump.sql.bz2 にMySQLの初期データに必要なデータを置いておきました。
          </p>
          <div class="well"> $ bzcat dump.sql.bz2 | mysql -uroot</div>
          <p>
            上のコマンドを実行し、データのインポートを行います。<br>
            数分かかります。気長に待ちましょう。<br>

            (ユーザー名・パスワードを変更した場合は適宜変更してください。)
          </p>
        </div>
        <div class="box">
          <h3>6. アプリケーションのライブラリの構築</h3>
          <ul>
            <li>Nodeだったらnpm</li>
            <li>Pythonだったらpip</li>
            <li>Rubyだったらgem</li>
          </ul>
          <p>以下Golangの例です。</p>
          <p>今回は問題ソースコードにセットアップ用スクリプト(<code>setup.sh</code>)にライブラリ構築用のコマンドが全部あるのでそれを実行します</p>
          <div class="well">
            $ cd /home/isucon/webapp/golang<br>
            $ sh setup.sh
          </div>
        </div>
        <div class="box">
          <h3>7. アプリケーションの実行</h3>
          <p>
            アプリケーションを実行してみてください。<br>
            環境変数等の設定はコードを見て行ってください。
          </p>
          <p><code>go run</code>でもいけますがGolangは実行ファイルを生成してくれるのでそちらで実行してみましょう</p>
          <div class="well">
            $ ./app
          </div>
        </div>
        <div class="box">
          <h3>8. アプリケーションの動作を確認</h3>
          <p>
            Team Infoに表示されているIPアドレスにブラウザでアクセスし、動作を確認してください。<br>
            デフォルトだと8000番で動くらしいので、http://xx.xx.xx.xx:8000(xx xxはIPアドレス)にブラウザでアクセスしてみましょう。
          </p>
          <img src="../../assets/manual.png" style="width: 80%; border: 1px solid;"/>
          <p>
            例として、「アカウント名」は mary、 「パスワード」は marymary を入力することでログインが行えます。
          </p>
        </div>
        <div class="box">
          <h3>9. Nginxのプロキシ設定</h3>
          <p>ベンチマークで正の点数を取るためには、普段ブラウザでアクセスした時に開かれる80番(HTTP)のポートでアプリが閲覧できるようにしなくてはいけません。
            Nginxを使うことで80番にきたリクエストを8000番に渡す（プロキシ）ことができます。 <code>/etc/nginx/conf.d/isucon.conf</code>に設定ファイルを書きましょう。</p>
          <div class="well">
            <pre>
server{
  listen 80 default_server;
  server_name _;

  location / {
    proxy_pass http://localhost:8000;
    proxy_set_header Host $host;
  }
}</pre>
          </div>
          <p>また<code>/etc/nginx/sites-enabled/default</code>を削除します。</p>
          <div class="well">
            $ sudo rm /etc/nginx/sites-enabled/default
          </div>
          <p>Nginxの設定を書き換えた後は必ず文法のチェックを行いましょう。</p>
          <div class="well">
            $ sudo nginx -t
          </div>
          <p>書けたらNginxの設定を再読込します。</p>
          <div class="well">
            $ sudo systemctl reload nginx
          </div>
          <p><code>./app</code>を起動したままhttp://xx.xx.xx.xxにアクセスしてみましょう。 8000番を指定しなくても表示されると思います。</p>
          <ul>
            <li><a href="https://nginx.org/en/docs/">nginx ドキュメント (英語)</a></li>
            <li><a href="http://mogile.web.fc2.com/nginx/index.html">nginx ドキュメント (日本語訳)</a>
              <ul>
                <li><a href="http://mogile.web.fc2.com/nginx/beginners_guide.html">ビギナーガイド (設定ファイルの構成・書き方)</a></li>
                <li><a href="http://mogile.web.fc2.com/nginx/http/ngx_http_proxy_module.html#proxy_pass">proxy_pass</a></li>
              </ul>
            </li>
          </ul>
        </div>
        <div class="box">
          <h3>10. 負荷走行を実行</h3>
          <p>
            ポータルのTeamInfoでベンチマークを行うボタンを押して、ベンチマークを開始してください。<br>
            キューはページ上部に表示されています。<br>
            この操作後、ポータルにてスコアが反映されているか確認して下さい。
          </p>
        </div>
      </vuestic-widget>
      <vuestic-widget headerText="レギュレーション">
        <div class="box">
          <p>基本スコアは以下のルールで算出されます。</p>
          <div class="well">
            成功レスポンス数(GET) x 1 + 成功レスポンス数(POST) x 2 + 成功レスポンス数(画像投稿) x 5 - (サーバエラー(error)レスポンス数 x 10 + リクエスト失敗(exception)数 x 20 + 遅延POSTレスポンス数 x 100)
          </div>
          <p>ただし、基本スコアと計測ツールの出すスコアが異なっている場合は、計測ツールの出すスコアが優先されます。</p>
          <h3>減点対象</h3>
          <p>以下の事項に抵触すると減点対象となります。</p>
          <ul>
            <li>存在するべきファイルへのアクセスが失敗する</li>
            <li>リクエスト失敗（通信エラー等）が発生する</li>
            <li>サーバエラー(Status 5xx)・クライアントエラー(Status 4xx)をアプリケーションが返す</li>
            <li>他、計測ツールのチェッカが検出したケース</li>
          </ul>
        </div>
        <div class="box">
          <h3>注意事項</h3>
          <ul>
            <li>リダイレクトはリダイレクト先が正しいレスポンスを返せた場合に、1回レスポンスが成功したと判断します</li>
            <li>POSTの失敗は大幅な減点対象です</li>
          </ul>
        </div>
        <div class="box">
          <h3>制約事項</h3>
          <p>以下の事項に抵触すると点数が無効となります。</p>
          <ul>
            <li>GET /initialize へのレスポンスが10秒以内に終わらない</li>
            <li>存在するべきDOM要素がレスポンスHTMLに存在しない</li>
          </ul>
        </div>
        <div class="box">
          <h3>禁止事項</h3>
          <p>以下の事項は特別に禁止する。</p>
          <ul>
            <li>他のチームへの妨害と主催者がみなす全ての行為</li>
          </ul>
        </div>
        <div class="box">
          <h3>ソフトウェア事項</h3>
          <p>コンテストにあたり、参加者は与えられたソフトウェア、もしくは自分で競技時間内に実装したソフトウェアを用いる。</p>
          <p>高速化対象のソフトウェアとして主催者からRuby, PHPによるWebアプリケーションが与えられる。ただし各々の性能が一致することを主催者は保証しない。どれをベースに用いてもよいし、独自で実装したものを用いてもよい。</p>
          <p>競技における高速化対象のアプリケーションとして与えられたアプリケーションから、以下の機能は変更しないこと。</p>
          <ul>
            <li>アクセス先のURI（ポート、およびHTTPリクエストパス）</li>
            <li>レスポンス(HTML)のDOM構造</li>
            <li>JavaScript/CSSファイルの内容</li>
            <li>画像および動画等のメディアファイルの内容</li>
          </ul>
          <p>各サーバにおけるソフトウェアの入れ替え、設定の変更、アプリケーションコードの変更および入れ替えなどは一切禁止しない。起動したインスタンス以外の外部リソースを利用する行為 (他のインスタンスに処理を委譲するなど) は禁止する。</p>
          <p>許可される事項には、例として以下のような作業が含まれる。</p>
          <ul>
            <li>DBスキーマの変更やインデックスの作成・削除</li>
            <li>キャッシュ機構の追加、jobqueue機構の追加による遅延書き込み</li>
            <li>他の言語による再実装</li>
          </ul>
          <p>ただし以下の事項に留意すること。</p>
          <ul>
            <li>コンテスト進行用のメンテナンスコマンドが正常に動作するよう互換性を保つこと</li>
            <li>各サーバの設定およびデータ構造は任意のタイミングでのサーバ再起動に耐えること</li>
            <li>サーバ再起動後にすべてのアプリケーションコードが正常動作する状態を維持すること</li>
            <li>ベンチマーク実行時にアプリケーションに書き込まれたデータは再起動後にも取得できること</li>
          </ul>
        </div>
        <div class="box">
          <h3>採点</h3>
          <p>採点は採点条件（後述）をクリアした参加者の間で、性能値（後述）の高さを競うものとする。</p>
          <p>採点条件として、以下の各チェックの検査を通過するものとする。</p>
          <ul>
            <li>負荷走行中、POSTしたデータが、POSTへのHTTPレスポンスを返してから即座に関連するURI GETのレスポンスデータに反映されていること</li>
            <li>レスポンスHTMLのDOM構造が変化していないこと</li>
            <li>ブラウザから対象アプリケーションにアクセスした結果、ページ上の表示および各種動作が正常であること</li>
          </ul>
          <p>性能値として、以下の指標を用いる。</p>
          <ul>
            <li>
              計測ツールの実行時間は1分間とする
              <ul>
                <li>細かい閾値ならびに配点についての詳細は当日のマニュアルに記載する</li>
              </ul>
            </li>
            <li>
              計測時間内のHTTPリクエスト成功数をベースとする
              <ul>
                <li>リクエストの種類毎に配点を変更する</li>
                <li>エラーの数により減点する</li>
              </ul>
            </li>
          </ul>
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

.well
  margin-bottom: 48px;
  padding: 1rem 1.5rem;

</style>

