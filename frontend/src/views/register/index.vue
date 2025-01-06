<template>
  <div class="register-container">
    <div class="register-box">
      <div class="register-header">
        <img src="@/assets/logo.svg" alt="logo" class="logo" />
        <h1 class="title">DATAX ADMIN</h1>
        <p class="subtitle">数据集成管理平台</p>
      </div>
      <a-form
        ref="formRef"
        :model="form"
        class="register-form"
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
        <a-form-item
          field="nickname"
          hide-label
        >
          <a-input
            v-model="form.nickname"
            placeholder="昵称"
            size="large"
          >
            <template #prefix>
              <icon-user-group />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item
          field="email"
          :rules="[{ type: 'email', message: '请输入正确的邮箱地址' }]"
          hide-label
        >
          <a-input
            v-model="form.email"
            placeholder="邮箱"
            size="large"
          >
            <template #prefix>
              <icon-email />
            </template>
          </a-input>
        </a-form-item>
        <a-button
          type="primary"
          html-type="submit"
          :loading="loading"
          long
          size="large"
        >
          注册
        </a-button>
        <div class="register-footer">
          <a-space>
            <a-link @click="$router.push('/login')">已有账号？立即登录</a-link>
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
import { IconUser, IconLock, IconUserGroup, IconEmail, IconMoonFill, IconSunFill } from '@arco-design/web-vue/es/icon';
import { useUserStore } from '@/stores/user';

const router = useRouter();
const userStore = useUserStore();
const formRef = ref();
const loading = ref(false);
const isDarkMode = ref(false);

const form = reactive({
  username: '',
  password: '',
  nickname: '',
  email: '',
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

const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    loading.value = true;
    await userStore.register(form.username, form.password, form.nickname, form.email);
    Message.success('注册成功，请登录');
    router.push('/login');
  } catch (error: any) {
    if (error.response?.data?.error) {
      Message.error(error.response.data.error);
    } else if (error.errors) {
      Message.error('请完善注册信息');
    } else {
      Message.error('注册失败');
    }
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  initTheme();
});
</script>

<style lang="less" scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--color-neutral-2);
  transition: all 0.2s;
}

.register-box {
  width: 100%;
  max-width: 480px;
  padding: 40px;
  margin: 0 20px;
  background: var(--color-bg-2);
  border-radius: 8px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
}

.register-header {
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
    font-weight: 600;
    letter-spacing: 2px;
    color: var(--color-text-1);
  }

  .subtitle {
    margin: 8px 0 0;
    font-size: 16px;
    color: var(--color-text-3);
  }
}

.register-form {
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

.register-footer {
  margin-top: 24px;
  text-align: center;
}

:deep([arco-theme='dark']) {
  .register-container {
    background: var(--color-bg-1);
  }

  .register-box {
    background: var(--color-bg-2);
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);
  }

  .register-form {
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
