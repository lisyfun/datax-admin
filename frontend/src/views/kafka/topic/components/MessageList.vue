<template>
  <div class="message-list">
    <a-card class="general-card">
      <template #title>
        <a-space>
          <span>消息列表</span>
          <a-tag>{{ topicName }}</a-tag>
          <a-tag type="success">消费者组: {{ searchForm.groupId }}</a-tag>
        </a-space>
      </template>

      <!-- 搜索表单 -->
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="24">
          <a-form :model="searchForm" layout="inline">
            <a-form-item label="分区">
              <a-select
                v-model="searchForm.partition"
                placeholder="请选择分区"
                style="width: 120px"
                @change="handlePartitionChange"
              >
                <a-option v-for="p in partitions" :key="p" :value="p">{{ p }}</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="偏移量">
              <a-input-number
                v-model="searchForm.offset"
                placeholder="偏移量"
                style="width: 160px"
              />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="fetchLatest">
                  <template #icon><icon-caret-right /></template>
                  最新消息
                </a-button>
                <a-button @click="fetchEarliest">
                  <template #icon><icon-caret-left /></template>
                  最早消息
                </a-button>
              </a-space>
            </a-form-item>
            <a-form-item label="消息数量">
              <a-input-number
                v-model="searchForm.count"
                :min="1"
                :max="100"
                placeholder="消息数量"
                style="width: 120px"
              />
            </a-form-item>
            <a-form-item label="Key 过滤">
              <a-input
                v-model="searchForm.keyFilter"
                placeholder="Key 模糊搜索"
                style="width: 200px"
                allow-clear
              />
            </a-form-item>
            <a-form-item label="Value 过滤">
              <a-input
                v-model="searchForm.valueFilter"
                placeholder="Value 模糊搜索"
                style="width: 200px"
                allow-clear
              />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="handleSearch">
                  <template #icon><icon-search /></template>
                  搜索
                </a-button>
                <a-button @click="handleReset">
                  <template #icon><icon-refresh /></template>
                  重置
                </a-button>
                <a-button @click="handleBack">
                  <template #icon><icon-left /></template>
                  返回
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-col>
      </a-row>

      <!-- 消息列表 -->
      <a-table
        :loading="loading"
        :data="messages"
        :pagination="false"
        :bordered="false"
        :scroll="{ y: 600 }"
      >
        <template #columns>
          <a-table-column
            title="偏移量"
            data-index="offset"
            :width="100"
            align="center"
          />
          <a-table-column
            title="Key"
            data-index="key"
            :width="200"
          />
          <a-table-column
            title="时间戳"
            data-index="timestamp"
            :width="180"
          >
            <template #cell="{ record }">
              {{ formatDateTime(record.timestamp) }}
            </template>
          </a-table-column>
          <a-table-column
            title="消息内容"
            data-index="value"
          >
            <template #cell="{ record }">
              <a-typography-paragraph
                :ellipsis="{ rows: 3, expandable: true }"
                style="margin-bottom: 0"
              >
                <pre style="margin: 0; white-space: pre-wrap;">{{ formatValue(record.value) }}</pre>
              </a-typography-paragraph>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { consumeMessages, getTopicPartitions } from '@/api/kafka';
import { formatDateTime } from '@/utils/date';
import {
  IconSearch,
  IconRefresh,
  IconLeft,
  IconCaretRight,
  IconCaretLeft,
} from '@arco-design/web-vue/es/icon';

const route = useRoute();
const router = useRouter();
const loading = ref(false);
const messages = ref([]);
const partitions = ref<number[]>([]);

// 从路由参数中获取信息
const clusterId = parseInt(route.params.clusterId as string, 10);
const topicName = route.params.topicName as string;

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
    } else {
      Message.error(res.data.message || '获取分区信息失败');
    }
  } catch (err: any) {
    Message.error(err.response?.data?.message || '获取分区信息失败');
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

    const res = await consumeMessages(clusterId, topicName, {
      partition: searchForm.partition,
      offset: searchForm.offset,
      count: searchForm.count,
      keyFilter: searchForm.keyFilter,
      valueFilter: searchForm.valueFilter,
      groupId: searchForm.groupId,
    });

    if (res.data.code === 0) {
      messages.value = res.data.data || [];
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
    Message.error(err.response?.data?.message || '获取消息失败，请检查网络连接或服务器状态');
  } finally {
    loading.value = false;
  }
};

// 分区变更
const handlePartitionChange = () => {
  searchForm.offset = -1; // 切换分区时重置偏移量为最新位置
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
  fetchMessages();
};

// 获取最新消息
const fetchLatest = () => {
  searchForm.offset = -1;
  fetchMessages();
};

// 获取最早消息
const fetchEarliest = () => {
  searchForm.offset = 0;
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
  padding: 0 20px 20px;
}

:deep(.arco-form-item) {
  margin-bottom: 16px;
}
</style>
