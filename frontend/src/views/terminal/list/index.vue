<template>
  <div class="terminals">
    <a-card>
      <template #title>终端管理</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchForm.name"
            placeholder="请输入终端名称"
            style="width: 200px"
            allow-clear
            @search="handleSearch"
            @press-enter="handleSearch"
          />
          <a-input-search
            v-model="searchForm.host"
            placeholder="请输入主机地址"
            style="width: 200px"
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
            <a-option value="online">
              <template #icon>
                <icon-check-circle-fill style="color: #00b42a" />
              </template>
              在线
            </a-option>
            <a-option value="offline">
              <template #icon>
                <icon-close-circle-fill style="color: #f53f3f" />
              </template>
              离线
            </a-option>
          </a-select>
          <a-button type="primary" @click="handleAdd">
            <template #icon><icon-plus /></template>
            新建终端
          </a-button>
          <a-button
            type="primary"
            status="success"
            @click="handleBatchUpload"
          >
            <template #icon><icon-upload /></template>
            批量上传 {{ selectedKeys.length ? `(${selectedKeys.length})` : '' }}
          </a-button>
          <a-button @click="() => fetchData()">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        row-key="id"
        :loading="loading"
        :data="tableData"
        :pagination="pagination"
        :bordered="false"
        :stripe="true"
        :hover="true"
        :scroll="{ x: '100%' }"
        :row-selection="{
          type: 'checkbox',
          showCheckedAll: true
        }"
        v-model:selected-keys="selectedKeys"
        @selection-change="onSelectionChange"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #columns>
          <a-table-column title="终端ID" data-index="id" :width="70" align="center">
            <template #cell="{ record }">
              <a-tag size="small" class="id-tag">{{ record.id }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="终端名称" data-index="name" :width="120">
            <template #cell="{ record }">
              <div class="name-cell">
                <icon-robot class="icon" />
                <span>{{ record.name }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="主机地址" data-index="host" :width="160">
            <template #cell="{ record }">
              <div class="host-cell">
                <icon-cloud class="icon" />
                <span>{{ record.host }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="端口" data-index="port" :width="80" align="center">
            <template #cell="{ record }">
              <a-tag size="small" color="arcoblue">{{ record.port }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="用户名" data-index="username" :width="100">
            <template #cell="{ record }">
              <div class="username-cell">
                <icon-user class="icon" />
                <span>{{ record.username }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="状态" data-index="status" :width="80" align="center">
            <template #cell="{ record }">
              <a-badge
                :status="record.status === 'online' ? 'success' : 'danger'"
                :text="record.status === 'online' ? '在线' : '离线'"
              />
            </template>
          </a-table-column>
          <a-table-column title="最后在线" data-index="lastSeen" :width="150">
            <template #cell="{ record }">
              <div class="time-cell">
                <icon-clock-circle class="icon" />
                <span>{{ formatDate(record.lastSeen) }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="创建时间" data-index="createdAt" :width="150">
            <template #cell="{ record }">
              <div class="time-cell">
                <icon-calendar class="icon" />
                <span>{{ formatDate(record.createdAt) }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="操作" align="center" :width="140">
            <template #cell="{ record }">
              <a-space size="mini">
                <a-tooltip content="连接终端">
                  <a-button
                    type="primary"
                    size="mini"
                    shape="circle"
                    @click="handleConnect(record)"
                    class="action-button"
                  >
                    <template #icon>
                      <icon-link />
                    </template>
                  </a-button>
                </a-tooltip>
                <a-tooltip content="上传文件">
                  <a-button
                    type="primary"
                    size="mini"
                    shape="circle"
                    status="success"
                    @click="handleUpload(record)"
                    class="action-button"
                  >
                    <template #icon>
                      <icon-upload />
                    </template>
                  </a-button>
                </a-tooltip>
                <a-tooltip content="编辑信息">
                  <a-button
                    type="primary"
                    size="mini"
                    shape="circle"
                    status="warning"
                    @click="handleEdit(record)"
                    class="action-button"
                  >
                    <template #icon>
                      <icon-edit />
                    </template>
                  </a-button>
                </a-tooltip>
                <a-tooltip content="删除终端">
                  <a-popconfirm
                    content="确定要删除这个终端吗？"
                    type="warning"
                    position="left"
                    @ok="handleDelete(record)"
                  >
                    <a-button
                      type="primary"
                      size="mini"
                      shape="circle"
                      status="danger"
                      class="action-button"
                    >
                      <template #icon>
                        <icon-delete />
                      </template>
                    </a-button>
                  </a-popconfirm>
                </a-tooltip>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>

    <!-- 终端表单对话框 -->
    <a-modal
      v-model:visible="visible"
      :title="formTitle"
      @ok="handleSubmit"
      @cancel="handleCancel"
      :mask-closable="false"
      :unmount-on-close="true"
      :width="480"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-position="right"
        :label-col-props="{ span: 5 }"
        :wrapper-col-props="{ span: 19 }"
      >
        <a-form-item field="name" label="终端名称" :rules="[{ required: true, message: '请输入终端名称' }]">
          <a-input
            v-model="formData.name"
            placeholder="请输入终端名称"
            allow-clear
          >
            <template #prefix>
              <icon-tag />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item field="host" label="主机地址" :rules="[{ required: true, message: '请输入主机地址' }]">
          <a-input
            v-model="formData.host"
            placeholder="请输入主机地址"
            allow-clear
          >
            <template #prefix>
              <icon-cloud />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item field="port" label="SSH端口" :rules="[{ required: true, message: '请输入SSH端口' }]">
          <a-input-number
            v-model="formData.port"
            placeholder="请输入SSH端口"
            :min="1"
            :max="65535"
            :default-value="22"
            style="width: 100%"
            mode="button"
          />
        </a-form-item>
        <a-form-item field="username" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input
            v-model="formData.username"
            placeholder="请输入用户名"
            allow-clear
          >
            <template #prefix>
              <icon-user />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password
            v-model="formData.password"
            :placeholder="formData.id ? '不修改请留空' : '请输入密码'"
            allow-clear
          >
            <template #prefix>
              <icon-lock />
            </template>
          </a-input-password>
          <template #extra>
            <span class="form-extra-tip">{{ formData.id ? '如不修改密码请留空' : '请输入SSH登录密码' }}</span>
          </template>
        </a-form-item>
      </a-form>
      <template #footer>
        <div class="modal-footer">
          <a-space>
            <a-button @click="handleCancel">取消</a-button>
            <a-button type="primary" :loading="submitLoading" @click="handleSubmit">
              {{ formData.id ? '更新' : '创建' }}
            </a-button>
          </a-space>
        </div>
      </template>
    </a-modal>

    <!-- 文件上传对话框 -->
    <a-modal
      v-model:visible="uploadVisible"
      title="上传文件"
      @cancel="handleUploadCancel"
      :mask-closable="false"
      :unmount-on-close="true"
      :width="600"
    >
      <a-form :model="{ path: uploadPath }" layout="vertical">
        <a-form-item field="path" label="上传路径">
          <a-input
            v-model="uploadPath"
            placeholder="请输入文件上传路径"
            allow-clear
          >
            <template #prefix>
              <icon-folder />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item field="files" label="选择文件">
          <a-upload
            v-model:file-list="fileList"
            :custom-request="(options: RequestOption) => ({ abort: () => {} })"
            :auto-upload="false"
            multiple
          >
            <template #upload-button>
              <a-button type="primary">
                <template #icon>
                  <icon-upload />
                </template>
                选择文件
              </a-button>
            </template>
          </a-upload>
        </a-form-item>
        <a-divider v-if="uploadRecords.length > 0">上传记录</a-divider>
        <div v-if="uploadRecords.length > 0" class="upload-records">
          <a-list :data="uploadRecords">
            <template #item="{ item }">
              <a-list-item>
                <div class="upload-record-item">
                  <span class="terminal-name">{{ item.terminalName }}</span>
                  <span class="upload-status">
                    <a-tag :color="getStatusColor(item.status)">
                      {{ getStatusText(item.status) }}
                    </a-tag>
                  </span>
                  <span v-if="item.message" class="upload-message">{{ item.message }}</span>
                </div>
              </a-list-item>
            </template>
          </a-list>
        </div>
      </a-form>
      <template #footer>
        <div class="modal-footer">
          <a-space>
            <a-button @click="handleUploadCancel">取消</a-button>
            <a-button
              type="primary"
              :loading="uploadLoading"
              @click="handleUploadSubmit"
              :disabled="uploadLoading"
            >
              {{ uploadLoading ? '上传中...' : '上传' }}
            </a-button>
          </a-space>
        </div>
      </template>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { FileItem, RequestOption } from '@arco-design/web-vue/es/upload/interfaces';
import {
  IconPlus,
  IconEdit,
  IconDelete,
  IconLink,
  IconSearch,
  IconRefresh,
  IconTag,
  IconCloud,
  IconUser,
  IconLock,
  IconCalendar,
  IconClockCircle,
  IconRobot,
  IconCheckCircleFill,
  IconCloseCircleFill,
  IconUpload,
  IconFolder,
} from '@arco-design/web-vue/es/icon';
import type { TerminalInfo } from '@/types/terminal';
import terminalApi from '@/api/terminal';
import { useRouter } from 'vue-router';

const loading = ref(false);
const submitLoading = ref(false);
const uploadVisible = ref(false);
const uploadLoading = ref(false);
const currentTerminal = ref<TerminalInfo | null>(null);
const fileList = ref<FileItem[]>([]);
const uploadPath = ref('/tmp');
const tableData = ref<TerminalInfo[]>([]);
const selectedKeys = ref<(string | number)[]>([]);
const uploadRecords = ref<{
  terminalId: number;
  terminalName: string;
  status: 'uploading' | 'success' | 'error';
  message?: string;
}[]>([]);
const router = useRouter();

const selectedCount = computed(() => selectedKeys.value.length);
const selectedIds = computed(() => selectedKeys.value.map(id => Number(id)));
const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showTotal: true,
  showJumper: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100],
});

// 搜索表单
const searchForm = reactive({
  name: '',
  host: '',
  status: '',
});

// 表单相关
const visible = ref(false);
const formTitle = ref('新建终端');
const formRef = ref();
const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  host: '',
  port: 22,
  username: '',
  password: '',
});

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入终端名称' }],
  host: [{ required: true, message: '请输入主机地址' }],
  port: [{ required: true, message: '请输入SSH端口' }],
  username: [{ required: true, message: '请输入用户名' }],
};

// 获取终端列表
const fetchData = async (page = 1, pageSize = pagination.pageSize) => {
  loading.value = true;
  try {
    const res = await terminalApi.getTerminalList({
      page,
      pageSize,
      ...searchForm,
    });
    tableData.value = res.data.list;
    pagination.total = res.data.total;
  } catch (error) {
    Message.error('获取终端列表失败');
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
  fetchData(current);
};

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  fetchData(1, pageSize);
};

// 格式化日期
const formatDate = (date: string) => {
  if (!date) return '-';
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  });
};

