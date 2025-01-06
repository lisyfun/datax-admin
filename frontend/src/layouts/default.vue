<template>
  <a-layout class="layout">
    <a-layout-header class="header">
      <div class="header-left">
        <div class="logo">
          <img src="../assets/logo.svg" alt="logo" />
          <h1>DataX Admin</h1>
        </div>
        <a-button class="menu-trigger" type="text" @click="collapsed = !collapsed">
          <template #icon>
            <icon-menu-fold v-if="collapsed" />
            <icon-menu-unfold v-else />
          </template>
        </a-button>
        <a-breadcrumb class="breadcrumb">
          <a-breadcrumb-item>首页</a-breadcrumb-item>
          <template v-for="(item, index) in breadcrumbItems" :key="index">
            <a-breadcrumb-item>{{ item }}</a-breadcrumb-item>
          </template>
        </a-breadcrumb>
      </div>
      <div class="header-right">
        <a-space :size="16">
          <a-button class="action-btn" type="text" @click="toggleTheme">
            <template #icon>
              <icon-moon-fill v-if="isDarkMode" />
              <icon-sun-fill v-else />
            </template>
          </a-button>
          <a-button class="action-btn" type="text">
            <template #icon><icon-notification /></template>
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
              <a-doption @click="handleProfile">
                <template #icon><icon-user /></template>
                个人信息
              </a-doption>
              <a-doption @click="handleSettings">
                <template #icon><icon-settings /></template>
                系统设置
              </a-doption>
              <a-divider style="margin: 4px 0" />
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
        :width="220"
        :hide-trigger="true"
        breakpoint="xl"
      >
        <a-menu
          :selected-keys="selectedKeys"
          :default-open-keys="openKeys"
          :style="{ width: '100%' }"
          @menu-item-click="handleMenuClick"
        >
          <a-menu-item key="/dashboard" class="menu-item">
            <template #icon><icon-home /></template>
            仪表盘
          </a-menu-item>
          <a-sub-menu key="system">
            <template #icon><icon-apps /></template>
            <template #title>系统管理</template>
            <a-menu-item key="/system/users" class="menu-item">
              <template #icon><icon-user /></template>
              用户管理
            </a-menu-item>
            <a-menu-item key="/system/roles" class="menu-item">
              <template #icon><icon-user-group /></template>
              角色管理
            </a-menu-item>
            <a-menu-item key="/system/permissions" class="menu-item">
              <template #icon><icon-safe /></template>
              权限管理
            </a-menu-item>
          </a-sub-menu>
          <a-sub-menu key="job">
            <template #icon><icon-calendar /></template>
            <template #title>任务管理</template>
            <a-menu-item key="/job/list" class="menu-item">
              <template #icon><icon-unordered-list /></template>
              任务列表
            </a-menu-item>
            <a-menu-item key="/job/history" class="menu-item">
              <template #icon><icon-clock-circle /></template>
              执行历史
            </a-menu-item>
          </a-sub-menu>
          <a-sub-menu key="monitor">
            <template #icon><icon-dashboard /></template>
            <template #title>系统监控</template>
            <a-menu-item key="/monitor/online" class="menu-item">
              <template #icon><icon-desktop /></template>
              在线用户
            </a-menu-item>
            <a-menu-item key="/monitor/server" class="menu-item">
              <template #icon><icon-cloud /></template>
              服务监控
            </a-menu-item>
            <a-menu-item key="/monitor/logs" class="menu-item">
              <template #icon><icon-file /></template>
              操作日志
            </a-menu-item>
          </a-sub-menu>
          <a-sub-menu key="system-tool">
            <template #icon><icon-bulb /></template>
            <template #title>系统工具</template>
            <a-menu-item key="/tool/generator" class="menu-item">
              <template #icon><icon-code /></template>
              代码生成
            </a-menu-item>
            <a-menu-item key="/tool/swagger" class="menu-item">
              <template #icon><icon-robot /></template>
              系统接口
            </a-menu-item>
          </a-sub-menu>
        </a-menu>
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
import { ref, computed, onMounted, watch, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { Message, Modal } from '@arco-design/web-vue';
import {
  IconHome,
  IconApps,
  IconUser,
  IconUserGroup,
  IconSafe,
  IconCalendar,
  IconUnorderedList,
  IconClockCircle,
  IconDashboard,
  IconDesktop,
  IconCloud,
  IconFile,
  IconBulb,
  IconCode,
  IconRobot,
  IconDown,
  IconNotification,
  IconFullscreen,
  IconFullscreenExit,
  IconMenuFold,
  IconMenuUnfold,
  IconExport,
  IconSettings,
  IconMoonFill,
  IconSunFill,
} from '@arco-design/web-vue/es/icon';
import * as userApi from '@/api/user';

const router = useRouter();
const route = useRoute();
const collapsed = ref(false);
const isDarkMode = ref(false);
const isFullscreen = ref(false);

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
  const paths = route.path.split('/').filter(Boolean);
  return paths.map(path => {
    switch (path) {
      case 'dashboard': return '仪表盘';
      case 'system': return '系统管理';
      case 'users': return '用户管理';
      case 'roles': return '角色管理';
      case 'permissions': return '权限管理';
      default: return path;
    }
  });
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

// 个人信息
const handleProfile = () => {
  router.push('/profile');
};

// 系统设置
const handleSettings = () => {
  router.push('/settings');
};

// 菜单相关
const selectedKeys = computed(() => [route.path]);
const openKeys = computed(() => {
  const path = route.path;
  if (path.startsWith('/system/')) {
    return ['system'];
  }
  return [];
});

const handleMenuClick = (key: string) => {
  router.push(key);
};

onMounted(() => {
  fetchUserInfo();
  initTheme();
  document.addEventListener('fullscreenchange', handleFullscreenChange);
});

// 组件卸载时移除事件监听
onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange);
});
</script>

<style scoped>
.layout {
  height: 100vh;
  background-color: var(--color-fill-2);
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
  font-weight: 500;
  white-space: nowrap;
  transition: all 0.2s ease;
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
}

.layout-sider {
  box-shadow: 2px 0 8px 0 rgba(29, 35, 41, 0.05);
  background: var(--color-bg-2);
  z-index: 99;
  transition: all 0.2s ease;
}

:deep(.arco-layout-sider-children) {
  overflow-y: auto;
  height: calc(100vh - 64px);
}

.layout-content {
  padding: 16px;
  overflow: auto;
  background: var(--color-neutral-2);
  transition: all 0.2s ease;
}

.content-wrapper {
  background: var(--color-bg-2);
  padding: 16px;
  border-radius: 4px;
  min-height: calc(100vh - 96px);
  box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
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
    background: var(--color-bg-1);
    box-shadow: 2px 0 8px 0 rgba(0, 0, 0, 0.15);
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
    background: var(--color-fill-3);
  }

  :deep(.arco-layout-sider-children)::-webkit-scrollbar-track,
  .layout-content::-webkit-scrollbar-track {
    background: var(--color-bg-2);
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
</style>
