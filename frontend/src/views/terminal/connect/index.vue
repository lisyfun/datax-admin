<template>
  <div class="terminal-connect" :class="{ 'fullscreen-mode': isFullscreen }">
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
          <a-button @click="toggleInfoPanel" v-if="connected">
            <template #icon>
              <icon-up v-if="showInfoPanel" />
              <icon-down v-else />
            </template>
            {{ showInfoPanel ? '收起信息' : '展开信息' }}
          </a-button>
          <a-button @click="toggleFullscreen">
            <template #icon><icon-fullscreen /></template>
            {{ isFullscreen ? '退出全屏' : '全屏' }}
          </a-button>
          <a-button @click="handleBack">
            <template #icon><icon-arrow-left /></template>
            返回列表
          </a-button>
        </a-space>
      </template>

      <div class="terminal-info" v-if="terminalInfo && showInfoPanel" :class="{ 'fullscreen-info': isFullscreen }">
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
  IconFullscreen,
  IconUp,
  IconDown,
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
const isFullscreen = ref(false);
const showInfoPanel = ref(true);
const terminalContainer = ref<HTMLElement>();
let terminal: Terminal | null = null;
let socket: WebSocket | null = null;
let resizeTimer: number | null = null;

// 处理F11键全屏切换
const handleF11KeyDown = (e: KeyboardEvent) => {
  if (e.key === 'F11') {
    e.preventDefault();
    toggleFullscreen();
  }
};

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
    fontSize: isFullscreen.value ? 16 : 14,
    fontFamily: 'Consolas, JetBrains Mono, Menlo, Monaco, "Courier New", monospace',
    scrollback: 1000,
    convertEol: true,
    lineHeight: 1.5,
    letterSpacing: 0.8,
  });

  // 添加插件
  const fitAddon = new FitAddon();
  terminal.loadAddon(fitAddon);
  terminal.loadAddon(new WebLinksAddon());

  // 挂载终端
  terminal.open(terminalElement);

  // 调整终端大小
  fitAddon.fit();

  // 全屏模式下调整xterm元素样式
  if (isFullscreen.value) {
    const xtermElement = terminalElement.querySelector('.xterm');
    if (xtermElement) {
      (xtermElement as HTMLElement).style.height = '100%';
    }
  }

  // 监听窗口大小变化
  const resizeObserver = new ResizeObserver(() => {
    // 防抖动处理
    if (resizeTimer) {
      clearTimeout(resizeTimer);
    }

    resizeTimer = window.setTimeout(() => {
      if (fitAddon) {
        fitAddon.fit();
      }

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
    }, 100);
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

  // 聚焦终端
  terminal.focus();

  // 全屏模式下需要额外调整布局
  if (isFullscreen.value && terminal) {
    setTimeout(() => {
      const fitAddon = new FitAddon();
      terminal?.loadAddon(fitAddon);
      fitAddon.fit();
    }, 100);
  }
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

    // 动态构建WebSocket URL
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = window.location.host;
    // 获取基础路径，例如 /datax
    const basePath = window.location.pathname.split('/')[1] || '';
    const wsUrl = `${wsProtocol}//${wsHost}/${basePath ? basePath + '/' : ''}ws/terminals/${terminalId.value}`;

    console.log('当前页面 URL:', window.location.href);
    console.log('基础路径:', basePath);
    console.log('WebSocket连接URL:', wsUrl);
    console.log('终端ID:', terminalId.value);

    socket = new WebSocket(wsUrl);
    console.log('WebSocket实例已创建');

    socket.onopen = () => {
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
      try {
        const data = JSON.parse(event.data);
        switch (data.type) {
          case 'output':
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

    // 动态构建WebSocket URL
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = window.location.host;
    // 获取基础路径，例如 /datax
    const basePath = window.location.pathname.split('/')[1] || '';
    const wsUrl = `${wsProtocol}//${wsHost}/${basePath ? basePath + '/' : ''}ws/terminals/${terminalId.value}`;

    // 创建测试连接
    const testSocket = new WebSocket(wsUrl);

    testSocket.onopen = () => {
      Message.success('WebSocket连接测试成功');
      canConnect.value = true;
      // 发送一个测试消息
      testSocket.send(JSON.stringify({ type: 'test', data: 'test connection' }));
      setTimeout(() => {
        testSocket.close(1000, '测试完成');
      }, 1000);
    };

    testSocket.onerror = (event: Event) => {
      canConnect.value = false;
      const wsEvent = event as WebSocketEventMap['error'];
      Message.error('WebSocket连接测试失败');
    };

    testSocket.onclose = (event: CloseEvent) => {
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
    canConnect.value = false;
    Message.error(`WebSocket连接测试出错: ${error instanceof Error ? error.message : String(error)}`);
  }
};

// 切换信息面板显示/隐藏
const toggleInfoPanel = () => {
  showInfoPanel.value = !showInfoPanel.value;

  nextTick(() => {
    if (terminal) {
      const fitAddon = new FitAddon();
      terminal.loadAddon(fitAddon);
      fitAddon.fit();
    }
  });
};

// 切换全屏
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value;

  // 全屏模式下默认隐藏信息面板
  if (isFullscreen.value) {
    showInfoPanel.value = false;
    // 添加全屏样式到body和html
    document.documentElement.classList.add('terminal-fullscreen');
    document.body.classList.add('terminal-fullscreen-body');
  } else {
    showInfoPanel.value = true;
    // 移除全屏样式
    document.documentElement.classList.remove('terminal-fullscreen');
    document.body.classList.remove('terminal-fullscreen-body');
  }

  nextTick(() => {
    if (terminal) {
      // 调整字体大小
      terminal.options.fontSize = isFullscreen.value ? 16 : 14;

      const fitAddon = new FitAddon();
      terminal.loadAddon(fitAddon);
      fitAddon.fit();

      // 全屏模式下调整xterm元素样式
      const terminalElement = document.getElementById('terminal');
      if (terminalElement) {
        const xtermElement = terminalElement.querySelector('.xterm');
        if (xtermElement) {
          (xtermElement as HTMLElement).style.height = '100%';
        }
      }
    }
  });
};

// 返回列表
const handleBack = () => {
  router.push('/terminal/list');
};

// 计算终端容器高度
const getTerminalHeight = () => {
  if (isFullscreen.value) {
    return showInfoPanel.value ? 'calc(100vh - 180px)' : 'calc(100vh - 75px)';
  } else {
    return showInfoPanel.value ? '550px' : '650px';
  }
};

// 组件挂载
onMounted(() => {
  fetchTerminalInfo();

  // 监听键盘快捷键 F11 切换全屏
  window.addEventListener('keydown', handleF11KeyDown);
});

// 组件卸载前清理
onBeforeUnmount(() => {
  // 退出全屏状态
  if (isFullscreen.value) {
    // 确保移除全屏样式
    document.documentElement.classList.remove('terminal-fullscreen');
    document.body.classList.remove('terminal-fullscreen-body');
  }

  // 移除键盘事件监听
  window.removeEventListener('keydown', handleF11KeyDown);

  if (socket) {
    socket.close();
  }
  if (terminal) {
    terminal.dispose();
  }
});
</script>

<style lang="less">
// 全局样式
:global(html.terminal-fullscreen) {
  overflow: hidden !important;
  margin: 0 !important;
  padding: 0 !important;
}

:global(body.terminal-fullscreen-body) {
  overflow: hidden !important;
  margin: 0 !important;
  padding: 0 !important;
  height: 100vh !important;
  width: 100vw !important;

  .arco-layout {
    height: 100vh !important;
    margin: 0 !important;
    padding: 0 !important;
  }

  .layout-content {
    padding: 0 !important;
    margin: 0 !important;
  }

  .content-wrapper {
    padding: 0 !important;
    margin: 0 !important;
    min-height: 100vh !important;
  }
}
</style>

<style lang="less" scoped>
.terminal-connect {
  padding: 16px;
  transition: padding 0.3s ease;

  :deep(.arco-card) {
    overflow: hidden;
  }

  &.fullscreen-mode {
    padding: 0;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 1000;
    background-color: var(--color-bg-1);
    margin: 0;
    border: none;
    border-radius: 0;
    overflow: hidden;

    .terminal-card {
      height: 100vh;
      border-radius: 0;
      box-shadow: none;
      margin: 0;
      border: none;

      :deep(.arco-card-body) {
        padding: 0;
        height: calc(100vh - 45px); // 减去头部高度
        display: flex;
        flex-direction: column;
        margin: 0;
        border: none;
      }

      :deep(.arco-card-header) {
        border-bottom-color: var(--color-border-2);
        background-color: var(--color-bg-2);
      }
    }

    .terminal-info {
      margin: 4px 8px;
      padding: 8px;
      flex-shrink: 0;
    }

    .terminal-container {
      border-radius: 0;
      flex: 1;
      margin: 0;
      height: auto !important; // 强制使用flex布局的高度

      &.connected {
        padding: 0;
      }

      .terminal-content {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
        overflow: hidden;
        background-color: #1e1e1e;
        padding: 0;
        margin: 0;

        :deep(.xterm) {
          height: 100% !important;
          display: flex;
          flex-direction: column;
          padding: 0;
          margin: 0;

          .xterm-viewport {
            flex: 1;
          }

          canvas {
            padding: 0;
            margin: 0;
          }
        }
      }
    }
  }

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
    transition: all 0.3s ease;

    &.fullscreen-info {
      margin-bottom: 8px;
      padding: 8px 16px;
      border-radius: 0;
    }
  }

  .terminal-container {
    height: v-bind('getTerminalHeight()');
    background-color: #1e1e1e;
    border-radius: 4px;
    overflow: hidden;
    transition: height 0.3s ease;

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
      display: flex;
      flex-direction: column;
    }
  }
}
</style>
