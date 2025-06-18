<template>
  <a-menu
    :selected-keys="selectedKeys"
    :open-keys="openKeys"
    :auto-open="false"
    @menu-item-click="handleMenuClick"
    @sub-menu-click="handleSubMenuClick"
  >
    <template v-for="menu in menuTree" :key="menu.id">
      <!-- 有子菜单的情况 -->
      <a-sub-menu v-if="menu.children?.length" :key="`sub-${menu.id}`">
        <template #icon>
          <component v-if="menu.icon && iconMap[menu.icon]" :is="iconMap[menu.icon]" />
        </template>
        <template #title>{{ menu.name }}</template>
        <template v-for="child in menu.children" :key="child.id">
          <a-menu-item
            v-if="!child.children?.length"
            :key="`item-${child.id}`"
          >
            <template #icon>
              <component v-if="child.icon && iconMap[child.icon]" :is="iconMap[child.icon]" />
            </template>
            {{ child.name }}
          </a-menu-item>
          <a-sub-menu v-else :key="`sub-${child.id}`">
            <template #icon>
              <component v-if="child.icon && iconMap[child.icon]" :is="iconMap[child.icon]" />
            </template>
            <template #title>{{ child.name }}</template>
            <a-menu-item
              v-for="item in child.children"
              :key="`item-${item.id}`"
            >
              <template #icon>
                <component v-if="item.icon && iconMap[item.icon]" :is="iconMap[item.icon]" />
              </template>
              {{ item.name }}
            </a-menu-item>
          </a-sub-menu>
        </template>
      </a-sub-menu>
      <!-- 没有子菜单的情况 -->
      <a-menu-item v-else :key="`item-${menu.id}`">
        <template #icon>
          <component v-if="menu.icon && iconMap[menu.icon]" :is="iconMap[menu.icon]" />
        </template>
        {{ menu.name }}
      </a-menu-item>
    </template>
  </a-menu>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useMenuStore } from '@/stores/menu';
import {
  IconDashboard,
  IconApps,
  IconUser,
  IconUserGroup,
  IconSafe,
  IconCalendar,
  IconUnorderedList,
  IconClockCircle,
  IconDesktop,
  IconCloud,
  IconFile,
  IconBulb,
  IconCode,
  IconRobot,
  IconCommon,
  IconCommand,
  IconLock,
  IconList,
} from '@arco-design/web-vue/es/icon';

const route = useRoute();
const router = useRouter();
const menuStore = useMenuStore();

// 图标映射
const iconMap: Record<string, any> = {
  'icon-dashboard': IconDashboard,
  'icon-apps': IconApps,
  'icon-user': IconUser,
  'icon-user-group': IconUserGroup,
  'icon-safe': IconSafe,
  'icon-calendar': IconCalendar,
  'icon-unordered-list': IconUnorderedList,
  'icon-clock-circle': IconClockCircle,
  'icon-desktop': IconDesktop,
  'icon-cloud': IconCloud,
  'icon-file': IconFile,
  'icon-bulb': IconBulb,
  'icon-code': IconCode,
  'icon-robot': IconRobot,
  'icon-common': IconCommon,
  'icon-command': IconCommand,
  'icon-lock': IconLock,
  'icon-menu': IconList,
};

// 当前选中的菜单项
const selectedKeys = ref<string[]>([]);
// 当前展开的子菜单
const openKeys = ref<string[]>([]);
// 标记是否是手动点击菜单
const isManualClick = ref(false);

// 设置默认展开的菜单
const setDefaultOpenKeys = () => {
  // 暂时禁用默认展开，避免影响菜单选中状态
  // 可以根据菜单名称或ID来设置默认展开的菜单
  // const defaultOpenMenus = ['系统管理']; // 可以添加多个菜单名称

  // 只有在没有当前路由匹配的菜单时，才设置默认展开
  if (selectedKeys.value.length === 0 && openKeys.value.length === 0) {
    console.log('没有选中的菜单，不设置默认展开');
  }
};

