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
          <a-button type="primary" @click="handleAdd" v-permission="'terminal.list.create'">
            <template #icon><icon-plus /></template>
            新建终端
          </a-button>
          <a-button
            type="primary"
            status="success"
            @click="handleBatchUpload"
            v-permission="'terminal.list.create'"
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
        :bordered="true"
        :stripe="true"
        :fixedHeader="true"
        :hover="true"
        :scroll="{ x: '100%', y: '100%' }"
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
                    v-permission="'terminal.list.connect'"
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
                    v-permission="'terminal.list.upload'"
                  >
                    <template #icon>
                      <icon-upload />
                    </template>
                  </a-button>
                </a-tooltip>
                <a-tooltip content="下载文件">
                  <a-button
                    type="primary"
                    size="mini"
                    shape="circle"
                    status="warning"
                    @click="handleDownload(record)"
                    class="action-button"
                    v-permission="'terminal.list.download'"
                  >
                    <template #icon>
                      <icon-download />
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
                    v-permission="'terminal.list.update'"
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
                      v-permission="'terminal.list.delete'"
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
        <a-form-item field="authType" label="认证方式" :rules="[{ required: true, message: '请选择认证方式' }]">
          <a-radio-group v-model="formData.authType">
            <a-radio value="password">密码认证</a-radio>
            <a-radio value="key">密钥文件认证</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item
          v-if="formData.authType === 'password'"
          field="password"
          label="密码"
        >
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
        <a-form-item
          v-if="formData.authType === 'key'"
          field="keyFile"
          label="密钥文件"
        >
          <a-upload
            :file-list="keyFileList"
            :show-file-list="false"
            :auto-upload="false"
            @change="handleKeyFileChange"
            accept=".pem,.key,.ppk"
          >
            <template #upload-button>
              <a-button>
                <template #icon>
                  <icon-upload />
                </template>
                选择密钥文件
              </a-button>
            </template>
          </a-upload>
          <div v-if="keyFileName" style="margin-top: 8px; color: #165dff;">
            已选择: {{ keyFileName }}
          </div>
        </a-form-item>
        <a-form-item
          v-if="formData.authType === 'key'"
          field="keyPassphrase"
          label="密钥密码"
        >
          <a-input-password
            v-model="formData.keyPassphrase"
            placeholder="如果密钥文件有密码请输入"
            allow-clear
          >
            <template #prefix>
              <icon-lock />
            </template>
          </a-input-password>
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

    <!-- 上传对话框 -->
    <a-modal
      v-model:visible="uploadVisible"
      title="上传文件"
      @cancel="handleUploadCancel"
      :mask-closable="false"
      :unmount-on-close="true"
      :footer="false"
      :width="520"
    >
      <a-form :model="{ path: uploadPath }" layout="vertical">
        <a-form-item field="path" label="上传路径">
          <a-input
            v-model="uploadPath"
            placeholder="请输入上传路径"
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
            :custom-request="customRequest"
            :drag="true"
            multiple
            @success="handleUploadSuccess"
            @error="handleUploadError"
            @progress="handleUploadProgress"
          >
          </a-upload>
        </a-form-item>
        <a-divider v-if="uploadRecords.length > 0">上传记录</a-divider>
        <div v-if="uploadRecords.length > 0" class="upload-records">
          <a-list :data="uploadRecords">
            <template #item="{ item }">
              <a-list-item>
                <div class="upload-record-item">
                  <span class="terminal-name">{{ item.terminalName }}</span>
                  <span class="file-name">{{ item.fileName }}</span>
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
    </a-modal>

    <!-- 文件列表对话框 -->
    <a-modal
      v-model:visible="fileListVisible"
      title="选择要下载的文件"
      @cancel="handleFileListCancel"
      :mask-closable="false"
      :unmount-on-close="true"
      :width="800"
      :modal-style="{ height: '80vh' }"
      :body-style="{ height: 'calc(80vh - 120px)', padding: '0' }"
    >
      <div class="file-list-container">
        <div class="file-list-header">
          <a-space>
            <a-button type="outline" @click="handleNavigateUp" :disabled="currentPath === '.'">
              <template #icon><icon-arrow-up /></template>
              返回上级
            </a-button>
            <a-input
              v-model="currentPath"
              placeholder="当前路径"
              style="width: 300px"
              @press-enter="handlePathChange"
            >
              <template #prefix>
                <icon-folder />
              </template>
            </a-input>
          </a-space>
        </div>
        <div class="file-list-content" @scroll="handleScroll" ref="fileListContentRef">
          <div class="file-list-info" v-if="totalFiles > 0">
            <span>共 {{ totalFiles }} 个文件，已显示 {{ serverFileList.length }} 个</span>
          </div>
          <a-table
            :data="serverFileList"
            :loading="fileListLoading"
            :pagination="false"
            :bordered="false"
          >
            <template #columns>
              <a-table-column title="文件名" data-index="name" :width="300">
                <template #cell="{ record }">
                  <a-space>
                    <icon-folder v-if="record.isDir" />
                    <icon-file v-else />
                    <a-button
                      type="text"
                      @click="handleFileClick(record)"
                      class="file-name-button"
                    >
                      <span class="file-name-text" :title="record.name">{{ record.name }}</span>
                    </a-button>
                  </a-space>
                </template>
              </a-table-column>
              <a-table-column title="大小" data-index="size" :width="120">
                <template #cell="{ record }">
                  {{ formatFileSize(record.size) }}
                </template>
              </a-table-column>
              <a-table-column title="修改时间" data-index="modTime" :width="180">
                <template #cell="{ record }">
                  {{ formatDate(record.modTime) }}
                </template>
              </a-table-column>
              <a-table-column title="操作" :width="120" align="center">
                <template #cell="{ record }">
                  <a-button
                    v-if="!record.isDir"
                    type="primary"
                    size="mini"
                    @click="handleDownloadFile(record)"
                  >
                    下载
                  </a-button>
                </template>
              </a-table-column>
            </template>
          </a-table>

          <!-- 加载更多指示器 -->
          <div v-if="loadingMore" class="loading-more">
            <a-spin size="small" />
            <span style="margin-left: 8px;">加载更多文件...</span>
          </div>

          <!-- 手动加载更多按钮 -->
          <div v-else-if="hasMore && !loadingMore" class="load-more-button">
            <a-button type="outline" @click="loadMoreFiles" :loading="loadingMore">
              加载更多 (剩余 {{ totalFiles - serverFileList.length }} 个文件)
            </a-button>
          </div>

          <!-- 没有更多数据提示 -->
          <div v-else-if="!hasMore && serverFileList.length > 0" class="no-more-data">
            <span>已显示全部文件</span>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { FileItem, RequestOption, UploadRequest } from '@arco-design/web-vue/es/upload/interfaces';
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
  IconDownload,
  IconArrowUp,
  IconFile,
} from '@arco-design/web-vue/es/icon';
import type { TerminalInfo } from '@/types/terminal';
import type { FileInfo } from '@/api/terminal';
import type { TreeNodeData } from '@arco-design/web-vue/es/tree/interface';
import terminalApi from '@/api/terminal';
import { useRouter } from 'vue-router';

