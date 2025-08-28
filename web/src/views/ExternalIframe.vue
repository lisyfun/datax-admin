<template>
  <div class="external-iframe-container">
    <div class="iframe-wrapper">
      <iframe
        :key="$route.fullPath"
        :src="url"
        class="external-iframe"
      ></iframe>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed} from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
/* 1️⃣ 用计算属性，每次点击都会带新的 _t，使 iframe 重新加载 */
const url = computed(() => {
  const raw = route.query.url as string;
  if (!raw) return '';
  const u = new URL(raw);          // 避免污染原始 url
  u.searchParams.set('_t', route.query._t as string || Date.now().toString());
  return u.toString();
});

/* 2️⃣ 重试：改一次 query 里的 _t 即可 */
const retry = () => {
  router.replace({
    name: route.name!,
    query: { ...route.query, _t: Date.now() }
  });
};

const goBack = () => {
  router.go(-1);
};

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