// 获取树形菜单数据
const menuTree = computed(() => {
  const tree = menuStore.menuTree;
  // 如果只有一个根节点且是"首页"，则直接返回其子菜单
  if (tree.length === 1 && tree[0].code === 'root' && tree[0].children) {
    return tree[0].children;
  }
  return tree;
});

// 提取菜单ID（去掉前缀）
const extractMenuId = (key: string): string => {
  if (key.startsWith('sub-') || key.startsWith('item-')) {
    return key.substring(key.indexOf('-') + 1);
  }
  return key;
};

// 处理菜单点击
const handleMenuClick = async (key: string) => {
  console.log('菜单点击事件触发, key:', key, 'key类型:', typeof key);

  // 标记为手动点击
  isManualClick.value = true;

  // 提取真实的菜单ID
  const menuId = extractMenuId(key);

  // 使用完整的菜单列表查找，而不是树形结构
  const fullMenuList = menuStore.flatMenuList;

  // 将 key 转换为数字进行比较
  const keyNum = parseInt(menuId);
  console.log('转换后的key:', keyNum);

  const menu = fullMenuList.find(m => m.id === keyNum);
  console.log('找到的菜单:', menu);

  if (menu && menu.path) {
    // 在跳转前，先保存当前的展开状态
    const currentParentKeys = getParentKeys(menuTree.value, menuId);
    console.log('点击前保存父级keys:', currentParentKeys);

    // 构建完整路径
    const fullPath = buildFullPath(menu);
    console.log('点击菜单:', menu.name, '原始路径:', menu.path, '构建路径:', fullPath);

    // 跳转路由
    await router.push(fullPath);

    // 设置选中状态
    selectedKeys.value = [key];

    // 确保父级菜单保持展开状态
    if (currentParentKeys.length > 0) {
      openKeys.value = currentParentKeys.map(k => `sub-${k}`);
      console.log('设置父级菜单展开状态:', openKeys.value);
    }

    console.log('设置选中菜单:', selectedKeys.value);

    // 重置标记
    setTimeout(() => {
      isManualClick.value = false;
    }, 100);
  } else {
    console.log('菜单或路径不存在, menu:', !!menu, 'path:', menu?.path);
    isManualClick.value = false;
  }
};

// 构建完整路径
const buildFullPath = (menu: any): string => {
  // 根据菜单数据结构，直接构建路径
  // 仪表盘: /dashboard
  // 任务列表: /job/list
  // 终端列表: /terminal/list
  // 用户管理: /system/users

  if (!menu.path) {
    return '/';
  }

  // 如果是根菜单的直接子菜单（parent_id = 1），直接使用路径
  if (menu.parent_id === 1) {
    return '/' + menu.path;
  }

  // 如果有父菜单，需要构建完整路径
  const fullMenuList = menuStore.flatMenuList;
  const parent = findMenuById(fullMenuList, menu.parent_id);

  if (parent && parent.parent_id === 1) {
    // 父菜单是根菜单的子菜单，构建二级路径
    return '/' + parent.path + '/' + menu.path;
  }

  // 默认情况，直接使用菜单路径
  return '/' + menu.path;
};

// 根据ID查找菜单项
const findMenuById = (menus: any[], id: number): any => {
  for (const menu of menus) {
    if (menu.id === id) {
      return menu;
    }
    if (menu.children?.length) {
      const found = findMenuById(menu.children, id);
      if (found) return found;
    }
  }
  return null;
};

// 获取菜单的层级
const getMenuLevel = (menus: any[], targetKey: string, level = 0): number => {
  for (const menu of menus) {
    if (menu.id.toString() === targetKey) {
      return level;
    }
    if (menu.children?.length) {
      const childLevel = getMenuLevel(menu.children, targetKey, level + 1);
      if (childLevel !== -1) {
        return childLevel;
      }
    }
  }
  return -1;
};



