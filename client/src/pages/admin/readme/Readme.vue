<template>
  <div>
    <va-card class="col-md-12" headerText="はじめに">
      <va-card-title>
        はじめに
      </va-card-title>
      <va-card-content>
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
          <h3>運営からのお願い</h3>
          <div>
            <li>悪いことはしないでね</li>
            <li>バグとかあったらこっそり教えてください</li>
          </div>
        </div>
      </va-card-content>
    </va-card>
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

