import { defineStore } from 'pinia';
import * as menuApi from '@/api/menu';

interface MenuItem {
  id: number;
  parent_id: number;
  name: string;
  path: string;
  component: string;
  icon: string;
  sort: number;
  status: number;
  hidden: boolean;
  cache: boolean;
  type: number;
  children?: MenuItem[];
}

interface MenuState {
  menuList: MenuItem[];
  loading: boolean;
}

export const useMenuStore = defineStore('menu', {
  state: (): MenuState => ({
    menuList: [],
    loading: false,
  }),

  getters: {
    menuTree(): MenuItem[] {
      const buildTree = (items: MenuItem[], parentId: number = 0): MenuItem[] => {
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
    async fetchUserMenus() {
      this.loading = true;
      try {
        const { data } = await menuApi.getUserMenus();
        this.menuList = data.list;
      } catch (error) {
        console.error('获取用户菜单列表失败:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  },
});
