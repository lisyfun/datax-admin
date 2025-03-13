<template>
  <div class="container">
    <Breadcrumb :items="['Kafka 管理', '主题管理']" />
    <a-card class="general-card" title="主题管理">
      <a-row>
        <a-col :flex="1">
          <a-form :model="searchForm" :label-col-props="{ span: 6 }" :wrapper-col-props="{ span: 18 }" label-align="left">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item field="search" label="搜索">
                  <a-input v-model="searchForm.search" placeholder="请输入主题名称" @press-enter="search" />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-space>
                  <a-button type="primary" @click="search">
                    <template #icon>
                      <icon-search />
                    </template>
                    搜索
                  </a-button>
                  <a-button @click="reset">
                    <template #icon>
                      <icon-refresh />
                    </template>
                    重置
                  </a-button>
                </a-space>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 32px" direction="vertical" />
        <a-col :flex="'86px'" style="text-align: right">
          <a-button type="primary" @click="openForm()">
            <template #icon>
              <icon-plus />
            </template>
            新增
          </a-button>
        </a-col>
      </a-row>
      <a-table
        row-key="name"
        :loading="loading"
        :pagination="pagination"
        :columns="columns"
        :data="renderData"
        :bordered="false"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #name="{ record }">
          <a-link @click="openConsumer(record)">{{ record.name }}</a-link>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openForm(record)">
              <icon-edit />
              编辑
            </a-button>
            <a-popconfirm content="确定要删除该主题吗？" @ok="handleDelete(record)">
              <a-button type="text" status="danger" size="small">
                <icon-delete />
                删除
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:visible="visible"
      :title="form.name ? '编辑主题' : '新增主题'"
      @ok="handleSubmit"
      @cancel="closeForm"
    >
      <a-form ref="formRef" :model="form" :rules="rules" :label-col-props="{ span: 6 }" :wrapper-col-props="{ span: 18 }">
        <a-form-item field="name" label="主题名称" :validate-trigger="['change', 'blur']">
          <a-input v-model="form.name" placeholder="请输入主题名称" :disabled="!!form.name" />
        </a-form-item>
        <a-form-item field="partitions" label="分区数" :validate-trigger="['change', 'blur']">
          <a-input-number v-model="form.partitions" :min="1" placeholder="请输入分区数" />
        </a-form-item>
        <a-form-item v-if="!form.name" field="replicas" label="副本数" :validate-trigger="['change', 'blur']">
          <a-input-number v-model="form.replicas" :min="1" placeholder="请输入副本数" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="consumerVisible"
      title="消费消息"
      ok-text="消费"
      @ok="handleConsume"
      @cancel="closeConsumer"
    >
      <a-form ref="consumerFormRef" :model="consumerForm" :rules="consumerRules" :label-col-props="{ span: 6 }" :wrapper-col-props="{ span: 18 }">
        <a-form-item field="partition" label="分区">
          <a-select v-model="consumerForm.partition">
            <a-option v-for="p in partitions" :key="p" :value="p">{{ p }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="offset" label="偏移量">
          <a-input-number v-model="consumerForm.offset" :min="-1" placeholder="请输入偏移量，-1 表示最新位置" />
        </a-form-item>
        <a-form-item field="count" label="消息数量">
          <a-input-number v-model="consumerForm.count" :min="1" :max="100" placeholder="请输入消息数量" />
        </a-form-item>
        <a-form-item field="keyFilter" label="消息 Key 过滤">
          <a-input v-model="consumerForm.keyFilter" placeholder="请输入消息 Key 过滤条件" />
        </a-form-item>
        <a-form-item field="valueFilter" label="消息内容过滤">
          <a-input v-model="consumerForm.valueFilter" placeholder="请输入消息内容过滤条件" />
        </a-form-item>
      </a-form>
      <template v-if="messages.length > 0">
        <a-divider />
        <a-table :columns="messageColumns" :data="messages" :pagination="false" :scroll="{ y: 400 }">
          <template #value="{ record }">
            <a-typography-paragraph :ellipsis="{ rows: 3, expandable: true }">
              {{ record.value }}
            </a-typography-paragraph>
          </template>
        </a-table>
      </template>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from 'vue';
import { useRoute } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import useLoading from '@/hooks/loading';
import { queryTopicList, createTopic, alterTopic, deleteTopic, getTopicPartitions, consumeMessages } from '@/api/kafka';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';

const route = useRoute();
const clusterId = computed(() => {
  const id = route.params.clusterId;
  return Array.isArray(id) ? parseInt(id[0], 10) : parseInt(id, 10);
});
const { loading, setLoading } = useLoading(true);
const renderData = ref([]);
const formRef = ref();
const visible = ref(false);
const consumerFormRef = ref();
const consumerVisible = ref(false);
const partitions = ref([]);
const messages = ref([]);

const searchForm = reactive({
  search: '',
});

const basePagination = {
  total: 0,
  current: 1,
  pageSize: 10,
};

const pagination = reactive({
  ...basePagination,
});

const form = reactive({
  name: '',
  partitions: 1,
  replicas: 1,
});

const consumerForm = reactive({
  name: '',
  partition: 0,
  offset: -1,
  count: 10,
  keyFilter: '',
  valueFilter: '',
});

const rules = {
  name: [{ required: true, message: '请输入主题名称' }],
  partitions: [{ required: true, message: '请输入分区数' }],
  replicas: [{ required: true, message: '请输入副本数' }],
};

const consumerRules = {
  partition: [{ required: true, message: '请选择分区' }],
  offset: [{ required: true, message: '请输入偏移量' }],
  count: [{ required: true, message: '请输入消息数量' }],
};

const columns = computed<TableColumnData[]>(() => [
  {
    title: '主题名称',
    dataIndex: 'name',
    slotName: 'name',
  },
  {
    title: '分区数',
    dataIndex: 'partitions',
  },
  {
    title: '副本数',
    dataIndex: 'replicas',
  },
  {
    title: '平均日志大小',
    dataIndex: 'avgLogSize',
  },
  {
    title: '总日志大小',
    dataIndex: 'logSize',
  },
  {
    title: '操作',
    dataIndex: 'operations',
    slotName: 'operations',
    width: 160,
  },
]);

const messageColumns = computed<TableColumnData[]>(() => [
  {
    title: '分区',
    dataIndex: 'partition',
    width: 100,
  },
  {
    title: '偏移量',
    dataIndex: 'offset',
    width: 100,
  },
  {
    title: '消息 Key',
    dataIndex: 'key',
    width: 200,
  },
  {
    title: '时间戳',
    dataIndex: 'timestamp',
    width: 200,
  },
  {
    title: '消息内容',
    dataIndex: 'value',
    slotName: 'value',
  },
]);

const fetchData = async () => {
  setLoading(true);
  try {
    const res = await queryTopicList(clusterId.value, {
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchForm.search,
    });
    if (res.data.code === 0) {
      renderData.value = res.data.data.items;
      pagination.total = res.data.data.total;
    } else {
      Message.error(res.data.message || '获取主题列表失败');
    }
  } catch (err: any) {
    Message.error(err.response?.data?.message || '获取主题列表失败');
  } finally {
    setLoading(false);
  }
};

const search = () => {
  pagination.current = 1;
  fetchData();
};

const reset = () => {
  searchForm.search = '';
  pagination.current = 1;
  fetchData();
};

const onPageChange = (current: number) => {
  pagination.current = current;
  fetchData();
};

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  fetchData();
};

