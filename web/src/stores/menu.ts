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

        // Helper function to convert MenuResponse to MenuItem recursively
        const convertToMenuItem = (item: menuApi.MenuResponse): MenuItem => {
          const menuItem: MenuItem = {
            id: item.id,
            parent_id: item.parent_id,
            name: item.name,
            path: item.path,
            component: item.component,
            icon: item.icon,
            sort: item.sort,
            status: item.status,
            hidden: Boolean(item.hidden),
            cache: Boolean(item.cache),
            type: item.type,
          };

          if (item.children && item.children.length > 0) {
            menuItem.children = item.children.map(convertToMenuItem);
          }

          return menuItem;
        };

        // Convert the whole list
        this.menuList = data.list.map(convertToMenuItem);
      } catch (error) {
        console.error('获取用户菜单列表失败:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  },
});
