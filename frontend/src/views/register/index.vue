<template>
  <div class="register-container">
    <div class="register-box">
      <div class="register-header">
        <img src="../../assets/logo.svg" alt="logo" class="logo" />
        <h1>DataX Admin</h1>
      </div>
      <a-form :model="form" @submit="handleSubmit">
        <a-form-item field="username" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model="form.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item field="password" label="密码" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model="form.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item field="nickname" label="昵称">
          <a-input v-model="form.nickname" placeholder="请输入昵称" />
        </a-form-item>
        <a-form-item field="email" label="邮箱" :rules="[{ type: 'email', message: '请输入正确的邮箱地址' }]">
          <a-input v-model="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item>
          <a-space :size="16" direction="vertical" fill>
            <a-button type="primary" html-type="submit" long :loading="loading">
              注册
            </a-button>
            <a-button type="text" long @click="$router.push('/login')">
              已有账号？立即登录
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { useUserStore } from '../../stores/user';

const router = useRouter();
const userStore = useUserStore();

const form = reactive({
  username: '',
  password: '',
  nickname: '',
  email: '',
});

const loading = ref(false);

const handleSubmit = async () => {
  try {
    loading.value = true;
    await userStore.register(form.username, form.password, form.nickname, form.email);
    Message.success('注册成功，请登录');
    router.push('/login');
  } catch (error: any) {
    Message.error(error.response?.data?.error || '注册失败');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--color-fill-2);
}

.register-box {
  background-color: var(--color-bg-1);
  border-radius: 4px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  padding: 24px;
  width: 400px;
}

.register-header {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
}

.logo {
  width: 40px;
  height: 40px;
  margin-right: 8px;
}

h1 {
  margin: 0;
  font-size: 24px;
  color: var(--color-text-1);
}
</style>