// 处理子菜单点击（手风琴效果）
const handleSubMenuClick = (key: string) => {
  // 提取真实的菜单ID
  const menuId = extractMenuId(key);
  const index = openKeys.value.indexOf(key);

  if (index > -1) {
    // 如果当前菜单已展开，则收起它及其所有子菜单
    const keysToRemove = [key];

    // 递归找到所有子菜单
    const findAllChildKeys = (menus: any[], parentKey: string) => {
      menus.forEach(menu => {
        if (menu.id.toString() === parentKey && menu.children?.length) {
          menu.children.forEach((child: any) => {
            if (child.children?.length) {
              keysToRemove.push(`sub-${child.id}`);
              findAllChildKeys(menu.children, child.id.toString());
            }
          });
        }
        if (menu.children?.length) {
          findAllChildKeys(menu.children, parentKey);
        }
      });
    };

    findAllChildKeys(menuTree.value, menuId);

    // 移除当前菜单及其所有子菜单
    openKeys.value = openKeys.value.filter(k => !keysToRemove.includes(k));
  } else {
    // 如果当前菜单未展开，实现手风琴效果
    const parentKeys = getParentKeys(menuTree.value, menuId);

    // 找到当前菜单的直接父菜单
    let directParent = null;
    const findDirectParent = (menus: any[], targetKey: string): any => {
      for (const menu of menus) {
        if (menu.children?.length) {
          for (const child of menu.children) {
            if (child.id.toString() === targetKey) {
              return menu;
            }
          }
          const found = findDirectParent(menu.children, targetKey);
          if (found) return found;
        }
      }
      return null;
    };

    directParent = findDirectParent(menuTree.value, menuId);

    // 找到需要收起的同级菜单
    const keysToRemove: string[] = [];

    if (directParent) {
      // 收起同级的其他有子菜单的菜单项
      directParent.children.forEach((sibling: any) => {
        if (sibling.id.toString() !== menuId && sibling.children?.length) {
          keysToRemove.push(`sub-${sibling.id}`);

          // 递归收起其所有子菜单
          const collectAllChildren = (children: any[]) => {
            children.forEach(child => {
              if (child.children?.length) {
                keysToRemove.push(`sub-${child.id}`);
                collectAllChildren(child.children);
              }
            });
          };

          collectAllChildren(sibling.children);
        }
      });
    } else {
      // 如果是顶级菜单，收起其他顶级菜单
      menuTree.value.forEach(menu => {
        if (menu.id.toString() !== menuId && menu.children?.length) {
          keysToRemove.push(`sub-${menu.id}`);

          // 递归收起其所有子菜单
          const collectAllChildren = (children: any[]) => {
            children.forEach(child => {
              if (child.children?.length) {
                keysToRemove.push(`sub-${child.id}`);
                collectAllChildren(child.children);
              }
            });
          };

          collectAllChildren(menu.children);
        }
      });
    }

    // 构建新的openKeys：保留所有父级菜单 + 当前菜单
    const newOpenKeys = [...parentKeys.map(k => `sub-${k}`), key];

    // 移除同级菜单，保留父级和当前菜单
    openKeys.value = newOpenKeys.filter(k => !keysToRemove.includes(k));
  }

  console.log('子菜单点击，key:', key, '父级菜单:', getParentKeys(menuTree.value, menuId), '当前展开的菜单:', openKeys.value);
};

// 根据 key 查找菜单项
const findMenuByKey = (menus: any[], key: string): any => {
  console.log('查找菜单, key:', key, '菜单列表:', menus);

  for (const menu of menus) {
    console.log('检查菜单:', menu.name, 'ID:', menu.id, '比较:', menu.id.toString() === key);

    if (menu.id.toString() === key) {
      return menu;
    }
    if (menu.children?.length) {
      const found = findMenuByKey(menu.children, key);
      if (found) return found;
    }
  }
  return null;
};

// 根据路由路径查找菜单项
const findMenuByPath = (menus: any[], path: string): any => {
  for (const menu of menus) {
    if (menu.path === path) {
      return menu;
    }
    if (menu.children?.length) {
      const found = findMenuByPath(menu.children, path);
      if (found) return found;
    }
  }
  return null;
};

