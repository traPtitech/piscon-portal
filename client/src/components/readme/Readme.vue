<template>
<div>
  <vuestic-widget class="col-md-12" headerText="はじめに">
    <div class="box">
      <h3>Q. Pisconとは？</h3>
      <p>A. <a href="https://trap.jp" target="_blank">東京工業大学デジタル創作同好会traP</a>の部内ISUCONのことです。</p>
    </div>
    <div class="box">
      <h3>Q. ISUCONとは？</h3>
      <p>Iikanjini SpeedUp CONtestの略です。
      <br>サーバーをチューニングして高いスコアを目指します。
      <br>毎年秋くらいにコンテストを開催しています。
      <br><a href="http://isucon.net/" target="_blank">ISUCON公式サイト</a>
      </p>
    </div>
    <div class="box">
      <h3>Q. どうやって始めるの？</h3>
      <p>(この回答はtraP内部向けです)</p>
      <p><router-link :to="{name: 'team-info'}">Team Info</router-link>にアクセスしてインスタンスを作成を押すとサーバーが作成されます。
      <br>今回は<b>環境構築も競技の一部になっている</b>ので作成されるインスタンスはすっからかんです。
      <br>まずは<router-link to="manual">環境構築マニュアル</router-link>を参考に環境構築をして正の点数を取るところから始めてみましょう！</p>
    </div>
    <div class="box">
      <h3>Q. 〇〇が分からない・〇〇について聞きたい</h3>
      <p><router-link :to="{name: 'q-and-a'}">Q & A</router-link>でじゃんじゃん質問しましょう！
      <br>運営が頑張って回答します。
      </p>
    </div>
    <div class="box">
      <h3>運営からのお願い</h3>
      <div>
        <li>悪いことはしないでね</li>
        <li>バグとかあったらこっそり教えてください</li>
        <li>どんどん質問していこう！</li>
      </div>
    </div>
  </vuestic-widget>
</div>
</template>

<script>
export default {
  name: 'readme',
}
</script>

<style lang="sass" scoped>
.box 
  margin-bottom: 48px;

a
  color: #4ae287;

</style>

