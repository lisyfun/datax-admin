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
            <a-option value="success">成功</a-option>
            <a-option value="failed">失败</a-option>
            <a-option value="running">运行中</a-option>
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
            <a-button type="text" size="small" @click="handleDownload(record)">
              <template #icon><icon-download /></template>
              下载日志
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 查看日志对话框 -->
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
import { ref, reactive } from 'vue';
import { Message } from '@arco-design/web-vue';
import { IconRefresh, IconEye, IconDownload } from '@arco-design/web-vue/es/icon';

const loading = ref(false);
const renderData = ref([]);
const showLogModal = ref(false);
const currentLog = ref('');

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
    width: 200,
    align: 'center' as const,
  },
];

const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showTotal: true,
  showJumper: true,
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

// 获取历史数据
const fetchData = async () => {
  try {
    loading.value = true;
    // TODO: 调用API获取历史数据
    // const { data } = await getJobHistory({
    //   page: pagination.current,
    //   page_size: pagination.pageSize,
    //   job_name: searchForm.jobName,
    //   status: searchForm.status || undefined,
    // });
    // renderData.value = data.items;
    // pagination.total = data.total;
  } catch (err) {
    Message.error('获取执行历史失败');
  } finally {
    loading.value = false;
  }
};

// 搜索
const handleSearch = () => {
  pagination.current = 1;
  fetchData();
};

// 分页变化
const onPageChange = (current: number) => {
  pagination.current = current;
  fetchData();
};

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  fetchData();
};

// 查看日志
const handleView = async (record: any) => {
  try {
    // TODO: 调用API获取日志内容
    // const { data } = await getJobLog(record.id);
    // currentLog.value = data.content;
    showLogModal.value = true;
  } catch (err) {
    Message.error('获取日志失败');
  }
};

// 下载日志
const handleDownload = async (record: any) => {
  try {
    // TODO: 调用API下载日志
    // await downloadJobLog(record.id);
    Message.success('日志下载成功');
  } catch (err) {
    Message.error('日志下载失败');
  }
};

// 关闭日志对话框
const handleLogModalClose = () => {
  showLogModal.value = false;
  currentLog.value = '';
};

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
