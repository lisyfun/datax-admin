<template>
  <div class="container">
    <a-card>
      <template #title>
        <div class="card-title">
          <div class="title-left">
            <div class="topic-name">
              <icon-message class="topic-icon" />
              <span>{{ topicName }}</span>
            </div>
            <a-space class="topic-stats">
              <span class="stat-item">
                <span class="stat-label">起始偏移量:</span>
                <span class="stat-value">{{ topicInfo.beginningOffset || '-' }}</span>
              </span>
              <span class="stat-item">
                <span class="stat-label">结束偏移量:</span>
                <span class="stat-value">{{ topicInfo.endOffset || '-' }}</span>
              </span>
              <span class="stat-item">
                <span class="stat-label">消息数量:</span>
                <span class="stat-value">{{ topicInfo.size || '-' }}</span>
              </span>
            </a-space>
          </div>
          <a-button @click="handleBack">
            <template #icon>
              <icon-left />
            </template>
            返回
          </a-button>
        </div>
      </template>

      <!-- 筛选条件 -->
      <div class="search-container">
        <div class="search-header">
          <div class="search-title">
            <icon-filter />
            <span>消息筛选</span>
          </div>
          <a-tooltip position="left">
            <template #content>
              <div class="tooltip-content">
                <p>- 分区：选择要查询的Kafka分区</p>
                <p>- 偏移量模式：选择从最新或最早的消息开始</p>
                <p>- 偏移量：指定具体的偏移量位置</p>
                <p>- 消息数量：限制返回的消息数量（最大100条）</p>
                <p>- 消息键/值：根据关键字过滤消息</p>
              </div>
            </template>
            <icon-question-circle style="cursor: pointer; color: var(--color-text-3);" />
          </a-tooltip>
        </div>

        <div class="search-content">
          <div class="search-row">
            <div class="filter-group">
              <div class="filter-label">分区:</div>
              <a-select
                v-model="searchForm.partition"
                placeholder="请选择分区"
                allow-clear
                @change="handlePartitionChange"
                style="width: 120px"
              >
                <a-option v-for="p in partitions" :key="p" :value="p">{{ p }}</a-option>
              </a-select>
            </div>

            <div class="filter-group">
              <div class="filter-label">偏移量模式:</div>
              <div class="offset-mode-switch">
                <span
                  class="mode-label"
                  :class="{ active: offsetReset === 'earliest' }"
                  @click="handleOffsetSwitchChange(false)"
                >最早</span>
                <a-switch
                  :model-value="offsetReset === 'latest'"
                  @change="handleOffsetSwitchChange"
                  type="round"
                />
                <span
                  class="mode-label"
                  :class="{ active: offsetReset === 'latest' }"
                  @click="handleOffsetSwitchChange(true)"
                >最新</span>
              </div>
            </div>

            <div class="filter-group">
              <div class="filter-label">偏移量:</div>
              <a-input-number
                v-model="searchForm.offset"
                placeholder="偏移量"
                style="width: 150px"
              />
            </div>

            <div class="filter-group">
              <div class="filter-label">消息数量:</div>
              <a-input-number
                v-model="searchForm.count"
                :min="1"
                :max="100"
                placeholder="消息数量"
                style="width: 120px"
              />
            </div>
          </div>

          <div class="search-row">
            <div class="filter-group">
              <div class="filter-label">消息键:</div>
              <a-input
                v-model="searchForm.keyFilter"
                placeholder="输入关键字过滤消息键"
                style="width: 220px"
                allow-clear
              >
                <template #prefix>
                  <icon-search />
                </template>
              </a-input>
            </div>

            <div class="filter-group">
              <div class="filter-label">消息值:</div>
              <a-input
                v-model="searchForm.valueFilter"
                placeholder="输入关键字过滤消息内容"
                style="width: 220px"
                allow-clear
              >
                <template #prefix>
                  <icon-search />
                </template>
              </a-input>
            </div>

            <div class="filter-actions">
              <a-tooltip content="拉取符合条件的消息">
                <a-button type="primary" @click="handleSearch" status="normal">
                  <template #icon>
                    <icon-search />
                  </template>
                  拉取消息
                </a-button>
              </a-tooltip>
              <a-tooltip content="重置筛选条件">
                <a-button @click="handleReset">
                  <template #icon>
                    <icon-refresh />
                  </template>
                  重置
                </a-button>
              </a-tooltip>
            </div>
          </div>
        </div>
      </div>

      <!-- 消息列表 -->
      <div class="message-container">
        <div v-for="(message, index) in messages" :key="index" class="message-item">
          <div class="message-meta">
            <span class="meta-item">
              <icon-apps class="meta-icon" />
              partition: <b>{{ message.partition }}</b>
            </span>
            <span class="meta-item">
              <icon-file class="meta-icon" />
              key: <b>{{ message.key || '-' }}</b>
            </span>
            <span class="meta-item">
              <icon-code class="meta-icon" />
              offset: <b>{{ message.offset }}</b>
            </span>
            <span class="meta-item">
              <icon-calendar class="meta-icon" />
              timestamp: <b>{{ formatDateTime(String(message.timestamp)) }}</b>
            </span>
          </div>
          <div class="message-content">
            <div v-if="!expandedMessages[index]" class="message-preview">
              <div class="toggle-button" @click="toggleMessageExpand(index)">
                <icon-expand />
              </div>
              <div class="content-container">
                <div class="content-actions">
                  <div class="copy-button" @click="copyMessageContent(message.value)" title="复制内容">
                    <icon-copy />
                  </div>
                </div>
                <pre class="preview-content">{{ getPreviewContent(message.value) }}</pre>
              </div>
            </div>
            <div v-else class="message-full">
              <div class="toggle-button" @click="toggleMessageExpand(index)">
                <icon-shrink />
              </div>
              <div class="content-container">
                <div class="content-actions">
                  <div class="copy-button" @click="copyMessageContent(message.value)" title="复制内容">
                    <icon-copy />
                  </div>
                </div>
                <pre>{{ formatValue(message.value) }}</pre>
              </div>
            </div>
          </div>
        </div>

        <div v-if="messages.length === 0 && !loading" class="empty-message">
          没有找到符合条件的消息
        </div>

        <div v-if="loading" class="loading-container">
          <a-spin />
        </div>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { consumeMessages, getTopicPartitions, getTopicInfo } from '@/api/kafka';
