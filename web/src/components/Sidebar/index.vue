<template>
  <div class="sidebar" :class="{ collapsed: isCollapsed, 'dark-mode': isDarkMode }">
    <!-- 主导航区域 -->
    <div class="navigation-section">
      <div class="nav-title" v-if="!isCollapsed">Main</div>

      <div class="nav-items">
        <template v-for="menu in menuTree" :key="menu.id">
          <!-- 普通菜单项 -->
          <div
            v-if="!menu.children?.length"
            class="nav-item"
            :class="{ active: isMenuActive(menu) }"
            @click="handleMenuClick(menu)"
          >
            <div class="nav-icon">
              <!-- 直接使用后端返回的 Arco Design 图标 -->
              <component v-if="menu.icon && iconMap[menu.icon]" :is="iconMap[menu.icon]" />
              <!-- 默认图标 -->
              <IconApps v-else />
            </div>
            <span class="nav-label" v-if="!isCollapsed">{{ menu.name }}</span>
          </div>

          <!-- 有子菜单的项目 -->
          <div v-else class="nav-item-group">
            <div
              class="nav-item parent-item"
              :class="{
                active: isMenuActive(menu),
                expanded: expandedMenus.includes(menu.id)
              }"
              @click="toggleSubmenu(menu)"
            >
              <div class="nav-icon">
                <!-- 直接使用后端返回的 Arco Design 图标 -->
                <component v-if="menu.icon && iconMap[menu.icon]" :is="iconMap[menu.icon]" />
                <!-- 默认图标 -->
                <IconApps v-else />
              </div>
              <span class="nav-label" v-if="!isCollapsed">{{ menu.name }}</span>
            </div>

            <!-- 子菜单 -->
            <div
              v-if="!isCollapsed && expandedMenus.includes(menu.id)"
              class="submenu-container"
            >
              <div class="submenu-items">
                <div
                  v-for="child in menu.children"
                  :key="child.id"
                  class="submenu-item"
                  :class="{ active: isMenuActive(child) }"
                  @click="handleMenuClick(child)"
                >
                  <span class="submenu-label">{{ child.name }}</span>
                  <div class="submenu-radius"></div>
                </div>
              </div>
              <div class="submenu-line"></div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, inject } from 'vue';
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
  IconLink,
  IconSettings,
} from '@arco-design/web-vue/es/icon';

const route = useRoute();
const router = useRouter();
const menuStore = useMenuStore();

// 从父组件注入的状态
const collapsed = inject('collapsed', ref(false));
const isDarkMode = inject('isDarkMode', ref(false));
const isCollapsed = computed(() => collapsed.value);

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
  'icon-link': IconLink,
  'icon-settings': IconSettings,
};

// 展开的菜单
const expandedMenus = ref<number[]>([]);

// 获取菜单树
const menuTree = computed(() => {
  const tree = menuStore.menuTree;
  if (tree.length === 1 && tree[0].code === 'root' && tree[0].children) {
    return tree[0].children;
  }
  return tree;
});

// 移除了 getMenuIcon 函数，直接使用后端返回的 menu.icon 字段

