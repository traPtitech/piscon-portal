import axios from 'axios'

// axios.defaults.withCredentials = true
axios.defaults.baseURL = process.env.NODE_ENV === 'development' ? 'http://localhost:8080' : ''

export default axios
