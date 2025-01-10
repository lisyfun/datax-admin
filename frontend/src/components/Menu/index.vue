<template>
  <a-menu
    :selected-keys="selectedKeys"
    :open-keys="openKeys"
    show-collapse-button
    :auto-open="false"
    @menu-item-click="handleMenuClick"
    @sub-menu-click="handleSubMenuClick"
  >
    <template v-for="menu in menuTree" :key="menu.id">
      <!-- 有子菜单的情况 -->
      <a-sub-menu v-if="menu.children?.length" :key="menu.id">
        <template #icon>
          <icon-font v-if="menu.icon" :type="menu.icon" />
        </template>
        <template #title>{{ menu.name }}</template>
        <template v-for="child in menu.children" :key="child.id">
          <a-menu-item
            v-if="!child.children?.length"
            :key="child.id"
          >
            <template #icon>
              <icon-font v-if="child.icon" :type="child.icon" />
            </template>
            {{ child.name }}
          </a-menu-item>
          <a-sub-menu v-else :key="child.id">
            <template #icon>
              <icon-font v-if="child.icon" :type="child.icon" />
            </template>
            <template #title>{{ child.name }}</template>
            <a-menu-item
              v-for="item in child.children"
              :key="item.id"
            >
              <template #icon>
                <icon-font v-if="item.icon" :type="item.icon" />
              </template>
              {{ item.name }}
            </a-menu-item>
          </a-sub-menu>
        </template>
      </a-sub-menu>
      <!-- 没有子菜单的情况 -->
      <a-menu-item v-else :key="menu.id">
        <template #icon>
          <icon-font v-if="menu.icon" :type="menu.icon" />
        </template>
        {{ menu.name }}
      </a-menu-item>
    </template>
  </a-menu>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useMenuStore } from '@/stores/modules/menu';

const route = useRoute();
const router = useRouter();
const menuStore = useMenuStore();

// 当前选中的菜单项
const selectedKeys = ref<string[]>([]);
// 当前展开的子菜单
const openKeys = ref<string[]>([]);

// 获取树形菜单数据
const menuTree = computed(() => menuStore.menuTree);

// 处理菜单点击
const handleMenuClick = (key: string) => {
  const menu = findMenuByKey(menuTree.value, key);
  if (menu && menu.path) {
    router.push(menu.path);
  }
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
  for (const menu of menus) {
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

// 监听路由变化，更新选中的菜单项
watch(
  () => route.path,
  (path) => {
    const menu = findMenuByPath(menuTree.value, path);
    if (menu) {
      const key = menu.id.toString();
      selectedKeys.value = [key];
      openKeys.value = getParentKeys(menuTree.value, key);
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
