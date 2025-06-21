import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { PermissionInfo } from '@/types/permission';
import * as permissionApi from '@/api/permission';

export const usePermissionStore = defineStore('permission', () => {
  const permissions = ref<PermissionInfo[]>([]);
  const menuTree = ref<PermissionInfo[]>([]);

  async function getUserPermissions() {
    const res = await permissionApi.getUserPermissions();
    permissions.value = res.data.permissions || [];
    // 过滤出菜单权限并构建树形结构
    const menus = permissions.value.filter(p => p.type === 'menu');
    menuTree.value = buildMenuTree(menus);
    console.log('用户权限已加载:', permissions.value.length, '个权限');
    console.log('菜单权限:', menus.length, '个');
    console.log('按钮权限:', permissions.value.filter(p => p.type === 'button').length, '个');
  }

  async function getPermissionTree() {
    const res = await permissionApi.getPermissionTree();
    return res.data.list;
  }

  function buildMenuTree(menus: PermissionInfo[]): PermissionInfo[] {
    const menuMap = new Map<number, PermissionInfo>();
    const result: PermissionInfo[] = [];

    // 创建菜单映射
    menus.forEach(menu => {
      menuMap.set(menu.id, { ...menu, children: [] });
    });

    // 构建树形结构
    menus.forEach(menu => {
      const node = menuMap.get(menu.id)!;
      if (menu.parent_id) {
        const parent = menuMap.get(menu.parent_id);
        if (parent) {
          parent.children?.push(node);
        }
      } else {
        result.push(node);
      }
    });

    return result;
  }

  function hasPermission(code: string): boolean {
    const permission = permissions.value.find(p => p.code === code);
    // 权限必须存在且不能被隐藏
    const result = permission && permission.hidden === 0;
    return !!result;
  }

  return {
    permissions,
    menuTree,
    getUserPermissions,
    getPermissionTree,
    hasPermission,
  };
});
