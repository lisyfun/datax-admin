import { computed } from 'vue';
import { usePermissionStore } from '@/stores/permission';

/**
 * 权限控制组合函数
 */
export function usePermission() {
  const permissionStore = usePermissionStore();

  /**
   * 检查是否有指定权限
   * @param code 权限编码
   * @returns 是否有权限
   */
  const hasPermission = (code: string): boolean => {
    return permissionStore.hasPermission(code);
  };

  /**
   * 检查是否有任意一个权限
   * @param codes 权限编码数组
   * @returns 是否有任意一个权限
   */
  const hasAnyPermission = (codes: string[]): boolean => {
    return codes.some(code => permissionStore.hasPermission(code));
  };

  /**
   * 检查是否有所有权限
   * @param codes 权限编码数组
   * @returns 是否有所有权限
   */
  const hasAllPermissions = (codes: string[]): boolean => {
    return codes.every(code => permissionStore.hasPermission(code));
  };

  /**
   * 响应式权限检查
   * @param code 权限编码
   * @returns 响应式的权限状态
   */
  const permission = (code: string) => {
    return computed(() => permissionStore.hasPermission(code));
  };

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    permission,
  };
}
