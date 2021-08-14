<template>
  <div>
    <div class="flex md12"></div>
    <div class="form-group">
      <div class="input-group">
        <va-input
          class="mb-4"
          v-model="betterize"
          type="textarea"
          placeholder="改善点を入力してください(記入しないとベンチマークを行えません)"
        />
      </div>
    </div>
    <div class="flex md12 my-2" v-for="i in team.instance.length" :key="i">
      <va-button
        :rounded="false"
        class="mr-4 item"
        @click="benchmark(i)"
        :disabled="benchmarkButton(i) || betterize === ''"
      >
        サーバ{{ i }}にベンチマークを行う
      </va-button>
      <va-button
        :rounded="false"
        class="mr-4 item"
        :color="instanceButtonColor(i)"
        @click="setOperationModal(i)"
        :disabled="instanceButton(i) || waiting"
      >
        {{ instanceButtonMessage(i) }}
      </va-button>
    </div>
    <div v-if="error" class="type-articles">
      {{ error }}
    </div>
    <va-modal hide-default-actions v-model="showOperationModal">
      <template #header>
        <h3>確認</h3>
      </template>
      <slot>
        <div>
          この操作は取り消せません。間違えて行わないように注意してください
        </div>
      </slot>
      <template #footer>
        <va-button @click="operationInstance(operationInstanceNumber)">
          実行
        </va-button>
        <va-button @click="showOperationModal = false"> キャンセル </va-button>
      </template>
    </va-modal>
    <va-modal v-model="showInfoModal">
      <va-content>
        <p>{{ infoModalMessage.better }}</p>
        <p>{{ infoModalMessage.message }}</p>
      </va-content>
    </va-modal>
  </div>
</template>
<script lang="ts">
import { defineComponent, ref, computed, PropType } from 'vue'
import { AxiosError } from 'axios'
import apis, {
  Instance,
  PostBenchmarkRequest,
  Response,
  Result
} from '../../../lib/apis'
import store from '../../../store'
export default defineComponent({
  name: 'Benchmark',
  props: {
    sortedInstance: {
      type: Array as PropType<Array<Instance>>,
      required: true
    },
    teamResults: { type: Array as PropType<Array<Result>>, required: true },
    showInfoModal: { type: Boolean, required: true },
    InfomodalMessage: { type: Object, required: true }
  },
  setup(props) {
    const betterize = ref<string>('')
    const error = ref<string>('')
    const showOperationModal = ref(false)
    const operationInstanceNumber = ref(0)
    const waiting = ref(false)
    const team = computed(() => store.state.Team)
    const benchmarkButton = (i: number) =>
      computed(
        () =>
          store.state.Queue?.find(
            que => que.team.name === store.state.Team?.name
          ) ||
          (!props.sortedInstance
            ? false
            : props.sortedInstance[i - 1].status !== 'ACTIVE')
      ).value
    const benchmark = (id: number) => {
      if (benchmarkButton(id) || !store.state.Team) {
        return
      }
      const req: PostBenchmarkRequest = {
        betterize: betterize.value
      }
      apis
        .benchmarkNameInstanceNumberPost(store.state.Team?.name, id, req)
        .then(() => {
          betterize.value = ''
          store.dispatch.getData()
        })
        .catch((err: AxiosError<Response>) => {
          error.value = !err.response?.data.message
            ? ''
            : err.response?.data.message
        })
    }
    const instanceButton = (i: number) => {
      return computed(
        () =>
          (!props.sortedInstance || !props.sortedInstance[i - 1]
            ? false
            : props.sortedInstance[i - 1].status !== 'ACTIVE') &&
          (!props.sortedInstance || !props.sortedInstance[i - 1]
            ? false
            : props.sortedInstance[i - 1].status !== 'NOT_EXIST')
      ).value
    }
    const instanceButtonColor = (i: number) =>
      computed(() => {
        if (!props.sortedInstance) {
          return
        }
        switch (props.sortedInstance[i - 1].status) {
          case 'ACTIVE':
            return 'danger'

          case 'NOT_EXIST':
            return 'info'

          default:
            return 'info'
        }
      }).value
    const instanceButtonMessage = (i: number) =>
      computed(() => {
        if (!props.sortedInstance) {
          return
        }
        switch (props.sortedInstance[i - 1].status) {
          case 'ACTIVE':
            return `インスタンス${
              props.sortedInstance[i - 1].instance_number
            }を削除する`

          case 'NOT_EXIST':
            return `インスタンス${
              props.sortedInstance[i - 1].instance_number
            }を作成する`

          case 'BUILDING':
            return `作成中`

          default:
            return props.sortedInstance[i - 1].status
        }
      }).value
    const setOperationModal = (id: number) => {
      showOperationModal.value = true
      operationInstanceNumber.value = id
    }
    const operationInstance = (id: number) => {
      showOperationModal.value = false
      waiting.value = true
      if (
        instanceButton(id) ||
        !props.sortedInstance ||
        !props.sortedInstance[id - 1] ||
        !store.state.Team
      ) {
        return
      }
      switch (props.sortedInstance[id - 1].status) {
        case 'ACTIVE':
          apis
            .instanceTeamIdInstanceNumberDelete(store.state.Team?.ID, id)
            .then(() => {
              store.dispatch.getData()
              waiting.value = false
            })
            .catch((err: AxiosError<Response>) => {
              error.value = !err.response?.data.message
                ? ''
                : err.response?.data.message
              waiting.value = false
            })
          break
        case 'NOT_EXIST':
          apis
            .instanceTeamIdInstanceNumberPost(store.state.Team.ID, id)
            .then(() => {
              store.dispatch.getData()
              waiting.value = false
            })
            .catch((err: AxiosError<Response>) => {
              error.value = !err.response?.data.message
                ? ''
                : err.response?.data.message
              waiting.value = false
            })
          break
      }
    }

    return {
      betterize,
      operationInstanceNumber,
      showOperationModal,
      team,
      benchmarkButton,
      benchmark,
      instanceButton,
      instanceButtonColor,
      instanceButtonMessage,
      setOperationModal,
      operationInstance
    }
  }
})
</script>