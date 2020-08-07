import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import "@/assets/icon/iconfont.css"
import 'element-ui/lib/theme-chalk/index.css';
import Vuex from 'vuex'
import $ from 'jquery'
import ElementUI from 'element-ui'

Vue.config.productionTip = false
Vue.use(ElementUI)
Vue.use(Vuex)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
