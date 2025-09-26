<template>
  <div class="job-history">
    <a-card>
      <template #title>执行历史</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchForm.jobName"
            placeholder="请输入任务名称"
            style="width: 300px"
            @search="handleSearch"
            allow-clear
            @clear="handleClear"
            @press-enter="handleSearch"
          />
          <a-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            style="width: 120px"
            allow-clear
            @change="handleSearch"
          >
            <a-option value="-1">执行中</a-option>
            <a-option value="1">成功</a-option>
            <a-option value="0">失败</a-option>
          </a-select>
          <a-button @click="fetchData">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
          <a-dropdown>
            <a-button>
              <template #icon><icon-delete /></template>
              清理日志
            </a-button>
            <template #content>
              <a-doption @click="handleClean(7)">清理一周前</a-doption>
              <a-doption @click="handleClean(15)">清理半个月前</a-doption>
              <a-doption @click="handleClean(30)">清理一个月前</a-doption>
              <a-doption @click="handleClean(90)">清理三个月前</a-doption>
              <a-doption @click="handleClean(180)">清理半年前</a-doption>
              <a-doption @click="handleClean(-1)">清理全部</a-doption>
            </template>
          </a-dropdown>
        </a-space>
      </template>

      <a-table
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        :columns="columns"
        :data="renderData"
        :scroll="{ x: '100%', y: '100%' }"
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

    <a-drawer
      v-model:visible="showLogModal"
      :title="`执行日志 - ${currentRecord?.job_name || ''}`"
      :footer="false"
      :mask-closable="true"
      :width="1200"
      placement="right"
      @cancel="handleLogModalClose"
    >
      <div class="log-header">
        <a-space>
          <a-tag :color="getStatusColor(currentRecord?.status || 0)">
            {{ getStatusText(currentRecord?.status || 0) }}
          </a-tag>
          <span v-if="currentRecord?.start_time">
            开始时间: {{ formatDateTime(currentRecord.start_time) }}
          </span>
          <span v-if="currentRecord?.end_time">
            结束时间: {{ formatDateTime(currentRecord.end_time) }}
          </span>
          <span v-if="currentRecord?.duration">
            耗时: {{ currentRecord.duration }}ms
          </span>
          <a-button
            type="primary"
            size="small"
            @click="refreshLog"
            :loading="refreshingLog"
            v-if="currentRecord?.status === -1"
          >
            <template #icon><icon-refresh /></template>
            刷新日志
          </a-button>
        </a-space>
      </div>
      <a-divider />
      <div class="log-content" ref="logContentRef">
        <pre>{{ currentLog }}</pre>
      </div>
    </a-drawer>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onBeforeUnmount, onMounted, nextTick } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import { IconRefresh, IconEye, IconClose, IconDelete } from '@arco-design/web-vue/es/icon';
import { getJobHistoryList, getJobHistoryDetail, cleanJobHistory } from '@/api/job';
import type { JobHistory } from '@/api/types';
import { formatDateTime } from '@/utils/date';
import { useRoute, useRouter } from 'vue-router';
import { usePageRefresh } from '@/composables/usePageRefresh';

const route = useRoute();
const router = useRouter();
const loading = ref(false);
const renderData = ref<JobHistory[]>([]);
const showLogModal = ref(false);
const currentLog = ref('');
const currentRecord = ref<JobHistory | null>(null);
const refreshingLog = ref(false);
const logContentRef = ref<HTMLElement>();
let isUnmounted = false;

interface SearchFormState {
  jobName: string;
  status: string;
  jobId?: number;
}

const searchForm = reactive<SearchFormState>({
  jobName: '',
  status: '',
  jobId: undefined,
});

// 使用页面刷新功能
usePageRefresh(() => {
  fetchData();
});

// 初始化查询参数
onMounted(() => {
  const { job_id, job_name, status } = route.query;
  if (job_id) {
    searchForm.jobId = parseInt(job_id as string);
  }
  if (job_name) {
    searchForm.jobName = job_name as string;
  }
  if (status) {
    searchForm.status = status as string;
  }
  fetchData();
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
  pageSize: 20,
  showTotal: true,
  showJumper: true,
  showPageSize: true,
});



