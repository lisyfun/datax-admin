<template>
  <a-layout class="layout">
    <a-layout-header class="header">
      <div class="header-left">
        <div class="logo">
          <img src="../assets/logo.svg" alt="logo" />
          <h1>数据集成平台</h1>
        </div>
        <a-button class="menu-trigger" type="text" @click="collapsed = !collapsed">
          <template #icon>
            <icon-menu-fold v-if="collapsed" />
            <icon-menu-unfold v-else />
          </template>
        </a-button>
        <a-breadcrumb class="breadcrumb">
          <template v-for="item in breadcrumbItems" :key="item">
            <a-breadcrumb-item>{{ item }}</a-breadcrumb-item>
          </template>
        </a-breadcrumb>
      </div>
      <div class="header-right">
        <a-space :size="16">
          <a-button class="action-btn" type="text" @click="toggleTheme">
            <template #icon>
              <icon-sun-fill v-if="isDarkMode" />
              <icon-moon-fill v-else />
            </template>
          </a-button>
          <a-button class="action-btn" type="text" @click="handleRefresh">
            <template #icon><icon-refresh /></template>
          </a-button>
          <a-button class="action-btn" type="text" @click="toggleFullscreen">
            <template #icon>
              <icon-fullscreen-exit v-if="isFullscreen" />
              <icon-fullscreen v-else />
            </template>
          </a-button>
          <a-dropdown trigger="click">
            <a-space :size="8" class="user-dropdown">
              <a-avatar :size="32" class="user-avatar">
                <img v-if="userInfo.avatar" :src="userInfo.avatar" alt="avatar" />
                <icon-user v-else />
              </a-avatar>
              <span class="username">{{ userInfo.nickname || userInfo.username }}</span>
              <icon-down />
            </a-space>
            <template #content>
              <a-doption @click="handleLogout">
                <template #icon><icon-export /></template>
                退出登录
              </a-doption>
            </template>
          </a-dropdown>
        </a-space>
      </div>
    </a-layout-header>
    <a-layout>
      <a-layout-sider
        class="layout-sider"
        :collapsed="collapsed"
        :width="200"
        :collapsed-width="50"
        :hide-trigger="true"
        breakpoint="xxl"
      >
        <Menu />
      </a-layout-sider>
      <a-layout-content class="layout-content">
        <div class="content-wrapper">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted, provide } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { Message, Modal } from '@arco-design/web-vue';
import {
  IconUser,
  IconDown,
  IconRefresh,
  IconFullscreen,
  IconFullscreenExit,
  IconMenuFold,
  IconMenuUnfold,
  IconExport,
  IconMoonFill,
  IconSunFill,
} from '@arco-design/web-vue/es/icon';
import * as userApi from '@/api/user';
import Menu from '@/components/Menu/index.vue';
import { useMenuStore } from '@/stores/menu';
import { usePermissionStore } from '@/stores/permission';

const router = useRouter();
const route = useRoute();
const menuStore = useMenuStore();
const permissionStore = usePermissionStore();
const collapsed = ref(false);
const isDarkMode = ref(false);
const isFullscreen = ref(false);
const isRefreshing = ref(false);

// 提供collapsed状态给子组件
provide('collapsed', collapsed);



// 初始化主题
const initTheme = () => {
  const theme = localStorage.getItem('theme') || 'light';
  isDarkMode.value = theme === 'dark';
  applyTheme(isDarkMode.value);
};

// 切换主题
const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  applyTheme(isDarkMode.value);
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light');
};

// 应用主题
const applyTheme = (dark: boolean) => {
  if (dark) {
    document.body.setAttribute('arco-theme', 'dark');
  } else {
    document.body.removeAttribute('arco-theme');
  }
};

// 切换全屏
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen();
    isFullscreen.value = true;
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen();
      isFullscreen.value = false;
    }
  }
};

// 监听全屏变化
const handleFullscreenChange = () => {
  isFullscreen.value = !!document.fullscreenElement;
};

