<template>
  <div class="terminal-connect">
    <a-card class="terminal-card">
      <template #title>
        <div class="card-title">
          <icon-link class="icon" />
          <span class="title">终端连接 - {{ terminalInfo?.name }}</span>
          <a-tag size="small" :color="terminalInfo?.status === 'online' ? 'green' : 'red'" style="margin-left: 8px">
            {{ terminalInfo?.status === 'online' ? '在线' : '离线' }}
          </a-tag>
        </div>
      </template>
      <template #extra>
        <a-space>
          <a-button type="primary" status="success" @click="handleConnect" :loading="connecting">
            <template #icon><icon-play-circle /></template>
            {{ connected ? '重新连接' : '连接' }}
          </a-button>
          <a-button @click="handleBack">
            <template #icon><icon-arrow-left /></template>
            返回列表
          </a-button>
        </a-space>
      </template>

      <div class="terminal-info" v-if="terminalInfo">
        <a-descriptions :column="2" :data="terminalInfoData" />
      </div>

      <div ref="terminalContainer" class="terminal-container" :class="{ connected }">
        <div v-if="!connected" class="terminal-placeholder">
          <icon-robot :style="{ fontSize: '48px', marginBottom: '16px' }" />
          <p>点击上方"连接"按钮开始连接终端</p>
        </div>
        <div v-else id="terminal" class="terminal-content"></div>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import {
  IconLink,
  IconPlayCircle,
  IconArrowLeft,
  IconRobot,
} from '@arco-design/web-vue/es/icon';
import type { TerminalInfo } from '@/types/terminal';
import terminalApi from '@/api/terminal';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import '@xterm/xterm/css/xterm.css';

const route = useRoute();
const router = useRouter();
const terminalId = computed(() => Number(route.params.id));

const terminalInfo = ref<TerminalInfo>();
const connecting = ref(false);
const connected = ref(false);
const terminalContainer = ref<HTMLElement>();
let terminal: Terminal | null = null;
let socket: WebSocket | null = null;

// 终端信息展示数据
const terminalInfoData = computed(() => {
  if (!terminalInfo.value) return [];
  return [
    {
      label: '终端名称',
      value: terminalInfo.value.name,
    },
    {
      label: '主机地址',
      value: terminalInfo.value.host,
    },
    {
      label: 'SSH端口',
      value: terminalInfo.value.port.toString(),
    },
    {
      label: '用户名',
      value: terminalInfo.value.username,
    },
    {
      label: '最后在线',
      value: formatDate(terminalInfo.value.lastSeen),
    },
    {
      label: '创建时间',
      value: formatDate(terminalInfo.value.createdAt),
    },
  ];
});

// 格式化日期
const formatDate = (date: string) => {
  if (!date) return '-';
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  });
};

// 获取终端信息
const fetchTerminalInfo = async () => {
  try {
    const res = await terminalApi.getTerminalById(terminalId.value);
    terminalInfo.value = res.data;
  } catch (error) {
    Message.error('获取终端信息失败');
    router.push('/terminal/list');
  }
};

// 初始化终端
const initTerminal = () => {
  if (!terminalContainer.value) return;

  // 创建终端实例
  terminal = new Terminal({
    cursorBlink: true,
    theme: {
      background: '#1e1e1e',
      foreground: '#ffffff',
    },
    fontSize: 14,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
  });

  // 添加插件
  const fitAddon = new FitAddon();
  terminal.loadAddon(fitAddon);
  terminal.loadAddon(new WebLinksAddon());

  // 挂载终端
  terminal.open(document.getElementById('terminal')!);
  fitAddon.fit();

  // 监听窗口大小变化
  const resizeObserver = new ResizeObserver(() => {
    fitAddon.fit();
  });
  resizeObserver.observe(terminalContainer.value);

  // 监听终端输入
  terminal.onData((data) => {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type: 'input', data }));
    }
  });
};

// 连接终端
const handleConnect = async () => {
  if (!terminalInfo.value) return;

  try {
    connecting.value = true;

    // 初始化终端
    if (!terminal) {
      initTerminal();
    }

    // 关闭已存在的连接
    if (socket) {
      socket.close();
    }

    // 创建WebSocket连接
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/terminals/connect/${terminalId.value}`;
    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      connected.value = true;
      Message.success('终端连接成功');
      terminal?.focus();
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'output') {
          terminal?.write(data.data);
        }
      } catch (error) {
        console.error('Failed to parse terminal message:', error);
      }
    };

    socket.onclose = () => {
      connected.value = false;
      Message.warning('终端连接已断开');
    };

    socket.onerror = () => {
      connected.value = false;
      Message.error('终端连接失败');
    };

  } catch (error) {
    Message.error('连接终端失败');
  } finally {
    connecting.value = false;
  }
};

// 返回列表
const handleBack = () => {
  router.push('/terminal/list');
};

// 组件挂载
onMounted(() => {
  fetchTerminalInfo();
});

// 组件卸载前清理
onBeforeUnmount(() => {
  if (socket) {
    socket.close();
  }
  if (terminal) {
    terminal.dispose();
  }
});
</script>

<style lang="less" scoped>
.terminal-connect {
  padding: 16px;

  .terminal-card {
    :deep(.arco-card-header) {
      border-bottom: 1px solid var(--color-border);
    }

    .card-title {
      display: flex;
      align-items: center;

      .icon {
        margin-right: 8px;
        font-size: 20px;
        color: rgb(var(--primary-6));
      }

      .title {
        font-size: 16px;
        font-weight: 500;
      }
    }
  }

  .terminal-info {
    margin-bottom: 16px;
    padding: 16px;
    background-color: var(--color-fill-2);
    border-radius: 4px;
  }

  .terminal-container {
    height: 500px;
    background-color: #1e1e1e;
    border-radius: 4px;
    overflow: hidden;

    &.connected {
      padding: 8px;
    }

    .terminal-placeholder {
      height: 100%;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      color: rgba(255, 255, 255, 0.5);
    }

    .terminal-content {
      width: 100%;
      height: 100%;
    }
  }
}
</style>
