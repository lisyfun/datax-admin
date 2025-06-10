<template>
  <a-modal
    :visible="visible"
    :title="formType === 'add' ? '新增用户' : '编辑用户'"
    @cancel="handleCancel"
    @ok="handleSubmit"
    :confirmLoading="loading"
  >
    <a-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-col="{ span: 6 }"
      wrapper-col="{ span: 18 }"
    >
      <a-form-item label="用户名" field="username">
        <a-input
          v-model="form.username"
          placeholder="请输入用户名"
          :disabled="formType === 'edit'"
        />
      </a-form-item>
      <a-form-item
        label="密码"
        field="password"
        :rules="formType === 'add' ? [{ required: true, message: '请输入密码' }] : undefined"
      >
        <a-input-password
          v-model="form.password"
          placeholder="请输入密码"
          allow-clear
        />
      </a-form-item>
      <a-form-item label="昵称" field="nickname">
        <a-input
          v-model="form.nickname"
          placeholder="请输入昵称"
          allow-clear
        />
      </a-form-item>
      <a-form-item label="邮箱" field="email">
        <a-input
          v-model="form.email"
          placeholder="请输入邮箱"
          allow-clear
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, watch } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { FormInstance } from '@arco-design/web-vue';
import * as userApi from '@/api/user';
import type { UserInfo } from '@/types/api';

const props = defineProps<{
  visible: boolean;
  formType: 'add' | 'edit';
  formData: UserInfo | null;
}>();

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void;
  (e: 'success'): void;
}>();

const formRef = ref<FormInstance>();
const loading = ref(false);

const form = reactive({
  username: '',
  password: '',
  nickname: '',
  email: '',
});

const rules = {
  username: [
    { required: true, message: '请输入用户名' },
    { minLength: 4, message: '用户名长度不能小于4个字符' },
  ],
};

watch(
  () => props.formData,
  (val) => {
    if (val) {
      form.username = val.username;
      form.nickname = val.nickname;
      form.email = val.email;
    } else {
      form.username = '';
      form.password = '';
      form.nickname = '';
      form.email = '';
    }
  },
  { immediate: true }
);

const handleCancel = () => {
  emit('update:visible', false);
};

const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    loading.value = true;

    if (props.formType === 'add') {
      await userApi.register({
        username: form.username,
        password: form.password,
        nickname: form.nickname,
        email: form.email,
      });
      Message.success('添加成功');
    } else {
      if (!props.formData) return;
      if (form.password) {
        await userApi.resetPassword(props.formData.id, { password: form.password });
      }
      await userApi.updateProfile({
        nickname: form.nickname,
        email: form.email,
      });
      Message.success('更新成功');
    }

    emit('success');
    emit('update:visible', false);
  } catch (error: any) {
    Message.error(error.response?.data?.error || '操作失败');
  } finally {
    loading.value = false;
  }
};
</script>
