<template>
  <div class="row">
    <div class="col-md-12">
      <vuestic-widget class="no-padding no-v-padding">
        <div class='vuestic-breadcrumbs'>
          <div class="que-list crumb">
            現在のキュー
          </div>
          <div  class="crumb que" :class="{benchNow: que.state==='benchmark'}" v-for="que in $store.state.Que">
            {{que.team.name}}
          </div>
        </div>
      </vuestic-widget>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'vuestic-breadcrumbs',
    props: {
      breadcrumbs: {
        type: Array,
        default: function () {
          return []
        }
      }
    }
  }
</script>

<style lang='scss' scoped>
  .vuestic-breadcrumbs {
    height: $breadcrumbs-height;
    display: flex;
    align-items: center;

    .crumb {
      &.disabled {
        pointer-events: none;
      }
      &:last-child::after {
        display: none;
      }

      &::after {
        padding: 0 5px;
        display: inline-block;
        content: $breadcrumbs-arrow-content;
        vertical-align: middle;
        color: $brand-primary;
        font-size: $breadcrumbs-arrow-font;
        font-family: FontAwesome;
      }
    }

    a {
      color: $breadcrumbs-gray;
      &:hover {
        color: $brand-primary;
      }
    }
  }

  .que {
    color: white;
    background: $brand-success;
    padding: 2px 4px;
    border-radius: 3px;
  }

  .benchNow{
    background: $brand-danger;
  }
</style>