// 获取所有父级菜单的 key
const getParentKeys = (menus: any[], targetKey: string, keys: string[] = []): string[] => {
  for (const menu of menus) {
    if (menu.children?.length) {
      if (menu.children.some((child: any) => child.id.toString() === targetKey)) {
        keys.push(menu.id.toString());
      }
      getParentKeys(menu.children, targetKey, keys);
    }
  }
  return keys;
};

// 根据构建的路径查找菜单项
const findMenuByConstructedPath = (menus: any[], targetPath: string): any => {
  for (const menu of menus) {
    const constructedPath = buildFullPath(menu);
    if (constructedPath === targetPath) {
      return menu;
    }
    if (menu.children?.length) {
      const found = findMenuByConstructedPath(menu.children, targetPath);
      if (found) return found;
    }
  }
  return null;
};

// 监听路由变化，更新选中的菜单项
watch(
  () => route.path,
  (path) => {
    console.log('=== 路由变化监听 ===');
    console.log('当前路由路径:', path);
    console.log('是否手动点击:', isManualClick.value);
    console.log('当前选中菜单:', selectedKeys.value);
    console.log('当前展开菜单:', openKeys.value);

    // 如果是手动点击菜单导致的路由变化，不重新设置展开状态
    if (isManualClick.value) {
      console.log('手动点击菜单，跳过展开状态重置');
      // 重置手动点击标记，确保下次路由变化能正常处理
      setTimeout(() => {
        isManualClick.value = false;
      }, 100);
      return;
    }

    // 简化匹配逻辑：根据路径直接匹配
    let matchedMenu = null;

    // 遍历所有菜单项，找到路径匹配的
    const findMatchingMenu = (menus: any[]): any => {
      for (const menu of menus) {
        const constructedPath = buildFullPath(menu);
        console.log('检查菜单:', menu.name, 'ID:', menu.id, '构建路径:', constructedPath, '当前路径:', path, '是否匹配:', constructedPath === path);

        // 精确匹配优先
        if (constructedPath === path) {
          console.log('✅ 找到精确匹配:', menu.name, 'ID:', menu.id);
          return menu;
        }

        // 递归检查子菜单
        if (menu.children?.length) {
          const found = findMatchingMenu(menu.children);
          if (found) return found;
        }
      }

      return null;
    };

    console.log('开始查找匹配菜单...');
    matchedMenu = findMatchingMenu(menuTree.value);

    if (matchedMenu) {
      const key = matchedMenu.id.toString();
      const parentKeys = getParentKeys(menuTree.value, key);
      const newSelectedKeys = [`item-${key}`];
      const newOpenKeys = parentKeys.map(k => `sub-${k}`);

      console.log('✅ 匹配到菜单:', matchedMenu.name, 'ID:', matchedMenu.id);
      console.log('设置选中菜单:', newSelectedKeys);
      console.log('设置展开菜单:', newOpenKeys);

      selectedKeys.value = newSelectedKeys;
      openKeys.value = newOpenKeys;
    } else {
      console.log('❌ 未找到匹配的菜单，清空选中状态');
      selectedKeys.value = [];
    }
    console.log('=== 路由变化监听结束 ===');
  },
  { immediate: true }
);

// 监听菜单数据变化
watch(
  () => menuStore.menuTree,
  () => {
    console.log('菜单数据已更新');
    // 菜单数据更新后设置默认展开
    setDefaultOpenKeys();
  },
  { immediate: true }
);


</script>

<style scoped>
:deep(.arco-menu-inline-header) {
  height: 40px;
  line-height: 40px;
}

/* 统一所有菜单项样式 */
:deep(.arco-menu-item),
:deep(.arco-sub-menu-title),
:deep(.arco-menu-inline-header) {
  height: 40px !important;
  line-height: 40px !important;
  margin: 2px 8px !important;
  border-radius: 6px !important;
  transition: all 0.2s ease;
  padding: 0 16px !important;
  display: flex !important;
  align-items: center !important;
  box-sizing: border-box !important;
}

/* 特别针对有图标的菜单项 */
:deep(.arco-menu-item.arco-menu-has-icon),
:deep(.arco-sub-menu-title) {
  padding: 0 16px !important;
}

