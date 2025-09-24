<template>
  <div class="container">
    <a-card class="general-card" title="集群管理">
      <a-row>
        <a-col :flex="1">
          <a-form :model="searchForm" :label-col-props="{ span: 0 }" :wrapper-col-props="{ span: 24 }" label-align="left">
            <a-row :gutter="16">
              <a-col :span="6">
                <a-form-item field="search" label="">
                  <a-input v-model="searchForm.search" placeholder="请输入集群名称" @press-enter="search" />
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
          <a-button type="primary" @click="openForm()" v-permission="'tools.kafka.cluster.create'">
            <template #icon>
              <icon-plus />
            </template>
            新增
          </a-button>
        </a-col>
      </a-row>
      <a-table
        row-key="id"
        :loading="loading"
        :pagination="{
          total: pagination.total,
          current: pagination.current,
          pageSize: pagination.pageSize,
          showTotal: true,
          showJumper: true,
          showPageSize: true,
        }"
        :columns="columns"
        :data="renderData"
        :bordered="false"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #name="{ record }">
          <a-link
            :disabled="!record.status"
            :style="{ cursor: record.status ? 'pointer' : 'not-allowed' }"
            @click="record.status && goToTopics(record)"
          >
            {{ record.name }}
          </a-link>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status ? 'green' : 'red'">
            {{ record.status ? '正常' : '不可用' }}
          </a-tag>
        </template>
        <template #topicCount="{ record }">
          <a-tag color="blue">{{ record.topicCount || 0 }}</a-tag>
        </template>
        <template #brokerCount="{ record }">
          <a-tag color="green">{{ record.brokerCount || 0 }}</a-tag>
        </template>
        <template #consumerGroupCount="{ record }">
          <a-tag color="purple">{{ record.consumerGroupCount || 0 }}</a-tag>
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openForm(record)" v-permission="'tools.kafka.cluster.update'">
              <icon-edit />
              编辑
            </a-button>
            <a-popconfirm content="确定要删除该集群吗？" @ok="handleDelete(record)">
              <a-button type="text" status="danger" size="small" v-permission="'tools.kafka.cluster.delete'">
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
      :title="form.id ? '编辑集群' : '新增集群'"
      @ok="handleSubmit"
      @cancel="closeForm"
    >
      <a-form ref="formRef" :model="form" :rules="rules" :label-col-props="{ span: 6 }" :wrapper-col-props="{ span: 18 }">
        <a-form-item field="name" label="集群名称" :validate-trigger="['change', 'blur']">
          <a-input v-model="form.name" placeholder="请输入集群名称" />
        </a-form-item>
        <a-form-item field="brokerServers" label="Broker 服务器" :validate-trigger="['change', 'blur']">
          <a-input v-model="form.brokerServers" placeholder="请输入 Broker 服务器地址，多个地址用逗号分隔" />
        </a-form-item>
        <a-form-item field="securityProtocol" label="安全协议" :validate-trigger="['change', 'blur']">
          <a-select v-model="form.securityProtocol" placeholder="请选择安全协议">
            <a-option value="PLAINTEXT">PLAINTEXT</a-option>
            <a-option value="SASL_PLAINTEXT">SASL_PLAINTEXT</a-option>
            <a-option value="SSL">SSL</a-option>
            <a-option value="SASL_SSL">SASL_SSL</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="saslMechanism" label="SASL 机制" :validate-trigger="['change', 'blur']">
          <a-select v-model="form.saslMechanism" placeholder="请选择 SASL 机制" :disabled="form.securityProtocol === 'PLAINTEXT'">
            <a-option value="PLAIN">PLAIN</a-option>
            <a-option value="SCRAM-SHA-256">SCRAM-SHA-256</a-option>
            <a-option value="SCRAM-SHA-512">SCRAM-SHA-512</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="username" label="用户名" :validate-trigger="['change', 'blur']">
          <a-input v-model="form.username" placeholder="请输入用户名" :disabled="form.securityProtocol === 'PLAINTEXT'" />
        </a-form-item>
        <a-form-item field="password" label="密码" :validate-trigger="['change', 'blur']">
          <a-input-password v-model="form.password" placeholder="请输入密码" :disabled="form.securityProtocol === 'PLAINTEXT'" />
        </a-form-item>
        <a-form-item field="description" label="描述" :validate-trigger="['change', 'blur']">
          <a-textarea v-model="form.description" placeholder="请输入描述" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import useLoading from '@/hooks/loading';
import { queryClusterList, createCluster, updateCluster, deleteCluster } from '@/api/kafka';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
interface KafkaCluster {
  id: number;
  name: string;
  brokerServers: string;
  securityProtocol: string;
  saslMechanism: string;
  username: string;
  password: string;
  description: string;
  topicCount: number;
  brokerCount: number;
  consumerGroupCount: number;
  status: boolean;
}

interface SearchForm {
  search: string;
}

