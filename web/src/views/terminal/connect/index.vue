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
            :disabled="connecting || !terminalInfo || connected"
          >
            <template #icon><icon-play-circle /></template>
            连接
          </a-button>
          <a-button type="primary" status="danger" @click="handleDisconnect" v-if="connected">
            <template #icon><icon-close /></template>
            关闭连接
          </a-button>
          <a-divider direction="vertical" />
          <a-tooltip content="减小字体">
            <a-button @click="decreaseFontSize" :disabled="!connected">
              <template #icon><icon-minus /></template>
            </a-button>
          </a-tooltip>
          <a-tooltip content="增大字体">
            <a-button @click="increaseFontSize" :disabled="!connected">
              <template #icon><icon-plus /></template>
            </a-button>
          </a-tooltip>
          <a-tooltip content="减小行高">
            <a-button @click="decreaseLineHeight" :disabled="!connected">
              <template #icon><icon-line-height /></template>
              <template #default>-</template>
            </a-button>
          </a-tooltip>
          <a-tooltip content="增大行高">
            <a-button @click="increaseLineHeight" :disabled="!connected">
              <template #icon><icon-line-height /></template>
              <template #default>+</template>
            </a-button>
          </a-tooltip>
          <a-divider direction="vertical" />
          <a-button @click="toggleInfoPanel" v-if="connected">
            <template #icon>
              <icon-up v-if="showInfoPanel" />
              <icon-down v-else />
            </template>
            {{ showInfoPanel ? '收起信息' : '展开信息' }}
          </a-button>
        </a-space>
      </template>

      <div class="terminal-info" v-if="terminalInfo && showInfoPanel" :class="{ 'fullscreen-info': isFullscreen }">
        <a-descriptions :column="2" :data="terminalInfoData" />
      </div>

      <div ref="terminalContainer" class="terminal-container" :class="{ connected }" @click="handleTerminalClick" @contextmenu.prevent="handleRightClick" tabindex="0" @focus="handleTerminalFocus">
        <template v-if="!connected">
          <div class="terminal-placeholder">
            <icon-robot :style="{ fontSize: '48px', marginBottom: '16px' }" />
            <p>点击上方"连接"按钮开始连接终端</p>
          </div>
        </template>
        <template v-else>
          <div id="terminal" class="terminal-content"></div>
        </template>

        <!-- 操作提示 -->
        <div class="action-feedback copy-feedback" v-if="showCopyTip">
          <icon-check-circle class="icon" /> 已复制到剪贴板
        </div>
        <div class="action-feedback paste-feedback" v-if="showPasteTip">
          <icon-check-circle class="icon" /> 已粘贴
        </div>
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
  IconRobot,
  IconUp,
  IconDown,
  IconClose,
  IconMinus,
  IconPlus,
  IconLineHeight,
  IconCheckCircle,
} from '@arco-design/web-vue/es/icon';
import type { TerminalInfo } from '@/types/terminal';
import terminalApi from '@/api/terminal';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { WebglAddon } from '@xterm/addon-webgl';
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
let reconnectTimer: number | null = null;

const RECONNECT_INTERVAL = 3000; // 3秒后尝试重连
const MAX_RECONNECT_ATTEMPTS = 3; // 最大重连次数

// 字体设置相关
const fontSize = ref(14);
const isBold = ref(false);
const lineHeight = ref(1.3);
const hasSelection = ref(false);

// 显示气泡提示的配置
const showCopyTip = ref(false);
const showPasteTip = ref(false);

// 从本地存储加载字体设置
const loadFontSettings = () => {
  if (terminalId.value) {
    const savedSettings = localStorage.getItem(`terminal_settings_${terminalId.value}`);
    if (savedSettings) {
      const settings = JSON.parse(savedSettings);
      fontSize.value = settings.fontSize || 14;
      isBold.value = settings.fontWeight === 'bold';
      lineHeight.value = settings.lineHeight || 1.3;
    }
  }
};

