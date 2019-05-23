<template>
  <div class="profile-dropdown col nav-item dropdown navbar-dropdown">
    <span class="profile-section-avatar-container">
      <slot></slot>
    </span>
  </div>
</template>

<script>
  export default {
    name: 'profile-section',

    props: {
      options: {
        type: Array,
        default: () => [
          {
            name: 'profile',
            redirectTo: '',
          },
          {
            name: 'logout',
            redirectTo: 'login'
          }
        ]
      }
    },
  }
</script>

<style lang="scss">
  .profile-dropdown {

    .profile-section-avatar-container {
      display: inline-block;
      width: 50px;
      height: 50px;
      background-color: white;
      border-radius: 50%;
      border: 2px solid $lighter-gray;
      overflow: hidden;
      position: absolute;
      right: 0;

      img {
        height: 100%;
        width: 100%;
      }
    }

  }

</style>
