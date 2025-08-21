<template>
  <div class="operation-management">
    <a-card>
      <template #title>操作管理</template>
      <template #extra>
        <a-space>
          <a-button type="primary" status="danger" @click="handleClearLogs">
            <template #icon><icon-delete /></template>
            清空记录
          </a-button>
          <a-button type="primary" status="warning" @click="handleBatchDelete" :disabled="selectedRowKeys.length === 0">
            <template #icon><icon-delete /></template>
            批量删除
          </a-button>
        </a-space>
      </template>

      <!-- 搜索表单 -->
      <a-form :model="searchForm" layout="inline" style="margin-bottom: 16px;">
        <a-form-item label="用户名">
          <a-input v-model="searchForm.username" placeholder="请输入用户名" style="width: 150px;" />
        </a-form-item>
        <a-form-item label="操作模块">
          <a-select v-model="searchForm.module" placeholder="请选择模块" style="width: 120px;" allow-clear>
            <a-option value="user">用户管理</a-option>
            <a-option value="role">角色管理</a-option>
            <a-option value="permission">权限管理</a-option>
            <a-option value="menu">菜单管理</a-option>
            <a-option value="job">任务管理</a-option>
            <a-option value="terminal">终端管理</a-option>
            <a-option value="kafka">Kafka管理</a-option>
            <a-option value="system">系统管理</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="操作动作">
          <a-select v-model="searchForm.action" placeholder="请选择动作" style="width: 120px;" allow-clear>
            <a-option value="create">创建</a-option>
            <a-option value="update">更新</a-option>
            <a-option value="delete">删除</a-option>
            <a-option value="view">查看</a-option>
            <a-option value="login">登录</a-option>
            <a-option value="logout">登出</a-option>
            <a-option value="export">导出</a-option>
            <a-option value="import">导入</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="操作状态">
          <a-select v-model="searchForm.status" placeholder="请选择状态" style="width: 100px;" allow-clear>
            <a-option :value="1">成功</a-option>
            <a-option :value="0">失败</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker v-model="searchForm.timeRange" style="width: 300px;" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">
            <template #icon><icon-search /></template>
            搜索
          </a-button>
          <a-button @click="handleReset" style="margin-left: 8px;">
            <template #icon><icon-refresh /></template>
            重置
          </a-button>
        </a-form-item>
      </a-form>

      <!-- 数据表格 -->
      <a-table
        :columns="columns"
        :data="logList"
        :loading="loading"
        :pagination="pagination"
        :row-selection="rowSelection"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      >
        <template #module="{ record }">
          <a-tag :color="getModuleColor(record.module)">
            {{ getModuleText(record.module) }}
          </a-tag>
        </template>
        <template #action="{ record }">
          <a-tag :color="getActionColor(record.action)">
            {{ getActionText(record.action) }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '成功' : '失败' }}
          </a-tag>
        </template>
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="handleViewDetail(record)">
            详情
          </a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">
            删除
          </a-button>
        </template>
      </a-table>
    </a-card>

    <!-- 详情对话框 -->
    <a-modal v-model:visible="detailVisible" title="操作详情" width="800px" :footer="false">
      <a-descriptions :data="detailData" :column="2" bordered />
    </a-modal>

    <!-- 清空记录对话框 -->
    <a-modal v-model:visible="clearVisible" title="清空操作记录" @ok="confirmClearLogs">
      <p>请选择要清空的操作记录范围：</p>
      <a-form :model="clearForm">
        <a-form-item label="清空天数前的记录">
          <a-input-number v-model="clearForm.beforeDays" :min="1" :max="365" placeholder="请输入天数" />
          <span style="margin-left: 8px; color: #999;">天</span>
        </a-form-item>
      </a-form>
      <p style="color: #ff4d4f; margin-top: 16px;">
        <icon-exclamation-circle-fill />
        此操作将删除 {{ clearForm.beforeDays }} 天前的所有操作记录，且无法恢复！
      </p>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import { Message, Modal } from '@arco-design/web-vue';
import {
  IconDelete,
  IconSearch,
  IconRefresh,
  IconExclamationCircleFill,
} from '@arco-design/web-vue/es/icon';
import { getOperationLogs, deleteOperationLog, batchDeleteOperationLogs, clearOperationLogs, type OperationLog } from '@/api/operation_log';

// 响应式数据
const loading = ref(false);
const logList = ref<OperationLog[]>([]);
const selectedRowKeys = ref<number[]>([]);
const detailVisible = ref(false);
const clearVisible = ref(false);
const detailData = ref<any[]>([]);

// 搜索表单
const searchForm = reactive({
  username: '',
  module: '',
  action: '',
  status: undefined,
  timeRange: [],
});

// 清空表单
const clearForm = reactive({
  beforeDays: 30,
});

// 分页配置
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showTotal: true,
});

