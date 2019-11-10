<template>
<div>
  <vuestic-widget v-for="(q, i) in questions" :key="i" :headerText="'Q.'+(i+1)">
    <div class="question">
      <div class="question-header">
        質問
      </div>
      <div class="question-body">
        {{q.question}}
      </div>
    </div>
    <div class="answer">
      <div class="answer-header">
        回答
      </div>
      <div v-if="q.answer" class="answer-body">
        {{q.answer}}
      </div>
      <div v-else>
        まだ回答されていません
      </div>
    </div>
    <div v-if="checkAdmin">
      <div>
        運営用回答欄
      </div>
      <div class="form-group">
      <div class="input-group">
        <textarea type="text" class="answer" name="answer" col="10" v-model="newA[i]"></textarea>
        <i class="bar"></i>
      </div>
    </div>
    <button class="btn btn-primary btn-small" @click="newAnswer(q.ID, i)">回答する</button>
    <button class="btn btn-danger btn-small" @click="deleteQuestion(q.ID, i)">削除する</button>
    </div>
  </vuestic-widget>

  <vuestic-widget headerText="質問する" v-if="$store.state.Me">
    <div>
      <p>バシバシ質問しましょう！</p>
      <p>バグ報告はこっそり@xecuaまでお願いします</p>
    </div>
    <div class="form-group">
      <div class="input-group">
        <textarea type="text" id="new" name="new" col="10" v-model="newQ"></textarea>
        <label class="control-label" for="new">質問文</label>
        <i class="bar"></i>
      </div>
    </div>
    <button class="btn btn-primary btn-small" @click="newQuestion">質問する</button>
  </vuestic-widget>
</div>
</template>

<script>
import axios from '../../services/axios'
export default {
  name: 'qa',
  data () {
    return {
      questions: [],
      newQ: '',
      newA: []
    }
  },
  created () {
    axios.get('/api/questions')
      .then(data => {
        this.questions = data.data
      })
  },
  computed: {
    checkAdmin () {
      const me = this.$store.state.Me
      if (!me) {
        return false
      }
      return me.name === 'nagatech' || me.name === 'to-hutohu' || me.name === 'xecua'
    }
  },
  methods: {
    async newQuestion () {
      const question = this.newQ
      if (!question) {
        return
      }
      await axios.post('/api/questions', {question: question})
        .then(() => {
          this.newQ = ''
        })
      this.questions.push({
        question: question,
        answer: ''
      })
    },
    async newAnswer (id, index) {
      const answer = this.newA[index]
      await axios.put(`/api/questions/${id}`, {answer: answer})
        .then(() => {
          this.newA[index] = ''
          this.questions[index].answer = answer
        })
    },
    async deleteQuestion (id, index) {
      if (!window.confirm('削除しますか？')) {
        return
      }
      await axios.delete(`/api/questions/${id}`)
        .then(() => {
          axios.get('/api/questions')
          .then(data => {
            this.questions = data.data
          })
        })
    }
  }
}
</script>

<style>

.question-body, .answer-body {
  word-wrap: break-word;
}

.question, .answer{
  margin-bottom: 12px;
}

.question-header, .answer-header {
  font-size: 1.4em;
}

</style>

