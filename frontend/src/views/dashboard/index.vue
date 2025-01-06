<template>
  <div class="dashboard">
    <a-card>
      <template #title>仪表盘</template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-card class="stat-card" :bordered="false">
            <div class="stat-header">
              <div class="stat-title">用户总数</div>
              <icon-user class="stat-icon" />
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.userCount }}</div>
              <div class="stat-label">总用户数量</div>
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card" :bordered="false">
            <div class="stat-header">
              <div class="stat-title">角色总数</div>
              <icon-user-group class="stat-icon" />
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.roleCount }}</div>
              <div class="stat-label">系统角色数量</div>
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card" :bordered="false">
            <div class="stat-header">
              <div class="stat-title">权限总数</div>
              <icon-safe class="stat-icon" />
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.permissionCount }}</div>
              <div class="stat-label">系统权限数量</div>
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card class="stat-card" :bordered="false">
            <div class="stat-header">
              <div class="stat-title">在线用户</div>
              <icon-user-add class="stat-icon" />
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.onlineCount }}</div>
              <div class="stat-label">当前在线用户</div>
            </div>
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16" style="margin-top: 16px">
        <a-col :span="12">
          <a-card class="list-card" :bordered="false">
            <template #title>
              <div class="card-title">
                <icon-history class="card-title-icon" />
                最近登录
              </div>
            </template>
            <a-table :data="recentLogins" :pagination="false" :bordered="false">
              <template #columns>
                <a-table-column title="用户名" data-index="username">
                  <template #cell="{ record }">
                    <a-space>
                      <icon-user />
                      {{ record.username }}
                    </a-space>
                  </template>
                </a-table-column>
                <a-table-column title="登录时间" data-index="loginTime">
                  <template #cell="{ record }">
                    <a-space>
                      <icon-calendar />
                      {{ record.loginTime }}
                    </a-space>
                  </template>
                </a-table-column>
                <a-table-column title="IP地址" data-index="ip">
                  <template #cell="{ record }">
                    <a-space>
                      <icon-wifi />
                      {{ record.ip }}
                    </a-space>
                  </template>
                </a-table-column>
              </template>
            </a-table>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card class="info-card" :bordered="false">
            <template #title>
              <div class="card-title">
                <icon-computer class="card-title-icon" />
                系统信息
              </div>
            </template>
            <a-descriptions :data="systemInfo" layout="inline-horizontal" bordered />
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import {
  IconUser,
  IconUserGroup,
  IconSafe,
  IconUserAdd,
  IconHistory,
  IconCalendar,
  IconWifi,
  IconComputer,
} from '@arco-design/web-vue/es/icon';

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
    value: 'DATAX ADMIN',
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
  /* background-color: var(--color-fill-2); */
  min-height: 100%;
}

.stat-card {
  height: 140px;
  background-color: var(--color-bg-1);
  transition: all 0.3s;
  border-radius: 8px;
  position: relative;
  overflow: hidden;
}

.stat-card::after {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  width: 80px;
  height: 80px;
  background: linear-gradient(45deg, transparent, rgba(var(--primary-6), 0.1));
  border-radius: 0 8px 0 50%;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  background-color: var(--color-bg-1);
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  position: relative;
  z-index: 1;
}

.stat-title {
  font-size: 16px;
  color: var(--color-text-1);
  font-weight: 500;
}

.stat-icon {
  font-size: 24px;
  color: rgb(var(--primary-6));
  background-color: rgba(var(--primary-6), 0.1);
  padding: 8px;
  border-radius: 8px;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
  position: relative;
  z-index: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: var(--color-text-1);
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: var(--color-text-3);
}

.list-card, .info-card {
  height: 100%;
  background-color: var(--color-bg-1);
  border-radius: 8px;
}

.card-title {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 500;
  padding: 4px 0;
}

.card-title-icon {
  margin-right: 8px;
  font-size: 18px;
  color: rgb(var(--primary-6));
  background-color: rgba(var(--primary-6), 0.1);
  padding: 6px;
  border-radius: 6px;
}

:deep(.arco-descriptions-item-label) {
  font-weight: 500;
  color: var(--color-text-2);
}

:deep(.arco-descriptions-item-value) {
  color: var(--color-text-1);
}

:deep(.arco-table-th) {
  background-color: var(--color-fill-2) !important;
  font-weight: 500;
}

:deep(.arco-table-tr:hover) {
  background-color: var(--color-fill-2) !important;
}

:deep(.arco-descriptions) {
  background-color: var(--color-bg-1);
}

:deep(.arco-descriptions-bordered) {
  border-radius: 8px;
  border: 1px solid var(--color-neutral-3);
}

:deep(.arco-descriptions-bordered .arco-descriptions-item) {
  padding: 12px 16px;
}

:deep(.arco-card) {
  border: 1px solid var(--color-neutral-3);
}

:deep(.arco-table-container) {
  border-radius: 8px;
}

:deep(.arco-table-cell) {
  padding: 12px 16px;
}
</style>
