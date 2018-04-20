import Vue from 'vue';
import Router from 'vue-router';
import Login from '@/components/Login';
import Regist from '@/components/Regist';
import Home from '@/components/Home';
import Role from '@/components/Role';
import User from '@/components/User';
import UserGroup from '@/components/UserGroup';
import Permission from '@/components/Permission';
import Task from '@/components/Task';
import Holiday from '@/components/Holiday';
import Leave from '@/components/Leave';

import store from '../store';

Vue.use(Router);

const router = new Router({
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/login',
      name: 'Login',
      component: Login,
    },
    {
      path: '/regist',
      name: 'Regist',
      component: Regist,
    },
    {
      path: '/task',
      name: 'Task',
      component: Task,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/holiday',
      name: 'Holiday',
      component: Holiday,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/leave',
      name: 'Leave',
      component: Leave,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/role',
      name: 'Role',
      component: Role,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/user',
      name: 'User',
      component: User,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/usergroup',
      name: 'UserGroup',
      component: UserGroup,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/permission',
      name: 'Permission',
      component: Permission,
      meta: {
        requiresAuth: true,
      },
    },
  ],
});

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!store.state.login) {
      next({
        path: '/login',
        query: { redirect: to.fullPath },
      });
    } else if (store.state.menu.indexOf(to.path) > -1) {
      // console.log('has permission');
      next();
    } else {
      // console.log('no permission');
      next(false);
    }
  } else if (store.state.login) {
    next({
      path: '/',
    });
  } else {
    next();
  }
});


export default router;
