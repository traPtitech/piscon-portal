import axios from 'axios'

// axios.defaults.withCredentials = true
axios.defaults.baseURL = process.env.NODE_ENV === 'development' ? 'http://localhost:8080' : process.env.VUE_APP_API_ENDPOINT

export default axios