// 新建终端
const handleAdd = () => {
  formTitle.value = '新建终端';
  formData.id = undefined;
  formData.name = '';
  formData.host = '';
  formData.port = 22;
  formData.username = '';
  formData.password = '';
  visible.value = true;
};

// 编辑终端
const handleEdit = (record: TerminalInfo) => {
  formTitle.value = '编辑终端';
  formData.id = record.id;
  formData.name = record.name;
  formData.host = record.host;
  formData.port = record.port;
  formData.username = record.username;
  formData.password = '';
  visible.value = true;
};

// 删除终端
const handleDelete = async (record: TerminalInfo) => {
  try {
    await terminalApi.deleteTerminal(record.id);
    Message.success('删除成功');
    if (tableData.value.length === 1 && pagination.current > 1) {
      pagination.current -= 1;
    }
    fetchData(pagination.current);
  } catch (error) {
    Message.error('删除失败');
  }
};

// 连接终端
const handleConnect = (record: TerminalInfo) => {
  router.push(`/terminal/connect/${record.id}`);
};

// 打开上传对话框
const handleUpload = (record: TerminalInfo) => {
  currentTerminal.value = record;
  uploadPath.value = '/tmp';
  fileList.value = [];
  uploadVisible.value = true;
};

// 处理选择变化
const onSelectionChange = (keys: (string | number)[]) => {
  selectedKeys.value = keys;
};