// 保存字体设置到本地存储
const saveFontSettings = () => {
  if (terminalId.value) {
    localStorage.setItem(
      `terminal_settings_${terminalId.value}`,
      JSON.stringify({
        fontSize: fontSize.value,
        fontWeight: isBold.value ? 'bold' : 'normal',
        lineHeight: lineHeight.value
      })
    );
  }
};

// 应用字体设置到终端
const applyFontSettings = () => {
  if (terminal) {
    terminal.options.fontSize = fontSize.value;
    terminal.options.fontWeight = isBold.value ? 'bold' : 'normal';
    terminal.options.lineHeight = lineHeight.value;

    // 重新适配终端大小
    const fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);
    fitAddon.fit();
  }
};

// 增大字体
const increaseFontSize = () => {
  if (fontSize.value < 32) {
    fontSize.value += 1;
    applyFontSettings();
    saveFontSettings();
  }
};

// 减小字体
const decreaseFontSize = () => {
  if (fontSize.value > 8) {
    fontSize.value -= 1;
    applyFontSettings();
    saveFontSettings();
  }
};

// 切换字体粗细
const toggleFontWeight = () => {
  isBold.value = !isBold.value;
  applyFontSettings();
  saveFontSettings();
};

// 增大行高
const increaseLineHeight = () => {
  if (lineHeight.value < 2.0) {
    lineHeight.value = Math.round((lineHeight.value + 0.1) * 10) / 10;
    applyFontSettings();
    saveFontSettings();
  }
};