const loading = ref(false);
const submitLoading = ref(false);
const uploadVisible = ref(false);
const uploadLoading = ref(false);
const currentTerminal = ref<TerminalInfo | null>(null);
const fileList = ref<FileItem[]>([]);
const keyFileList = ref<FileItem[]>([]);
const keyFileName = ref('');
const uploadPath = ref('/tmp');
const tableData = ref<TerminalInfo[]>([]);
const selectedKeys = ref<(string | number)[]>([]);
const uploadRecords = ref<{
  terminalId: number;
  terminalName: string;
  fileName: string;
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
  authType: 'password' as 'password' | 'key',
  password: '',
  keyFile: '',
  keyPassphrase: '',
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

// 处理密钥文件上传
const handleKeyFileChange = (fileList: FileItem[]) => {
  if (fileList.length > 0) {
    const file = fileList[0];
    keyFileName.value = file.name || '';

    // 读取文件内容
    const reader = new FileReader();
    reader.onload = (e) => {
      formData.keyFile = e.target?.result as string;
    };
    reader.readAsText(file.file as File);
  } else {
    keyFileName.value = '';
    formData.keyFile = '';
  }
};

// 新建终端
const handleAdd = () => {
  formTitle.value = '新建终端';
  formData.id = undefined;
  formData.name = '';
  formData.host = '';
  formData.port = 22;
  formData.username = '';
  formData.authType = 'password';
  formData.password = '';
  formData.keyFile = '';
  formData.keyPassphrase = '';
  keyFileList.value = [];
  keyFileName.value = '';
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
  formData.authType = record.authType || 'password';
  formData.password = '';
  formData.keyFile = record.keyFile || '';
  formData.keyPassphrase = record.keyPassphrase || '';
  keyFileList.value = [];
  if (record.keyFile) {
    keyFileName.value = '已设置密钥文件';
  } else {
    keyFileName.value = '';
  }
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
  const terminalPath = `/datax/#/terminal/connect/${record.id}`;
  const terminalUrl = `${window.location.origin}${terminalPath}`;
  // 设置打开新窗口标题
  const windowTitle = `终端连接 - ${record.name}`;
  let newWindow = window.open(terminalUrl, '', 'width=1024,height=768,menubar=no,toolbar=no,location=no,status=no');
  if (newWindow) {
    setTimeout(() => newWindow!.document.title = windowTitle, 300);
  }
};


// 打开上传对话框
const handleUpload = (record: TerminalInfo) => {
  selectedKeys.value = [record.id];
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

// 上传记录类型
interface UploadRecord {
  terminalId: number;
  terminalName: string;
  fileName: string;
  status: 'uploading' | 'success' | 'error';
  message?: string;
}

// 处理文件上传请求
const customRequest = (option: RequestOption): UploadRequest => {
  const { fileItem, onProgress, onSuccess, onError } = option;
  if (!fileItem.file) {
    onError();
    return {
      abort: () => {}
    };
  }

  const file = fileItem.file;

  // 为每个选中的终端创建上传记录
  selectedKeys.value.forEach(terminalId => {
    const terminal = tableData.value.find(t => t.id === Number(terminalId));
    uploadRecords.value.push({
      terminalId: Number(terminalId),
      terminalName: terminal?.name || `终端${terminalId}`,
      fileName: file.name,
      status: 'uploading'
    });
  });

  // 准备上传数据
  const formData = new FormData();
  formData.append('path', uploadPath.value);
  formData.append('files', file);

  // 发送上传请求
  selectedKeys.value.forEach(async (terminalId) => {
    try {
      await terminalApi.uploadFiles(Number(terminalId), formData);
      // 更新上传记录状态为成功
      const record = uploadRecords.value.find(
        r => r.terminalId === Number(terminalId) && r.fileName === file.name
      );
      if (record) {
        record.status = 'success';
      }
      onSuccess();
    } catch (error) {
      // 更新上传记录状态为失败
      const record = uploadRecords.value.find(
        r => r.terminalId === Number(terminalId) && r.fileName === file.name
      );
      if (record) {
        record.status = 'error';
        record.message = error instanceof Error ? error.message : '上传失败';
      }
      onError();
    }
  });

  return {
    abort: () => {
      console.log('上传被取消');
    }
  };
};

// 处理上传成功
const handleUploadSuccess = (fileItem: FileItem) => {
  console.log('上传成功:', fileItem);
};

// 处理上传失败
const handleUploadError = (fileItem: FileItem) => {
  console.log('上传失败:', fileItem);
};

// 处理上传进度
const handleUploadProgress = (fileItem: FileItem, ev?: ProgressEvent) => {
  if (ev) {
    console.log('上传进度:', fileItem, ev);
  }
};

// 取消上传
const handleUploadCancel = () => {
  uploadVisible.value = false;
  fileList.value = [];
  uploadPath.value = '/tmp';
  uploadRecords.value = [];
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
        authType: formData.authType,
        password: formData.password || undefined,
        keyFile: formData.keyFile,
        keyPassphrase: formData.keyPassphrase,
      });
      Message.success('更新成功');
    } else {
      await terminalApi.createTerminal({
        name: formData.name,
        host: formData.host,
        port: formData.port,
        username: formData.username,
        authType: formData.authType,
        password: formData.password,
        keyFile: formData.keyFile,
        keyPassphrase: formData.keyPassphrase,
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
  keyFileList.value = [];
  keyFileName.value = '';
  formData.authType = 'password';
  formData.keyFile = '';
  formData.keyPassphrase = '';
};

// 文件列表相关
const fileListVisible = ref(false);
const fileListLoading = ref(false);
const loadingMore = ref(false);
const serverFileList = ref<FileInfo[]>([]);
const currentPath = ref('.');
const currentPage = ref(1);
const pageSize = ref(100);
const totalFiles = ref(0);
const hasMore = ref(false);
const fileListContentRef = ref<HTMLElement>();

// 格式化文件大小
const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B';
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB';
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(2) + ' MB';
  return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
};

// 获取文件列表（重置列表）
const fetchFileList = async (path: string) => {
  if (!currentTerminal.value) return;

  fileListLoading.value = true;
  currentPage.value = 1;

  try {
    const res = await terminalApi.getFileList(currentTerminal.value.id, path, 1, pageSize.value);
    console.log('API响应:', res);
    console.log('文件列表数据:', res.data.data);
    console.log('文件数量:', res.data.data?.length || 0);

    serverFileList.value = res.data.data || [];
    totalFiles.value = res.data.total || 0;
    hasMore.value = res.data.hasMore || false;
    currentPage.value = res.data.page || 1;
  } catch (error) {
    console.error('获取文件列表失败:', error);
    Message.error('获取文件列表失败');
    serverFileList.value = [];
    totalFiles.value = 0;
    hasMore.value = false;
  } finally {
    fileListLoading.value = false;
  }
};

// 加载更多文件
const loadMoreFiles = async () => {
  if (!currentTerminal.value || !hasMore.value || loadingMore.value) return;

  loadingMore.value = true;
  const nextPage = currentPage.value + 1;

  try {
    const res = await terminalApi.getFileList(currentTerminal.value.id, currentPath.value, nextPage, pageSize.value);

    // 追加新数据到现有列表
    serverFileList.value.push(...(res.data.data || []));
    hasMore.value = res.data.hasMore || false;
    currentPage.value = res.data.page || nextPage;

    console.log(`加载第${nextPage}页，新增${res.data.data?.length || 0}个文件`);
  } catch (error) {
    console.error('加载更多文件失败:', error);
    Message.error('加载更多文件失败');
  } finally {
    loadingMore.value = false;
  }
};

// 处理滚动事件
const handleScroll = (event: Event) => {
  const target = event.target as HTMLElement;
  const { scrollTop, scrollHeight, clientHeight } = target;

  // 当滚动到底部附近时（距离底部50px以内）触发加载更多
  if (scrollHeight - scrollTop - clientHeight < 50) {
    loadMoreFiles();
  }
};

// 处理文件点击
const handleFileClick = (file: FileInfo) => {
  if (file.isDir) {
    currentPath.value = file.path;
    fetchFileList(file.path);
  }
};

// 处理路径变化
const handlePathChange = () => {
  fetchFileList(currentPath.value);
};

// 返回上级目录
const handleNavigateUp = () => {
  const parentPath = currentPath.value.split('/').slice(0, -1).join('/') || '.';
  currentPath.value = parentPath;
  fetchFileList(parentPath);
};

// 打开文件列表对话框
const handleDownload = async (record: TerminalInfo) => {
  currentTerminal.value = record;
  currentPath.value = '.';
  fileListVisible.value = true;
  await fetchFileList('.');
};

// 下载文件
const handleDownloadFile = async (file: FileInfo) => {
  if (!currentTerminal.value) return;

  try {
    Message.info({
      content: `开始下载文件: ${file.name}`,
      position: 'top'
    });

    const response = await terminalApi.downloadFile(currentTerminal.value.id, file.path);

    // 创建 Blob 对象
    const blob = new Blob([response.data], { type: 'application/octet-stream' });

    // 创建下载链接
    const downloadUrl = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = file.name;

    // 触发下载
    document.body.appendChild(link);
    link.click();

    // 清理
    window.URL.revokeObjectURL(downloadUrl);
    document.body.removeChild(link);

    Message.success({
      content: `文件 ${file.name} 下载成功`,
      position: 'top'
    });
  } catch (error) {
    console.error('下载文件失败:', error);
    Message.error({
      content: `下载失败: ${error instanceof Error ? error.message : '未知错误'}`,
      position: 'top'
    });
  }
};

// 关闭文件列表对话框
const handleFileListCancel = () => {
  fileListVisible.value = false;
  serverFileList.value = [];
  currentPath.value = '.';
  currentTerminal.value = null;
  // 重置分页状态
  currentPage.value = 1;
  totalFiles.value = 0;
  hasMore.value = false;
  loadingMore.value = false;
};

// 初始加载数据
fetchData();
</script>

<style lang="less" scoped>
.terminals {
  padding: 16px;
  height: calc(100vh - 80px); /* 减去header和padding的高度 */
  display: flex;
  flex-direction: column;
  overflow: hidden;

  :deep(.arco-card) {
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .arco-card-header {
      border-bottom: 1px solid var(--color-border);
      flex-shrink: 0;
      padding: 16px 20px;
    }
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
  gap: 12px;

  .terminal-name {
    flex: 0 0 120px;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .file-name {
    flex: 1;
    color: var(--color-text-2);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .upload-status {
    flex: 0 0 auto;
  }

  .upload-message {
    flex: 0 0 auto;
    color: var(--color-text-3);
    font-size: 12px;
  }
}

.arco-upload-slide {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 150px;
  border: 2px dashed var(--color-border);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease-in-out;

  &:hover {
    border-color: rgb(var(--primary-6));
    background-color: var(--color-fill-2);
  }

  .arco-upload-slide-text {
    display: flex;
    flex-direction: column;
    align-items: center;
    color: var(--color-text-3);

    .icon {
      font-size: 24px;
      margin-bottom: 8px;
    }

    p {
      margin: 0;
    }
  }
}

.file-list-container {
  height: 100%;
  display: flex;
  flex-direction: column;

  .file-list-header {
    padding: 16px;
    border-bottom: 1px solid var(--color-border);
  }

  .file-list-content {
    flex: 1;
    overflow: auto;
    padding: 16px;

    :deep(.arco-table) {
      .arco-table-th {
        background-color: var(--color-fill-2);
      }
    }

    .file-list-info {
      padding: 8px 16px;
      background-color: var(--color-fill-1);
      border-radius: 4px;
      margin-bottom: 16px;
      font-size: 12px;
      color: var(--color-text-3);
    }

    .loading-more {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 16px;
      color: var(--color-text-3);
    }

    .load-more-button {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 16px;
    }

    .no-more-data {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 16px;
      color: var(--color-text-4);
      font-size: 12px;
    }

    .file-name-button {
      width: 100%;
      text-align: left;
      padding: 0;

      .file-name-text {
        display: block;
        width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        max-width: 240px; /* 留出图标和间距的空间 */
      }
    }
  }
}
</style>
