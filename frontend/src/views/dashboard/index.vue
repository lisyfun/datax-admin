<template>
  <div class="dashboard">
    <a-card>
      <template #title>仪表盘</template>
      <a-row type="flex" :gutter="20">
        <a-col :span="4">
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
        <a-col :span="4">
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
        <a-col :span="4">
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
        <a-col :span="5">
          <a-card class="stat-card" :bordered="false">
            <div class="stat-header">
              <div class="stat-title">成功任务</div>
              <icon-check-circle class="stat-icon success" />
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.successCount }}</div>
              <div class="stat-label">成功执行任务</div>
            </div>
          </a-card>
        </a-col>
        <a-col :span="5">
          <a-card class="stat-card" :bordered="false">
            <div class="stat-header">
              <div class="stat-title">失败任务</div>
              <icon-close-circle class="stat-icon error" />
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.failedCount }}</div>
              <div class="stat-label">失败执行任务</div>
            </div>
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16" style="margin-top: 16px">
        <a-col :span="12">
          <a-card class="chart-card" :bordered="false">
            <template #title>
              <div class="card-title">
                <icon-bar-chart class="card-title-icon" />
                任务执行趋势
              </div>
            </template>
            <div ref="chartRef" class="chart-container"></div>
          </a-card>
        </a-col>
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
                      {{ formatDateTime(record.loginTime) }}
                    </a-space>
                  </template>
                </a-table-column>
                <a-table-column title="IP地址" data-index="ip">
                  <template #cell="{ record }">
                    <a-space>
                      <icon-wifi />
                      <a-tooltip content="登录IP">
                        {{ record.ip }}
                      </a-tooltip>
                    </a-space>
                  </template>
                </a-table-column>
              </template>
            </a-table>
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16" style="margin-top: 16px">
        <a-col :span="24">
          <a-card class="info-card" :bordered="false">
            <template #title>
              <div class="card-title">
                <icon-computer class="card-title-icon" />
                系统信息
              </div>
            </template>
            <a-descriptions :column="3" :data="systemInfoData" layout="inline-horizontal" bordered />
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, computed, nextTick, inject } from 'vue';
import {
  IconUser,
  IconUserGroup,
  IconSafe,
  IconUserAdd,
  IconHistory,
  IconCalendar,
  IconWifi,
  IconComputer,
  IconCheckCircle,
  IconCloseCircle,
  IconBarChart,
} from '@arco-design/web-vue/es/icon';
import * as echarts from 'echarts';
import { getDashboard } from '@/api/dashboard';
import { formatDateTime } from '@/utils/date';

interface Stats {
  userCount: number;
  roleCount: number;
  permissionCount: number;
  successCount: number;
  failedCount: number;
}

interface RecentLogin {
  username: string;
  loginTime: string;
  ip: string;
}

interface JobExecutionTrend {
  date: string;
  successCount: number;
  failedCount: number;
}

interface SystemInfo {
  systemName: string;
  version: string;
  os: string;
  goVersion: string;
  dbVersion: string;
  uptime: string;
  cpuUsage: string;
  memoryTotal: string;
  memoryUsed: string;
  memoryUsage: string;
  diskTotal: string;
  diskUsed: string;
  diskUsage: string;
  numGoroutine: number;
  numCPU: number;
}

const stats = ref<Stats>({
  userCount: 0,
  roleCount: 0,
  permissionCount: 0,
  successCount: 0,
  failedCount: 0,
});

const recentLogins = ref<RecentLogin[]>([]);
const trendData = ref<JobExecutionTrend[]>([]);
const systemInfo = ref<SystemInfo>({
  systemName: '',
  version: '',
  os: '',
  goVersion: '',
  dbVersion: '',
  uptime: '',
  cpuUsage: '',
  memoryTotal: '',
  memoryUsed: '',
  memoryUsage: '',
  diskTotal: '',
  diskUsed: '',
  diskUsage: '',
  numGoroutine: 0,
  numCPU: 0,
});

const systemInfoData = computed(() => [
  { label: '系统名称', value: systemInfo.value.systemName },
  { label: '系统版本', value: systemInfo.value.version },
  { label: '操作系统', value: systemInfo.value.os },
  { label: 'Go版本', value: systemInfo.value.goVersion },
  { label: '数据库版本', value: systemInfo.value.dbVersion },
  { label: '运行时间', value: systemInfo.value.uptime },
  { label: 'CPU使用率', value: systemInfo.value.cpuUsage },
  { label: '内存总量', value: systemInfo.value.memoryTotal },
  { label: '已用内存', value: systemInfo.value.memoryUsed },
  { label: '内存使用率', value: systemInfo.value.memoryUsage },
  { label: '磁盘总量', value: systemInfo.value.diskTotal },
  { label: '已用磁盘', value: systemInfo.value.diskUsed },
  { label: '磁盘使用率', value: systemInfo.value.diskUsage },
  { label: 'Goroutine数量', value: String(systemInfo.value.numGoroutine) },
  { label: 'CPU核心数', value: String(systemInfo.value.numCPU) },
]);

const chartRef = ref<HTMLElement>();
let chart: echarts.ECharts | null = null;

const getThemeColor = (cssVar: string) => {
  const el = document.documentElement;
  const varValue = getComputedStyle(el).getPropertyValue(cssVar).trim();
  return varValue.startsWith('rgb') ? varValue : `rgb(${varValue})`;
};

const getRgbaColor = (cssVar: string, alpha: number) => {
  const el = document.documentElement;
  const varValue = getComputedStyle(el).getPropertyValue(cssVar).trim();
  if (varValue.startsWith('rgb')) {
    return varValue.replace('rgb', 'rgba').replace(')', `, ${alpha})`);
  }
  return `rgba(${varValue}, ${alpha})`;
};