// 减小行高
const decreaseLineHeight = () => {
  if (lineHeight.value > 1.0) {
    lineHeight.value = Math.round((lineHeight.value - 0.1) * 10) / 10;
    applyFontSettings();
    saveFontSettings();
  }
};

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
    return;
  }

  // 创建终端实例
  terminal = new Terminal({
    cursorBlink: true,
    theme: {
      background: '#1e1e1e',
      foreground: '#ffffff',
    },
    fontSize: fontSize.value,
    fontFamily: 'Consolas, Menlo, Monaco, "Courier New", monospace',
    fontWeight: isBold.value ? 'bold' : 'normal',
    scrollback: 1000,
    convertEol: true,
    lineHeight: lineHeight.value,
    letterSpacing: 0.5,
    allowTransparency: true,
    rightClickSelectsWord: true, // 右键选中单词
  });

  // 添加插件
  const fitAddon = new FitAddon();
  terminal.loadAddon(fitAddon);
  terminal.loadAddon(new WebLinksAddon());

  // 添加WebGL渲染插件提高性能
  try {
    const webglAddon = new WebglAddon();
    terminal.loadAddon(webglAddon);
  } catch (e) {
    // WebGL渲染初始化失败，将使用Canvas渲染
  }

  // 监听选择事件，自动复制选中的内容
  terminal.onSelectionChange(() => {
    const hasTextSelected = terminal?.hasSelection() || false;
    hasSelection.value = hasTextSelected;

    if (hasTextSelected) {
      copySelectionToClipboard();
    }
  });

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
  if (!terminalInfo.value) {
    return;
  }

  try {
    connecting.value = true;

    // 如果已经连接，先断开现有连接
    if (connected.value) {
      if (socket) {
        socket.close(1000, '重新连接');
      }
      if (terminal) {
        terminal.clear();
        terminal.dispose();
        terminal = null;
      }
      connected.value = false;
      await nextTick();
    }

    // 先设置连接状态，让DOM元素显示出来
    connected.value = true;

    // 等待DOM更新
    await nextTick();

    // 重新初始化终端
    initTerminal();
    if (!terminal) {
      throw new Error('终端初始化失败');
    }

    // 关闭已存在的连接
    if (socket) {
      socket.close();
    }

    // 动态构建WebSocket URL
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = window.location.host;
    const basePath = window.location.pathname.split('/')[1] || '';
    const wsUrl = `${wsProtocol}//${wsHost}/${basePath ? basePath + '/' : ''}ws/terminals/${terminalId.value}`;

    socket = new WebSocket(wsUrl);

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
            Message.error(data.data);
            break;
          default:
            // 收到未知类型的消息，静默处理
            break;
        }
      } catch (error) {
        // 解析WebSocket消息失败，静默处理
      }
    };

    socket.onclose = (event) => {
      // 只有在非用户主动关闭的情况下才处理重连
      if (event.code !== 1000) {
        connected.value = false;
        Message.warning('终端连接已断开，正在尝试重新连接...');
        handleReconnect();
      } else {
        connected.value = false;
        if (event.reason !== '用户主动关闭连接') {
          Message.warning('终端连接已断开');
        }
      }
    };

    socket.onerror = (error) => {
      connected.value = false;
      Message.error('终端连接失败');
    };

  } catch (error) {
    Message.error('连接终端失败');
    connected.value = false;  // 连接失败时重置状态
  } finally {
    connecting.value = false;
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
      terminal.options.fontSize = isFullscreen.value ? 14 : 12;

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

// 关闭终端连接
const handleDisconnect = () => {
  // 清除重连定时器
  if (reconnectTimer) {
    clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
  if (socket) {
    socket.close(1000, '用户主动关闭连接');
    socket = null;
  }
  if (terminal) {
    terminal.clear();
    terminal = null;
  }
  connected.value = false;
  Message.success('终端连接已关闭');
};

// 添加重连机制
const handleReconnect = () => {
  let reconnectAttempts = 0;

  const attemptReconnect = () => {
    if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
      Message.error('重连失败，请手动重新连接');
      return;
    }

    reconnectAttempts++;

    // 重新创建WebSocket连接
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = window.location.host;
    const basePath = window.location.pathname.split('/')[1] || '';
    const wsUrl = `${wsProtocol}//${wsHost}/${basePath ? basePath + '/' : ''}ws/terminals/${terminalId.value}`;

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      connected.value = true;
      reconnectAttempts = 0;
      Message.success('终端重新连接成功');
      terminal?.focus();
    };

    socket.onclose = (event) => {
      if (event.code !== 1000) {
        // 如果不是主动关闭，继续尝试重连
        setTimeout(attemptReconnect, RECONNECT_INTERVAL);
      }
    };

    socket.onerror = () => {
      // 重连过程中发生错误，静默处理
    };
  };

  attemptReconnect();
};

// 自动复制选中内容到剪贴板
const copySelectionToClipboard = () => {
  if (!terminal || !terminal.hasSelection()) return;

  const selection = terminal.getSelection();
  if (!selection) return;

  // 优先使用 Clipboard API
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(selection)
      .then(() => {
        showCopiedFeedback();
      })
      .catch(() => {
        // 降级处理
        fallbackCopyTextToClipboard(selection);
      });
  } else {
    // 不支持 Clipboard API，直接降级
    fallbackCopyTextToClipboard(selection);
  }
};

// 降级方案
function fallbackCopyTextToClipboard(text: string) {
  const textArea = document.createElement("textarea");
  textArea.value = text;
  // 避免页面滚动
  textArea.style.position = "fixed";
  textArea.style.top = "0";
  textArea.style.left = "0";
  textArea.style.width = "2em";
  textArea.style.height = "2em";
  textArea.style.padding = "0";
  textArea.style.border = "none";
  textArea.style.outline = "none";
  textArea.style.boxShadow = "none";
  textArea.style.background = "transparent";
  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    const successful = document.execCommand('copy');
    if (successful) {
      showCopiedFeedback();
    } else {
      Message.error('复制失败，请手动复制');
    }
  } catch (err) {
    Message.error('复制失败，请手动复制');
  }
  document.body.removeChild(textArea);
}

