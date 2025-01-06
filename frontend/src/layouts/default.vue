<template>
  <a-layout class="layout">
    <a-layout-header class="header">
      <div class="logo">
        <img src="../assets/logo.svg" alt="logo" />
        <h1>DataX Admin</h1>
      </div>
      <div class="header-right">
        <a-dropdown>
          <a-button>
            <template #icon><icon-user /></template>
            {{ userInfo.nickname || userInfo.username }}
          </a-button>
          <template #content>
            <a-doption @click="handleLogout">
              <template #icon><icon-export /></template>
              退出登录
            </a-doption>
          </template>
        </a-dropdown>
      </div>
    </a-layout-header>
    <a-layout>
      <a-layout-sider
        collapsible
        :collapsed="collapsed"
        :width="200"
        :style="{ background: 'var(--color-bg-2)' }"
      >
        <a-menu
          :selected-keys="selectedKeys"
          :default-open-keys="openKeys"
          :style="{ width: '100%' }"
          @menu-item-click="handleMenuClick"
        >
          <a-menu-item key="/dashboard">
            <template #icon><icon-dashboard /></template>
            仪表盘
          </a-menu-item>
          <a-sub-menu key="system">
            <template #icon><icon-settings /></template>
            <template #title>系统管理</template>
            <a-menu-item key="/system/users">
              <template #icon><icon-user /></template>
              用户管理
            </a-menu-item>
            <a-menu-item key="/system/roles">
              <template #icon><icon-user-group /></template>
              角色管理
            </a-menu-item>
            <a-menu-item key="/system/permissions">
              <template #icon><icon-safe /></template>
              权限管理
            </a-menu-item>
          </a-sub-menu>
        </a-menu>
      </a-layout-sider>
      <a-layout-content>
        <router-view></router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import * as userApi from '@/api/user';

const router = useRouter();
const route = useRoute();
const collapsed = ref(false);

// 用户信息
const userInfo = ref({
  username: '',
  nickname: '',
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
const handleLogout = async () => {
  try {
    await userApi.logout();
    Message.success('退出登录成功');
    router.push('/login');
  } catch (error: any) {
    Message.error('退出登录失败');
  }
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
});
</script>

<style scoped>
.layout {
  height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border);
}

.logo {
  display: flex;
  align-items: center;
}

.logo img {
  width: 32px;
  height: 32px;
  margin-right: 8px;
}

.logo h1 {
  margin: 0;
  color: var(--color-text-1);
  font-size: 18px;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
}

:deep(.arco-layout-sider-children) {
  overflow-y: auto;
}

.arco-layout-content {
  overflow-y: auto;
  background: var(--color-fill-2);
}
</style>