import { formatDateTime } from '@/utils/date';
import {
  IconLeft,
  IconSearch,
  IconRefresh,
  IconFilter,
  IconQuestionCircle,
  IconApps,
  IconExpand,
  IconShrink,
  IconFile,
  IconCode,
  IconCalendar,
  IconCopy,
} from '@arco-design/web-vue/es/icon';

const route = useRoute();
const router = useRouter();
const loading = ref(false);
const messages = ref<KafkaMessage[]>([]);
const partitions = ref<number[]>([]);
const offsetReset = ref('latest');
const expandedMessages = ref<boolean[]>([]);

// 从路由参数中获取信息
const clusterId = parseInt(route.params.clusterId as string, 10);
const topicName = route.params.topicName as string;

// 主题信息
const topicInfo = reactive({
  beginningOffset: 0,
  endOffset: 0,
  size: 0,
});

// 搜索表单
const searchForm = reactive({
  partition: parseInt(route.query.partition as string, 10) || 0,
  offset: parseInt(route.query.offset as string, 10) || -1,
  count: 10,
  keyFilter: '',
  valueFilter: '',
  groupId: 'datax-admin', // 添加默认消费者组
});

// 格式化消息内容
const formatValue = (value: string) => {
  try {
    const parsed = JSON.parse(value);
    return JSON.stringify(parsed, null, 2);
  } catch {
    return value;
  }
};

// 定义消息类型接口
interface KafkaMessage {
  partition: number;
  offset: number;
  key: string;
  value: string;
  timestamp: number;
}

// 获取主题分区信息
const fetchPartitions = async () => {
  try {
    const res = await getTopicPartitions(clusterId, topicName);
    if (res.data.code === 0) {
      // 直接使用后端返回的分区数组
      partitions.value = res.data.data;
      console.log('获取到的分区列表:', partitions.value);
      if (!searchForm.partition && partitions.value.length > 0) {
        searchForm.partition = partitions.value[0];
      }

      // 获取主题信息
      fetchTopicInfo();
    } else {
      Message.error(res.data.message || '获取分区信息失败');
    }
  } catch (err: any) {
    Message.error(err.response?.data?.message || '获取分区信息失败');
  }
};

// 获取主题信息（起始偏移量、结束偏移量、消息数量）
const fetchTopicInfo = async () => {
  try {
    // 假设后端提供了获取主题信息的API
    // 这里需要根据实际API调整
    const res = await getTopicInfo(clusterId, topicName);
    if (res.data.code === 0) {
      const data = res.data.data;
      topicInfo.beginningOffset = data.beginningOffset;
      topicInfo.endOffset = data.endOffset;
      topicInfo.size = data.size;
    } else {
      Message.error(res.data.message || '获取主题信息失败');
    }
  } catch (err: any) {
    console.error('获取主题信息失败:', err);
    Message.error(err.response?.data?.message || '获取主题信息失败');
  }
};

