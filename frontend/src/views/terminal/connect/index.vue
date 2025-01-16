<template>
  <div class="container">
    <a-card class="terminal-card">
      <template #title>
        <a-space>
          <span>{{ terminalInfo?.name || '终端连接' }}</span>
          <a-tag :color="connected ? 'green' : 'red'">
            {{ connected ? '已连接' : '未连接' }}
          </a-tag>
        </a-space>
      </template>
      <template #extra>
        <a-space>
          <a-button type="text" @click="handleBack">
            <template #icon>
              <icon-left />
            </template>
            返回列表
          </a-button>
          <a-button
            type="primary"
            status="success"
            v-if="!connected"
            :loading="connecting"
            @click="handleConnect"
          >
            <template #icon>
              <icon-play-circle />
            </template>
            连接终端
          </a-button>
          <a-button
            type="primary"
            status="danger"
            v-else
            :loading="disconnecting"
            @click="handleDisconnect"
          >
            <template #icon>
              <icon-pause-circle />
            </template>
            断开连接
          </a-button>
        </a-space>
      </template>
      <div ref="terminalContainer" class="terminal-container"></div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { IconLeft, IconPlayCircle, IconPauseCircle } from '@arco-design/web-vue/es/icon';
import * as terminalApi from '@/api/terminal';
import type { Terminal } from '@/api/terminal';

const route = useRoute();
const router = useRouter();
const terminalContainer = ref<HTMLElement>();
const terminalInfo = ref<Terminal>();
const connected = ref(false);
const connecting = ref(false);
const disconnecting = ref(false);

// 获取终端信息
const fetchTerminalInfo = async () => {
  const id = route.params.id as string;
  try {
    const res = await terminalApi.getTerminalInfo(id);
    terminalInfo.value = res.data;
  } catch (error) {
    Message.error('获取终端信息失败');
    router.push('/terminal/list');
  }
};

// 连接终端
const handleConnect = async () => {
  if (!terminalInfo.value?.id) return;

  connecting.value = true;
  try {
    await terminalApi.connectTerminal(terminalInfo.value.id);
    connected.value = true;
    Message.success('连接成功');
    // TODO: 建立WebSocket连接，实现终端交互
  } catch (error) {
    Message.error('连接失败');
  } finally {
    connecting.value = false;
  }
};

// 断开连接
const handleDisconnect = async () => {
  if (!terminalInfo.value?.id) return;

  disconnecting.value = true;
  try {
    await terminalApi.disconnectTerminal(terminalInfo.value.id);
    connected.value = false;
    Message.success('已断开连接');
    // TODO: 关闭WebSocket连接
  } catch (error) {
    Message.error('断开连接失败');
  } finally {
    disconnecting.value = false;
  }
};

// 返回列表
const handleBack = () => {
  router.push('/terminal/list');
};

// 组件挂载时获取终端信息
onMounted(() => {
  fetchTerminalInfo();
});

// 组件卸载前断开连接
onBeforeUnmount(async () => {
  if (connected.value) {
    await handleDisconnect();
  }
});
</script>

<style scoped>
.container {
  padding: 0 20px 20px 20px;
}

.terminal-card {
  margin-top: 16px;
}

.terminal-container {
  width: 100%;
  height: calc(100vh - 200px);
  background-color: #1e1e1e;
  border-radius: 4px;
  padding: 8px;
}
</style>
