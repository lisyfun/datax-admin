import { defineStore } from 'pinia';
import { getUserMenus, getMenuList } from '@/api/menu';
import type { MenuResponse } from '@/api/menu';

interface MenuState {
  menuList: MenuResponse[];
  loading: boolean;
}

export const useMenuStore = defineStore('menu', {
  state: (): MenuState => ({
    menuList: [],
    loading: false,
  }),

  getters: {
    // 获取树形菜单
    menuTree(): MenuResponse[] {
      const buildTree = (items: MenuResponse[], parentId: number = 0): MenuResponse[] => {
        return items
          .filter(item => item.parent_id === parentId)
          .map(item => ({
            ...item,
            children: buildTree(items, item.id),
          }))
          .sort((a, b) => a.sort - b.sort);
      };

      return buildTree(this.menuList);
    },
  },

  actions: {
    // 获取用户菜单列表
    async fetchUserMenus() {
      this.loading = true;
      try {
        const { data } = await getUserMenus();
        this.menuList = data.list;
      } catch (error) {
        console.error('获取用户菜单列表失败:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // 获取所有菜单列表（用于菜单管理）
    async fetchMenuList() {
      this.loading = true;
      try {
        const { data } = await getMenuList();
        this.menuList = data.list;
      } catch (error) {
        console.error('获取菜单列表失败:', error);
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