const getCurrentTheme = () => {
  return document.body.hasAttribute('arco-theme') ? 'dark' : 'light';
};

const initChart = () => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value, getCurrentTheme());
    window.addEventListener('resize', () => {
      if (chart) {
        chart.resize();
      }
    });
    updateChart();
  }
};

const updateChart = () => {
  if (!chart) return;

  const option = {
    color: ['#67c23a', '#f56c6c'],
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'line'
      }
    },
    legend: {
      data: ['成功', '失败'],
      top: 10
    },
    grid: {
      left: '8%',
      right: '5%',
      bottom: '8%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: trendData.value.map(item => item.date),
      axisLabel: {
        rotate: 0,
        fontSize: 12,
        margin: 8
      },
      boundaryGap: false
    },
    yAxis: {
      type: 'value',
      splitLine: {
        lineStyle: {
          type: 'dashed'
        }
      },
      axisLabel: {
        fontSize: 12
      },
      minInterval: 1,
      min: 0,
      max: function(value: { max: number }) {
        return Math.ceil(value.max * 1.1);
      }
    },
    series: [
      {
        name: '成功',
        type: 'line',
        data: trendData.value.map(item => item.successCount),
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        itemStyle: {
          color: '#67c23a'
        },
        lineStyle: {
          width: 2.5,
          color: '#67c23a'
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [{
              offset: 0,
              color: 'rgba(103, 194, 58, 0.3)'
            }, {
              offset: 1,
              color: 'rgba(103, 194, 58, 0.05)'
            }]
          }
        },
        emphasis: {
          focus: 'series',
          lineStyle: {
            width: 3
          }
        },
        showSymbol: true
      },
      {
        name: '失败',
        type: 'line',
        data: trendData.value.map(item => item.failedCount),
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        itemStyle: {
          color: '#f56c6c'
        },
        lineStyle: {
          width: 2.5,
          color: '#f56c6c'
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [{
              offset: 0,
              color: 'rgba(245, 108, 108, 0.3)'
            }, {
              offset: 1,
              color: 'rgba(245, 108, 108, 0.05)'
            }]
          }
        },
        emphasis: {
          focus: 'series',
          lineStyle: {
            width: 3
          }
        },
        showSymbol: true
      }
    ]
  };

  chart.setOption(option, true);
};

// 获取仪表盘数据
const fetchDashboardData = async () => {
  try {
    const res = await getDashboard();
    stats.value = res.data.stats;
    recentLogins.value = res.data.recentLogins;
    trendData.value = res.data.trendData;
    systemInfo.value = res.data.systemInfo;
    await nextTick();
    initChart();
  } catch (error) {
    console.error('获取仪表盘数据失败:', error);
  }
};

// 监听刷新事件
const handlePageRefresh = () => {
  fetchDashboardData();
};

// 监听主题变化
const observer = new MutationObserver(() => {
  if (chart) {
    const theme = getCurrentTheme();
    chart.dispose();
    chart = echarts.init(chartRef.value, theme);
    updateChart();
  }
});

onMounted(() => {
  fetchDashboardData();
  // 添加刷新事件监听
  window.addEventListener('page-refresh', handlePageRefresh);
  // 监听 body 的属性变化
  observer.observe(document.body, {
    attributes: true,
    attributeFilter: ['arco-theme']
  });
});

onUnmounted(() => {
  if (chart) {
    window.removeEventListener('resize', () => chart?.resize());
    chart.dispose();
    chart = null;
  }
  // 停止监听
  observer.disconnect();
  // 移除刷新事件监听
  window.removeEventListener('page-refresh', handlePageRefresh);
});
</script>

<style scoped>
.dashboard {
  padding: 16px;
  min-height: 100%;
}

.stat-card {
  height: 140px;
  background-color: var(--color-bg-2);
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
  background: linear-gradient(45deg, transparent, rgba(var(--primary-6), 0.15));
  border-radius: 0 8px 0 50%;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  background-color: var(--color-bg-3);
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

.stat-icon.success {
  color: rgb(var(--success-6));
  background-color: rgba(var(--success-6), 0.1);
}

.stat-icon.error {
  color: rgb(var(--danger-6));
  background-color: rgba(var(--danger-6), 0.1);
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

.list-card,
.info-card,
.chart-card {
  height: 100%;
  background-color: var(--color-bg-2);
  border-radius: 8px;
}

.chart-container {
  height: 360px;
  width: 100%;
  min-height: 360px;
  padding: 0 0 16px 0;
}

.chart-card {
  height: 100%;
  min-height: 420px;
  background-color: var(--color-bg-2);
  border-radius: 8px;
}

.list-card {
  height: 100%;
  min-height: 420px;
  background-color: var(--color-bg-2);
  border-radius: 8px;
  :deep(.arco-table-container) {
    height: 360px;
    overflow: auto;
  }
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
  background-color: var(--color-fill-3) !important;
  font-weight: 500;
}

:deep(.arco-table-tr:hover) {
  background-color: var(--color-fill-3) !important;
}

:deep(.arco-descriptions) {
  background-color: var(--color-bg-2);
}

:deep(.arco-descriptions-bordered) {
  border-radius: 8px;
  border: 1px solid var(--color-neutral-4);
}

:deep(.arco-descriptions-bordered .arco-descriptions-item) {
  padding: 12px 16px;
}

:deep(.arco-card) {
  background-color: var(--color-bg-2);
  border: 1px solid var(--color-neutral-4);
}

:deep(.arco-table-container) {
  border-radius: 8px;
  border: 1px solid var(--color-neutral-4);
}

:deep(.arco-table-cell) {
  padding: 12px 16px;
}
</style>
