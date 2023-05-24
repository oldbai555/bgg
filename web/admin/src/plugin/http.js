import Vue from 'vue'
import axios from 'axios'

let Url = 'http://localhost:8003'

axios.defaults.baseURL = Url

axios.interceptors.request.use(config => {
  config.headers.Authorization = `${window.sessionStorage.getItem('sid')}`
  return config
})

Vue.prototype.$http = axios

export { Url }
