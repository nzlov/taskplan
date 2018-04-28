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
    usergroup: 0,
    realname: '',
    permissions: {},
    editall: false, // 允许编辑所有的任务
    menu: [],
    users: [],
    roles: [],
    holidays: [],
    usergroups: [],
    constmenu: [
      {
        icon: 'home',
        path: '/',
        title: '主页',
      },
      {
        icon: 'assignment',
        path: '/task',
        title: '任务',
      },
      {
        icon: 'bubble_chart',
        path: '/task_gantt',
        title: '任务甘特图',
      },
      {
        icon: 'person',
        path: '/user',
        title: '用户',
      },
      {
        icon: 'bubble_chart',
        path: '/leave',
        title: '请假',
      },
      {
        icon: 'bubble_chart',
        path: '/holiday',
        title: '假期',
      },
      {
        icon: 'account_box',
        path: '/role',
        title: '角色',
      },
      {
        icon: 'group',
        path: '/usergroup',
        title: '用户组',
      },
      {
        icon: 'bubble_chart',
        path: '/permission',
        title: '权限',
      },
    ],
  },
  mutations: {
    init(state) {
      console.dir('init');
      const datastr = sessionStorage.session;
      if (datastr != null) {
        const data = JSON.parse(datastr);
        if (data.login) {
          console.dir('check');
          HttpUtil.LGet(data, '/login').then((resp) => {
            switch (resp.data.code) {
              case 1000: {
                state.id = data.id;
                state.user = data.user;
                state.token = data.token;
                state.usergroup = data.usergroup;
                state.realname = data.realname;
                state.permissions = data.permissions;
                state.menu = data.menu;
                state.editall = data.permissions['list.all'] !== undefined;
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
      state.usergroup = data.User.UserGroupID;
      state.token = data.token;
      state.realname = data.User.RealName;
      state.permissions = data.permission;
      state.editall = (data.permission['list.all'] !== undefined);
      state.menu = data.menu;
      state.login = true;
      sessionStorage.session = JSON.stringify(state);
    },
    logout(state) {
      console.dir('logout');
      state.login = false;
      state.id = 0;
      state.usergroup = 0;
      state.user = '';
      state.realname = '';
      state.token = '';
      state.permissions = {};
      state.menu = [];
      sessionStorage.clear();
    },
    reloadusers(state) {
      new Promise((resolve) => {
        const objs = [{
          name: '无',
          id: 0,
        }];
        HttpUtil.LGet(state, '/user').then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                objs.push({
                  id: element.ID,
                  name: element.RealName,
                });
              });
              break;
            }
            default:
              console.dir('服务器报错');
          }
          resolve({
            objs,
          });
        }).catch((e) => {
          console.dir(e);
          console.dir('服务器报错');
          resolve({
            objs,
          });
        });
      }).then((data) => {
        state.users = data.objs;
      });
    },
    reloadusergroups(state) {
      new Promise((resolve) => {
        const objs = [{
          name: '无',
          id: 0,
        }];
        HttpUtil.LGet(state, '/usergroup').then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                objs.push({
                  id: element.ID,
                  name: element.Name,
                });
              });
              break;
            }
            default:
              console.dir('服务器报错');
          }
          resolve({
            objs,
          });
        }).catch(() => {
          console.dir('服务器报错');
          resolve({
            objs,
          });
        });
      }).then((data) => {
        state.usergroups = data.objs;
      });
    },
    reloadroles(state) {
      new Promise((resolve) => {
        const objs = [{
          name: '无',
          id: 0,
        }];
        HttpUtil.LGet(state, '/role').then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                objs.push({
                  id: element.ID,
                  name: element.Name,
                });
              });
              break;
            }
            default:
              console.dir('服务器报错');
          }
          resolve({
            objs,
          });
        }).catch(() => {
          console.dir('服务器报错');
          resolve({
            objs,
          });
        });
      }).then((data) => {
        state.roles = data.objs;
      });
    },
    reloadholidays(state) {
      const formatDate = (d) => {
        if (d <= 0) {
          return '';
        }
        const now = new Date(d * 1000);
        const year = now.getFullYear();
        let month = now.getMonth() + 1;
        let date = now.getDate();
        if (month <= 9) {
          month = `0${month}`;
        }
        if (date <= 9) {
          date = `0${date}`;
        }
        return `${year}-${month}-${date}`;
      };
      new Promise((resolve) => {
        const objs = [];
        HttpUtil.LGet(state, '/holiday').then((resp) => {
          switch (resp.data.code) {
            case 0: {
              resp.data.data.data.forEach((element) => {
                objs.push(formatDate(element.Day));
              });
              break;
            }
            default:
              console.dir('服务器报错');
          }
          resolve({
            objs,
          });
        }).catch(() => {
          console.dir('服务器报错');
          resolve({
            objs,
          });
        });
      }).then((data) => {
        state.holidays = data.objs;
      });
    },
  },
});
export default store;