interface Pagination {
  total: number;
  current: number;
  pageSize: number;
}

const router = useRouter();
const { loading, setLoading } = useLoading(true);
const renderData = ref<KafkaCluster[]>([]);
const formRef = ref();
const visible = ref(false);

const searchForm = reactive<SearchForm>({
  search: '',
});

const basePagination: Pagination = {
  total: 0,
  current: 1,
  pageSize: 10,
};

const pagination = reactive<Pagination>({
  ...basePagination,
});

const form = reactive<KafkaCluster>({
  id: 0,
  name: '',
  brokerServers: '',
  securityProtocol: 'PLAINTEXT',
  saslMechanism: 'PLAIN',
  username: '',
  password: '',
  description: '',
  topicCount: 0,
  brokerCount: 0,
  consumerGroupCount: 0,
  status: true,
});

const rules = {
  name: [{ required: true, message: '请输入集群名称' }],
  brokerServers: [{ required: true, message: '请输入 Broker 服务器地址' }],
  securityProtocol: [{ required: true, message: '请选择安全协议' }],
  saslMechanism: [{
    validator: (value: string, callback: (error?: string) => void) => {
      if (form.securityProtocol !== 'PLAINTEXT' && !value) {
        callback('请选择 SASL 机制');
      } else {
        callback();
      }
    },
  }],
  username: [{
    validator: (value: string, callback: (error?: string) => void) => {
      if (form.securityProtocol !== 'PLAINTEXT' && !value) {
        callback('请输入用户名');
      } else {
        callback();
      }
    },
  }],
  password: [{
    validator: (value: string, callback: (error?: string) => void) => {
      if (form.securityProtocol !== 'PLAINTEXT' && !value) {
        callback('请输入密码');
      } else {
        callback();
      }
    },
  }],
};

const columns = computed<TableColumnData[]>(() => [
  {
    title: '集群名称',
    dataIndex: 'name',
    slotName: 'name',
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 100,
    align: 'center',
    slotName: 'status',
  },
  {
    title: 'Broker 服务器',
    dataIndex: 'brokerServers',
  },
  {
    title: '主题数',
    dataIndex: 'topicCount',
    width: 100,
    align: 'center',
    slotName: 'topicCount',
  },
  {
    title: 'Broker 数',
    dataIndex: 'brokerCount',
    width: 100,
    align: 'center',
    slotName: 'brokerCount',
  },
  {
    title: '安全协议',
    dataIndex: 'securityProtocol',
  },
  {
    title: 'SASL 机制',
    dataIndex: 'saslMechanism',
  },
  {
    title: '操作',
    dataIndex: 'operations',
    slotName: 'operations',
    width: 160,
    fixed: 'right',
  },
]);

const fetchData = async () => {
  setLoading(true);
  try {
    const res = await queryClusterList({
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchForm.search,
    });
    if (res.data.code === 0) {
      renderData.value = res.data.data.items;
      pagination.total = res.data.data.total;
    } else {
      Message.error(res.data.message);
    }
  } catch (err) {
    Message.error('获取数据失败');
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

const openForm = (record?: KafkaCluster) => {
  if (record) {
    Object.assign(form, record);
  }
  visible.value = true;
};

const closeForm = () => {
  form.id = 0;
  form.name = '';
  form.brokerServers = '';
  form.securityProtocol = 'PLAINTEXT';
  form.saslMechanism = 'PLAIN';
  form.username = '';
  form.password = '';
  form.description = '';
  form.topicCount = 0;
  form.brokerCount = 0;
  form.consumerGroupCount = 0;
  form.status = true;
  visible.value = false;
};

const handleSubmit = async () => {
  const res = await formRef.value?.validate();
  if (!res) {
    try {
      if (form.id) {
        await updateCluster(form);
        Message.success('更新成功');
      } else {
        await createCluster(form);
        Message.success('创建成功');
      }
      closeForm();
      fetchData();
    } catch (err: any) {
      Message.error(err.response?.data?.message || '操作失败');
    }
  }
};

const handleDelete = async (record: KafkaCluster) => {
  try {
    await deleteCluster(record.id);
    Message.success('删除成功');
    fetchData();
  } catch (err: any) {
    Message.error(err.response?.data?.message || '删除失败');
  }
};

const goToTopics = (record: KafkaCluster) => {
  router.push({
    name: 'KafkaTopic',
    params: { clusterId: record.id },
    query: { clusterName: record.name }
  });
};

// 监听安全协议变更
watch(() => form.securityProtocol, (newValue) => {
  // 重新验证相关字段
  formRef.value?.validateField(['saslMechanism', 'username', 'password']);

  // 如果切换到 PLAINTEXT，清空认证相关字段
  if (newValue === 'PLAINTEXT') {
    form.saslMechanism = '';
    form.username = '';
    form.password = '';
  }
});

fetchData();
</script>

<style scoped lang="less">
.container {
  padding: 16px;
}
</style>
