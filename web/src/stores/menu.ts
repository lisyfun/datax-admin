import { defineStore } from 'pinia';
import * as permissionApi from '@/api/permission';
import type { PermissionInfo } from '@/types/permission';

interface MenuState {
  menuList: PermissionInfo[];
  loading: boolean;
}

export const useMenuStore = defineStore('menu', {
  state: (): MenuState => ({
    menuList: [],
    loading: false,
  }),

  getters: {
    // 获取树形菜单
    menuTree(): PermissionInfo[] {
      // 如果menuList已经是树形结构，直接返回
      return this.menuList;
    },

    // 获取扁平化的菜单列表（用于路径构建）
    flatMenuList(): PermissionInfo[] {
      const flatten = (menus: PermissionInfo[]): PermissionInfo[] => {
        let result: PermissionInfo[] = [];
        for (const menu of menus) {
          result.push(menu);
          if (menu.children && menu.children.length > 0) {
            result = result.concat(flatten(menu.children));
          }
        }
        return result;
      };
      return flatten(this.menuList);
    },
  },

  actions: {
    // 获取用户菜单列表
    async fetchUserMenus() {
      this.loading = true;
      try {
        const { data } = await permissionApi.getUserMenus();
        this.menuList = data.list;
      } catch (error) {
        console.error('获取用户菜单列表失败:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // 重置状态
    resetMenuState() {
      this.menuList = [];
      this.loading = false;
    },
  },
});
