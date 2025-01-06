import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import DefaultLayout from '@/layouts/default.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/index.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: DefaultLayout,
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'system/users',
        name: 'Users',
        component: () => import('@/views/system/users/index.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'system/roles',
        name: 'Roles',
        component: () => import('@/views/system/roles/index.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'system/permissions',
        name: 'Permissions',
        component: () => import('@/views/system/permissions/index.vue'),
        meta: { requiresAuth: true }
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
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
