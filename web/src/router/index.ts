import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import type { AppRouteRecordRaw } from './types';
import DefaultLayout from '@/layouts/default.vue';
import { useUserStore } from '@/stores/user';
import { usePermissionStore } from '@/stores/permission';

declare module 'vue-router' {
  interface RouteMeta {
    title?: string;
    requiresAuth?: boolean;
    icon?: string;
    order?: number;
    hideInMenu?: boolean;
    roles?: string[];
    permission?: string;
  }
}

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
        path: 'job',
        name: 'Job',
        component: () => import('@/views/job/index.vue'),
        redirect: '/job/list',
        meta: {
          title: '任务管理',
          requiresAuth: true,
          icon: 'icon-calendar',
          order: 1,
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
              permission: 'job.list',
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
              permission: 'job.history',
            },
          },
        ],
      },
      {
        path: 'terminal',
        name: 'Terminal',
        component: () => import('@/views/terminal/index.vue'),
        redirect: '/terminal/list',
        meta: {
          title: '终端管理',
          requiresAuth: true,
          icon: 'icon-command',
          order: 2,
        },
        children: [
          {
            path: 'list',
            name: 'TerminalList',
            component: () => import('@/views/terminal/list/index.vue'),
            meta: {
              title: '终端列表',
              requiresAuth: true,
              icon: 'icon-desktop',
              roles: ['*'],
              permission: 'terminal.list',
            },
          },
          {
            path: 'connect/:id',
            name: 'TerminalConnect',
            component: () => import('@/views/terminal/connect/index.vue'),
            meta: {
              title: '终端连接',
              requiresAuth: true,
              hideInMenu: true,
              roles: ['*'],
              permission: 'terminal.connect',
            },
          },
        ],
      },
      {
            path: 'kafka',
            name: 'Kafka',
            component: () => import('@/views/kafka/index.vue'),
            meta: {
              title: 'KAFKA 管理',
              requiresAuth: true,
              icon: 'icon-apps',
               order: 3,
            },
            redirect: '/kafka/cluster',
            children: [
              {
                path: 'cluster',
                name: 'KafkaCluster',
                component: () => import('@/views/kafka/cluster/index.vue'),
                meta: {
                  title: '集群管理',
                  requiresAuth: true,
                  roles: ['*'],
                  icon: 'icon-apps',
                  permission: 'tools.kafka.cluster',
                },
              },
              {
                path: 'cluster/:clusterId/topic',
                name: 'KafkaTopic',
                component: () => import('@/views/kafka/topic/index.vue'),
                meta: {
                  title: '主题管理',
                  requiresAuth: true,
                  roles: ['*'],
                  hideInMenu: true,
                  permission: 'tools.kafka.topic',
                },
              },
              {
                path: 'clusters/:clusterId/topics/:topicName/messages',
                name: 'KafkaMessage',
                component: () => import('@/views/kafka/topic/components/MessageList.vue'),
                meta: {
                  title: '消息列表',
                  requiresAuth: true,
                  roles: ['*'],
                  hideInMenu: true,
                  permission: 'tools.kafka.message',
                },
              },
            ],
          },
      {
        path: 'tools',
        name: 'Tools',
        component: () => import('@/views/tools/index.vue'),
        redirect: '/tools/json-formatter',
        meta: {
          title: '工具箱',
          requiresAuth: true,
          icon: 'icon-common',
          order: 4,
        },
        children: [
          {
            path: 'json-formatter',
            name: 'JsonFormatter',
            component: () => import('@/views/tools/JsonFormatter.vue'),
            meta: {
              title: 'JSON 格式化',
              requiresAuth: true,
              icon: 'icon-code',
              roles: ['*'],
              permission: 'tools.json-formatter',
            },
          },
          {
            path: 'crypto',
            name: 'Crypto',
            component: () => import('@/views/tools/Crypto.vue'),
            meta: {
              title: '加解密工具',
              requiresAuth: true,
              icon: 'icon-lock',
              roles: ['*'],
              permission: 'tools.crypto',
            },
          },

          {
            path: 'redis',
            name: 'Redis',
            component: () => import('@/views/tools/redis/index.vue'),
            meta: {
              title: 'Redis 管理',
              requiresAuth: true,
              icon: 'icon-database',
              roles: ['*'],
              permission: 'tools.redis',
            },
          },

        ],
      },
      {
        path: 'external-iframe',
        name: 'ExternalIframe',
        component: () => import('@/views/ExternalIframe.vue'),
        meta: {
          title: '外部页面',
          requiresAuth: true,
          hideInMenu: true,
          roles: ['*'],
        },
      },
      {
        path: 'system',
        name: 'System',
        component: () => import('@/views/system/index.vue'),
        redirect: '/system/users',
        meta: {
          title: '系统管理',
          requiresAuth: true,
          icon: 'icon-apps',
          order: 999,
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
              permission: 'system.users',
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
              permission: 'system.roles',
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
              permission: 'system.permissions',
            },
          },
          {
            path: 'logs',
            name: 'OperationLogs',
            component: () => import('@/views/system/logs/index.vue'),
            meta: {
              title: '操作管理',
              requiresAuth: true,
              icon: 'icon-file',
              roles: ['admin'],
              permission: 'system.logs',
            },
          },
        ],
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/404/index.vue'),
    meta: {
      title: '页面未找到',
      requiresAuth: false,
      hideInMenu: true,
    },
  },
];

const router = createRouter({
  history: createWebHashHistory('/datax/'),
  routes: appRoutes as RouteRecordRaw[],
});

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false);

  // 不需要认证的路由直接放行
  if (!requiresAuth) {
    next();
    return;
  }

  const userStore = useUserStore();

  try {
    // 尝试获取用户信息，如果成功表示已登录
    await userStore.getUserInfo();
  } catch (err) {
    // 获取用户信息失败，表示未登录
    if (to.path !== '/login') {
      next({ path: '/login', query: { redirect: to.fullPath } });
    } else {
      next();
    }
    return;
  }

  // 已登录用户访问登录页，重定向到首页
  if (to.path === '/login') {
    next({ path: '/' });
    return;
  }

  // 路由级权限校验：没有对应权限码时不允许直接通过地址栏访问
  const permissionCode = to.meta.permission;
  if (permissionCode) {
    const permissionStore = usePermissionStore();
    try {
      await permissionStore.ensurePermissions(userStore.userInfo?.id);
    } catch (err) {
      // 权限加载失败时按无权限处理，避免绕过路由限制
      console.error('加载用户权限失败:', err);
    }

    if (!permissionStore.hasPermissionCode(permissionCode)) {
      Message.warning('您没有访问该页面的权限');
      next({ path: '/dashboard', replace: true });
      return;
    }
  }

  next();
});

export default router;
