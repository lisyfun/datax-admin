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
      <a-form-item field="password" label="密码" :rules="[{ required: !props.isEdit, message: '请输入密码' }]">
        <a-input-password v-model="formData.password" placeholder="请输入密码" allow-clear />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, defineProps, defineEmits, watch } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { FieldRule } from '@arco-design/web-vue/es/form/interface';
import terminalApi from '@/api/terminal';
import type { TerminalInfo, CreateTerminalData, UpdateTerminalData } from '@/types/terminal';

const props = defineProps<{
  visible: boolean;
  isEdit: boolean;
  data?: TerminalInfo;
}>();

const emit = defineEmits(['update:visible', 'success']);

const formRef = ref();
const formData = reactive<CreateTerminalData>({
  name: '',
  host: '',
  port: 22,
  username: '',
  password: '',
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
  password: [
    { required: !props.isEdit, message: '请输入密码' },
  ],
};

// 监听数据变化
watch(
  () => props.data,
  (newVal) => {
    if (newVal) {
      const { name, host, port, username } = newVal;
      Object.assign(formData, { name, host, port, username });
    }
  },
  { immediate: true },
);

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    if (props.isEdit && props.data) {
      const updateData: UpdateTerminalData = {
        name: formData.name,
        host: formData.host,
        port: formData.port,
        username: formData.username,
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
  emit('update:visible', false);
};
</script>

<style scoped>
:deep(.arco-form-item-label-col) {
  min-width: 80px;
}
</style>
