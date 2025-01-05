import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import Layout from '@/layout/index.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/index.vue'),
    meta: {
      title: '注册',
      requiresAuth: false
    }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '仪表盘',
          requiresAuth: true
        }
      },
      {
        path: 'system',
        name: 'System',
        meta: {
          title: '系统管理',
          requiresAuth: true
        },
        children: [
          {
            path: 'permissions',
            name: 'Permissions',
            component: () => import('@/views/system/permissions/index.vue'),
            meta: {
              title: '权限管理',
              requiresAuth: true
            }
          },
          {
            path: 'roles',
            name: 'Roles',
            component: () => import('@/views/system/roles/index.vue'),
            meta: {
              title: '角色管理',
              requiresAuth: true
            }
          },
          {
            path: 'users',
            name: 'Users',
            component: () => import('@/views/system/users/index.vue'),
            meta: {
              title: '用户管理',
              requiresAuth: true
            }
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title} - DataX Admin`

  // 检查是否需要登录认证
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('token')
    if (!token) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }
  }

  next()
})

export default router