// 获取消息列表
const fetchMessages = async () => {
  loading.value = true;
  messages.value = []; // 清空之前的消息
  expandedMessages.value = []; // 清空展开状态

  try {
    console.log('开始获取消息，参数:', {
      clusterId,
      topicName,
      partition: searchForm.partition,
      offset: searchForm.offset,
      count: searchForm.count,
      keyFilter: searchForm.keyFilter,
      valueFilter: searchForm.valueFilter,
      groupId: searchForm.groupId,
    });

    // 设置请求超时
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 15000); // 15秒超时

    const res = await consumeMessages(clusterId, topicName, {
      partition: searchForm.partition,
      offset: searchForm.offset,
      count: searchForm.count,
      keyFilter: searchForm.keyFilter,
      valueFilter: searchForm.valueFilter,
      groupId: searchForm.groupId,
    }, { signal: controller.signal });

    clearTimeout(timeoutId);

    if (res.data.code === 0) {
      messages.value = res.data.data || [] as KafkaMessage[];
      // 初始化所有消息为未展开状态
      expandedMessages.value = new Array(messages.value.length).fill(false);

      if (messages.value.length === 0) {
        Message.info('没有找到符合条件的消息');
      } else {
        Message.success(`成功获取 ${messages.value.length} 条消息`);
      }
    } else {
      Message.error(res.data.message || '获取消息失败');
    }
  } catch (err: any) {
    console.error('获取消息失败:', err);
    if (err.name === 'AbortError') {
      Message.error('请求超时，请尝试减少消息数量或使用更精确的过滤条件');
    } else {
      Message.error(err.response?.data?.message || '获取消息失败，请检查网络连接或服务器状态');
    }
  } finally {
    loading.value = false;
  }
};

// 分区变更
const handlePartitionChange = () => {
  searchForm.offset = -1; // 切换分区时重置偏移量为最新位置
};

// 偏移量重置选项变更
const handleOffsetSwitchChange = (checked: boolean) => {
  offsetReset.value = checked ? 'latest' : 'earliest';
  searchForm.offset = checked ? -1 : 0;
};

// 搜索
const handleSearch = () => {
  fetchMessages();
};

// 重置
const handleReset = () => {
  offsetReset.value = 'latest';
  searchForm.offset = -1;
  searchForm.count = 10;
  searchForm.keyFilter = '';
  searchForm.valueFilter = '';
  fetchMessages();
};

// 返回主题列表
const handleBack = () => {
  router.push({
    name: 'KafkaTopic',
    params: { clusterId },
    query: { clusterName: route.query.clusterName as string }
  });
};

// 获取消息预览内容
const getPreviewContent = (value: string) => {
  try {
    // 尝试解析 JSON
    const parsed = JSON.parse(value);
    const formatted = JSON.stringify(parsed, null, 2);

    // 只返回第一行，不添加省略号
    const lines = formatted.split('\n');
    if (lines.length > 1) {
      return lines[0];
    }
    return formatted;
  } catch {
    // 非 JSON 内容，返回前 50 个字符，不添加省略号
    if (value.length > 50) {
      return value.slice(0, 50);
    }
    return value;
  }
};

// 展开或收起消息内容
const toggleMessageExpand = (index: number) => {
  // 确保expandedMessages数组已初始化
  if (!expandedMessages.value) {
    expandedMessages.value = new Array(messages.value.length).fill(false);
  }
  // 切换指定索引的消息展开状态
  expandedMessages.value[index] = !expandedMessages.value[index];
};

// 复制消息内容（支持不安全环境）
const copyMessageContent = (content: string) => {
  // 尝试使用现代Clipboard API
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(content)
      .then(() => {
        Message.success('消息内容已复制到剪贴板');
      })
      .catch(() => {
        // 如果Clipboard API失败，回退到传统方法
        fallbackCopyTextToClipboard(content);
      });
  } else {
    // 在不安全上下文中使用传统方法
    fallbackCopyTextToClipboard(content);
  }
};

