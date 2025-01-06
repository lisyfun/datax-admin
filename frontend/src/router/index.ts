import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import type { AppRouteRecordRaw } from './types';
import DefaultLayout from '@/layouts/default.vue';

export const appRoutes: AppRouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      requiresAuth: false,
      hideInMenu: true,
    },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/index.vue'),
    meta: {
      title: '注册',
      requiresAuth: false,
      hideInMenu: true,
    },
  },
  {
    path: '/',
    name: 'Root',
    component: DefaultLayout,
    redirect: '/dashboard',
    meta: {
      title: '首页',
      requiresAuth: true,
    },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '仪表盘',
          requiresAuth: true,
          icon: 'icon-dashboard',
          order: 0,
          roles: ['*'],
        },
      },
      {
        path: 'system',
        name: 'System',
        component: () => import('@/views/system/index.vue'),
        meta: {
          title: '系统管理',
          requiresAuth: true,
          icon: 'icon-apps',
          order: 1,
        },
        children: [
          {
            path: 'users',
            name: 'Users',
            component: () => import('@/views/system/users/index.vue'),
            meta: {
              title: '用户管理',
              requiresAuth: true,
              icon: 'icon-user',
              roles: ['admin'],
            },
          },
          {
            path: 'roles',
            name: 'Roles',
            component: () => import('@/views/system/roles/index.vue'),
            meta: {
              title: '角色管理',
              requiresAuth: true,
              icon: 'icon-user-group',
              roles: ['admin'],
            },
          },
          {
            path: 'permissions',
            name: 'Permissions',
            component: () => import('@/views/system/permissions/index.vue'),
            meta: {
              title: '权限管理',
              requiresAuth: true,
              icon: 'icon-safe',
              roles: ['admin'],
            },
          },
        ],
      },
      {
        path: 'job',
        name: 'Job',
        component: () => import('@/views/job/index.vue'),
        meta: {
          title: '任务管理',
          requiresAuth: true,
          icon: 'icon-calendar',
          order: 2,
        },
        children: [
          {
            path: 'list',
            name: 'JobList',
            component: () => import('@/views/job/list/index.vue'),
            meta: {
              title: '任务列表',
              requiresAuth: true,
              icon: 'icon-unordered-list',
              roles: ['*'],
            },
          },
          {
            path: 'history',
            name: 'JobHistory',
            component: () => import('@/views/job/history/index.vue'),
            meta: {
              title: '执行历史',
              requiresAuth: true,
              icon: 'icon-clock-circle',
              roles: ['*'],
            },
          },
        ],
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes: appRoutes as RouteRecordRaw[],
});

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false);

  if (requiresAuth && !token) {
    // 需要认证但未登录，重定向到登录页
    next({ path: '/login', query: { redirect: to.fullPath } });
  } else if (token && (to.path === '/login' || to.path === '/register')) {
    // 已登录但访问登录或注册页，重定向到首页
    next({ path: '/' });
  } else {
    next();
  }
});

export default router;
