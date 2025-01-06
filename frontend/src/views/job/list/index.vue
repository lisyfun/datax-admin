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
          <a-button type="primary" status="success" @click="handleBatchExecute" :disabled="!selectedKeys.length">
            <template #icon><icon-play-circle /></template>
            批量执行一次
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
          {{ record.type === 1 ? 'Shell脚本' : 'DataX任务' }}
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
    <a-modal
      v-model:visible="showForm"
      :title="isEdit ? '编辑任务' : '新建任务'"
      @ok="handleFormSubmit"
      @cancel="handleFormCancel"
    >
      <a-form ref="formRef" :model="form" :rules="rules">
        <a-form-item field="name" label="任务名称" :rules="[{ required: true, message: '请输入任务名称' }]">
          <a-input v-model="form.name" placeholder="请输入任务名称" />
        </a-form-item>
        <a-form-item field="type" label="任务类型" :rules="[{ required: true, message: '请选择任务类型' }]">
          <a-radio-group v-model="form.type">
            <a-radio :value="1">Shell脚本</a-radio>
            <a-radio :value="2">DataX任务</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item field="description" label="任务描述">
          <a-textarea v-model="form.description" placeholder="请输入任务描述" />
        </a-form-item>
        <a-form-item
          field="command"
          :label="form.type === 1 ? '执行命令' : 'DataX配置'"
          :rules="[{ required: true, message: form.type === 1 ? '请输入执行命令' : '请输入DataX配置' }]"
        >
          <a-textarea
            v-model="form.command"
            :placeholder="form.type === 1 ? '请输入要执行的命令' : '请输入DataX任务配置JSON'"
            :auto-size="form.type === 2 ? { minRows: 10, maxRows: 20 } : undefined"
          />
          <template v-if="form.type === 2" #extra>
            <div class="datax-tools">
              <a-space>
                <a-button type="text" @click="handleFormatJson">
                  <template #icon><icon-code /></template>
                  格式化JSON
                </a-button>
                <a-button type="text" @click="handleLoadTemplate">
                  <template #icon><icon-file /></template>
                  加载模板
                </a-button>
              </a-space>
            </div>
          </template>
        </a-form-item>
        <a-form-item field="working_dir" label="工作目录" v-if="form.type === 1">
          <a-input v-model="form.working_dir" placeholder="请输入工作目录，默认为当前目录" />
        </a-form-item>
        <a-form-item field="cron_expr" label="Cron 表达式" :rules="[{ required: true, message: '请输入 Cron 表达式' }]">
          <a-input v-model="form.cron_expr" placeholder="请输入 Cron 表达式">
            <template #append>
              <a-tooltip content="Cron 表达式在线生成">
                <a-button @click="openCronHelper">
                  <template #icon><icon-question-circle /></template>
                </a-button>
              </a-tooltip>
            </template>
          </a-input>
        </a-form-item>
        <a-form-item field="timeout" label="超时时间(秒)">
          <a-input-number
            v-model="form.timeout"
            placeholder="请输入超时时间"
            :min="0"
            :max="86400"
            :step="1"
          />
        </a-form-item>
        <a-form-item field="retry_times" label="重试次数">
          <a-input-number
            v-model="form.retry_times"
            placeholder="请输入重试次数"
            :min="0"
            :max="10"
            :step="1"
          />
        </a-form-item>
        <a-form-item field="retry_interval" label="重试间隔(秒)">
          <a-input-number
            v-model="form.retry_interval"
            placeholder="请输入重试间隔"
            :min="1"
            :max="3600"
            :step="1"
          />
        </a-form-item>
      </a-form>
    </a-modal>
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
  IconQuestionCircle,
  IconCode,
  IconFile,
} from '@arco-design/web-vue/es/icon';
import type { Job } from '@/api/job';
import {
  getJobList,
  startJob,
  stopJob,
  deleteJob,
  executeJob,
  executeJobs,
  createJob,
  updateJob,
} from '@/api/job';

const loading = ref(false);
const renderData = ref<Job[]>([]);
const selectedKeys = ref<number[]>([]);

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
    title: '任务类型',
    dataIndex: 'type',
    slotName: 'type',
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

// 表单相关
const showForm = ref(false);
const isEdit = ref(false);
const formRef = ref();
const form = reactive({
  id: 0,
  name: '',
  description: '',
  type: 1,
  command: '',
  working_dir: '',
  cron_expr: '',
  timeout: 0,
  retry_times: 0,
  retry_interval: 0,
});

