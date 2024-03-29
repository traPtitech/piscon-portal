<template>
  <va-card :title="$t('tables.infiniteScroll')">
    <div class="data-table-infinite-scroll--container" ref="scrollable" @scroll="onScroll">
      <va-data-table
        :fields="fields"
        :data="users"
        api-mode
        no-pagination
      >
        <template #marker="props">
          <va-icon name="fa fa-circle" :color="props.rowData.color" size="8px" />
        </template>
      </va-data-table>

      <div class="flex-center ma-3">
        <spring-spinner
          v-if="loading"
          :animation-duration="2000"
          :size="60"
          :color="theme.primary"
        />
      </div>
    </div>
  </va-card>
</template>

<script>
import { SpringSpinner } from 'epic-spinners'
import users from '../data/users.json'
import { useGlobalConfig } from 'vuestic-ui'

export default {
  components: {
    SpringSpinner,
  },
  data () {
    return {
      users: [],
      loading: false,
      offset: 0,
    }
  },
  computed: {
    theme() {
      return useGlobalConfig().getGlobalConfig().colors;
    },
    fields () {
      return [{
        name: '__slot:marker',
        width: '30px',
        dataClass: 'text-center',
      }, {
        name: 'fullName',
        title: this.$t('tables.headings.name'),
      }, {
        name: 'email',
        title: this.$t('tables.headings.email'),
      }, {
        name: 'country',
        title: this.$t('tables.headings.country'),
      }]
    },
  },
  created () {
    this.loadMore()
  },
  methods: {
    async loadMore () {
      this.loading = true

      const users = await this.readUsers()
      this.users = this.users.concat(users)
      this.loading = false
    },
    async readUsers () {
      await new Promise(resolve => setTimeout(resolve, 600))
      return users.slice(0, 10)
    },
    onScroll (e) {
      if (this.loading) {
        return
      }

      const { target } = e

      if (target.offsetHeight + target.scrollTop === target.scrollHeight) {
        this.loadMore()
      }
    },
  },
}
</script>

<style lang="scss">
  .data-table-infinite-scroll--container {
    height: 300px;
    overflow-y: auto;
  }
</style>
