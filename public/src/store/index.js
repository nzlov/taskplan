import Vue from 'vue';
import Vuex from 'vuex';
import HttpUtil from '../utils/http';

Vue.use(Vuex);


const store = new Vuex.Store({
  state: {
    login: false,
    id: 0,
    token: '',
    user: '',
    realname: '',
    permissions: {},
    menu: [],
  },
  mutations: {
    init(state) {
      console.dir('init');
      const datastr = sessionStorage.session;
      if (datastr != null) {
        const data = JSON.parse(datastr);
        if (data.login) {
          HttpUtil.LGet(data, '/login').then((resp) => {
            switch (resp.data.code) {
              case 1000: {
                console.dir('check');
                state.id = data.id;
                state.user = data.user;
                state.token = data.token;
                state.realname = data.realname;
                state.permissions = data.permissions;
                state.menu = data.menu;
                state.login = true;
                break;
              }
              default:
                sessionStorage.clear();
            }
          }).catch(() => {
            sessionStorage.clear();
          });
        }
      }
    },
    login(state, data) {
      console.dir('login');
      state.id = data.User.ID;
      state.user = data.User.Name;
      state.token = data.token;
      state.realname = data.User.RealName;
      state.permissions = data.permission;
      state.menu = data.menu;
      state.login = true;
      sessionStorage.session = JSON.stringify(state);
      console.dir(sessionStorage.session);
    },
    logout(state) {
      console.dir('logout');
      state.login = false;
      state.id = 0;
      state.user = '';
      state.realname = '';
      state.token = '';
      state.permissions = {};
      state.menu = [];
      sessionStorage.clear();
    },
  },
});
export default store;