// 显示复制成功的反馈
const showCopiedFeedback = () => {
  // 显示一个小提示，表明复制成功
  showCopyTip.value = true;
  setTimeout(() => {
    showCopyTip.value = false;
  }, 1000);
};

// 处理右键点击，直接粘贴
const handleRightClick = async (e: MouseEvent) => {
  if (!connected.value || !terminal) return;

  // 1. 优先 Clipboard API
  if (navigator.clipboard && navigator.clipboard.readText) {
    try {
      const text = await navigator.clipboard.readText();
      if (text && socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'input', data: text }));
        showPasteTip.value = true;
        setTimeout(() => { showPasteTip.value = false; }, 1000);
      } else {
        Message.warning('剪贴板为空，请先复制内容');
      }
      return;
    } catch (err) {
      // 失败降级
    }
  }

  // 2. 尝试 execCommand（大概率无效，但可尝试）
  let pasted = false;
  try {
    const textArea = document.createElement('textarea');
    textArea.style.position = 'fixed';
    textArea.style.opacity = '0';
    document.body.appendChild(textArea);
    textArea.focus();
    if (document.execCommand('paste')) {
      const text = textArea.value;
      if (text && socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'input', data: text }));
        showPasteTip.value = true;
        setTimeout(() => { showPasteTip.value = false; }, 1000);
        pasted = true;
      }
    }
    document.body.removeChild(textArea);
  } catch (e) {
    // 失败
  }

  if (!pasted) {
    // 3. 最终提示用户手动粘贴
    Message.warning('浏览器限制，无法自动粘贴，请使用快捷键 shift+insert 粘贴');
    terminal.focus();
  }
};

// 粘贴事件处理
const handlePasteEvent = async (event: ClipboardEvent) => {
  if (!connected.value || !terminal) return;
  
  event.preventDefault();
  event.stopPropagation();
  
  try {
    // 优先从事件中获取数据
    let text = event.clipboardData?.getData('text') || '';
    
    // 如果事件中没有数据，尝试其他方法
    if (!text || !text.trim()) {
      text = await readClipboard();
    }
    
    if (text && text.trim()) {
      if (socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'input', data: text }));
        showPasteTip.value = true;
        setTimeout(() => { showPasteTip.value = false; }, 1000);
        console.log('右键粘贴成功:', text.length, '个字符');
      } else {
        Message.warning('终端连接已断开');
      }
    } else {
      Message.warning('剪贴板内容为空');
    }
  } catch (e) {
    console.error('粘贴事件处理失败:', e);
    Message.warning('粘贴失败');
  }
};

// 请求剪贴板权限
const requestClipboardPermission = async () => {
  try {
    if (navigator.permissions) {
      const permission = await navigator.permissions.query({ name: 'clipboard-read' as PermissionName });
      return permission.state === 'granted';
    }
    return false;
  } catch (e) {
    console.warn('无法检查剪贴板权限:', e);
    return false;
  }
};

