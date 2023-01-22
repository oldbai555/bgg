import Vue from 'vue'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify'
import day from 'dayjs'
import {formatDate} from './plugins/time'


import './plugins/http'

Vue.filter('dateformat', function(indate, outdate) {
  return day(indate).format(outdate)
})

Vue.filter('formatDate', function(timestamp) {
  return formatDate(timestamp)
})

Vue.config.productionTip = false

new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
