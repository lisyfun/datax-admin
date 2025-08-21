<template>
  <div class="external-iframe-container">
    <div class="iframe-header">
      <h3>{{ title }}</h3>
      <a-button type="text" @click="goBack">
        <template #icon><IconClose /></template>
        关闭
      </a-button>
    </div>
    <div class="iframe-wrapper">
      <iframe
        :src="url"
        class="external-iframe"
        @load="onIframeLoad"
        @error="onIframeError"
      ></iframe>
      <div v-if="loading" class="iframe-loading">
        <a-spin size="large">
          <template #indicator>
            <IconLoading />
          </template>
        </a-spin>
        <p>正在加载页面...</p>
      </div>
      <div v-if="error" class="iframe-error">
        <a-result status="error" title="加载失败">
          <template #extra>
            <a-button type="primary" @click="retry">重试</a-button>
            <a-button @click="goBack">返回</a-button>
          </template>
        </a-result>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  IconClose,
  IconLoading,
} from '@arco-design/web-vue/es/icon';

const route = useRoute();
const router = useRouter();

const url = ref('');
const title = ref('外部页面');
const loading = ref(true);
const error = ref(false);

const onIframeLoad = () => {
  loading.value = false;
  error.value = false;
};

const onIframeError = () => {
  loading.value = false;
  error.value = true;
};

const retry = () => {
  error.value = false;
  loading.value = true;
  // 重新加载iframe
  const iframe = document.querySelector('.external-iframe') as HTMLIFrameElement;
  if (iframe) {
    iframe.src = iframe.src;
  }
};

const goBack = () => {
  router.go(-1);
};

onMounted(() => {
  // 从路由参数中获取URL和标题
  const query = route.query;
  if (query.url) {
    url.value = query.url as string;
  }
  if (query.title) {
    title.value = query.title as string;
  }

  // 如果没有URL，显示错误
  if (!url.value) {
    error.value = true;
    loading.value = false;
  }
});
</script>

<style scoped>
.external-iframe-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-1);
}

.iframe-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--color-border-2);
  background: var(--color-bg-1);
}

.iframe-header h3 {
  margin: 0;
  color: var(--color-text-1);
  font-weight: 600;
}

.iframe-wrapper {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.external-iframe {
  width: 100%;
  height: 100%;
  border: none;
  display: block;
}

.iframe-loading,
.iframe-error {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
}

.iframe-loading p {
  margin-top: 16px;
  color: var(--color-text-2);
}

.iframe-error {
  width: 100%;
  max-width: 400px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .iframe-header {
    padding: 12px 16px;
  }

  .iframe-header h3 {
    font-size: 16px;
  }
}
</style>