import { DEFAULT_LAYOUT } from '@/router/base';
import { AppRouteRecordRaw } from '@/router/types';

const SYSTEM: AppRouteRecordRaw = {
  path: '/system',
  name: 'system',
  component: DEFAULT_LAYOUT,
  meta: {
    title: 'menu.system',
    requiresAuth: true,
    icon: 'icon-settings',
    order: 0,
  },
  children: [
    {
      path: 'menu',
      name: 'Menu',
      component: () => import('@/views/system/menu/index.vue'),
      meta: {
        title: 'menu.system.menu',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    // ... 其他系统管理路由
  ],
};

export default SYSTEM;
