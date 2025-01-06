<template>
  <div class="permissions">
    <a-card>
      <template #title>权限管理</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchKeyword"
            placeholder="请输入权限名称或编码"
            style="width: 300px"
            @search="handleSearch"
          />
          <a-button type="primary" @click="() => handleAdd()">
            <template #icon><IconPlus /></template>
            新增权限
          </a-button>
        </a-space>
      </template>

      <a-table
        :data="permissions"
        :loading="loading"
        :pagination="false"
        row-key="id"
        :tree-props="{
          children: 'children'
        }"
      >
        <template #columns>
          <a-table-column title="权限名称" data-index="name" />
          <a-table-column title="权限编码" data-index="code" />
          <a-table-column title="类型" align="center">
            <template #cell="{ record }">
              <a-tag :color="record.type === 'menu' ? 'blue' : 'green'">
                {{ record.type === 'menu' ? '菜单' : '按钮' }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="路径" data-index="path" />
          <a-table-column title="组件" data-index="component" />
          <a-table-column title="图标" align="center">
            <template #cell="{ record }">
              <i v-if="record.icon" :class="record.icon"></i>
              <span v-else>-</span>
            </template>
          </a-table-column>
          <a-table-column title="排序" data-index="sort" align="center" />
          <a-table-column title="状态" align="center">
            <template #cell="{ record }">
              <a-switch
                :model-value="record.status === 1"
                @update:model-value="(value) => handleStatusChange(record, Boolean(value))"
              />
            </template>
          </a-table-column>
          <a-table-column title="操作" align="center">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="() => handleAdd(record)">
                  <template #icon><IconPlus /></template>
                  新增子权限
                </a-button>
                <a-button type="text" size="small" @click="() => handleEdit(record)">
                  <template #icon><IconEdit /></template>
                  编辑
                </a-button>
                <a-popconfirm
                  content="确定要删除该权限吗？"
                  @ok="() => handleDelete(record)"
                >
                  <a-button type="text" status="danger" size="small">
                    <template #icon><IconDelete /></template>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>

    <!-- 权限表单对话框 -->
    <a-modal
      v-model:visible="showPermissionForm"
      :title="isEdit ? '编辑权限' : '新增权限'"
      @ok="handlePermissionFormSubmit"
      @cancel="handlePermissionFormCancel"
    >
      <a-form ref="formRef" :model="formData" :rules="formRules">
        <a-form-item field="name" label="权限名称">
          <a-input v-model="formData.name" placeholder="请输入权限名称" />
        </a-form-item>
        <a-form-item field="code" label="权限编码">
          <a-input v-model="formData.code" placeholder="请输入权限编码" :disabled="isEdit" />
        </a-form-item>
        <a-form-item field="type" label="类型">
          <a-radio-group v-model="formData.type">
            <a-radio value="menu">菜单</a-radio>
            <a-radio value="button">按钮</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item field="parent_id" label="上级权限">
          <a-tree-select
            v-model="formData.parent_id"
            :data="permissionTree"
            placeholder="请选择上级权限"
            allow-clear
          />
        </a-form-item>
        <template v-if="formData.type === 'menu'">
          <a-form-item field="path" label="路径">
            <a-input v-model="formData.path" placeholder="请输入路径" />
          </a-form-item>
          <a-form-item field="component" label="组件">
            <a-input v-model="formData.component" placeholder="请输入组件路径，例如：views/system/user/index" />
          </a-form-item>
          <a-form-item field="icon" label="图标">
            <a-input v-model="formData.icon" placeholder="请输入图标类名" />
          </a-form-item>
        </template>
        <a-form-item field="sort" label="排序">
          <a-input-number v-model="formData.sort" placeholder="请输入排序" :min="0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { FormInstance, FieldRule } from '@arco-design/web-vue';
import { PermissionInfo } from '@/types/permission';
import { convertToTreeData } from '@/types/permission';
import * as permissionApi from '@/api/permission';
import {
  IconPlus,
  IconEdit,
  IconDelete,
} from '@arco-design/web-vue/es/icon';

// 表格数据
const permissions = ref<PermissionInfo[]>([]);
const permissionTree = computed(() => convertToTreeData(permissions.value));
const loading = ref(false);

// 表单相关
const formRef = ref<FormInstance>();
const isEdit = ref(false);
const showPermissionForm = ref(false);
const formData = reactive({
  id: 0,
  name: '',
  code: '',
  type: 'menu' as 'menu' | 'button',
  parent_id: undefined as number | undefined,
  path: '',
  component: '',
  icon: '',
  sort: 0,
});

const formRules: Record<string, FieldRule[]> = {
  name: [
    { required: true, message: '请输入权限名称' }
  ],
  code: [
    { required: true, message: '请输入权限编码' },
    { validator: (value: string) => /^[a-z][a-z0-9_]*$/.test(value), message: '权限编码只能包含小写字母、数字和下划线，且必须以字母开头' }
  ],
  type: [
    { required: true, message: '请选择类型' }
  ],
  path: [
    { validator: (value: string) => !value || value.startsWith('/'), message: '路径必须以/开头' }
  ],
  sort: [
    { required: true, type: 'number', message: '请输入排序' },
    { type: 'number', min: 0, message: '排序不能小于0' }
  ]
};

// 添加搜索相关变量
const searchKeyword = ref('');

// 添加搜索处理函数
const handleSearch = () => {
  fetchPermissions();
};

// 获取权限列表
const fetchPermissions = async () => {
  try {
    loading.value = true;
    const response = await permissionApi.getPermissionTree();
    permissions.value = response.data.list;
  } catch (error) {
    Message.error('获取权限列表失败');
  } finally {
    loading.value = false;
  }
};

const handleAdd = (parent?: PermissionInfo) => {
  isEdit.value = false;
  formData.id = 0;
  formData.name = '';
  formData.code = '';
  formData.type = 'menu';
  formData.parent_id = parent?.id;
  formData.path = '';
  formData.component = '';
  formData.icon = '';
  formData.sort = 0;
  showPermissionForm.value = true;
};

const handleStatusChange = async (record: PermissionInfo, value: boolean) => {
  try {
    await permissionApi.updatePermission(record.id, {
      name: record.name,
      type: record.type,
      status: value ? 1 : 0,
      parent_id: record.parent_id,
      path: record.path,
      component: record.component,
      icon: record.icon,
      sort: record.sort
    });
    Message.success('状态更新成功');
    await fetchPermissions();
  } catch (error) {
    Message.error('状态更新失败');
  }
};

// 编辑权限
const handleEdit = (record: PermissionInfo) => {
  isEdit.value = true;
  formData.id = record.id;
  formData.name = record.name;
  formData.code = record.code;
  formData.type = record.type;
  formData.parent_id = record.parent_id;
  formData.path = record.path || '';
  formData.component = record.component || '';
  formData.icon = record.icon || '';
  formData.sort = record.sort;
  showPermissionForm.value = true;
};

// 删除权限
const handleDelete = async (record: PermissionInfo) => {
  try {
    await permissionApi.deletePermission(record.id);
    Message.success('删除成功');
    await fetchPermissions();
  } catch (error: any) {
    Message.error(error.response?.data?.error || '删除失败');
  }
};

// 提交权限表单
const handlePermissionFormSubmit = async () => {
  try {
    if (!formRef.value) return;
    await formRef.value.validate();

    if (isEdit.value) {
      await permissionApi.updatePermission(formData.id, {
        name: formData.name,
        code: formData.code,
        type: formData.type,
        parent_id: formData.parent_id,
        path: formData.path,
        component: formData.component,
        icon: formData.icon,
        sort: formData.sort,
      });
    } else {
      await permissionApi.createPermission({
        name: formData.name,
        code: formData.code,
        type: formData.type,
        parent_id: formData.parent_id,
        path: formData.path,
        component: formData.component,
        icon: formData.icon,
        sort: formData.sort,
      });
    }
    Message.success(isEdit.value ? '编辑成功' : '新增成功');
    showPermissionForm.value = false;
    await fetchPermissions();
  } catch (error: any) {
    if (error.response?.data?.error) {
      Message.error(error.response.data.error);
    } else if (error.errors) {
      Message.error('表单验证失败，请检查输入');
    } else {
      Message.error(isEdit.value ? '编辑失败' : '新增失败');
    }
  }
};

// 取消权限表单
const handlePermissionFormCancel = () => {
  showPermissionForm.value = false;
};

onMounted(() => {
  fetchPermissions();
});
</script>

<style scoped>
.permissions {
  padding: 16px;
}
</style>
