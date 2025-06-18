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
      <a-sub-menu v-if="menu.children?.length" :key="menu.id">
        <template #icon>
          <component v-if="menu.icon && iconMap[menu.icon]" :is="iconMap[menu.icon]" />
        </template>
        <template #title>{{ menu.name }}</template>
        <template v-for="child in menu.children" :key="child.id">
          <a-menu-item
            v-if="!child.children?.length"
            :key="child.id"
          >
            <template #icon>
              <component v-if="child.icon && iconMap[child.icon]" :is="iconMap[child.icon]" />
            </template>
            {{ child.name }}
          </a-menu-item>
          <a-sub-menu v-else :key="child.id">
            <template #icon>
              <component v-if="child.icon && iconMap[child.icon]" :is="iconMap[child.icon]" />
            </template>
            <template #title>{{ child.name }}</template>
            <a-menu-item
              v-for="item in child.children"
              :key="item.id"
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
      <a-menu-item v-else :key="menu.id">
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

// 获取树形菜单数据
const menuTree = computed(() => {
  const tree = menuStore.menuTree;
  // 如果只有一个根节点且是"首页"，则直接返回其子菜单
  if (tree.length === 1 && tree[0].code === 'root' && tree[0].children) {
    return tree[0].children;
  }
  return tree;
});

// 处理菜单点击
const handleMenuClick = async (key: string) => {
  console.log('菜单点击事件触发, key:', key, 'key类型:', typeof key);

  // 标记为手动点击
  isManualClick.value = true;

  // 使用完整的菜单列表查找，而不是树形结构
  const fullMenuList = menuStore.flatMenuList;

  // 将 key 转换为数字进行比较
  const keyNum = parseInt(key.toString());
  console.log('转换后的key:', keyNum);

  const menu = fullMenuList.find(m => m.id === keyNum);
  console.log('找到的菜单:', menu);

  if (menu && menu.path) {
    // 在跳转前，先保存当前的展开状态
    const currentParentKeys = getParentKeys(menuTree.value, key);
    console.log('点击前保存父级keys:', currentParentKeys);

    // 构建完整路径
    const fullPath = buildFullPath(menu);
    console.log('点击菜单:', menu.name, '原始路径:', menu.path, '构建路径:', fullPath);

    // 跳转路由
    await router.push(fullPath);

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

// 处理子菜单点击
const handleSubMenuClick = (key: string) => {
  const index = openKeys.value.indexOf(key);
  if (index > -1) {
    openKeys.value.splice(index, 1);
  } else {
    openKeys.value.push(key);
  }
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
    console.log('当前路由路径:', path, '是否手动点击:', isManualClick.value);

    // 如果是手动点击菜单导致的路由变化，不重新设置展开状态
    if (isManualClick.value) {
      console.log('手动点击菜单，跳过展开状态重置');
      return;
    }

    // 简化匹配逻辑：根据路径直接匹配
    let matchedMenu = null;

    // 遍历所有菜单项，找到路径匹配的
    const findMatchingMenu = (menus: any[]): any => {
      for (const menu of menus) {
        const constructedPath = buildFullPath(menu);
        console.log('检查菜单:', menu.name, '构建路径:', constructedPath);

        if (constructedPath === path) {
          return menu;
        }

        if (menu.children?.length) {
          const found = findMatchingMenu(menu.children);
          if (found) return found;
        }
      }
      return null;
    };

    matchedMenu = findMatchingMenu(menuTree.value);

    if (matchedMenu) {
      const key = matchedMenu.id.toString();
      const parentKeys = getParentKeys(menuTree.value, key);
      selectedKeys.value = [key];
      openKeys.value = parentKeys;
      console.log('匹配到菜单:', matchedMenu.name, 'key:', key, '父级keys:', parentKeys);
      console.log('当前展开的菜单:', openKeys.value);
    } else {
      console.log('未找到匹配的菜单');
    }
  },
  { immediate: true }
);


</script>

<style scoped>
:deep(.arco-menu-inline-header) {
  height: 40px;
  line-height: 40px;
}

:deep(.arco-menu-item) {
  height: 40px;
  line-height: 40px;
}

:deep(.arco-icon) {
  font-size: 16px;
}
</style>
