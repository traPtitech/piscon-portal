import axios from 'axios'

// axios.defaults.withCredentials = true
axios.defaults.baseURL = process.env.NODE_ENV === 'development' ? 'http://localhost:4000' : 'http://160.251.13.26/'
// axios.defaults.baseURL = 'http://118.27.18.240'

export default axios