// 打开批量上传对话框
const handleBatchUpload = () => {
  if (!selectedKeys.value.length) {
    Message.warning('请先选择要上传的终端');
    return;
  }
  uploadPath.value = '/tmp';
  fileList.value = [];
  uploadVisible.value = true;
};

// 处理文件上传
const handleUploadSubmit = async () => {
  if (!selectedKeys.value.length) {
    Message.warning('请选择要上传的终端');
    return;
  }

  if (fileList.value.length === 0) {
    Message.warning('请选择要上传的文件');
    return;
  }

  try {
    uploadLoading.value = true;
    const formData = new FormData();
    formData.append('path', uploadPath.value);
    fileList.value.forEach(file => {
      if (file.file) {
        formData.append('files', file.file);
      }
    });

    // 初始化上传记录
    uploadRecords.value = selectedKeys.value.map(id => {
      const terminal = tableData.value.find(t => t.id === Number(id));
      return {
        terminalId: Number(id),
        terminalName: terminal?.name || `终端${id}`,
        status: 'uploading'
      };
    });

    // 并行上传文件
    await Promise.all(
      uploadRecords.value.map(async record => {
        try {
          await terminalApi.uploadFiles(record.terminalId, formData);
          record.status = 'success';
          record.message = '上传成功';
        } catch (error) {
          record.status = 'error';
          record.message = error instanceof Error ? error.message : '上传失败';
        }
      })
    );

    // 检查是否全部成功
    const hasError = uploadRecords.value.some(record => record.status === 'error');
    if (hasError) {
      Message.warning('部分终端上传失败，请查看详情');
    } else {
      Message.success('所有文件上传成功');
    }
  } catch (error) {
    Message.error('文件上传失败');
  } finally {
    uploadLoading.value = false;
  }
};

