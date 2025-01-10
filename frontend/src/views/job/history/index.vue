<template>
  <div class="history">
    <a-card>
      <template #title>执行历史</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchForm.jobName"
            placeholder="请输入任务名称"
            style="width: 300px"
            @search="handleSearch"
          />
          <a-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            style="width: 120px"
            allow-clear
            @change="handleSearch"
          >
            <a-option :value="1">成功</a-option>
            <a-option :value="0">失败</a-option>
          </a-select>
          <a-button @click="fetchData">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        :columns="columns"
        :data="renderData"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleView(record)">
              <template #icon><icon-eye /></template>
              查看
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:visible="showLogModal"
      title="执行日志"
      :footer="false"
      :mask-closable="true"
      @cancel="handleLogModalClose"
    >
      <div class="log-content">
        <pre>{{ currentLog }}</pre>
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onBeforeUnmount } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import { IconRefresh, IconEye } from '@arco-design/web-vue/es/icon';
import { getJobHistoryList } from '@/api/job';
import type { JobHistory } from '@/api/types';
import { formatDateTime } from '@/utils/date';

const loading = ref(false);
const renderData = ref<JobHistory[]>([]);
const showLogModal = ref(false);
const currentLog = ref('');
let isUnmounted = false;

interface SearchFormState {
  jobName: string;
  status: string;
}

const searchForm = reactive<SearchFormState>({
  jobName: '',
  status: '',
});

const columns = [
  {
    title: '任务名称',
    dataIndex: 'job_name',
  },
  {
    title: '开始时间',
    dataIndex: 'start_time',
    render: (data: any) => formatDateTime(data?.record?.start_time),
  },
  {
    title: '结束时间',
    dataIndex: 'end_time',
    render: (data: any) => formatDateTime(data?.record?.end_time),
  },
  {
    title: '执行时长',
    dataIndex: 'duration',
    render: (data: any) => data?.record?.duration ? `${data.record.duration}ms` : '-',
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
    width: 200,
    align: 'center' as const,
  },
] as TableColumnData[];

const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showTotal: true,
  showJumper: true,
});

const getStatusColor = (status: number) => {
  switch (status) {
    case 1:
      return 'green';
    case 0:
      return 'red';
    default:
      return 'gray';
  }
};

const getStatusText = (status: number) => {
  switch (status) {
    case 1:
      return '成功';
    case 0:
      return '失败';
    default:
      return '未知';
  }
};

// 获取历史数据
const fetchData = async () => {
  if (isUnmounted) return;

  try {
    loading.value = true;
    const { data } = await getJobHistoryList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: searchForm.jobName || undefined,
      status: searchForm.status ? parseInt(searchForm.status) : undefined,
    });

    if (!isUnmounted) {
      renderData.value = data.items;
      pagination.total = data.total;
    }
  } catch (err) {
    if (!isUnmounted) {
      Message.error('获取执行历史失败');
    }
  } finally {
    if (!isUnmounted) {
      loading.value = false;
    }
  }
};

// 搜索
const handleSearch = () => {
  if (isUnmounted) return;
  pagination.current = 1;
  fetchData();
};

// 分页变化
const onPageChange = (current: number) => {
  if (isUnmounted) return;
  pagination.current = current;
  fetchData();
};

const onPageSizeChange = (pageSize: number) => {
  if (isUnmounted) return;
  pagination.pageSize = pageSize;
  fetchData();
};

// 查看日志
const handleView = async (record: JobHistory) => {
  if (isUnmounted) return;
  currentLog.value = record.output || '无输出';
  if (record.error) {
    currentLog.value += `\n\n错误信息：\n${record.error}`;
  }
  showLogModal.value = true;
};

// 关闭日志对话框
const handleLogModalClose = () => {
  if (isUnmounted) return;
  showLogModal.value = false;
  currentLog.value = '';
};

onBeforeUnmount(() => {
  isUnmounted = true;
  renderData.value = [];
  currentLog.value = '';
  showLogModal.value = false;
});

// 初始化加载数据
fetchData();
</script>

<style scoped>
.history {
  padding: 16px;
}

.log-content {
  max-height: 500px;
  overflow-y: auto;
  background-color: var(--color-fill-2);
  padding: 16px;
  border-radius: 4px;
}

.log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: monospace;
  font-size: 14px;
  line-height: 1.5;
  color: var(--color-text-2);
}
</style>
