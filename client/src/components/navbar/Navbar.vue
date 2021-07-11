<template>
  <div class="app-layout__navbar">
    <va-navbar color="white">
      <template v-slot:left>
        <div class="left">
          <va-icon-menu-collapsed
            @click="isSidebarMinimized = !isSidebarMinimized"
            :class="{ 'x-flip': isSidebarMinimized }"
            class="va-navbar__item"
            :color="colors.primary"
          />
          <router-link to="/dashboard">
            <piscon-logo class="logo" />
          </router-link>
        </div>
      </template>
      <template #right>
        <div v-if="user">
          {{ user.name }}
        </div>
        <div v-else>
          <router-link to="/auth/login"> Signin with traQ </router-link>
        </div>
      </template>
    </va-navbar>
  </div>
</template>

<script>
import { useColors } from 'vuestic-ui'
import { computed } from 'vue'
import PisconLogo from '@/components/piscon-logo'
import VaIconMenuCollapsed from '@/components/icons/VaIconMenuCollapsed'
import store from '@/store'

export default {
  components: { PisconLogo, VaIconMenuCollapsed },
  setup() {
    const { getColors } = useColors()
    const colors = computed(() => getColors())
    const user = computed(() => store.state.User)

    const isSidebarMinimized = computed({
      get: () => store.state.isSidebarMinimized,
      set: value => store.commit.updateSidebarCollapsedState(value)
    })

    const userName = computed(() => store.state.userName)
    return {
      colors,
      isSidebarMinimized,
      userName,
      user
    }
  }
}
</script>

<style lang="scss" scoped>
.va-navbar {
  box-shadow: var(--va-box-shadow);
  z-index: 2;
  &__center {
    @media screen and (max-width: 1200px) {
      .app-navbar__github-button {
        display: none;
      }
    }
    @media screen and (max-width: 950px) {
      .app-navbar__text {
        display: none;
      }
    }
  }

  @media screen and (max-width: 950px) {
    .left {
      width: 100%;
    }
    .app-navbar__actions {
      width: 100%;
      display: flex;
      justify-content: space-between;
    }
  }
}

.left {
  display: flex;
  align-items: center;
  & > * {
    margin-right: 1.5rem;
  }
  & > *:last-child {
    margin-right: 0;
  }
}

.x-flip {
  transform: scaleX(-100%);
}

.app-navbar__text > * {
  margin-right: 0.5rem;
  &:last-child {
    margin-right: 0;
  }
}
</style>