:deep(.arco-menu-item:hover),
:deep(.arco-sub-menu-title:hover),
:deep(.arco-menu-inline-header:hover) {
  background-color: rgba(var(--primary-6), 0.1);
}

:deep(.arco-menu-item.arco-menu-selected),
:deep(.arco-sub-menu-title.arco-sub-menu-open),
:deep(.arco-menu-inline-header.arco-menu-selected) {
  background-color: rgba(var(--primary-6), 0.15);
  border-radius: 6px;
  color: rgb(var(--primary-6));
  font-weight: 500;
}

:deep(.arco-menu-item.arco-menu-selected::before),
:deep(.arco-sub-menu-title.arco-sub-menu-open::before) {
  display: none;
}

:deep(.arco-sub-menu-inline .arco-menu-item) {
  margin: 2px 16px 2px 24px;
}

/* 收起状态下的菜单项样式 - 统一所有类型 */
:deep(.arco-menu-collapse .arco-menu-item),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-has-icon),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop-header),
:deep(.arco-menu-collapse .arco-sub-menu-title),
:deep(.arco-menu-collapse .arco-menu-inline-header) {
  margin: 2px 4px !important;
  padding: 0 !important;
  justify-content: center !important;
  display: flex !important;
  align-items: center !important;
  border-radius: 6px !important;
  height: 40px !important;
  width: calc(100% - 8px) !important;
  box-sizing: border-box !important;
}

/* 收起状态下的菜单项内容样式 */
:deep(.arco-menu-collapse .arco-menu-item-inner),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-has-icon .arco-menu-item-inner),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop .arco-menu-item-inner),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop-header .arco-menu-item-inner),
:deep(.arco-menu-collapse .arco-sub-menu-title),
:deep(.arco-menu-collapse .arco-menu-inline-header) {
  justify-content: center !important;
  padding: 0 !important;
  width: 100% !important;
  display: flex !important;
  align-items: center !important;
}

/* 收起状态下隐藏文字 - 统一所有类型 */
:deep(.arco-menu-collapse .arco-sub-menu-title .arco-menu-title),
:deep(.arco-menu-collapse .arco-menu-item .arco-menu-title),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop .arco-menu-title),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop-header .arco-menu-title),
:deep(.arco-menu-collapse .arco-menu-inline-header .arco-menu-title) {
  display: none !important;
}

:deep(.arco-menu-collapse .arco-sub-menu-suffix) {
  display: none;
}

/* 确保图标在收起状态下居中 - 统一所有类型 */
:deep(.arco-menu-collapse .arco-menu-icon),
:deep(.arco-menu-collapse .arco-sub-menu-icon) {
  margin: 0 !important;
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
  width: 100% !important;
  height: 100% !important;
}

/* 收起状态下的悬停效果 */
:deep(.arco-menu-collapse .arco-menu-item:hover),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-has-icon:hover),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop:hover),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop-header:hover),
:deep(.arco-menu-collapse .arco-sub-menu-title:hover),
:deep(.arco-menu-collapse .arco-menu-inline-header:hover) {
  background-color: rgba(var(--primary-6), 0.1) !important;
}

/* 收起状态下的选中效果 */
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-selected),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-has-icon.arco-menu-selected),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop.arco-menu-selected),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop-header.arco-menu-selected),
:deep(.arco-menu-collapse .arco-sub-menu-title.arco-sub-menu-open),
:deep(.arco-menu-collapse .arco-menu-inline-header.arco-menu-selected) {
  background-color: rgba(var(--primary-6), 0.15) !important;
  color: rgb(var(--primary-6)) !important;
}

/* 强制重置收起状态下的所有菜单项样式 - 最高优先级 */
:deep(.arco-menu.arco-menu-collapse .arco-menu-item),
:deep(.arco-menu.arco-menu-collapse .arco-menu-item.arco-menu-has-icon),
:deep(.arco-menu.arco-menu-collapse .arco-menu-item.arco-menu-pop),
:deep(.arco-menu.arco-menu-collapse .arco-menu-item.arco-menu-pop-header) {
  margin: 2px 4px !important;
  padding: 0 !important;
  justify-content: center !important;
  display: flex !important;
  align-items: center !important;
  border-radius: 6px !important;
  height: 40px !important;
  width: calc(100% - 8px) !important;
  box-sizing: border-box !important;
}

