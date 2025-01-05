<template>
  <div class="dashboard">
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card class="stat-card">
          <template #title>用户总数</template>
          <div class="stat-value">
            <icon-user class="stat-icon" />
            <span>{{ stats.userCount }}</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <template #title>角色总数</template>
          <div class="stat-value">
            <icon-user-group class="stat-icon" />
            <span>{{ stats.roleCount }}</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <template #title>权限总数</template>
          <div class="stat-value">
            <icon-safe class="stat-icon" />
            <span>{{ stats.permissionCount }}</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <template #title>在线用户</template>
          <div class="stat-value">
            <icon-user-add class="stat-icon" />
            <span>{{ stats.onlineCount }}</span>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="12">
        <a-card title="最近登录">
          <a-table :data="recentLogins" :pagination="false">
            <template #columns>
              <a-table-column title="用户名" data-index="username" />
              <a-table-column title="登录时间" data-index="loginTime" />
              <a-table-column title="IP地址" data-index="ip" />
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="系统信息">
          <a-descriptions :data="systemInfo" layout="inline-horizontal" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';

interface Stats {
  userCount: number;
  roleCount: number;
  permissionCount: number;
  onlineCount: number;
}

interface RecentLogin {
  username: string;
  loginTime: string;
  ip: string;
}

const stats = ref<Stats>({
  userCount: 0,
  roleCount: 0,
  permissionCount: 0,
  onlineCount: 0,
});

const recentLogins = ref<RecentLogin[]>([]);

const systemInfo = [
  {
    label: '系统名称',
    value: 'DataX Admin',
  },
  {
    label: '系统版本',
    value: 'v1.0.0',
  },
  {
    label: '服务器操作系统',
    value: 'CentOS 7.9',
  },
  {
    label: 'Go 版本',
    value: 'go1.22.2',
  },
  {
    label: '数据库版本',
    value: 'MySQL 8.0',
  },
];

// TODO: 从后端获取统计数据
const fetchStats = async () => {
  // 模拟数据
  stats.value = {
    userCount: 100,
    roleCount: 10,
    permissionCount: 50,
    onlineCount: 5,
  };
};

// TODO: 从后端获取最近登录记录
const fetchRecentLogins = async () => {
  // 模拟数据
  recentLogins.value = [
    {
      username: 'admin',
      loginTime: '2024-01-05 23:30:00',
      ip: '127.0.0.1',
    },
    {
      username: 'user1',
      loginTime: '2024-01-05 23:15:00',
      ip: '192.168.1.100',
    },
  ];
};

onMounted(() => {
  fetchStats();
  fetchRecentLogins();
});
</script>

<style scoped>
.dashboard {
  padding: 16px;
}

.stat-card {
  background-color: var(--color-bg-2);
}

.stat-value {
  display: flex;
  align-items: center;
  font-size: 24px;
  font-weight: 500;
}

.stat-icon {
  font-size: 32px;
  margin-right: 16px;
  color: var(--color-primary-light-4);
}
</style>
