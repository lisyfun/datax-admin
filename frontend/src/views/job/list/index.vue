<template>
  <div class="container">
    <div class="header">
      <a-typography-title :heading="6" style="margin: 0">
        任务列表
      </a-typography-title>
      <div class="toolbar">
        <a-space>
          <a-input-search
            v-model="searchForm.name"
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
            <a-option :value="1">运行中</a-option>
            <a-option :value="0">已停止</a-option>
          </a-select>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建任务
          </a-button>
          <a-button @click="fetchData">
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
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '运行中' : '已停止' }}
          </a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleEdit(record)">
              <template #icon><icon-edit /></template>
              编辑
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
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { Message } from '@arco-design/web-vue';
import {
  IconPlus,
  IconRefresh,
  IconEdit,
  IconDelete,
  IconPlayCircle,
  IconPauseCircle,
} from '@arco-design/web-vue/es/icon';
import type { Job } from '@/api/job';
import {
  getJobList,
  startJob,
  stopJob,
  deleteJob,
} from '@/api/job';

const loading = ref(false);
const renderData = ref<Job[]>([]);

const searchForm = reactive({
  name: '',
  status: '',
});

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
    dataIndex: 'cron_expr',
  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
  },
  {
    title: '创建时间',
    dataIndex: 'create_time',
  },
  {
    title: '更新时间',
    dataIndex: 'update_time',
  },
  {
    title: '操作',
    dataIndex: 'operations',
    slotName: 'operations',
    width: 250,
  },
];

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
      name: searchForm.name,
      status: searchForm.status ? Number(searchForm.status) : undefined,
    });
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
  // TODO: 打开创建任务弹窗
};

// 编辑任务
const handleEdit = (record: Job) => {
  // TODO: 打开编辑任务弹窗
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

// 初始化加载数据
fetchData();
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
