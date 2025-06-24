<template>
  <a-modal
    v-model:visible="visible"
    :title="isEdit ? '编辑终端' : '新建终端'"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-align="right"
      :style="{ width: '100%' }"
      @submit="handleSubmit"
    >
      <a-form-item field="name" label="终端名称" :rules="[{ required: true, message: '请输入终端名称' }]">
        <a-input v-model="formData.name" placeholder="请输入终端名称" allow-clear />
      </a-form-item>
      <a-form-item field="host" label="主机地址" :rules="[{ required: true, message: '请输入主机地址' }]">
        <a-input v-model="formData.host" placeholder="请输入主机地址，例如：192.168.1.100" allow-clear />
      </a-form-item>
      <a-form-item field="port" label="SSH端口" :rules="[{ required: true, message: '请输入SSH端口' }]">
        <a-input-number
          v-model="formData.port"
          placeholder="请输入SSH端口"
          :min="1"
          :max="65535"
          :default-value="22"
        />
      </a-form-item>
      <a-form-item field="username" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
        <a-input v-model="formData.username" placeholder="请输入用户名" allow-clear />
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
        :rules="[{ required: formData.authType === 'password' && !props.isEdit, message: '请输入密码' }]"
      >
        <a-input-password v-model="formData.password" placeholder="请输入密码" allow-clear />
      </a-form-item>
      <a-form-item
        v-if="formData.authType === 'key'"
        field="keyFile"
        label="密钥文件"
        :rules="[{ required: formData.authType === 'key', message: '请上传密钥文件' }]"
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
        <a-input-password v-model="formData.keyPassphrase" placeholder="如果密钥文件有密码请输入" allow-clear />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, defineProps, defineEmits, watch } from 'vue';
import { Message } from '@arco-design/web-vue';
import { IconUpload } from '@arco-design/web-vue/es/icon';
import type { FieldRule } from '@arco-design/web-vue/es/form/interface';
import type { FileItem } from '@arco-design/web-vue/es/upload/interfaces';
import terminalApi from '@/api/terminal';
import type { TerminalInfo, CreateTerminalData, UpdateTerminalData } from '@/types/terminal';

const props = defineProps<{
  visible: boolean;
  isEdit: boolean;
  data?: TerminalInfo;
}>();

const emit = defineEmits(['update:visible', 'success']);

const formRef = ref();
const keyFileList = ref<FileItem[]>([]);
const keyFileName = ref('');

const formData = reactive<CreateTerminalData>({
  name: '',
  host: '',
  port: 22,
  username: '',
  authType: 'password',
  password: '',
  keyFile: '',
  keyPassphrase: '',
});

// 表单验证规则
const rules: Record<string, FieldRule[]> = {
  name: [
    { required: true, message: '请输入终端名称' },
  ],
  host: [
    { required: true, message: '请输入主机地址' },
  ],
  port: [
    { required: true, message: '请输入SSH端口' },
    { type: 'number', min: 1, max: 65535, message: '端口范围为1-65535' } as FieldRule,
  ],
  username: [
    { required: true, message: '请输入用户名' },
  ],
  authType: [
    { required: true, message: '请选择认证方式' },
  ],
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

// 监听数据变化
watch(
  () => props.data,
  (newVal) => {
    if (newVal) {
      const { name, host, port, username, authType, keyFile, keyPassphrase } = newVal;
      Object.assign(formData, {
        name,
        host,
        port,
        username,
        authType: authType || 'password',
        keyFile: keyFile || '',
        keyPassphrase: keyPassphrase || ''
      });
      if (keyFile) {
        keyFileName.value = '已设置密钥文件';
      }
    }
  },
  { immediate: true },
);

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate();

    // 验证认证方式对应的字段
    if (formData.authType === 'password' && !formData.password && !props.isEdit) {
      Message.error('请输入密码');
      return;
    }
    if (formData.authType === 'key' && !formData.keyFile) {
      Message.error('请选择密钥文件');
      return;
    }

    if (props.isEdit && props.data) {
      const updateData: UpdateTerminalData = {
        name: formData.name,
        host: formData.host,
        port: formData.port,
        username: formData.username,
        authType: formData.authType,
        keyFile: formData.keyFile,
        keyPassphrase: formData.keyPassphrase,
      };
      if (formData.password) {
        updateData.password = formData.password;
      }
      await terminalApi.updateTerminal(props.data.id, updateData);
      Message.success('更新成功');
    } else {
      await terminalApi.createTerminal(formData);
      Message.success('创建成功');
    }
    emit('success');
    handleCancel();
  } catch (error) {
    // 表单验证失败或请求失败
  }
};

// 取消
const handleCancel = () => {
  formRef.value.resetFields();
  keyFileList.value = [];
  keyFileName.value = '';
  formData.authType = 'password';
  formData.keyFile = '';
  formData.keyPassphrase = '';
  emit('update:visible', false);
};
</script>

<style scoped>
:deep(.arco-form-item-label-col) {
  min-width: 80px;
}
</style>