// 用户信息
const userInfo = ref({
  username: '',
  nickname: '',
  avatar: '',
});

// 面包屑导航
const breadcrumbItems = computed(() => {
  const matched = route.matched;
  return matched.map(record => record.meta?.title || record.name);
});

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const res = await userApi.getUserInfo();
    userInfo.value = res.data;
  } catch (error: any) {
    Message.error('获取用户信息失败');
  }
};

// 退出登录
const handleLogout = () => {
  Modal.confirm({
    title: '确认退出',
    content: '确定要退出登录吗？',
    onOk: async () => {
      try {
        await userApi.logout();
        Message.success('退出登录成功');
        router.push('/login');
      } catch (error: any) {
        Message.error('退出登录失败');
      }
    },
  });
};



// 处理页面刷新
const handleRefresh = async () => {
  if (isRefreshing.value) return;
  isRefreshing.value = true;

  try {
    // 触发当前页面的刷新方法
    const refreshEvent = new CustomEvent('page-refresh');
    window.dispatchEvent(refreshEvent);
    Message.success('刷新成功');
  } catch (error) {
    Message.error('刷新失败');
  } finally {
    isRefreshing.value = false;
  }
};

// 提供刷新方法给所有子组件
provide('triggerRefresh', handleRefresh);

onMounted(async () => {
  await fetchUserInfo();
  await menuStore.fetchUserMenus(); // 加载用户菜单
  await permissionStore.getUserPermissions(); // 加载用户权限
  initTheme();
  document.addEventListener('fullscreenchange', handleFullscreenChange);
});

// 组件卸载时移除事件监听
onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange);
});
</script>

<style scoped>
/* 全局字体设置 */
:root {
  --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
    'Noto Color Emoji';
  --font-size-base: 14px;
  --line-height-base: 1.5715;
}

/* 应用全局字体设置 */
.layout {
  height: 100vh;
  background-color: var(--color-fill-2);
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  line-height: var(--line-height-base);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  position: relative;
  overflow: hidden;
}