// 获取状态颜色
const getStatusColor = (status: 'uploading' | 'success' | 'error') => {
  switch (status) {
    case 'uploading':
      return 'blue';
    case 'success':
      return 'green';
    case 'error':
      return 'red';
  }
};

// 获取状态文本
const getStatusText = (status: 'uploading' | 'success' | 'error') => {
  switch (status) {
    case 'uploading':
      return '上传中';
    case 'success':
      return '成功';
    case 'error':
      return '失败';
  }
};

// 取消上传
const handleUploadCancel = () => {
  uploadVisible.value = false;
  fileList.value = [];
  uploadPath.value = '/tmp';
  uploadRecords.value = [];
};

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    submitLoading.value = true;

    if (formData.id) {
      await terminalApi.updateTerminal(formData.id, {
        name: formData.name,
        host: formData.host,
        port: formData.port,
        username: formData.username,
        password: formData.password || undefined,
      });
      Message.success('更新成功');
    } else {
      await terminalApi.createTerminal({
        name: formData.name,
        host: formData.host,
        port: formData.port,
        username: formData.username,
        password: formData.password,
      });
      Message.success('创建成功');
    }
    visible.value = false;
    fetchData(pagination.current);
  } catch (error: unknown) {
    // 表单验证错误
    if (error && typeof error === 'object' && 'name' in error && error.name === 'FormValidateError') {
      return;
    }
    Message.error(formData.id ? '更新失败' : '创建失败');
  } finally {
    submitLoading.value = false;
  }
};

// 取消表单
const handleCancel = () => {
  visible.value = false;
  formRef.value.resetFields();
};

// 初始加载数据
fetchData();
</script>

<style lang="less" scoped>
.terminals {
  padding: 16px;

  :deep(.arco-card) {
    .arco-card-header {
      border-bottom: 1px solid var(--color-border);
    }
  }

  .id-tag {
    background-color: var(--color-fill-2);
    border: none;
    border-radius: 10px;
  }

  .name-cell,
  .host-cell,
  .username-cell,
  .time-cell {
    display: flex;
    align-items: center;

    .icon {
      margin-right: 8px;
      font-size: 16px;
      color: rgb(var(--primary-6));
    }
  }

  .action-button {
    transition: all 0.2s ease-in-out;

    &:hover {
      transform: scale(1.1);
    }
  }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid var(--color-border);
  background-color: var(--color-bg-2);
}

:deep(.arco-modal) {
  .arco-form-item-label-col {
    opacity: 0.8;
  }

  .arco-input-wrapper,
  .arco-select-view,
  .arco-input-number {
    transition: all 0.2s ease-in-out;

    &:hover {
      border-color: rgb(var(--primary-6));
    }

    &:focus-within {
      border-color: rgb(var(--primary-6));
      box-shadow: 0 0 0 2px rgba(var(--primary-6), 0.2);
    }
  }

  .form-extra-tip {
    color: var(--color-text-3);
    font-size: 12px;
  }
}

.upload-records {
  max-height: 200px;
  overflow-y: auto;
  margin: 0 -20px;
  padding: 0 20px;

  :deep(.arco-list-item) {
    padding: 8px 0;
  }
}

.upload-record-item {
  display: flex;
  align-items: center;
  width: 100%;

  .terminal-name {
    flex: 1;
    font-weight: 500;
  }

  .upload-status {
    margin: 0 12px;
  }

  .upload-message {
    color: var(--color-text-3);
    font-size: 12px;
  }
}
</style>
