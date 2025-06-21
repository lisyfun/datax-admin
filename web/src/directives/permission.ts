import type { App, DirectiveBinding } from 'vue';
import { usePermissionStore } from '@/stores/permission';
import { watch } from 'vue';

/**
 * 权限指令
 * 用法：v-permission="'permission.code'"
 * 或者：v-permission="['permission.code1', 'permission.code2']"
 */
export default {
  install(app: App) {
    app.directive('permission', {
      mounted(el: HTMLElement, binding: DirectiveBinding) {
        const permissionStore = usePermissionStore();
        const { value } = binding;

        console.log('v-permission directive mounted:', value);

        if (!value) {
          console.warn('v-permission directive requires a permission code');
          return;
        }

        // 检查权限的函数
        const checkPermission = () => {
          let hasPermission = false;

          if (Array.isArray(value)) {
            // 如果是数组，检查是否有任意一个权限
            hasPermission = value.some(code => permissionStore.hasPermission(code));
          } else {
            // 如果是字符串，检查单个权限
            hasPermission = permissionStore.hasPermission(value);
          }

          console.log(`Permission check for "${value}":`, hasPermission);

          if (!hasPermission) {
            // 如果没有权限，隐藏元素
            el.style.display = 'none';
            console.log(`Element hidden due to missing permission: ${value}`);
          } else {
            // 如果有权限，显示元素
            el.style.display = '';
            console.log(`Element shown with permission: ${value}`);
          }
        };

        // 立即检查一次
        checkPermission();

        // 监听权限数据变化
        const stopWatcher = watch(
          () => permissionStore.permissions,
          () => {
            console.log('Permissions changed, rechecking...');
            checkPermission();
          },
          { deep: true }
        );

        // 将停止监听的函数存储到元素上，以便在卸载时清理
        (el as any)._permissionWatcher = stopWatcher;
      },

      updated(el: HTMLElement, binding: DirectiveBinding) {
        const permissionStore = usePermissionStore();
        const { value } = binding;

        if (!value) {
          return;
        }

        let hasPermission = false;

        if (Array.isArray(value)) {
          hasPermission = value.some(code => permissionStore.hasPermission(code));
        } else {
          hasPermission = permissionStore.hasPermission(value);
        }

        if (hasPermission) {
          el.style.display = '';
        } else {
          el.style.display = 'none';
        }
      },

      unmounted(el: HTMLElement) {
        // 清理监听器
        if ((el as any)._permissionWatcher) {
          (el as any)._permissionWatcher();
          delete (el as any)._permissionWatcher;
        }
      }
    });
  }
};
