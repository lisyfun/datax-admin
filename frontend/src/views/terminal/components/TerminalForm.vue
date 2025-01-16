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
      <a-form-item field="password" label="密码" :rules="[{ required: true, message: '请输入密码' }]">
        <a-input-password v-model="formData.password" placeholder="请输入密码" allow-clear />
      </a-form-item>
      <a-form-item field="description" label="描述">
        <a-textarea
          v-model="formData.description"
          placeholder="请输入终端描述"
          :max-length="200"
          show-word-limit
          allow-clear
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, defineProps, defineEmits, watch } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { FieldRule } from '@arco-design/web-vue/es/form/interface';
import * as terminalApi from '@/api/terminal';
import type { Terminal } from '@/api/terminal';

const props = defineProps<{
  visible: boolean;
  isEdit: boolean;
  data?: Terminal;
}>();

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void;
  (e: 'success'): void;
}>();

const formRef = ref();
const formData = reactive<Omit<Terminal, 'id' | 'status' | 'created_at' | 'updated_at'>>({
  name: '',
  host: '',
  port: 22,
  username: '',
  password: '',
  description: '',
});

// 表单校验规则
const rules = {
  name: [
    { required: true, message: '请输入终端名称' },
    { maxLength: 50, message: '终端名称不能超过50个字符' },
  ],
  host: [
    { required: true, message: '请输入主机地址' },
    {
      match: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$|^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$/,
      message: '请输入正确的IP地址或域名'
    },
  ],
  port: [
    { required: true, message: '请输入SSH端口' },
    { validator: (value: number) => value >= 1 && value <= 65535, message: '端口范围为1-65535' },
  ],
  username: [
    { required: true, message: '请输入用户名' },
    { maxLength: 32, message: '用户名不能超过32个字符' },
  ],
  password: [
    { required: true, message: '请输入密码' },
    { minLength: 6, message: '密码不能少于6个字符' },
    { maxLength: 32, message: '密码不能超过32个字符' },
  ],
} as Record<string, FieldRule[]>;

// 重置表单
const resetForm = () => {
  formData.name = '';
  formData.host = '';
  formData.port = 22;
  formData.username = '';
  formData.password = '';
  formData.description = '';
  if (formRef.value) {
    formRef.value.resetFields();
  }
};

// 监听数据变化
const updateFormData = () => {
  if (props.data) {
    const { name, host, port, username, description } = props.data;
    Object.assign(formData, { name, host, port, username, description });
  } else {
    resetForm();
  }
};

// 监听visible变化
const visible = ref(props.visible);
watch(() => props.visible, (val) => {
  visible.value = val;
  if (val) {
    updateFormData();
  }
});

watch(visible, (val) => {
  emit('update:visible', val);
  if (!val) {
    resetForm();
  }
});

const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    if (props.isEdit && props.data?.id) {
      await terminalApi.updateTerminal(props.data.id, formData);
    } else {
      await terminalApi.createTerminal(formData);
    }
    Message.success(`${props.isEdit ? '更新' : '创建'}成功`);
    visible.value = false;
    emit('success');
  } catch (error) {
    if (error instanceof Error) {
      Message.error(error.message);
    }
  }
};

const handleCancel = () => {
  visible.value = false;
};
</script>

<style scoped>
:deep(.arco-form-item-label-col) {
  min-width: 80px;
}
</style>
