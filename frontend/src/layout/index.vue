<template>
  <div class="layout">
    <a-layout style="height: 100%">
      <a-layout-sider
        :collapsed="collapsed"
        :collapsible="true"
        :width="220"
        breakpoint="xl"
        @collapse="handleCollapse"
      >
        <div class="logo">
          <img src="../assets/logo.svg" alt="logo" />
          <h1 v-show="!collapsed">DataX Admin</h1>
        </div>
        <a-menu
          :selected-keys="selectedKeys"
          :open-keys="openKeys"
          show-collapse-button
          :style="{ width: '100%' }"
          @menu-item-click="handleMenuClick"
          @sub-menu-click="handleSubMenuClick"
        >
          <a-menu-item key="dashboard">
            <template #icon><icon-dashboard /></template>
            仪表盘
          </a-menu-item>
          <a-sub-menu key="system">
            <template #icon><icon-settings /></template>
            <template #title>系统管理</template>
            <a-menu-item key="system/permissions">
              <template #icon><icon-safe /></template>
              权限管理
            </a-menu-item>
            <a-menu-item key="system/roles">
              <template #icon><icon-user-group /></template>
              角色管理
            </a-menu-item>
            <a-menu-item key="system/users">
              <template #icon><icon-user /></template>
              用户管理
            </a-menu-item>
          </a-sub-menu>
        </a-menu>
      </a-layout-sider>
      <a-layout>
        <a-layout-header>
          <div class="header-content">
            <a-button
              class="collapse-btn"
              type="text"
              @click="collapsed = !collapsed"
            >
              <template #icon>
                <icon-menu-fold v-if="collapsed" />
                <icon-menu-unfold v-else />
              </template>
            </a-button>
            <div class="header-right">
              <a-dropdown trigger="click">
                <a-avatar :size="32">
                  <icon-user />
                </a-avatar>
                <template #content>
                  <a-doption>
                    <template #icon><icon-user /></template>
                    个人信息
                  </a-doption>
                  <a-doption>
                    <template #icon><icon-poweroff /></template>
                    退出登录
                  </a-doption>
                </template>
              </a-dropdown>
            </div>
          </div>
        </a-layout-header>
        <a-layout-content>
          <div class="content">
            <router-view></router-view>
          </div>
        </a-layout-content>
      </a-layout>
    </a-layout>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const collapsed = ref(false)
const selectedKeys = ref<string[]>([])
const openKeys = ref<string[]>([])

// 根据当前路由设置选中的菜单项
watch(
  () => route.path,
  (path) => {
    const paths = path.split('/').filter(Boolean)
    selectedKeys.value = [paths.join('/')]
    openKeys.value = paths.length > 1 ? [paths[0]] : []
  },
  { immediate: true }
)

const handleCollapse = (value: boolean) => {
  collapsed.value = value
}

const handleMenuClick = (key: string) => {
  router.push(`/${key}`)
}

const handleSubMenuClick = (key: string) => {
  const index = openKeys.value.indexOf(key)
  if (index > -1) {
    openKeys.value.splice(index, 1)
  } else {
    openKeys.value.push(key)
  }
}
</script>

<style scoped>
.layout {
  height: 100%;
}

.logo {
  height: 64px;
  padding: 16px;
  display: flex;
  align-items: center;
  background: var(--color-bg-2);
}

.logo img {
  width: 32px;
  height: 32px;
  margin-right: 12px;
}

.logo h1 {
  margin: 0;
  color: var(--color-text-1);
  font-weight: 500;
  font-size: 18px;
  line-height: 32px;
  white-space: nowrap;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 16px;
  background: var(--color-bg-2);
}

.collapse-btn {
  margin-right: 16px;
}

.header-right {
  display: flex;
  align-items: center;
}

.content {
  padding: 16px;
  background: var(--color-bg-2);
  height: 100%;
}
</style>