// 构建完整路径
const buildFullPath = (menu: any): string => {
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

// 根据ID查找菜单
const findMenuById = (menus: any[], id: number): any => {
  for (const menu of menus) {
    if (menu.id === id) {
      return menu;
    }
  }
  return null;
};

// 判断菜单是否激活
const isMenuActive = (menu: any) => {
  const fullPath = buildFullPath(menu);
  return route.path === fullPath || route.path.startsWith(fullPath + '/');
};

// 切换子菜单展开状态
const toggleSubmenu = (menu: any) => {
  const index = expandedMenus.value.indexOf(menu.id);
  if (index > -1) {
    expandedMenus.value.splice(index, 1);
  } else {
    expandedMenus.value.push(menu.id);
  }
};

// 处理菜单点击
const handleMenuClick = async (menu: any) => {
  console.log('菜单点击事件触发, menu:', menu.name, 'path:', menu.path);

  // 检查是否为外部链接
  if (menu.is_external && menu.external_url) {
    // 根据open_type决定打开方式
    if (menu.open_type === 0) {
      // 内嵌显示
      await router.push({
        path: '/external-iframe',
        query: {
          url: menu.external_url,
        }
      });
    } else {
      // 新窗口打开（默认）
      window.open(menu.external_url, '_blank', 'noopener,noreferrer');
    }
  } else if (menu.path) {
    // 构建完整路径
    const fullPath = buildFullPath(menu);
    console.log('点击菜单:', menu.name, '原始路径:', menu.path, '构建路径:', fullPath);

    // 跳转路由
    await router.push(fullPath);
  } else {
    console.log('菜单路径不存在, menu:', menu);
  }
};


// 切换收起状态
const toggleCollapse = () => {
  collapsed.value = !collapsed.value;
};

// 组件挂载时初始化
onMounted(async () => {
  // 获取菜单数据
  await menuStore.fetchUserMenus();
});
</script>

<style lang="less" scoped>
.sidebar {
  width: 220px;
  height: 100vh;
  background: #ffffff;
  padding: 24px 16px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  position: relative;
  border-right: 1px solid #e1e4e8;

  &.collapsed {
    width: 64px;
    padding: 16px 8px;
  }
}

.navigation-section {
  flex: 1;
}

.nav-title {
  font-size: 12px;
  color: #8b949e;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 16px;
  padding: 0 8px;
}

.nav-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  min-height: 44px;
  box-sizing: border-box;

  &:hover {
    background: #f6f8fa;
  }

  &.active {
    background: #e7f3ff;
    color: #0969da;

    .nav-icon img {
      filter: brightness(0) saturate(100%) invert(27%) sepia(51%) saturate(2878%) hue-rotate(346deg) brightness(104%) contrast(97%);
    }
  }

  &.parent-item {
    &.expanded .nav-chevron img {
      transform: rotate(90deg);
    }
  }

  &.logout-item {
    color: #d1242f;

    &:hover {
      background: #ffebee;
    }

    .nav-icon img {
      filter: brightness(0) saturate(100%) invert(15%) sepia(94%) saturate(7482%) hue-rotate(3deg) brightness(90%) contrast(96%);
    }
  }

  .nav-icon {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    flex-shrink: 0;

    img {
      width: 16px;
      height: 16px;
      transition: all 0.2s ease;
    }
  }

  .nav-label {
    flex: 1;
    font-size: 14px;
    font-weight: 500;
    line-height: 1.2;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .nav-chevron {
    width: 16px;
    height: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 8px;

    img {
      width: 10px;
      height: 6px;
      transition: transform 0.2s ease;
      opacity: 0.6;
    }
  }
}

.submenu-container {
  margin-top: 4px;
  position: relative;
}

.submenu-items {
  padding-left: 32px;
  position: relative;
}

.submenu-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  min-height: 36px;

  &:hover {
    background: #f6f8fa;
    transform: translateX(2px);
  }

  &.active {
    background: #e7f3ff;
    color: #0969da;
    font-weight: 500;
  }

  .submenu-label {
    font-size: 13px;
    font-weight: 400;
    line-height: 1.2;
  }

  .submenu-radius {
    position: absolute;
    left: -16px;
    top: 50%;
    transform: translateY(-50%);
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background: #8b949e;
    opacity: 0.6;
    transition: all 0.2s ease;
  }

  &:hover .submenu-radius,
  &.active .submenu-radius {
    background: #0969da;
    width: 5px;
    height: 5px;
    opacity: 1;
  }
}

.submenu-line {
  position: absolute;
  left: 16px;
  top: 0;
  bottom: 0;
  width: 1px;
  background: #e1e4e8;
}

.bottom-section {
  margin-top: auto;
}

.collapse-button {
  position: absolute;
  right: -12px;
  top: 50%;
  transform: translateY(-50%);
  width: 24px;
  height: 24px;
  background: #ffffff;
  border: 1px solid #e1e4e8;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  &:hover {
    background: #f6f8fa;
    border-color: #d1d9e0;
  }

  .collapse-icon {
    width: 12px;
    height: 12px;
    display: flex;
    align-items: center;
    justify-content: center;

    img {
      width: 8px;
      height: 5px;
      transition: transform 0.2s ease;
    }
  }

  .collapsed & .collapse-icon img {
    transform: rotate(180deg);
  }
}


.sidebar.collapsed {
  .nav-item {
    justify-content: center;
    padding: 12px;

    .nav-icon {
      margin-right: 0;
    }
  }

  .collapse-button .collapse-icon img {
    transform: rotate(180deg);
  }
}

// 暗色主题适配
.sidebar.dark-mode {
  background: linear-gradient(180deg, rgba(30, 30, 30, 0.98), rgba(20, 20, 20, 0.95));
  border-right: 1px solid #30363d;

  .nav-title {
    color: #8b949e;
  }

  .nav-item {
    color: #e6edf3;

    &:hover {
      background: #21262d;
    }

    &.active {
      background: #1f2937;
      color: #58a6ff;
    }

    &.logout-item {
      color: #f85149;

      &:hover {
        background: #2d1617;
      }
    }
  }

  .submenu-item {
    color: #e6edf3;

    &:hover {
      background: #21262d;
    }

    &.active {
      background: #1f2937;
      color: #58a6ff;
    }

    .submenu-radius {
      background: #8b949e;
    }

    &:hover .submenu-radius,
    &.active .submenu-radius {
      background: #58a6ff;
    }
  }

  .submenu-line {
    background: #30363d;
  }

  .collapse-button {
    background: #21262d;
    border-color: #30363d;

    &:hover {
      background: #30363d;
      border-color: #484f58;
    }
  }
}
</style>
