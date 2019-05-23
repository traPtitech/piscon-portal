<template>
<div>
  <vuestic-widget v-for="(q, i) in questions" :key="i" :headerText="'Q.'+(i+1)">
    <div class="question">
      <div class="question-header">
        質問
      </div>
      <div class="question-body">{{q.q}}</div>
    </div>
    <div class="answer">
      <div class="answer-header">
        回答
      </div>
      <div class="answer-body">{{q.a}}</div>
    </div>
  </vuestic-widget>

  <vuestic-widget headerText="質問する" v-if="$store.state.Me.name !== '-'">
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
      newQ: ''
    }
  },
  created () {
    axios.get('/api/questions')
      .then(data => {
        this.questions = data.data
        console.log(this.questions)
      })
  },
  methods: {
    newQuestion () {
      axios.post('/api/new', {q: this.newQ})
        .then(() => {
          this.newQ = ''
        })
    }
  }
}
</script>

<style>

.question-body, .answer-body {
  word-wrap: break-word;
  white-space: pre-line;
}

.question, .answer{
  margin-bottom: 12px;
}

.question-header, .answer-header {
  font-size: 1.4em;
}

</style>