// 表格列配置
const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '用户名', dataIndex: 'username', width: 120 },
  { title: '操作模块', dataIndex: 'module', slotName: 'module', width: 120 },
  { title: '操作动作', dataIndex: 'action', slotName: 'action', width: 100 },
  { title: '操作描述', dataIndex: 'description', width: 200, ellipsis: true, tooltip: true },
  { title: 'IP地址', dataIndex: 'ip', width: 120 },
  { title: '状态', dataIndex: 'status', slotName: 'status', width: 80 },
  { title: '操作时间', dataIndex: 'created_at', slotName: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

// 行选择配置
const rowSelection = reactive({
  type: 'checkbox',
  showCheckedAll: true,
  onSelectionChange: (rowKeys: string[]) => {
    selectedRowKeys.value = rowKeys.map(key => parseInt(key));
  },
});

// 获取模块颜色
const getModuleColor = (module: string) => {
  const colors: Record<string, string> = {
    user: 'blue',
    role: 'green',
    permission: 'orange',
    menu: 'purple',
    job: 'cyan',
    terminal: 'magenta',
    kafka: 'red',
    system: 'gray',
  };
  return colors[module] || 'gray';
};

// 获取模块文本
const getModuleText = (module: string) => {
  const texts: Record<string, string> = {
    user: '用户管理',
    role: '角色管理',
    permission: '权限管理',
    menu: '菜单管理',
    job: '任务管理',
    terminal: '终端管理',
    kafka: 'Kafka管理',
    system: '系统管理',
  };
  return texts[module] || module;
};

// 获取动作颜色
const getActionColor = (action: string) => {
  const colors: Record<string, string> = {
    create: 'green',
    update: 'blue',
    delete: 'red',
    view: 'gray',
    login: 'cyan',
    logout: 'orange',
    export: 'purple',
    import: 'magenta',
  };
  return colors[action] || 'gray';
};

// 获取动作文本
const getActionText = (action: string) => {
  const texts: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    view: '查看',
    login: '登录',
    logout: '登出',
    export: '导出',
    import: '导入',
  };
  return texts[action] || action;
};

// 格式化时间
const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN');
};

// 获取操作记录列表
const fetchLogs = async () => {
  loading.value = true;
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      username: searchForm.username,
      module: searchForm.module,
      action: searchForm.action,
      status: searchForm.status,
      start_time: searchForm.timeRange[0] ? new Date(searchForm.timeRange[0]).toISOString() : '',
      end_time: searchForm.timeRange[1] ? new Date(searchForm.timeRange[1]).toISOString() : '',
    };

    const { data } = await getOperationLogs(params);
    logList.value = data.list || [];
    pagination.total = data.total || 0;
  } catch (error: any) {
    Message.error(error.message || '获取操作记录失败');
  } finally {
    loading.value = false;
  }
};

// 搜索
const handleSearch = () => {
  pagination.current = 1;
  fetchLogs();
};

// 重置
const handleReset = () => {
  Object.assign(searchForm, {
    username: '',
    module: '',
    action: '',
    status: undefined,
    timeRange: [],
  });
  pagination.current = 1;
  fetchLogs();
};

// 分页变化
const handlePageChange = (page: number) => {
  pagination.current = page;
  fetchLogs();
};

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  pagination.current = 1;
  fetchLogs();
};

// 查看详情
const handleViewDetail = (record: any) => {
  detailData.value = [
    { label: 'ID', value: record.id },
    { label: '用户名', value: record.username },
    { label: '操作模块', value: getModuleText(record.module) },
    { label: '操作动作', value: getActionText(record.action) },
    { label: '操作描述', value: record.description },
    { label: 'IP地址', value: record.ip },
    { label: '用户代理', value: record.user_agent },
    { label: '请求数据', value: record.request_data },
    { label: '状态', value: record.status === 1 ? '成功' : '失败' },
    { label: '错误信息', value: record.error_msg || '无' },
    { label: '操作时间', value: formatTime(record.created_at) },
  ];
  detailVisible.value = true;
};

// 删除单个记录
const handleDelete = (record: any) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这条操作记录吗？',
    onOk: async () => {
      try {
        await deleteOperationLog(record.id);
        Message.success('删除成功');
        fetchLogs();
      } catch (error: any) {
        Message.error(error.message || '删除失败');
      }
    },
  });
};

// 批量删除
const handleBatchDelete = () => {
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 条操作记录吗？`,
    onOk: async () => {
      try {
        await batchDeleteOperationLogs({ ids: selectedRowKeys.value });
        Message.success('批量删除成功');
        selectedRowKeys.value = [];
        fetchLogs();
      } catch (error: any) {
        Message.error(error.message || '批量删除失败');
      }
    },
  });
};

// 清空操作记录
const handleClearLogs = () => {
  clearVisible.value = true;
};

const confirmClearLogs = async () => {
  try {
    await clearOperationLogs({ before_days: clearForm.beforeDays });
    Message.success('清空操作记录成功');
    clearVisible.value = false;
    fetchLogs();
  } catch (error: any) {
    Message.error(error.message || '清空操作记录失败');
  }
};

// 页面加载时获取数据
onMounted(() => {
  fetchLogs();
});
</script>

<style scoped>
.operation-management {
  padding: 16px;
}
</style>