// 更可靠的剪贴板读取函数
const readClipboard = async (): Promise<string> => {
  let text = '';
  
  try {
    // 方法1: 现代剪贴板API（需要HTTPS或localhost）
    if (navigator.clipboard && navigator.clipboard.readText) {
      // 检查权限
      const hasPermission = await requestClipboardPermission();
      if (hasPermission || location.protocol === 'https:' || location.hostname === 'localhost') {
        text = await navigator.clipboard.readText();
        if (text) return text;
      }
    }
  } catch (e) {
    console.warn('现代剪贴板API失败:', e);
  }
  
  try {
    // 方法2: 使用execCommand降级方案
    const textArea = document.createElement('textarea');
    textArea.style.position = 'fixed';
    textArea.style.left = '-999999px';
    textArea.style.top = '-999999px';
    textArea.style.opacity = '0';
    textArea.style.pointerEvents = 'none';
    document.body.appendChild(textArea);
    
    textArea.focus();
    textArea.select();
    
    // 尝试执行粘贴命令
    const success = document.execCommand('paste');
    if (success && textArea.value) {
      text = textArea.value;
    }
    
    document.body.removeChild(textArea);
    
    if (text) return text;
  } catch (e) {
    console.warn('execCommand粘贴失败:', e);
  }
  
  // 方法3: 尝试从事件中获取（如果有的话）
  try {
    // 创建一个隐藏的可编辑元素来捕获粘贴事件
    const hiddenInput = document.createElement('div');
    hiddenInput.contentEditable = 'true';
    hiddenInput.style.position = 'fixed';
    hiddenInput.style.left = '-999999px';
    hiddenInput.style.top = '-999999px';
    hiddenInput.style.opacity = '0';
    hiddenInput.style.pointerEvents = 'none';
    document.body.appendChild(hiddenInput);
    
    return new Promise((resolve) => {
      const pasteHandler = (e: ClipboardEvent) => {
        e.preventDefault();
        const clipboardText = e.clipboardData?.getData('text') || '';
        hiddenInput.removeEventListener('paste', pasteHandler);
        document.body.removeChild(hiddenInput);
        resolve(clipboardText);
      };
      
      hiddenInput.addEventListener('paste', pasteHandler);
      hiddenInput.focus();
      
      // 模拟Ctrl+V
      const pasteEvent = new KeyboardEvent('keydown', {
        key: 'v',
        ctrlKey: true,
        bubbles: true
      });
      hiddenInput.dispatchEvent(pasteEvent);
      
      // 超时处理
      setTimeout(() => {
        hiddenInput.removeEventListener('paste', pasteHandler);
        if (document.body.contains(hiddenInput)) {
          document.body.removeChild(hiddenInput);
        }
        resolve('');
      }, 1000);
    });
  } catch (e) {
    console.warn('事件捕获粘贴失败:', e);
  }
  
  return '';
};

// 监听 Ctrl+V/⌘+V 粘贴
const handleKeyDown = async (event: KeyboardEvent) => {
  if (!connected.value || !terminal) return;

  // Windows/Linux: Ctrl+V，Mac: Meta+V (Cmd+V)
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'v') {
    event.preventDefault();
    event.stopPropagation();
    
    try {
      const text = await readClipboard();
      
      if (text && text.trim()) {
        if (socket?.readyState === WebSocket.OPEN) {
          socket.send(JSON.stringify({ type: 'input', data: text }));
          showPasteTip.value = true;
          setTimeout(() => { showPasteTip.value = false; }, 1000);
          console.log('粘贴成功:', text.length, '个字符');
        } else {
          Message.warning('终端连接已断开');
        }
      } else {
        Message.warning('剪贴板内容为空或无法访问，请尝试右键粘贴');
      }
    } catch (e) {
      console.error('粘贴失败:', e);
      Message.warning('粘贴失败，请尝试右键粘贴或检查浏览器权限');
    }
  }
  
  // 支持 Shift+Insert 粘贴（Linux/Windows 传统快捷键）
  if (event.shiftKey && event.key === 'Insert') {
    event.preventDefault();
    event.stopPropagation();
    
    try {
      const text = await readClipboard();
      if (text && text.trim() && socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'input', data: text }));
        showPasteTip.value = true;
        setTimeout(() => { showPasteTip.value = false; }, 1000);
      } else if (!text || !text.trim()) {
        Message.warning('剪贴板内容为空');
      }
    } catch (e) {
      console.error('Shift+Insert 粘贴失败:', e);
      Message.warning('无法读取剪贴板内容，请检查浏览器权限');
    }
  }
};

// 组件挂载
onMounted(async () => {
  await fetchTerminalInfo();
  loadFontSettings(); // 加载字体设置

  // 自动触发连接
  if (terminalInfo.value && terminalInfo.value.status === 'online') {
    handleConnect();
  }

  // 监听粘贴事件
  if (terminalContainer.value) {
    terminalContainer.value.addEventListener('paste', handlePasteEvent);
    terminalContainer.value.addEventListener('keydown', handleKeyDown);
  }
});

