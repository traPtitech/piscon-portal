<template>
<div>
  <vuestic-widget class="col-md-12" headerText="はじめに">
    <div class="box">
      <h3>Q. Pisconとは？</h3>
      <p>A. traP内のISUCONのことです</p>
    </div>
    <div class="box">
      <h3>Q. ISUCONとは？</h3>
      <p>iikanjini speedup contestの略です。<br/>サーバーをチューニングして高いスコアを目指します。</p>
    </div>
    <div class="box">
      <h3>Q. どうやって始めるの？</h3>
      <p><router-link style="color:#4ae387;" :to="{name: 'team-info'}">Team Info</router-link>にアクセスしてインスタンスを作成を押すとサーバーが作成されます。
      <br>今回は環境構築も競技の一部になっているので作成されるインスタンスはすっからかんです。
      <br><a style="color:#4ae387;" href="https://git.trapti.tech/to-hutohu/isucon">このリポジトリ</a>のREADME.mdを参考に環境構築をして正の点数を取るところから始めてみましょう！</p>
    </div>
    <div class="box">
      <h3>運営からのお願い</h3>
      <div>
        <li>悪いことはしないでね</li>
        <li>バグとかあったら教えてください</li>
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

<style lang="sass">
.box 
  margin-bottom: 48px;

</style>

