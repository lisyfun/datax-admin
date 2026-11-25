import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { PermissionInfo } from '@/types/permission';
import * as permissionApi from '@/api/permission';
import { useUserStore } from '@/stores/user';

export const usePermissionStore = defineStore('permission', () => {
  const permissions = ref<PermissionInfo[]>([]);
  const menuTree = ref<PermissionInfo[]>([]);
  // 当前权限数据所属的用户，用于避免切换账号后复用上一个用户的权限
  const loadedUserId = ref<number | null>(null);

  async function getUserPermissions() {
    const res = await permissionApi.getUserPermissions();
    permissions.value = res.data.permissions || [];
    // 过滤出菜单权限并构建树形结构
    const menus = permissions.value.filter(p => p.type === 'menu');
    menuTree.value = buildMenuTree(menus);
    const userStore = useUserStore();
    loadedUserId.value = userStore.userInfo?.id ?? null;
    console.log('用户权限已加载:', permissions.value.length, '个权限');
    console.log('菜单权限:', menus.length, '个');
    console.log('按钮权限:', permissions.value.filter(p => p.type === 'button').length, '个');
  }

  // 确保当前用户的权限已加载（路由守卫使用）
  async function ensurePermissions(userId?: number) {
    if (loadedUserId.value !== null && (userId === undefined || loadedUserId.value === userId)) {
      return;
    }
    await getUserPermissions();
  }

  // 重置权限状态（退出登录时调用）
  function resetPermissions() {
    permissions.value = [];
    menuTree.value = [];
    loadedUserId.value = null;
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

  // 仅判断权限码是否存在（用于路由守卫，隐藏菜单型权限同样有效）
  function hasPermissionCode(code: string): boolean {
    return permissions.value.some(p => p.code === code);
  }

  return {
    permissions,
    menuTree,
    loadedUserId,
    getUserPermissions,
    ensurePermissions,
    resetPermissions,
    getPermissionTree,
    hasPermission,
    hasPermissionCode,
  };
});
