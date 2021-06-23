import axios from 'axios'

// axios.defaults.withCredentials = true
axios.defaults.baseURL = process.env.NODE_ENV === 'development' ? 'http://localhost:4000' : 'https://piscon-portal.trap.jp/'
// axios.defaults.baseURL = 'https://piscon-portal.trap.jp/'

export default axios