// 传统复制方法（兼容性支持）
const fallbackCopyTextToClipboard = (text: string) => {
  try {
    // 创建临时文本区域
    const textArea = document.createElement('textarea');
    textArea.value = text;

    // 设置样式使其不可见
    textArea.style.position = 'fixed';
    textArea.style.top = '0';
    textArea.style.left = '0';
    textArea.style.width = '2em';
    textArea.style.height = '2em';
    textArea.style.padding = '0';
    textArea.style.border = 'none';
    textArea.style.outline = 'none';
    textArea.style.boxShadow = 'none';
    textArea.style.background = 'transparent';

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    // 尝试执行复制命令
    const successful = document.execCommand('copy');
    if (successful) {
      Message.success('消息内容已复制到剪贴板');
    } else {
      Message.error('复制失败，请手动复制');
    }

    // 清理
    document.body.removeChild(textArea);
  } catch (err) {
    Message.error('复制失败，请手动复制');
    console.error('复制失败:', err);
  }
};

onMounted(() => {
  fetchPartitions();
  fetchMessages();
});
</script>

<style scoped lang="less">
.container {
  padding: 16px;
}

.message-list {
  padding: 20px;
  background-color: var(--color-bg-1);
  min-height: calc(100vh - 60px);
}

.topic-header {
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--color-border-2);
}

.back-button {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 500;
  color: var(--color-text-1);
  margin-bottom: 16px;
  cursor: pointer;

  .arco-icon {
    margin-right: 8px;
  }
}

.topic-info {
  display: flex;
  gap: 40px;
}

.info-item {
  display: flex;
  flex-direction: column;
}

.info-label {
  font-size: 13px;
  color: var(--color-text-3);
  margin-bottom: 4px;
}

.info-value {
  font-size: 16px;
  font-weight: 500;
  color: var(--color-text-1);
}

.search-container {
  margin-bottom: 20px;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  border: 1px solid var(--color-border-2);
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  }
}

.search-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background-color: var(--color-fill-2);
  border-bottom: 1px solid var(--color-border-2);
}

.search-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);

  .arco-icon {
    color: var(--color-primary);
    font-size: 16px;
  }
}

.search-content {
  padding: 16px;
  background-color: var(--color-bg-1);
}

.search-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 20px;
  margin-bottom: 16px;

  &:last-child {
    margin-bottom: 0;
    padding-top: 16px;
    border-top: 1px dashed var(--color-border-2);
  }

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;

    .filter-group {
      width: 100%;
    }

    .filter-actions {
      width: 100%;
      margin-top: 8px;
      justify-content: flex-end;
    }
  }
}

.filter-group {
  display: flex;
  align-items: center;
  transition: all 0.2s ease;

  &:hover {
    .filter-label {
      color: var(--color-text-1);
    }
  }

  .offset-mode-switch {
    display: flex;
    align-items: center;
    gap: 8px;

    .mode-label {
      font-size: 13px;
      color: var(--color-text-3);
      transition: color 0.2s ease;
      user-select: none;
      cursor: pointer;

      &.active {
        color: var(--color-primary);
        font-weight: 500;
      }
    }

    :deep(.arco-switch) {
      min-width: 44px;
      height: 22px;
      background-color: var(--color-fill-4);
      transition: all 0.2s ease-in-out;

      &:hover {
        background-color: var(--color-fill-3);
      }

      .arco-switch-handle {
        width: 18px;
        height: 18px;
        top: 2px;
        transition: all 0.2s ease-in-out;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
      }

      &.arco-switch-checked {
        background-color: var(--color-primary-light-1);

        &:hover {
          background-color: var(--color-primary-light-2);
        }

        .arco-switch-handle {
          left: calc(100% - 18px - 2px);
        }
      }
    }
  }

  :deep(.arco-select) {
    .arco-select-view {
      border-radius: 4px;
      transition: all 0.2s ease;

      &:hover {
        border-color: var(--color-primary-light-3);
      }
    }
  }

  :deep(.arco-input-wrapper) {
    border-radius: 4px;
    transition: all 0.2s ease;

    &:hover {
      border-color: var(--color-primary-light-3);
    }

    .arco-input-prefix {
      color: var(--color-text-3);
    }
  }

  :deep(.arco-input-number) {
    border-radius: 4px;
    transition: all 0.2s ease;

    &:hover {
      border-color: var(--color-primary-light-3);
    }
  }
}

.filter-label {
  margin-right: 8px;
  font-size: 13px;
  color: var(--color-text-2);
  white-space: nowrap;
  font-weight: 500;
  transition: color 0.2s ease;
}

.filter-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;

  :deep(.arco-btn) {
    border-radius: 4px;
    padding: 0 16px;
    height: 32px;
    font-weight: 500;
    transition: all 0.2s ease;

    .arco-icon {
      margin-right: 4px;
    }

    &:hover {
      transform: translateY(-1px);
    }

    &:active {
      transform: translateY(0);
    }
  }
}

.message-container {
  margin-top: 16px;
}

