import axios from 'axios'

// axios.defaults.withCredentials = true
axios.defaults.baseURL = process.env.NODE_ENV === 'development' ? 'http://localhost:4000' : 'https://portal.emoine.tech'

export default axios