// 表单规则
const rules = {
  name: [
    { required: true, message: '请输入任务名称' },
    { maxLength: 50, message: '任务名称最多 50 个字符' },
  ],
  command: [
    { required: true, message: '请输入执行命令' },
  ],
  cron_expr: [
    { required: true, message: '请输入 Cron 表达式' },
    {
      pattern: /^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*\/([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])) (\*|([0-9]|1[0-9]|2[0-3])|\*\/([0-9]|1[0-9]|2[0-3])) (\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*\/([1-9]|1[0-9]|2[0-9]|3[0-1])) (\*|([1-9]|1[0-2])|\*\/([1-9]|1[0-2])) (\*|([0-6])|\*\/([0-6]))$/,
      message: '请输入正确的 Cron 表达式',
    },
  ],
  timeout: [
    { type: 'number' as const, min: 0, max: 86400, message: '超时时间范围为 0-86400 秒' },
  ],
  retry_times: [
    { type: 'number' as const, min: 0, max: 10, message: '重试次数范围为 0-10 次' },
  ],
  retry_interval: [
    { type: 'number' as const, min: 1, max: 3600, message: '重试间隔范围为 1-3600 秒' },
  ],
};

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
  isEdit.value = false;
  form.id = 0;
  form.name = '';
  form.description = '';
  form.type = 1;
  form.command = '';
  form.working_dir = '';
  form.cron_expr = '';
  form.timeout = 0;
  form.retry_times = 0;
  form.retry_interval = 0;
  showForm.value = true;
};

// 编辑任务
const handleEdit = (record: Job) => {
  isEdit.value = true;
  form.id = record.id;
  form.name = record.name;
  form.description = record.description || '';
  form.type = record.type;
  form.command = record.command;
  form.working_dir = record.working_dir || '';
  form.cron_expr = record.cron_expr;
  form.timeout = record.timeout || 0;
  form.retry_times = record.retry_times || 0;
  form.retry_interval = record.retry_interval || 0;
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

// 提交表单
const handleFormSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();

    const data = {
      name: form.name,
      description: form.description,
      type: form.type,
      command: form.command,
      working_dir: form.working_dir,
      cron_expr: form.cron_expr,
      timeout: form.timeout,
      retry_times: form.retry_times,
      retry_interval: form.retry_interval,
    };

    if (isEdit.value) {
      await updateJob(form.id, data);
      Message.success('编辑任务成功');
    } else {
      await createJob(data);
      Message.success('创建任务成功');
    }

    showForm.value = false;
    fetchData();
  } catch (err: any) {
    if (err.response?.data?.error) {
      Message.error(err.response.data.error);
    } else if (err.errors) {
      Message.error('表单验证失败，请检查输入');
    } else {
      Message.error(isEdit.value ? '编辑任务失败' : '创建任务失败');
    }
  }
};

// 取消表单
const handleFormCancel = () => {
  showForm.value = false;
  if (formRef.value) {
    formRef.value.resetFields();
  }
};

// 打开 Cron 表达式帮助
const openCronHelper = () => {
  window.open('https://crontab.guru/', '_blank');
};

// 格式化 JSON
const handleFormatJson = () => {
  try {
    const jsonObj = JSON.parse(form.command);
    form.command = JSON.stringify(jsonObj, null, 2);
  } catch (err) {
    Message.error('JSON 格式错误，无法格式化');
  }
};

// 加载模板
const handleLoadTemplate = () => {
  const template = {
    "job": {
      "setting": {
        "speed": {
          "channel": 3
        }
      },
      "content": [
        {
          "reader": {
            "name": "mysqlreader",
            "parameter": {
              "username": "root",
              "password": "password",
              "column": ["*"],
              "connection": [
                {
                  "table": ["table"],
                  "jdbcUrl": ["jdbc:mysql://localhost:3306/database"]
                }
              ]
            }
          },
          "writer": {
            "name": "mysqlwriter",
            "parameter": {
              "username": "root",
              "password": "password",
              "column": ["*"],
              "connection": [
                {
                  "table": ["table"],
                  "jdbcUrl": "jdbc:mysql://localhost:3306/database"
                }
              ]
            }
          }
        }
      ]
    }
  };
  form.command = JSON.stringify(template, null, 2);
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

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.general-card {
  border-radius: 4px;
}

.datax-tools {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed var(--color-neutral-3);
}

.datax-tools :deep(.arco-btn) {
  padding: 0 4px;
  font-size: 14px;
  color: var(--color-text-3);
}

.datax-tools :deep(.arco-btn:hover) {
  color: rgb(var(--primary-6));
}
</style>
