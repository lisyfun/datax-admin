<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <img src="../../assets/logo.svg" alt="logo" class="logo" />
        <h1>DataX Admin</h1>
      </div>
      <a-card class="login-card" :style="{ width: '400px' }">
        <template #title>
          <div class="login-title">系统登录</div>
        </template>
        <a-form :model="form" @submit="handleSubmit">
          <a-form-item field="username" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
            <a-input v-model="form.username" placeholder="请输入用户名" />
          </a-form-item>
          <a-form-item field="password" label="密码" :rules="[{ required: true, message: '请输入密码' }]">
            <a-input-password v-model="form.password" placeholder="请输入密码" />
          </a-form-item>
          <a-form-item>
            <a-space :size="16" direction="vertical" fill>
              <a-button type="primary" html-type="submit" long :loading="loading">
                登录
              </a-button>
              <a-button type="text" long @click="$router.push('/register')">
                还没有账号？立即注册
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>
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
});

const loading = ref(false);

const handleSubmit = async () => {
  try {
    loading.value = true;
    await userStore.login(form.username, form.password);
    Message.success('登录成功');
    router.push('/');
  } catch (error: any) {
    Message.error(error.response?.data?.error || '登录失败');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--color-fill-2);
}

.login-box {
  background-color: var(--color-bg-1);
  border-radius: 4px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  padding: 24px;
}

.login-header {
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

.login-title {
  font-size: 24px;
  font-weight: 500;
  color: var(--color-text-1);
  text-align: center;
}
</style>
