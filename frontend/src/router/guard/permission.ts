import type { Router, RouteRecordRaw } from 'vue-router';
import { useUserStore } from '@/stores/modules/user';
import { useMenuStore } from '@/stores/modules/menu';
import { WHITE_LIST } from '../constants';

export default function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore();
    const menuStore = useMenuStore();
    const token = userStore.token;

    // 白名单路由直接放行
    if (WHITE_LIST.includes(to.name as string)) {
      next();
      return;
    }

    // 未登录跳转到登录页
    if (!token) {
      next({
        name: 'login',
        query: {
          redirect: to.fullPath,
        },
      });
      return;
    }

    // 已经获取过菜单列表，直接放行
    if (menuStore.menuList.length > 0) {
      next();
      return;
    }

    try {
      // 获取用户菜单列表
      await menuStore.fetchUserMenus();
      // 生成动态路由
      const routes = generateRoutesFromMenu(menuStore.menuTree);
      // 添加动态路由
      routes.forEach((route) => {
        router.addRoute(route);
      });
      // 重定向到目标路由
      next({ ...to, replace: true });
    } catch (error) {
      // 获取菜单失败，可能是 token 过期等原因
      userStore.logout();
      next({
        name: 'login',
        query: {
          redirect: to.fullPath,
        },
      });
    }
  });
}

// 根据菜单生成路由配置
function generateRoutesFromMenu(menus: any[]): RouteRecordRaw[] {
  return menus
    .filter((menu) => menu.type === 1 && !menu.hidden) // 只处理显示的菜单类型
    .map((menu) => {
      const route: RouteRecordRaw = {
        path: menu.path,
        name: menu.name,
        component: loadComponent(menu.component),
        meta: {
          title: menu.name,
          icon: menu.icon,
          keepAlive: menu.cache,
        },
      };

      if (menu.children?.length) {
        route.children = generateRoutesFromMenu(menu.children);
      }

      return route;
    });
}

// 动态加载组件
function loadComponent(component: string) {
  // 这里假设所有的组件都在 views 目录下
  return () => import(`@/views/${component}.vue`);
}
