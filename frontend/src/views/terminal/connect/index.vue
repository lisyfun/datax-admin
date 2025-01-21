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
          <a-button
            type="primary"
            status="success"
            @click="handleConnect"
            :loading="connecting"
            :disabled="!canConnect || connected"
          >
            <template #icon><icon-play-circle /></template>
            {{ connected ? '重新连接' : '连接' }}
          </a-button>
          <a-button type="primary" @click="testWebSocket" :disabled="connecting">
            <template #icon><icon-bug /></template>
            测试连接
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
        <template v-if="!connected">
          <div class="terminal-placeholder">
            <icon-robot :style="{ fontSize: '48px', marginBottom: '16px' }" />
            <p>点击上方"连接"按钮开始连接终端</p>
          </div>
        </template>
        <template v-else>
          <div id="terminal" class="terminal-content"></div>
        </template>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import {
  IconLink,
  IconPlayCircle,
  IconArrowLeft,
  IconRobot,
  IconBug,
} from '@arco-design/web-vue/es/icon';
import type { TerminalInfo } from '@/types/terminal';
import terminalApi from '@/api/terminal';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import '@xterm/xterm/css/xterm.css';
import { backendConfig } from '@/config';

const route = useRoute();
const router = useRouter();
const terminalId = computed(() => Number(route.params.id));

const terminalInfo = ref<TerminalInfo>();
const connecting = ref(false);
const connected = ref(false);
const canConnect = ref(false);
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

  // 等待DOM元素准备好
  const terminalElement = document.getElementById('terminal');
  if (!terminalElement) {
    console.error('终端DOM元素未找到');
    return;
  }

  // 创建终端实例
  terminal = new Terminal({
    cursorBlink: true,
    theme: {
      background: '#1e1e1e',
      foreground: '#ffffff',
    },
    fontSize: 14,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    scrollback: 1000,
    convertEol: true,
    lineHeight: 1.3,
    letterSpacing: 0.5,
  });

  // 添加插件
  const fitAddon = new FitAddon();
  terminal.loadAddon(fitAddon);
  terminal.loadAddon(new WebLinksAddon());

  // 挂载终端
  terminal.open(terminalElement);
  fitAddon.fit();

  // 监听窗口大小变化
  const resizeObserver = new ResizeObserver(() => {
    fitAddon.fit();
    // 发送新的终端大小到服务器
    if (socket?.readyState === WebSocket.OPEN && terminal) {
      socket.send(JSON.stringify({
        type: 'resize',
        data: JSON.stringify({
          cols: terminal.cols,
          rows: terminal.rows,
        }),
      }));
    }
  });
  resizeObserver.observe(terminalContainer.value);

  // 监听终端输入
  terminal.onData((data) => {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type: 'input', data }));
    }
  });

  // 监听终端大小变化
  terminal.onResize((size) => {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({
        type: 'resize',
        data: JSON.stringify({
          cols: size.cols,
          rows: size.rows,
        }),
      }));
    }
  });
};

