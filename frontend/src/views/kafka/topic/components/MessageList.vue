<template>
  <div class="message-list">
    <!-- 主题信息头部 -->
    <div class="topic-header">
      <div class="back-button" @click="handleBack">
        <icon-left />
        <span>{{ topicName }}</span>
      </div>
      <div class="topic-info">
        <div class="info-item">
          <div class="info-label">起始偏移量</div>
          <div class="info-value">{{ topicInfo.beginningOffset || '-' }}</div>
        </div>
        <div class="info-item">
          <div class="info-label">结束偏移量</div>
          <div class="info-value">{{ topicInfo.endOffset || '-' }}</div>
        </div>
        <div class="info-item">
          <div class="info-label">消息数量</div>
          <div class="info-value">{{ topicInfo.size || '-' }}</div>
        </div>
      </div>
    </div>

    <!-- 筛选条件 -->
    <div class="filter-bar">
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
        <div class="filter-label">自动偏移量重置:</div>
        <a-select
          v-model="offsetReset"
          style="width: 120px"
          @change="handleOffsetResetChange"
        >
          <a-option value="latest">最新</a-option>
          <a-option value="earliest">最早</a-option>
        </a-select>
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
        <div class="filter-label">数量:</div>
        <a-input-number
          v-model="searchForm.count"
          :min="1"
          :max="100"
          placeholder="消息数量"
          style="width: 120px"
        />
      </div>
    </div>

    <div class="filter-bar">
      <div class="filter-group">
        <div class="filter-label">key:</div>
        <a-input
          v-model="searchForm.keyFilter"
          placeholder="filter key"
          style="width: 220px"
          allow-clear
        />
      </div>

      <div class="filter-group">
        <div class="filter-label">value:</div>
        <a-input
          v-model="searchForm.valueFilter"
          placeholder="filter value"
          style="width: 220px"
          allow-clear
        />
      </div>

      <div class="filter-group">
        <a-button type="primary" @click="handleSearch">
          拉取
        </a-button>
        <a-button @click="handleReset" style="margin-left: 8px">
          重置
        </a-button>
      </div>
    </div>

    <!-- 消息列表 -->
    <div class="message-container">
      <div v-for="(message, index) in messages" :key="index" class="message-item">
        <div class="message-meta">
          <span class="meta-item">partition: <b>{{ message.partition }}</b></span>
          <span class="meta-item">key: <b>{{ message.key || '-' }}</b></span>
          <span class="meta-item">offset: <b>{{ message.offset }}</b></span>
          <span class="meta-item">timestamp: <b>{{ formatDateTime(String(message.timestamp)) }}</b></span>
        </div>
        <div class="message-content">
          <a-typography-paragraph :ellipsis="{rows: 3, expandable: true, showTooltip: true}">
            <pre>{{ formatValue(message.value) }}</pre>
          </a-typography-paragraph>
        </div>
      </div>

      <div v-if="messages.length === 0 && !loading" class="empty-message">
        没有找到符合条件的消息
      </div>

      <div v-if="loading" class="loading-container">
        <a-spin />
      </div>
    </div>
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
} from '@arco-design/web-vue/es/icon';

const route = useRoute();
const router = useRouter();
const loading = ref(false);
const messages = ref<KafkaMessage[]>([]);
const partitions = ref<number[]>([]);
const offsetReset = ref('latest');

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
const handleOffsetResetChange = () => {
  if (offsetReset.value === 'latest') {
    searchForm.offset = -1;
  } else if (offsetReset.value === 'earliest') {
    searchForm.offset = 0;
  }
};

// 搜索
const handleSearch = () => {
  fetchMessages();
};

// 重置
const handleReset = () => {
  searchForm.offset = -1;
  searchForm.count = 10;
  searchForm.keyFilter = '';
  searchForm.valueFilter = '';
  offsetReset.value = 'latest';
  fetchMessages();
};

// 返回主题列表
const handleBack = () => {
  router.push({
    name: 'KafkaTopic',
    params: { clusterId },
  });
};

onMounted(() => {
  fetchPartitions();
  fetchMessages();
});
</script>

<style scoped lang="less">
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

.filter-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
  padding: 12px 16px;
  background-color: var(--color-bg-2);
  border-radius: 4px;
}

.filter-group {
  display: flex;
  align-items: center;
}

.filter-label {
  margin-right: 8px;
  font-size: 13px;
  color: var(--color-text-2);
  white-space: nowrap;
}

.message-container {
  margin-top: 20px;
  position: relative;
  min-height: 200px;
}

.message-item {
  margin-bottom: 16px;
  padding: 12px 16px;
  background-color: var(--color-bg-2);
  border-radius: 4px;
  border-left: 3px solid var(--color-primary-light-4);
}

.message-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px dashed var(--color-border-2);
}

.meta-item {
  font-size: 12px;
  color: var(--color-text-3);

  b {
    color: var(--color-text-1);
    font-weight: 500;
  }
}

.message-content {
  :deep(.arco-typography) {
    margin-bottom: 0;
  }

  pre {
    margin: 0;
    padding: 8px;
    background-color: var(--color-fill-1);
    border-radius: 2px;
    font-size: 13px;
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 500px;
    overflow-y: auto;
  }

  :deep(.arco-typography-operation-expand) {
    margin-left: 8px;
    color: var(--color-primary);
    font-size: 12px;
  }
}

.empty-message {
  text-align: center;
  padding: 40px 0;
  color: var(--color-text-3);
}

.loading-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.5);
}
</style>