// 处理终端容器点击事件
const handleTerminalClick = () => {
  if (terminal && connected.value) {
    terminal.focus();
  }
  // 确保容器获得焦点以接收键盘事件
  if (terminalContainer.value) {
    terminalContainer.value.focus();
  }
};

// 处理终端容器焦点事件
const handleTerminalFocus = () => {
  if (terminal && connected.value) {
    terminal.focus();
  }
};

// 组件卸载前清理
onBeforeUnmount(() => {
  // 清除重连定时器
  if (reconnectTimer) {
    clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }

  if (socket) {
    socket.close(1000, '组件卸载');
  }
  if (terminal) {
    terminal.dispose();
  }
  // 移除粘贴事件监听
  if (terminalContainer.value) {
    terminalContainer.value.removeEventListener('paste', handlePasteEvent);
    terminalContainer.value.removeEventListener('keydown', handleKeyDown);
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

  // 隐藏顶部栏和左侧导航栏
  .arco-layout-header,
  .arco-layout-sider {
    display: none !important;
  }

  .arco-layout-content {
    margin-left: 0 !important;
    margin-top: 0 !important;
  }
}
</style>

<style lang="less" scoped>
.terminal-connect {
  padding: 0;
  height: 100vh;
  width: 100vw;
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
  display: flex;
  flex-direction: column;

  :deep(.arco-card) {
    height: 100vh;
    border-radius: 0;
    box-shadow: none;
    margin: 0;
    border: none;
    overflow: hidden;
    display: flex;
    flex-direction: column;

    .arco-card-header {
      border-bottom-color: var(--color-border-2);
      background-color: var(--color-bg-2);
      padding: 12px 20px;
      flex-shrink: 0;
    }

    .arco-card-body {
      padding: 0;
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }
  }

  .terminal-info {
    padding: 12px 20px;
    background-color: var(--color-fill-2);
    border-bottom: 1px solid var(--color-border);
    flex-shrink: 0;

    :deep(.arco-descriptions) {
      .arco-descriptions-item {
        padding: 8px 0;

        .arco-descriptions-item-label {
          color: var(--color-text-2);
          font-weight: 500;
        }

        .arco-descriptions-item-value {
          color: var(--color-text-1);
        }
      }
    }
  }

  .terminal-container {
    flex: 1;
    position: relative;
    background-color: #1e1e1e;
    overflow: hidden;
    display: flex;
    flex-direction: column;

    &.connected {
      padding: 8px;
    }

    .terminal-placeholder {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      text-align: center;
      color: rgba(255, 255, 255, 0.5);

      .icon {
        font-size: 48px;
        margin-bottom: 16px;
      }

      p {
        font-size: 16px;
        margin: 0;
      }
    }

    .terminal-content {
      width: 100%;
      height: 100%;
      display: flex;
      flex-direction: column;
      overflow: hidden;
      background-color: #1e1e1e;

      :deep(.xterm) {
        height: 100% !important;
        display: flex;
        flex-direction: column;
        padding-left: 8px;

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

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;

  .icon {
    font-size: 20px;
    color: rgb(var(--primary-6));
  }

  .title {
    font-size: 16px;
    font-weight: 500;
  }
}

// 添加字体控制按钮样式
:deep(.arco-btn) {
  &[disabled] {
    opacity: 0.5;
  }
}

// 添加复制粘贴操作反馈的样式
.action-feedback {
  position: absolute;
  background-color: rgba(0, 0, 0, 0.75);
  color: #fff;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 6px;
  z-index: 100;

  &.copy-feedback {
    top: 20px;
    right: 20px;
  }

  &.paste-feedback {
    bottom: 20px;
    right: 20px;
  }

  .icon {
    color: #52c41a;
  }
}
</style>