const openForm = (record?: any) => {
  if (record) {
    form.name = record.name;
    form.partitions = record.partitions;
  }
  visible.value = true;
};

const closeForm = () => {
  form.name = '';
  form.partitions = 1;
  form.replicas = 1;
  visible.value = false;
};

const handleSubmit = async () => {
  const res = await formRef.value?.validate();
  if (!res) {
    try {
      if (form.name) {
        await alterTopic(clusterId.value, form.name, {
          partitions: form.partitions,
        });
        Message.success('更新成功');
      } else {
        await createTopic(clusterId.value, {
          name: form.name,
          partitions: form.partitions,
          replicas: form.replicas,
        });
        Message.success('创建成功');
      }
      closeForm();
      fetchData();
    } catch (err) {
      // handle error
    }
  }
};

const handleDelete = async (record: any) => {
  try {
    await deleteTopic(clusterId.value, record.name);
    Message.success('删除成功');
    fetchData();
  } catch (err) {
    // handle error
  }
};

const openConsumer = async (record: any) => {
  try {
    const { data } = await getTopicPartitions(clusterId.value, record.name);
    partitions.value = data;
    consumerForm.name = record.name;
    consumerForm.partition = data[0] || 0;
    messages.value = [];
    consumerVisible.value = true;
  } catch (err) {
    // handle error
  }
};

const closeConsumer = () => {
  consumerForm.name = '';
  consumerForm.partition = 0;
  consumerForm.offset = -1;
  consumerForm.count = 10;
  consumerForm.keyFilter = '';
  consumerForm.valueFilter = '';
  messages.value = [];
  consumerVisible.value = false;
};

const handleConsume = async () => {
  const res = await consumerFormRef.value?.validate();
  if (!res) {
    try {
      const { data } = await consumeMessages(clusterId.value, consumerForm.name, {
        partition: consumerForm.partition,
        offset: consumerForm.offset,
        count: consumerForm.count,
        keyFilter: consumerForm.keyFilter,
        valueFilter: consumerForm.valueFilter,
      });
      messages.value = data;
      if (data.length === 0) {
        Message.info('没有找到符合条件的消息');
      }
    } catch (err) {
      // handle error
    }
  }
};

fetchData();
</script>

<style scoped lang="less">
.container {
  padding: 0 20px 20px;
}
</style>