// 连接终端
const handleConnect = async () => {
  if (!terminalInfo.value || !canConnect.value) {
    console.error('终端信息未获取到或未通过连接测试');
    return;
  }

  try {
    connecting.value = true;
    console.log('终端信息:', terminalInfo.value);

    // 先设置连接状态，让DOM元素显示出来
    connected.value = true;

    // 等待DOM更新
    await nextTick();

    // 初始化终端
    if (!terminal) {
      console.log('初始化终端实例');
      initTerminal();
      if (!terminal) {
        throw new Error('终端初始化失败');
      }
    }

    // 关闭已存在的连接
    if (socket) {
      console.log('关闭已存在的WebSocket连接');
      socket.close();
    }

    // 创建WebSocket连接
    const wsUrl = `${backendConfig.wsBaseUrl}/datax/ws/terminals/${terminalId.value}`;
    console.log('WebSocket连接URL:', wsUrl);
    console.log('终端ID:', terminalId.value);
    console.log('后端服务器:', backendConfig.wsBaseUrl);

    socket = new WebSocket(wsUrl);
    console.log('WebSocket实例已创建');

    socket.onopen = () => {
      console.log('WebSocket连接已建立');
      connected.value = true;
      Message.success('终端连接成功');
      terminal?.focus();

      // 发送终端大小
      if (terminal && socket) {
        const resizeData = {
          type: 'resize',
          data: JSON.stringify({
            cols: terminal.cols,
            rows: terminal.rows,
          }),
        };
        console.log('发送终端大小数据:', resizeData);
        socket.send(JSON.stringify(resizeData));
      }
    };

    socket.onmessage = (event) => {
      console.log('收到WebSocket消息:', event.data);
      try {
        const data = JSON.parse(event.data);
        switch (data.type) {
          case 'output':
            console.log('收到终端输出:', data.data);
            terminal?.write(data.data);
            break;
          case 'error':
            console.error('服务器返回错误:', data.data);
            Message.error(data.data);
            break;
        }
      } catch (error) {
        console.error('解析WebSocket消息失败:', error, event.data);
      }
    };

    socket.onclose = (event) => {
      console.log('WebSocket连接已关闭:', event.code, event.reason);
      connected.value = false;
      Message.warning('终端连接已断开');
    };

    socket.onerror = (error) => {
      console.error('WebSocket连接错误:', error);
      connected.value = false;
      Message.error('终端连接失败');
    };

  } catch (error) {
    console.error('连接终端失败:', error);
    Message.error('连接终端失败');
    connected.value = false;  // 连接失败时重置状态
  } finally {
    connecting.value = false;
  }
};

// 测试WebSocket连接
const testWebSocket = async () => {
  try {
    if (!terminalInfo.value) {
      Message.error('终端信息未获取到');
      return;
    }

    const wsUrl = `${backendConfig.wsBaseUrl}/datax/ws/terminals/${terminalId.value}`;

    console.log('开始测试WebSocket连接');
    console.log('终端ID:', terminalId.value);
    console.log('终端信息:', {
      name: terminalInfo.value.name,
      host: terminalInfo.value.host,
      port: terminalInfo.value.port,
      username: terminalInfo.value.username,
      status: terminalInfo.value.status
    });
    console.log('WebSocket URL:', wsUrl);
    console.log('后端服务器:', backendConfig.wsBaseUrl);

    // 创建测试连接
    console.log('正在创建WebSocket实例...');
    const testSocket = new WebSocket(wsUrl);

    testSocket.onopen = () => {
      console.log('WebSocket连接已建立');
      Message.success('WebSocket连接测试成功');
      canConnect.value = true;
      // 发送一个测试消息
      testSocket.send(JSON.stringify({ type: 'test', data: 'test connection' }));
      setTimeout(() => {
        console.log('正在关闭测试连接...');
        testSocket.close(1000, '测试完成');
      }, 1000);
    };

    testSocket.onerror = (event: Event) => {
      console.error('WebSocket连接失败:', event);
      canConnect.value = false;
      const wsEvent = event as WebSocketEventMap['error'];
      Message.error('WebSocket连接测试失败');
    };

    testSocket.onclose = (event: CloseEvent) => {
      console.log('WebSocket连接已关闭:', {
        code: event.code,
        reason: event.reason,
        wasClean: event.wasClean
      });
      if (event.code !== 1000) {
        canConnect.value = false;
        Message.warning(`WebSocket连接已关闭: ${event.code}`);
      }
    };

    testSocket.onmessage = (event: MessageEvent) => {
      console.log('收到WebSocket消息:', event.data);
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'error') {
          Message.error(data.data);
        }
      } catch (error) {
        console.error('解析WebSocket消息失败:', error);
      }
    };
  } catch (error) {
    console.error('测试连接出错:', error);
    canConnect.value = false;
    Message.error(`WebSocket连接测试出错: ${error instanceof Error ? error.message : String(error)}`);
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
