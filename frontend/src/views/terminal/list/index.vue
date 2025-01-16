<template>
  <div class="container">
    <a-card class="general-card" title="终端列表">
      <template #extra>
        <a-button type="primary" @click="handleAdd">
          <template #icon>
            <IconPlus />
          </template>
          新建终端
        </a-button>
      </template>
      <a-table
        :loading="loading"
        :data="tableData"
        :pagination="pagination"
        @page-change="onPageChange"
      >
        <template #columns>
          <a-table-column title="终端ID" data-index="id" />
          <a-table-column title="名称" data-index="name" />
          <a-table-column title="主机" data-index="host" />
          <a-table-column title="状态" data-index="status">
            <template #cell="{ record }">
              <a-tag :color="record.status === 'online' ? 'green' : 'red'">
                {{ record.status === 'online' ? '在线' : '离线' }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="创建时间" data-index="created_at" />
          <a-table-column title="操作" align="center">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="handleConnect(record)">
                  <template #icon>
                    <IconLink />
                  </template>
                  连接
                </a-button>
                <a-button type="text" size="small" @click="handleEdit(record)">
                  <template #icon>
                    <IconEdit />
                  </template>
                  编辑
                </a-button>
                <a-popconfirm
                  content="确定要删除这个终端吗？"
                  @ok="handleDelete(record)"
                >
                  <a-button type="text" size="small" status="danger">
                    <template #icon>
                      <IconDelete />
                    </template>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { Message } from '@arco-design/web-vue';
import { IconPlus, IconEdit, IconDelete, IconLink } from '@arco-design/web-vue/es/icon';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';

const loading = ref(false);
const tableData = ref([]);
const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
});

const fetchData = async (page = 1) => {
  loading.value = true;
  try {
    // TODO: 调用后端API获取终端列表数据
    // const res = await api.getTerminalList({ page, pageSize: pagination.pageSize });
    // tableData.value = res.data;
    // pagination.total = res.total;
  } catch (error) {
    Message.error('获取终端列表失败');
  } finally {
    loading.value = false;
  }
};

const onPageChange = (current: number) => {
  pagination.current = current;
  fetchData(current);
};

const handleAdd = () => {
  // TODO: 实现新建终端逻辑
};

const handleEdit = (record: any) => {
  // TODO: 实现编辑终端逻辑
};

const handleDelete = async (record: any) => {
  // TODO: 实现删除终端逻辑
};

const handleConnect = (record: any) => {
  // TODO: 实现连接终端逻辑
};

// 初始加载数据
fetchData();
</script>

<style scoped>
.container {
  padding: 0 20px 20px 20px;
}

.general-card {
  margin-top: 16px;
}
</style>
