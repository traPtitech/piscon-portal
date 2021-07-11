<template>
  <div>
    <va-card v-for="(q, i) in questions" :key="i">
      <va-card-title>
        {{ 'Q.' + (i + 1) }}
      </va-card-title>
      <va-card-content>
        <div class="question">
          <div class="question-header">質問</div>
          <div class="question-body">
            {{ q.question }}
          </div>
        </div>
        <div class="answer">
          <div class="answer-header">回答</div>
          <div v-if="q.answer" class="answer-body">
            {{ q.answer }}
          </div>
          <div v-else>まだ回答されていません</div>
        </div>
        <div v-if="checkAdmin">
          <div>運営用回答欄</div>
          <div class="form-group">
            <div class="input-group">
              <textarea
                type="text"
                class="answer"
                name="answer"
                col="10"
                v-model="newA[i]"
              ></textarea>
              <i class="bar"></i>
            </div>
          </div>
          <button class="btn btn-primary btn-small" @click="newAnswer(q.ID, i)">
            回答する
          </button>
          <button
            class="btn btn-danger btn-small"
            @click="deleteQuestion(q.ID, i)"
          >
            削除する
          </button>
        </div>
      </va-card-content>
    </va-card>

    <va-card headerText="質問する" v-if="user">
      <div>
        <p>バシバシ質問しましょう！</p>
        <p>バグ報告はこっそり@hijiki51までお願いします</p>
      </div>
      <div class="form-group">
        <div class="input-group">
          <textarea
            type="text"
            id="new"
            name="new"
            col="10"
            v-model="newQ"
          ></textarea>
          <label class="control-label" for="new">質問文</label>
          <i class="bar"></i>
        </div>
      </div>
      <button class="btn btn-primary btn-small" @click="newQuestion">
        質問する
      </button>
    </va-card>
  </div>
</template>

<script lang="ts">
import apis, { Questions } from '../../..//lib/apis'
import store from '../../../store'
import { ref } from 'vue'
export default {
  name: 'qa',
  setup() {
    const newQ = ref('')
    const getQuestions = async () => {
      const newQuestions = await apis.questionsGet().then(data => data.data)
      return newQuestions
    }
    const questions = ref([] as Questions[])
    getQuestions().then(data => (questions.value = data))
    const newA = ref([] as string[])

    const user = store.state.User
    const checkAdmin = () => {
      if (!user) {
        return false
      }
      return (
        user.name === 'nagatech' ||
        user.name === 'to-hutohu' ||
        user.name === 'xecua' ||
        user.name === 'hosshii' ||
        user.name === 'hijiki51'
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
      await apis.questionsPost(question).then(() => (newQ.value = ''))
      questions.value.push(question)
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
    return {
      questions,
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