/* 禁用菜单文字选中 */
:deep(.arco-menu) {
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

:deep(.arco-menu-item),
:deep(.arco-sub-menu-title),
:deep(.arco-menu-title),
:deep(.arco-breadcrumb) {
  font-size: 14px;
  letter-spacing: 0.3px;
  height: 46px;
  line-height: 46px;
  font-weight: 500;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

:deep(.arco-menu-icon),
:deep(.arco-sub-menu-icon),
:deep(.arco-menu-item-icon) {
  font-size: 18px;
  stroke-width: 2.5;
  stroke: currentColor;
}

:deep(.arco-menu-selected) {
  font-weight: 600 !important;
  background-color: var(--color-fill-2) !important;
}

:deep(.arco-menu-item:hover),
:deep(.arco-sub-menu-title:hover) {
  background-color: var(--color-fill-2) !important;
}

:deep(.arco-menu-light .arco-menu-selected) {
  background-color: var(--color-fill-2) !important;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  background: var(--color-bg-2);
  box-shadow: 0 1px 4px 0 rgba(0, 21, 41, 0.08);
  z-index: 100;
  transition: all 0.2s ease;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  height: 64px;
  margin-right: 16px;
}

.logo img {
  width: 32px;
  height: 32px;
  margin-right: 8px;
  transition: all 0.2s ease;
}

.logo h1 {
  margin: 0;
  color: var(--color-text-1);
  font-size: 18px;
  font-weight: 600;
  white-space: nowrap;
  transition: all 0.2s ease;
  letter-spacing: 0.3px;
}

.menu-trigger {
  margin-right: 16px;
  font-size: 18px;
  transition: all 0.2s ease;
}

.menu-trigger:hover {
  color: rgb(var(--primary-6));
}

.breadcrumb {
  margin-left: 8px;
}

.header-right {
  display: flex;
  align-items: center;
}

.action-btn {
  font-size: 16px;
  transition: all 0.2s ease;
}

.action-btn:hover {
  color: rgb(var(--primary-6));
  background-color: var(--color-fill-3);
}

.user-dropdown {
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.user-dropdown:hover {
  background-color: var(--color-fill-3);
}

.user-avatar {
  transition: transform 0.2s ease;
}

.user-avatar:hover {
  transform: scale(1.1);
}

.username {
  color: var(--color-text-1);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.2px;
}

.layout-sider {
  box-shadow: 2px 0 8px 0 rgba(29, 35, 41, 0.05);
  background: var(--color-bg-2);
  z-index: 99;
  transition: all 0.2s ease;
  position: fixed;
  top: 64px;
  left: 0;
  bottom: 0;
  height: calc(100vh - 64px);
}

:deep(.arco-layout-sider-children) {
  overflow-y: auto;
  height: 100%;
  width: 100% !important;
}

:deep(.arco-menu) {
  width: 100% !important;
}

:deep(.arco-menu-collapse) {
  width: 100% !important;

  .arco-menu-item,
  .arco-menu-item-inner,
  .arco-sub-menu-title {
    padding: 0 !important;
    justify-content: center !important;
    width: 100%;
    display: flex;
    align-items: center;
  }

  .arco-menu-icon,
  .arco-sub-menu-icon {
    margin: 0 !important;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
  }

  .arco-menu-title {
    display: none;
  }

  .arco-sub-menu-suffix {
    display: none;
  }
}

/* 隐藏所有父菜单的箭头 */
:deep(.arco-sub-menu-suffix),
:deep(.arco-sub-menu-arrow),
:deep(.arco-icon-down),
:deep(.arco-icon-right) {
  display: none !important;
}


/* 基础菜单样式 - 统一所有菜单项 */
:deep(.arco-menu-item),
:deep(.arco-sub-menu-title),
:deep(.arco-menu-inline-header) {
  display: flex !important;
  align-items: center !important;
  height: 40px !important;
  line-height: 40px !important;
  text-align: left;
  padding: 0 16px !important;
  box-sizing: border-box !important;
}

/* 特别针对有图标的菜单项 */
:deep(.arco-menu-item.arco-menu-has-icon) {
  padding: 0 16px !important;
}

:deep(.arco-menu-item-inner) {
  width: 100% !important;
  display: flex !important;
  align-items: center !important;
  justify-content: flex-start !important;
  padding: 0 !important;
}

:deep(.arco-menu-title) {
  text-align: left;
  justify-content: flex-start;
}

/* 确保收起状态下图标完全居中 - 统一所有类型 */
:deep(.arco-menu-collapse .arco-menu-item),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-has-icon),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop),
:deep(.arco-menu-collapse .arco-menu-item.arco-menu-pop-header),
:deep(.arco-menu-collapse .arco-sub-menu-title),
:deep(.arco-menu-collapse .arco-menu-inline-header) {
  padding: 0 !important;
  margin: 2px 4px !important;
  justify-content: center !important;
  border-radius: 6px !important;
  height: 40px !important;
  width: calc(100% - 8px) !important;
  box-sizing: border-box !important;
  display: flex !important;
  align-items: center !important;
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

:deep(.arco-menu-collapse .arco-menu-icon),
:deep(.arco-menu-collapse .arco-sub-menu-icon) {
  margin: 0 !important;
  width: auto !important;
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
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

/* 确保所有菜单图标对齐 */
:deep(.arco-menu-icon),
:deep(.arco-sub-menu-icon) {
  margin-right: 8px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  width: 16px !important;
  height: 16px !important;
  flex-shrink: 0 !important;
}

.layout-content {
  padding: 16px;
  overflow: auto;
  background: var(--color-neutral-2);
  transition: all 0.2s ease;
  margin-left: 220px;
  margin-top: 64px;
  height: calc(100vh - 64px);
}

.content-wrapper {
  background: var(--color-bg-2);
  padding: 16px;
  border-radius: 4px;
  min-height: calc(100vh - 96px);
  box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
  font-size: 14px;
  line-height: 1.6;
  letter-spacing: 0.2px;
}

/* 当侧边栏收起时的样式 */
.layout-sider.collapsed + .layout-content {
  margin-left: 50px;
}

/* 路由切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 滚动条样式 */
:deep(.arco-layout-sider-children)::-webkit-scrollbar,
.layout-content::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

:deep(.arco-layout-sider-children)::-webkit-scrollbar-thumb,
.layout-content::-webkit-scrollbar-thumb {
  border-radius: 4px;
  background: var(--color-fill-4);
}

:deep(.arco-layout-sider-children)::-webkit-scrollbar-track,
.layout-content::-webkit-scrollbar-track {
  border-radius: 4px;
  background: var(--color-fill-2);
}

/* 暗黑模式下的样式调整 */
:deep([arco-theme='dark']) {
  .layout {
    background-color: var(--color-bg-1);
  }

  .header {
    background: var(--color-bg-1);
    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.2);
  }

  .layout-sider {
    background: linear-gradient(180deg,
      rgba(30, 30, 30, 0.98),
      rgba(20, 20, 20, 0.95));
    box-shadow: 4px 0 20px 0 rgba(0, 0, 0, 0.2),
                2px 0 8px 0 rgba(0, 0, 0, 0.1);
    border-right: 1px solid rgba(60, 60, 60, 0.3);
  }

  .layout-content {
    background: var(--color-neutral-1);
  }

  .content-wrapper {
    background: var(--color-bg-1);
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.15);
  }

  .action-btn:hover {
    background-color: var(--color-fill-2);
  }

  .user-dropdown:hover {
    background-color: var(--color-fill-2);
  }

  /* 滚动条暗色适配 */
  :deep(.arco-layout-sider-children)::-webkit-scrollbar-thumb,
  .layout-content::-webkit-scrollbar-thumb {
    background: linear-gradient(180deg,
      rgba(80, 80, 80, 0.8),
      rgba(100, 100, 100, 0.9));
    border: 1px solid rgba(120, 120, 120, 0.3);
  }

  :deep(.arco-layout-sider-children)::-webkit-scrollbar-thumb:hover,
  .layout-content::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(180deg,
      rgba(120, 120, 120, 0.9),
      rgba(140, 140, 140, 1));
  }

  :deep(.arco-layout-sider-children)::-webkit-scrollbar-track,
  .layout-content::-webkit-scrollbar-track {
    background: rgba(20, 20, 20, 0.5);
  }

  /* 菜单项暗色适配 */
  .menu-item {
    color: var(--color-text-2);

    &:hover {
      color: rgb(var(--primary-6));
      background-color: var(--color-fill-2);
    }
  }

  /* 面包屑暗色适配 */
  .breadcrumb {
    color: var(--color-text-2);
  }

  /* Logo 文字暗色适配 */
  .logo h1 {
    color: var(--color-text-1);
  }

  /* 用户名暗色适配 */
  .username {
    color: var(--color-text-1);
  }
}

/* 主题切换动画 */
.arco-icon {
  transition: transform 0.3s ease;
}

.action-btn:hover .arco-icon {
  transform: rotate(30deg);
}

/* 菜单项字体优化 */
:deep(.arco-menu-item) {
  font-size: 14px;
  letter-spacing: 0.2px;
}


/* 面包屑字体优化 */
:deep(.arco-breadcrumb-item) {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.2px;
}

/* 按钮文字优化 */
:deep(.arco-btn) {
  font-weight: 500;
  letter-spacing: 0.2px;
}

/* 标题文字优化 */
:deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
  font-weight: 600;
  letter-spacing: 0.3px;
  line-height: 1.4;
}

.layout-content {
  padding: 16px;
  overflow: auto;
  background: var(--color-neutral-2);
  transition: all 0.2s ease;
  margin-left: v-bind('collapsed ? "50px" : "200px"');
  margin-top: 64px;
  height: calc(100vh - 64px);
}

</style>
