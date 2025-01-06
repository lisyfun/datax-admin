<template>
  <div class="container">
    <div class="header">
      <a-typography-title :heading="6" style="margin: 0">
        任务列表
      </a-typography-title>
      <div class="toolbar">
        <a-space>
          <a-button type="primary">
            <template #icon><icon-plus /></template>
            新建任务
          </a-button>
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
          <a-tag :color="record.status === 'running' ? 'green' : 'red'">
            {{ record.status === 'running' ? '运行中' : '已停止' }}
          </a-tag>
        </template>
        <template #operations>
          <a-space>
            <a-button type="text" size="small">
              <template #icon><icon-edit /></template>
              编辑
            </a-button>
            <a-button type="text" size="small">
              <template #icon><icon-play-circle /></template>
              启动
            </a-button>
            <a-button type="text" size="small" status="danger">
              <template #icon><icon-delete /></template>
              删除
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { IconPlus, IconRefresh, IconEdit, IconDelete, IconPlayCircle } from '@arco-design/web-vue/es/icon';

const loading = ref(false);
const renderData = ref([]);

const columns = [
  {
    title: '任务名称',
    dataIndex: 'name',
  },
  {
    title: '任务描述',
    dataIndex: 'description',
  },
  {
    title: 'Cron 表达式',
    dataIndex: 'cron',
  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
  },
  {
    title: '更新时间',
    dataIndex: 'updateTime',
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

// TODO: 实现任务列表的 CRUD 操作
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
