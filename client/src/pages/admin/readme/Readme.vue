<template>
  <va-content class="typography content">
    <div>
      <va-card class="col-md-12">
        <va-card-title> <h3 class="h-fix">はじめに</h3> </va-card-title>
        <va-card-content>
          <div class="mb-4">
            <h4>Q. Pisconとは？</h4>
            <p>
              A.
              <a href="https://trap.jp" target="_blank"
              >東京工業大学デジタル創作同好会traP</a
              >の部内ISUCONのことです。
            </p>
          </div>
          <div class="mb-4">
            <h4>Q. ISUCONとは？</h4>
            <p>
              Iikanjini SpeedUp CONtestの略です。
              <br />サーバーをチューニングして高いスコアを目指します。
              <br />毎年秋くらいにコンテストを開催しています。 <br /><a
                href="http://isucon.net/"
                target="_blank"
              >ISUCON公式サイト</a
              >
            </p>
          </div>
          <div class="mb-4">
            <h4>運営からのお願い</h4>
            <div>
              <li>悪いことはしないでね</li>
              <li>バグとかあったらこっそり教えてください</li>
            </div>
          </div>
        </va-card-content>
      </va-card>
    </div>
  </va-content>
</template>

<script>
export default {
  name: 'readme'
}
</script>

<style lang="sass" scoped>
.box
  margin-bottom: 48px

a
  color: #4ae287

.h-fix
  margin-bottom: auto
  margin-top: auto
</style>

