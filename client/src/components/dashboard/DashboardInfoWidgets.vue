<template>
  <div class="row dashboard-info-widgets">
    <div class="col-md-6 col-xl-3">
      <vuestic-widget class="info-widget" headerText="ベンチ回数">
        <div class="info-widget-inner">
          <div class="stats">
            <div class="stats-number">
              <i class="ion ion-md-arrow-up text-primary stats-icon"></i>
              {{$store.getters.resultCount}}
            </div>
            <div class="stats-title">回</div>
          </div>
        </div>
      </vuestic-widget>
    </div>
    <div class="col-md-6 col-xl-3">
      <vuestic-widget class="info-widget brand-info" headerText="参加人数">
        <div class="info-widget-inner">
          <div class="stats">
            <div class="stats-number">
              <i class="ion ion-md-people stats-icon icon-wide"></i>
              {{$store.state.AllResults.length}}
            </div>
            <div class="stats-title">Team</div>
          </div>
        </div>
      </vuestic-widget>
    </div>
    <div class="col-md-6 col-xl-3">
      <vuestic-widget class="info-widget" headerText="何したら？な人へ">
        <div class="info-widget-inner">
          <div class="stats">
            <div class="stats-number">
              <a href="http://dsas.blog.klab.org/archives/pixiv-isucon2016-1.html" target="_blank">
                <span style="color:#4ae387;">攻略記事</span>
              </a>
            </div>
            <div class="stats-title">これやるだけでも<br>力付きます</div>
          </div>
        </div>
      </vuestic-widget>
    </div>
    <div class="col-md-6 col-xl-3">
      <vuestic-widget class="info-widget" headerText="困ったら">
        <div class="info-widget-inner">
          <div class="stats">
            <div class="stats-number">
              <router-link
                 :to="{name: 'q-and-a'}"
              >
              <span style="color:#4ae387;">Q&A</span>
              </router-link>
            </div>
            <div class="stats-title">バシバシ質問！</div>
          </div>
        </div>
      </vuestic-widget>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'dashboard-info-widgets',

  }
</script>

<style lang="scss" scoped>
  .stats-number, .stats-title {
    line-height: 1;
  }

  .info-widget-inner {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    position: relative;
    width: 100%;

    &.has-chart {
      justify-content: space-between;
    }

    .stats {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      height: 100%;
    }
  }

  .stats-number {
    position: relative;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    font-size: 2.625rem;
    margin-bottom: 0.5rem;

    .stats-icon {
      font-size: 1.5625rem;
      position: absolute;
      top: 0.625rem;
      left: -1.25rem;

      &.icon-wide {
        left: -1.875rem;
      }
    }
  }

</style>
