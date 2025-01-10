<template>
  <div class="jobs">
    <a-card>
      <template #title>任务列表</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchForm.name"
            placeholder="请输入任务名称"
            style="width: 300px"
            allow-clear
            @search="handleSearch"
            @press-enter="handleSearch"
          />
          <a-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            style="width: 120px"
            allow-clear
            @change="handleSearch"
          >
            <a-option :value="1">运行中</a-option>
            <a-option :value="0">已停止</a-option>
          </a-select>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建任务
          </a-button>
          <a-button type="primary" status="success" @click="handleBatchExecute" :disabled="!selectedKeys.length">
            <template #icon><icon-play-circle /></template>
            批量执行一次
          </a-button>
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
        :row-selection="rowSelection"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '运行中' : '已停止' }}
          </a-tag>
        </template>
        <template #type="{ record }">
          <a-tag :color="getJobTypeColor(record.type)">
            {{ getJobTypeText(record.type) }}
          </a-tag>
        </template>
        <template #created_at="{ record }">
          {{ formatDateTime(record.created_at) }}
        </template>
        <template #updated_at="{ record }">
          {{ formatDateTime(record.updated_at) }}
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleEdit(record)">
              <template #icon><icon-edit /></template>
              编辑
            </a-button>
            <a-button
              type="text"
              size="small"
              status="success"
              @click="handleExecute(record)"
            >
              <template #icon><icon-play-circle /></template>
              执行一次
            </a-button>
            <a-button
              v-if="record.status === 0"
              type="text"
              size="small"
              @click="handleStart(record)"
            >
              <template #icon><icon-play-circle /></template>
              启动
            </a-button>
            <a-button
              v-else
              type="text"
              size="small"
              status="warning"
              @click="handleStop(record)"
            >
              <template #icon><icon-pause-circle /></template>
              停止
            </a-button>
            <a-button
              type="text"
              size="small"
              @click="handleHistory(record)"
            >
              <template #icon><icon-history /></template>
              执行历史
            </a-button>
            <a-popconfirm
              content="确定要删除该任务吗？"
              @ok="handleDelete(record)"
            >
              <a-button type="text" size="small" status="danger">
                <template #icon><icon-delete /></template>
                删除
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 任务表单对话框 -->
    <job-form
      v-model:visible="showForm"
      :is-edit="isEdit"
      :job-data="currentJob"
      @success="handleFormSuccess"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useRouter } from 'vue-router';
import {
  IconPlus,
  IconPlayCircle,
  IconPauseCircle,
  IconRefresh,
  IconEdit,
  IconDelete,
  IconHistory,
} from '@arco-design/web-vue/es/icon';
import type { Job, JobStatus, JobListRequest } from '@/api/types';
import {
  getJobList,
  startJob,
  stopJob,
  deleteJob,
  executeJob,
  executeJobs,
} from '@/api/job';
import JobForm from './components/JobForm.vue';

const router = useRouter();

const loading = ref(false);
const renderData = ref<Job[]>([]);
const selectedKeys = ref<number[]>([]);
const showForm = ref(false);
const isEdit = ref(false);
const currentJob = ref<Job>();

interface SearchFormState {
  name: string;
  status: JobStatus | '';
}

const searchForm = reactive<SearchFormState>({
  name: '',
  status: '',
});

const columns = [
  {
    title: '任务名称',
    dataIndex: 'name',
    align: 'center' as const,
  },
  {
    title: '任务类型',
    dataIndex: 'type',
    slotName: 'type',
    align: 'center' as const,
  },
  {
    title: '任务描述',
    dataIndex: 'description',
    align: 'center' as const,
  },
  {
    title: 'Cron 表达式',
    dataIndex: 'cron_expr',
    width: 200,
    align: 'center' as const,
  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
    align: 'center' as const,
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    slotName: 'created_at',
    align: 'center' as const,
  },
  {
    title: '操作',
    dataIndex: 'operations',
    slotName: 'operations',
    width: 320,
    align: 'center' as const,
  },
];

const rowSelection = {
  type: 'checkbox' as const,
  showCheckedAll: true,
  onlyCurrent: false,
  onChange: (selectedRowKeys: (string | number)[]) => {
    selectedKeys.value = selectedRowKeys as number[];
  },
};

const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
});

// 获取任务列表数据
const fetchData = async () => {
  try {
    loading.value = true;
    const { data } = await getJobList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: searchForm.name,
      status: searchForm.status ? searchForm.status as JobStatus : undefined,
    } as JobListRequest);
    renderData.value = data.items;
    pagination.total = data.total;
  } catch (err) {
    Message.error('获取任务列表失败');
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

// 创建任务
const handleCreate = () => {
  isEdit.value = false;
  currentJob.value = undefined;
  showForm.value = true;
};

// 编辑任务
const handleEdit = (record: Job) => {
  isEdit.value = true;
  currentJob.value = record;
  showForm.value = true;
};

// 执行一次任务
const handleExecute = async (record: Job) => {
  try {
    await executeJob(record.id);
    Message.success('执行任务成功');
  } catch (err) {
    Message.error('执行任务失败');
  }
};

// 批量执行任务
const handleBatchExecute = async () => {
  if (!selectedKeys.value.length) {
    Message.warning('请选择要执行的任务');
    return;
  }

  try {
    await executeJobs(selectedKeys.value);
    Message.success('批量执行任务成功');
  } catch (err) {
    Message.error('批量执行任务失败');
  }
};

// 启动任务
const handleStart = async (record: Job) => {
  try {
    await startJob(record.id);
    Message.success('启动任务成功');
    fetchData();
  } catch (err) {
    Message.error('启动任务失败');
  }
};

// 停止任务
const handleStop = async (record: Job) => {
  try {
    await stopJob(record.id);
    Message.success('停止任务成功');
    fetchData();
  } catch (err) {
    Message.error('停止任务失败');
  }
};

// 删除任务
const handleDelete = async (record: Job) => {
  try {
    await deleteJob(record.id);
    Message.success('删除任务成功');
    fetchData();
  } catch (err) {
    Message.error('删除任务失败');
  }
};

// 表单提交成功
const handleFormSuccess = () => {
  fetchData();
};

// 添加任务类型相关的工具函数
const getJobTypeText = (type: string) => {
  switch (type) {
    case 'shell':
      return 'Shell脚本';
    case 'http':
      return 'HTTP请求';
    case 'datax':
      return 'DataX任务';
    default:
      return '未知类型';
  }
};

const getJobTypeColor = (type: string) => {
  switch (type) {
    case 'shell':
      return 'blue';
    case 'http':
      return 'green';
    case 'datax':
      return 'purple';
    default:
      return 'gray';
  }
};

// 添加时间格式化函数
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  });
};

// 跳转到执行历史
const handleHistory = (record: Job) => {
  router.push({
    path: '/job/history',
    query: {
      job_id: record.id.toString(),
      job_name: record.name
    }
  });
};

// 初始化加载数据
fetchData();
</script>

<style scoped>
.jobs {
  padding: 16px;
}
</style>
