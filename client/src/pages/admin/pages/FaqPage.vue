<template>
  <va-content>
    <va-card class="flex md12 item mb-3" v-for="(q, i) in questions" :key="i">
      <va-card-title>
        <h5>{{ 'Q.' + (i + 1) }}</h5>
      </va-card-title>
      <va-card-content>
        <h5 class="mb-4">質問</h5>
        <div class="text-block">
          <p>
            {{ q.question }}
          </p>
        </div>
        <h5 class="mb-4">回答</h5>
        <div v-if="q.answer" class="mb-4 text-block">
          {{ q.answer }}
        </div>
        <div v-else class="mb-4 text-block text--secondary">
          まだ回答されていません
        </div>
        <div v-if="checkAdmin">
          <strong>運営用回答欄</strong>
          <va-input
            class="mb-4"
            v-model="newA[i]"
            type="textarea"
            placeholder="回答"
          />
          <va-button :rounded="false" class="mr-4" @click="newAnswer(q.ID, i)">
            回答する
          </va-button>
          <va-button
            :rounded="false"
            class="mr-4"
            @click="deleteQuestion(q.ID, i)"
          >
            削除する
          </va-button>
        </div>
      </va-card-content>
    </va-card>

    <va-card v-if="user">
      <va-card-title><h5>質問する</h5></va-card-title>
      <va-card-content>
        <div class="mb-4">
          <p>バシバシ質問しましょう！</p>
          <p>バグ報告はこっそり@hijiki51までお願いします</p>
        </div>
        <va-input
          class="mb-4"
          v-model="newQ"
          type="textarea"
          placeholder="質問文"
        />
        <va-button :rounded="false" class="mr-4" @click="newQuestion">
          質問する
        </va-button>
      </va-card-content>
    </va-card>
  </va-content>
</template>

<script lang="ts">
import apis, { Questions } from '../../../lib/apis'
import store from '../../../store'
import { computed, ref } from 'vue'
export default {
  name: 'qa',
  setup() {
    const newQ = ref('')
    const questions = ref([] as Questions[])
    const newA = ref([] as string[])
    const getQuestions = async () => {
      const newQuestions = await apis.questionsGet().then(data => data.data)
      questions.value = newQuestions
      newA.value = newQuestions.map(v => (v.answer ? v.answer : ''))
    }
    const user = computed(() => store.state.User)
    const checkAdmin = () => {
      if (!user.value) {
        return false
      }
      return (
        user.value.name === 'nagatech' ||
        user.value.name === 'to-hutohu' ||
        user.value.name === 'xecua' ||
        user.value.name === 'hosshii' ||
        user.value.name === 'hijiki51'
      ) //TODO
    }
    const newQuestion = async () => {
      const questionText = newQ.value
      if (!questionText) {
        return
      }
      const question: Questions = {
        question: questionText
      }
      await apis.questionsPost(question)
      newQ.value = ''
      await getQuestions()
    }
    const newAnswer = async (id: number, index: number) => {
      const answerText = newA.value[index]
      const answer: Questions = {
        answer: answerText
      }
      await apis.questionsIdPut(id, answer).then(() => {
        newA.value[index] = ''
        questions.value[index].answer = answerText
      })
      await getQuestions()
    }
    const deleteQuestion = async (id: number) => {
      if (!window.confirm('削除しますか？')) {
        return
      }
      await apis.questionsIdDelete(id).then(() => {
        apis.questionsGet().then(data => {
          questions.value = data.data
        })
      })
    }

    getQuestions()
    return {
      questions,
      user,
      newQ,
      newA,
      checkAdmin,
      newQuestion,
      newAnswer,
      deleteQuestion
    }
  }
}
</script>

<style>
.question-body,
.answer-body {
  word-wrap: break-word;
}

.question,
.answer {
  margin-bottom: 12px;
}

.question-header,
.answer-header {
  font-size: 1.4em;
}
</style>