.message-item {
  margin-bottom: 16px;
  padding: 16px;
  background-color: var(--color-bg-2);
  border-radius: 6px;
  border-left: 3px solid var(--color-primary-light-4);
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.message-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-1px);
}

.message-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--color-border-2);
  font-size: 14px;
  color: var(--color-text-2);
  align-items: center;
}

.meta-item {
  display: flex;
  align-items: center;
  white-space: nowrap;
  background-color: var(--color-fill-2);
  padding: 4px 10px;
  border-radius: 4px;
  transition: all 0.2s ease;

  :deep(.arco-icon) {
    margin-right: 6px;
    font-size: 16px;
    color: var(--color-primary-light-3);
  }
}

.meta-item:hover {
  background-color: var(--color-fill-3);
}

.meta-item b {
  color: var(--color-text-1);
  font-weight: 500;
  margin-left: 4px;
}

.message-content {
  position: relative;
  width: 100%;
  max-width: 100%;
  display: flex;
  align-items: flex-start;
}

.message-preview, .message-full {
  display: flex;
  width: 100%;
  align-items: flex-start;
}

.toggle-button {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--color-primary);
  font-size: 16px;
  transition: all 0.2s ease;
  margin-right: 10px;
  margin-top: 10px;
  width: 28px;
  height: 28px;
  border-radius: 4px;
  background-color: var(--color-fill-2);
  flex-shrink: 0;
  border: 1px solid var(--color-border-2);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.toggle-button:hover {
  color: var(--color-primary-light-3);
  background-color: var(--color-fill-3);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.content-container {
  position: relative;
  background-color: var(--color-fill-1);
  border-radius: 4px;
  padding: 12px;
  width: 100%;
  overflow-x: auto;
  border: 1px solid var(--color-border-2);
  margin-top: 4px;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.05);
}

.content-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  gap: 8px;
  z-index: 2;
}

.copy-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 4px;
  background-color: var(--color-fill-2);
  color: var(--color-text-2);
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid var(--color-border-2);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  opacity: 0.7;
}

.copy-button:hover {
  opacity: 1;
  color: var(--color-primary);
  background-color: var(--color-fill-3);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.preview-content {
  max-height: 40px;
  overflow-y: hidden;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  position: relative;
  line-height: 1.4;
  color: var(--color-text-1);
}

.preview-content::after {
  display: none;
}

.message-full {
  width: 100%;
  max-width: 100%;
}

.message-full pre {
  max-height: 500px;
  overflow-y: auto;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.4;
  width: 100%;
}

.empty-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: var(--color-text-3);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  color: var(--color-text-4);
}

.loading-container {
  display: flex;
  justify-content: center;
  padding: 40px 0;
}

.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;

  .title-left {
    display: flex;
    align-items: center;
    gap: 24px;
    flex-wrap: wrap;

    .topic-name {
      display: flex;
      align-items: center;
      font-size: 18px;
      font-weight: 600;
      color: var(--color-primary);
      background: linear-gradient(to right, var(--color-primary-light-4), var(--color-primary-light-3));
      padding: 3px 9px;
      border-radius: 6px;
      position: relative;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
      border: 1px solid var(--color-primary-light-2);
      max-width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;

      .topic-icon {
        margin-right: 10px;
        font-size: 20px;
        color: var(--color-primary);
      }

      span {
        position: relative;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      &::before {
        content: '';
        position: absolute;
        left: 0;
        top: 0;
        height: 100%;
        width: 4px;
        background-color: var(--color-primary);
        border-radius: 6px 0 0 6px;
      }

      @media (max-width: 768px) {
        font-size: 16px;
        padding: 6px 12px;

        .topic-icon {
          font-size: 18px;
          margin-right: 8px;
        }
      }

      @media (max-width: 480px) {
        width: 100%;
      }
    }

    .topic-stats {
      @media (max-width: 768px) {
        margin-top: 8px;
      }

      @media (max-width: 480px) {
        width: 100%;
        justify-content: space-between;
      }

      .stat-item {
        font-size: 14px;

        .stat-label {
          color: var(--color-text-3);
          margin-right: 8px;
        }

        .stat-value {
          color: var(--color-text-1);
          font-weight: 500;
        }

        @media (max-width: 480px) {
          font-size: 12px;
        }
      }
    }
  }
}

:deep(.tooltip-content) {
  max-width: 300px;

  p {
    margin: 4px 0;
    font-size: 12px;
    line-height: 1.5;

    &:first-child {
      margin-top: 0;
    }

    &:last-child {
      margin-bottom: 0;
    }
  }
}
</style>