// 获取历史数据
const fetchData = async () => {
  if (isUnmounted) return;

  try {
    loading.value = true;
    const { data } = await getJobHistoryList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: searchForm.jobName || undefined,
      status: searchForm.status !== '' ? Number(searchForm.status) : undefined,
      job_id: searchForm.jobId,
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

// 清除筛选条件
const handleClear = () => {
  searchForm.status = '';
  searchForm.jobId = undefined;
  // 清除URL中的查询参数
  router.replace({
    path: route.path
  });
  fetchData();
};

// 搜索
const handleSearch = () => {
  if (isUnmounted) return;
  pagination.current = 1;
  // 更新URL查询参数
  const query: Record<string, string> = {};
  if (searchForm.jobId) {
    query.job_id = searchForm.jobId.toString();
  }
  if (searchForm.jobName) {
    query.job_name = searchForm.jobName;
  }
  if (searchForm.status) {
    query.status = searchForm.status;
  }
  router.replace({
    path: route.path,
    query
  });
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

// 获取状态颜色
const getStatusColor = (status: number) => {
  switch (status) {
    case -1: return 'blue';   // 执行中
    case 1: return 'green';   // 成功
    case 0: return 'red';     // 失败
    default: return 'gray';
  }
};

// 获取状态文本
const getStatusText = (status: number) => {
  switch (status) {
    case -1: return '执行中';
    case 1: return '成功';
    case 0: return '失败';
    default: return '未知';
  }
};

// 查看日志
const handleView = async (record: JobHistory) => {
  if (isUnmounted) return;
  currentRecord.value = record;
  currentLog.value = record.output || '无输出';
  if (record.error) {
    currentLog.value += `\n\n错误信息：\n${record.error}`;
  }
  showLogModal.value = true;

  // 打开弹框后滚动到底部
  scrollToBottom();
};

// 滚动到日志底部
const scrollToBottom = () => {
  nextTick(() => {
    if (logContentRef.value) {
      logContentRef.value.scrollTop = logContentRef.value.scrollHeight;
    }
  });
};

// 刷新日志
const refreshLog = async () => {
  if (!currentRecord.value || refreshingLog.value) return;

  try {
    refreshingLog.value = true;
    const { data } = await getJobHistoryDetail(currentRecord.value.id);

    // 更新当前记录
    currentRecord.value = data;

    // 更新日志内容
    currentLog.value = data.output || '无输出';
    if (data.error) {
      currentLog.value += `\n\n错误信息：\n${data.error}`;
    }

    // 如果任务已完成，同时刷新列表数据
    if (data.status !== -1) {
      fetchData();
    }

    // 滚动到底部
    scrollToBottom();

    Message.success('日志已刷新');
  } catch (err) {
    Message.error('刷新日志失败');
  } finally {
    refreshingLog.value = false;
  }
};

// 关闭日志对话框
const handleLogModalClose = () => {
  if (isUnmounted) return;
  showLogModal.value = false;
  currentLog.value = '';
  currentRecord.value = null;
  refreshingLog.value = false;
};

// 清理日志
const handleClean = async (days: number) => {
  try {
    await cleanJobHistory(days);
    Message.success('清理成功');
    fetchData();
  } catch (err) {
    Message.error('清理失败');
  }
};

onBeforeUnmount(() => {
  isUnmounted = true;
  renderData.value = [];
  currentLog.value = '';
  currentRecord.value = null;
  refreshingLog.value = false;
  showLogModal.value = false;
});
</script>

<style scoped>
.job-history {
  padding: 16px;
  height: calc(100vh - 80px); /* 减去header和padding的高度 */
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 确保卡片占满剩余空间 */
:deep(.arco-card) {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

:deep(.arco-card-body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0; /* 重要：允许flex子项收缩 */
}

/* 表格容器样式 */
:deep(.arco-table-container) {
  flex: 1;
  overflow: auto;
  min-height: 0; /* 重要：允许flex子项收缩 */
}

/* 表格主体区域可以直接滚动 */
:deep(.arco-table-body) {
  overflow: auto;
}

:deep(.arco-table-tbody) {
  overflow: visible;
}

/* 表格包装器 */
:deep(.arco-table-wrapper) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 表格主体 */
:deep(.arco-table) {
  flex: 1;
  overflow: hidden;
}

/* 分页器固定在底部 */
:deep(.arco-pagination) {
  margin-top: 16px;
  flex-shrink: 0;
  padding: 8px 0;
}

/* 搜索区域样式 */
:deep(.arco-card-header) {
  flex-shrink: 0;
  padding: 16px 20px;
}

/* 搜索表单样式 */
:deep(.arco-form) {
  flex-shrink: 0;
  margin-bottom: 16px;
}

.log-header {
  padding: 12px 0;
  background-color: var(--color-fill-1);
  border-radius: 4px;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.log-content {
  height: calc(100vh - 200px);
  overflow-y: auto;
  background-color: var(--color-fill-2);
  border-radius: 4px;
  border: 1px solid var(--color-border);
}

.log-content pre {
  margin: 0;
  padding: 16px;
  font-family: Monaco, Menlo, Consolas, "Courier New", monospace;
  font-size: 13px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: var(--color-text-1);
}
</style>
