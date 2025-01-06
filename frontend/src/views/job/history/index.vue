<template>
  <div class="container">
    <div class="header">
      <a-typography-title :heading="6" style="margin: 0">
        执行历史
      </a-typography-title>
      <div class="toolbar">
        <a-space>
          <a-button>
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </a-space>
      </div>
    </div>
    <a-card class="general-card" :bordered="false">
      <a-table
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        :columns="columns"
        :data="renderData"
      >
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #operations>
          <a-space>
            <a-button type="text" size="small">
              <template #icon><icon-eye /></template>
              查看
            </a-button>
            <a-button type="text" size="small">
              <template #icon><icon-download /></template>
              下载日志
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { IconRefresh, IconEye, IconDownload } from '@arco-design/web-vue/es/icon';

const loading = ref(false);
const renderData = ref([]);

const columns = [
  {
    title: '任务名称',
    dataIndex: 'jobName',
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
  },
  {
    title: '结束时间',
    dataIndex: 'endTime',
  },
  {
    title: '执行时长',
    dataIndex: 'duration',
  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
  },
  {
    title: '操作',
    dataIndex: 'operations',
    slotName: 'operations',
  },
];

const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
});

const getStatusColor = (status: string) => {
  switch (status) {
    case 'success':
      return 'green';
    case 'failed':
      return 'red';
    case 'running':
      return 'blue';
    default:
      return 'gray';
  }
};

const getStatusText = (status: string) => {
  switch (status) {
    case 'success':
      return '成功';
    case 'failed':
      return '失败';
    case 'running':
      return '运行中';
    default:
      return '未知';
  }
};

// TODO: 实现任务历史的查询和日志下载功能
</script>

<style scoped>
.container {
  padding: 0 20px 20px 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

:deep(.general-card) {
  border-radius: 4px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