:deep(.arco-menu.arco-menu-collapse .arco-menu-item .arco-menu-item-inner),
:deep(.arco-menu.arco-menu-collapse .arco-menu-item.arco-menu-has-icon .arco-menu-item-inner),
:deep(.arco-menu.arco-menu-collapse .arco-menu-item.arco-menu-pop .arco-menu-item-inner),
:deep(.arco-menu.arco-menu-collapse .arco-menu-item.arco-menu-pop-header .arco-menu-item-inner) {
  justify-content: center !important;
  padding: 0 !important;
  width: 100% !important;
  display: flex !important;
  align-items: center !important;
}

/* 通用的收起状态菜单项样式 - 覆盖所有可能的类组合 */
:deep(.arco-menu-collapse [class*="arco-menu-item"]) {
  margin: 2px 4px !important;
  padding: 0 !important;
  justify-content: center !important;
  display: flex !important;
  align-items: center !important;
  border-radius: 6px !important;
  height: 40px !important;
  width: calc(100% - 8px) !important;
  box-sizing: border-box !important;
}

:deep(.arco-menu-collapse [class*="arco-menu-item"] .arco-menu-item-inner) {
  justify-content: center !important;
  padding: 0 !important;
  width: 100% !important;
  display: flex !important;
  align-items: center !important;
}

/* 最强制的样式 - 使用更高的特异性 */
:deep(.arco-menu.arco-menu-collapse > .arco-menu-item),
:deep(.arco-menu.arco-menu-collapse > [class*="arco-menu-item"]) {
  margin: 2px 4px !important;
  padding: 0 !important;
  justify-content: center !important;
  display: flex !important;
  align-items: center !important;
  border-radius: 6px !important;
  height: 40px !important;
  width: calc(100% - 8px) !important;
  box-sizing: border-box !important;
  background: transparent !important;
}

:deep(.arco-menu.arco-menu-collapse > .arco-menu-item .arco-menu-item-inner),
:deep(.arco-menu.arco-menu-collapse > [class*="arco-menu-item"] .arco-menu-item-inner) {
  justify-content: center !important;
  padding: 0 !important;
  width: 100% !important;
  display: flex !important;
  align-items: center !important;
}

/* 终极解决方案 - 覆盖所有直接子元素 */
:deep(.arco-menu-collapse > *) {
  margin: 2px 4px !important;
  padding: 0 !important;
  justify-content: center !important;
  display: flex !important;
  align-items: center !important;
  border-radius: 6px !important;
  height: 40px !important;
  width: calc(100% - 8px) !important;
  box-sizing: border-box !important;
}

:deep(.arco-menu-collapse > * .arco-menu-item-inner) {
  justify-content: center !important;
  padding: 0 !important;
  width: 100% !important;
  display: flex !important;
  align-items: center !important;
}

/* 确保菜单项内容对齐 */
:deep(.arco-menu-item-inner) {
  display: flex !important;
  align-items: center !important;
  width: 100% !important;
  padding: 0 !important;
}

/* 特别处理单独的菜单项（如仪表盘）确保与子菜单标题对齐 */
:deep(.arco-menu-item:not(.arco-menu-indent-1):not(.arco-menu-indent-2)) {
  margin: 2px 8px !important;
  border-radius: 6px !important;
  height: 40px !important;
  line-height: 40px !important;
  padding: 0 16px !important;
}

:deep(.arco-menu-icon) {
  margin-right: 8px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  width: 16px !important;
  height: 16px !important;
}

:deep(.arco-sub-menu-icon) {
  margin-right: 8px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  width: 16px !important;
  height: 16px !important;
}

:deep(.arco-icon) {
  font-size: 16px;
}
</style>
