// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import 'vuetify/dist/vuetify.min.css';

import Vue from 'vue';
import Vuetify from 'vuetify';
import App from './App';
import store from './store';
import router from './router';
import HttpUtil from './utils/http';

Vue.use(Vuetify);

Vue.config.productionTip = false;
Vue.prototype.$http = HttpUtil;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  store,
  router,
  components: { App },
  template: '<App/>',
});
