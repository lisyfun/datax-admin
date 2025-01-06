<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <img src="@/assets/logo.svg" alt="logo" class="logo" />
        <h1 class="title">DataX Admin</h1>
        <p class="subtitle">数据集成管理平台</p>
      </div>
      <a-form
        ref="formRef"
        :model="form"
        class="login-form"
        @submit="handleSubmit"
      >
        <a-form-item
          field="username"
          :rules="[{ required: true, message: '请输入用户名' }]"
          hide-label
        >
          <a-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
          >
            <template #prefix>
              <icon-user />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item
          field="password"
          :rules="[{ required: true, message: '请输入密码' }]"
          hide-label
        >
          <a-input-password
            v-model="form.password"
            placeholder="密码"
            size="large"
            allow-clear
          >
            <template #prefix>
              <icon-lock />
            </template>
          </a-input-password>
        </a-form-item>
        <div class="login-options">
          <a-checkbox v-model="rememberMe">记住我</a-checkbox>
          <a-link>忘记密码？</a-link>
        </div>
        <a-button
          type="primary"
          html-type="submit"
          :loading="loading"
          long
          size="large"
        >
          登录
        </a-button>
        <div class="login-footer">
          <a-space>
            <a-link @click="handleRegister">注册账号</a-link>
            <a-divider direction="vertical" />
            <a-button
              type="text"
              size="small"
              @click="toggleTheme"
            >
              <template #icon>
                <icon-moon-fill v-if="isDarkMode" />
                <icon-sun-fill v-else />
              </template>
              {{ isDarkMode ? '暗色' : '亮色' }}主题
            </a-button>
          </a-space>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { IconUser, IconLock, IconMoonFill, IconSunFill } from '@arco-design/web-vue/es/icon';
import { login } from '@/api/user';

const router = useRouter();
const formRef = ref();
const loading = ref(false);
const rememberMe = ref(false);
const isDarkMode = ref(false);

const form = reactive({
  username: '',
  password: '',
});

// 初始化主题
const initTheme = () => {
  const theme = localStorage.getItem('theme') || 'light';
  isDarkMode.value = theme === 'dark';
  applyTheme(isDarkMode.value);
};

// 切换主题
const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  applyTheme(isDarkMode.value);
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light');
};

// 应用主题
const applyTheme = (dark: boolean) => {
  if (dark) {
    document.body.setAttribute('arco-theme', 'dark');
  } else {
    document.body.removeAttribute('arco-theme');
  }
};

// 处理登录
const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    loading.value = true;

    const res = await login({
      username: form.username,
      password: form.password,
    });

    if (rememberMe.value) {
      localStorage.setItem('username', form.username);
    } else {
      localStorage.removeItem('username');
    }

    localStorage.setItem('token', res.data.token);
    Message.success('登录成功');
    router.push('/');
  } catch (err: any) {
    if (err.response?.data?.error) {
      Message.error(err.response.data.error);
    } else if (err.errors) {
      Message.error('请输入用户名和密码');
    } else {
      Message.error('登录失败');
    }
  } finally {
    loading.value = false;
  }
};

// 跳转注册页
const handleRegister = () => {
  router.push('/register');
};

onMounted(() => {
  initTheme();
  // 如果记住了用户名，自动填充
  const savedUsername = localStorage.getItem('username');
  if (savedUsername) {
    form.username = savedUsername;
    rememberMe.value = true;
  }
});
</script>

<style lang="less" scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--color-neutral-2);
  transition: all 0.2s;
}

.login-box {
  width: 100%;
  max-width: 480px;
  padding: 40px;
  margin: 0 20px;
  background: var(--color-bg-2);
  border-radius: 8px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;

  .logo {
    width: 80px;
    height: 80px;
    margin-bottom: 16px;
  }

  .title {
    margin: 0;
    font-size: 28px;
    font-weight: 500;
    color: var(--color-text-1);
  }

  .subtitle {
    margin: 8px 0 0;
    font-size: 16px;
    color: var(--color-text-3);
  }
}

.login-form {
  :deep(.arco-form-item) {
    margin-bottom: 24px;
  }

  :deep(.arco-input-wrapper) {
    background-color: var(--color-fill-2);
    border: 1px solid var(--color-neutral-3);
    transition: all 0.2s;

    &:hover {
      background-color: var(--color-fill-3);
      border-color: rgb(var(--primary-6));
    }

    &.arco-input-focus {
      background-color: var(--color-bg-2);
      border-color: rgb(var(--primary-6));
      box-shadow: 0 0 0 2px rgba(var(--primary-6), 0.2);
    }
  }
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.login-footer {
  margin-top: 24px;
  text-align: center;
}

:deep([arco-theme='dark']) {
  .login-container {
    background: var(--color-bg-1);
  }

  .login-box {
    background: var(--color-bg-2);
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);
  }

  .login-form {
    :deep(.arco-input-wrapper) {
      background-color: var(--color-fill-2);
      border-color: var(--color-neutral-4);

      &:hover {
        background-color: var(--color-fill-3);
        border-color: rgb(var(--primary-6));
      }

      &.arco-input-focus {
        background-color: var(--color-bg-2);
        border-color: rgb(var(--primary-6));
        box-shadow: 0 0 0 2px rgba(var(--primary-6), 0.2);
      }
    }
  }
}
</style>
