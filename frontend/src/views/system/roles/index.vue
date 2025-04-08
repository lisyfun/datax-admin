<template>
  <div class="roles" v-if="initialized">
    <a-card>
      <template #title>角色管理</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchKeyword"
            placeholder="请输入角色名称或编码"
            style="width: 300px"
            @search="handleSearch"
          />
          <a-button type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            新增角色
          </a-button>
        </a-space>
      </template>

      <a-table
        :data="roles"
        :loading="loading"
        :pagination="{
          total,
          current: page,
          pageSize: pageSize,
          showTotal: true,
          showJumper: true,
        }"
        @page-change="handlePageChange"
      >
        <template #columns>
          <a-table-column title="角色名称" data-index="name" />
          <a-table-column title="角色编码" data-index="code" />
          <a-table-column title="描述" data-index="description" />
          <a-table-column title="状态" align="center">
            <template #cell="{ record }">
              <a-switch
                :model-value="record.status === 1"
                @change="(value) => handleStatusChange(record, Boolean(value))"
              />
            </template>
          </a-table-column>
          <a-table-column title="操作" align="center">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="handleEdit(record)">
                  <template #icon><IconEdit /></template>
                  编辑
                </a-button>
                <a-button type="text" size="small" @click="handleAssignPermissions(record)">
                  <template #icon><IconSafe /></template>
                  分配权限
                </a-button>
                <a-popconfirm
                  content="确定要删除该角色吗？"
                  @ok="handleDelete(record)"
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

    <!-- 角色表单对话框 -->
    <a-modal
      v-model:visible="showRoleForm"
      :title="isEdit ? '编辑角色' : '新增角色'"
      @ok="handleRoleFormSubmit"
      @cancel="handleRoleFormCancel"
      :unmount-on-close="true"
    >
      <a-form ref="roleFormRef" :model="roleForm">
        <a-form-item field="name" label="角色名称" :rules="[{ required: true, message: '请输入角色名称' }]">
          <a-input v-model="roleForm.name" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item field="code" label="角色编码" :rules="[{ required: true, message: '请输入角色编码' }]">
          <a-input v-model="roleForm.code" placeholder="请输入角色编码" :disabled="isEdit" />
        </a-form-item>
        <a-form-item field="description" label="描述">
          <a-textarea v-model="roleForm.description" placeholder="请输入描述" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 分配权限对话框 -->
    <a-modal
      v-model:visible="showPermissionForm"
      title="分配权限"
      @ok="handlePermissionFormSubmit"
      @cancel="handlePermissionFormCancel"
      :unmount-on-close="true"
      :mask-closable="false"
      :footer="showPermissionForm ? undefined : false"
    >
      <template #default>
        <div v-if="showPermissionForm && permissionTree.length > 0">
          <a-spin :loading="permissionLoading">
            <a-tree
              v-model:checked-keys="permissionForm.permissionIds"
              :data="permissionTree"
              :field-names="{
                key: 'id',
                title: 'name',
                children: 'children',
              }"
              checkable
              :check-strictly="false"
              :default-expand-all="true"
            />
          </a-spin>
        </div>
        <div v-else-if="showPermissionForm">
          <a-empty description="暂无权限数据" />
        </div>
      </template>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { TreeNodeData } from '@arco-design/web-vue';
import type { RoleInfo, UpdateRoleParams, UpdateRolePermissionsParams } from '@/types/role';
import type { PermissionInfo } from '@/types/permission';
import { convertToTreeData } from '@/types/permission';
import * as roleApi from '@/api/role';
import * as permissionApi from '@/api/permission';
import { usePageRefresh } from '@/composables/usePageRefresh';
import {
  IconPlus,
  IconEdit,
  IconDelete,
  IconSafe,
} from '@arco-design/web-vue/es/icon';

// 组件状态控制
const isUnmounted = ref(false);

// 表格数据
const roles = ref<RoleInfo[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);
const loading = ref(false);
const searchKeyword = ref('');

// 角色表单
const showRoleForm = ref(false);
const isEdit = ref(false);
const roleFormRef = ref();
const roleForm = reactive({
  id: 0,
  name: '',
  code: '',
  description: '',
});

// 权限表单
const showPermissionForm = ref(false);
const permissionTree = ref<TreeNodeData[]>([]);
const permissionForm = reactive({
  roleId: 0,
  permissionIds: [] as number[],
});

// 添加权限加载状态
const permissionLoading = ref(false);

// 初始化标志
const initialized = ref(false);

// 获取角色列表
const fetchRoles = async () => {
  if (isUnmounted.value) return;
  try {
    loading.value = true;
    const res = await roleApi.getRoleList({
      page: page.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value,
    });
    if (!isUnmounted.value) {
      roles.value = res.data.items || [];
      total.value = res.data.total || 0;
    }
  } catch (error: any) {
    if (!isUnmounted.value) {
      Message.error(error.response?.data?.error || '获取角色列表失败');
      roles.value = [];
      total.value = 0;
    }
  } finally {
    if (!isUnmounted.value) {
      loading.value = false;
    }
  }
};

// 获取权限树
const fetchPermissionTree = async () => {
  if (!initialized.value || isUnmounted.value) return;
  try {
    permissionLoading.value = true;
    const res = await permissionApi.getPermissionTree();
    if (!isUnmounted.value) {
      const treeData = res.data.list || [];
      // 递归处理树形数据
      const processTreeData = (items: any[]): TreeNodeData[] => {
        return items.map(item => ({
          id: item.id,
          name: item.name,
          children: item.children ? processTreeData(item.children) : [],
          key: item.id,
          title: item.name,
        }));
      };
      permissionTree.value = processTreeData(treeData);
    }
  } catch (error: any) {
    if (!isUnmounted.value) {
      Message.error(error.response?.data?.error || '获取权限树失败');
      permissionTree.value = [];
    }
  } finally {
    if (!isUnmounted.value) {
      permissionLoading.value = false;
    }
  }
};

// 获取角色权限
const fetchRolePermissions = async (roleId: number) => {
  if (!initialized.value || isUnmounted.value) return;
  try {
    permissionLoading.value = true;
    const res = await roleApi.getRolePermissions(roleId);
    if (!isUnmounted.value && res.data) {
      permissionForm.permissionIds = res.data.permissions || [];
    }
  } catch (error: any) {
    if (!isUnmounted.value) {
      Message.error(error.response?.data?.error || '获取角色权限失败');
      permissionForm.permissionIds = [];
    }
  } finally {
    if (!isUnmounted.value) {
      permissionLoading.value = false;
    }
  }
};

// 搜索
const handleSearch = () => {
  if (!initialized.value || isUnmounted.value) return;
  page.value = 1;
  fetchRoles();
};

// 分页
const handlePageChange = (current: number) => {
  if (!initialized.value || isUnmounted.value) return;
  page.value = current;
  fetchRoles();
};

// 新增角色
const handleAdd = () => {
  if (!initialized.value || isUnmounted.value) return;
  isEdit.value = false;
  roleForm.id = 0;
  roleForm.name = '';
  roleForm.code = '';
  roleForm.description = '';
  showRoleForm.value = true;
};

// 编辑角色
const handleEdit = (record: RoleInfo) => {
  if (!initialized.value || isUnmounted.value) return;
  isEdit.value = true;
  roleForm.id = record.id;
  roleForm.name = record.name || '';
  roleForm.code = record.code || '';
  roleForm.description = record.description || '';
  showRoleForm.value = true;
};

// 分配权限
const handleAssignPermissions = async (record: RoleInfo) => {
  if (!initialized.value || isUnmounted.value) return;
  try {
    loading.value = true;
    permissionForm.roleId = record.id;
    permissionForm.permissionIds = []; // 重置选中的权限

    // 确保权限树已加载
    if (permissionTree.value.length === 0) {
      await fetchPermissionTree();
    }

    // 获取角色权限
    await fetchRolePermissions(record.id);

    if (!isUnmounted.value) {
      showPermissionForm.value = true;
      await nextTick();
    }
  } catch (error: any) {
    if (!isUnmounted.value) {
      Message.error(error.response?.data?.error || '加载权限数据失败');
    }
  } finally {
    if (!isUnmounted.value) {
      loading.value = false;
    }
  }
};

// 更新状态
const handleStatusChange = async (record: RoleInfo, value: boolean) => {
  if (!initialized.value || isUnmounted.value) return;
  try {
    loading.value = true;
    const updateData: UpdateRoleParams = {
      name: record.name,
      description: record.description,
      status: value ? 1 : 0
    };
    await roleApi.updateRole(record.id, updateData);
    if (!isUnmounted.value) {
      Message.success('状态更新成功');
      await fetchRoles();
    }
  } catch (error: any) {
    if (!isUnmounted.value) {
      Message.error(error.response?.data?.error || '状态更新失败');
    }
  } finally {
    if (!isUnmounted.value) {
      loading.value = false;
    }
  }
};

// 删除角色
const handleDelete = async (record: RoleInfo) => {
  if (!initialized.value || isUnmounted.value) return;
  try {
    loading.value = true;
    await roleApi.deleteRole(record.id);
    if (!isUnmounted.value) {
      Message.success('删除成功');
      await fetchRoles();
    }
  } catch (error: any) {
    if (!isUnmounted.value) {
      Message.error(error.response?.data?.error || '删除失败');
    }
  } finally {
    if (!isUnmounted.value) {
      loading.value = false;
    }
  }
};

// 提交角色表单
const handleRoleFormSubmit = async () => {
  if (!initialized.value || isUnmounted.value || !roleFormRef.value) return;
  try {
    loading.value = true;
    await roleFormRef.value.validate();

    if (isEdit.value) {
      await roleApi.updateRole(roleForm.id, {
        name: roleForm.name,
        description: roleForm.description,
        status: roles.value.find(r => r.id === roleForm.id)?.status || 1
      });
      if (!isUnmounted.value) {
        Message.success('编辑成功');
      }
    } else {
      await roleApi.createRole({
        name: roleForm.name,
        code: roleForm.code,
        description: roleForm.description
      });
      if (!isUnmounted.value) {
        Message.success('新增成功');
      }
    }
    if (!isUnmounted.value) {
      showRoleForm.value = false;
      await fetchRoles();
    }
  } catch (error: any) {
    if (!isUnmounted.value) {
      if (error.response?.data?.error) {
        Message.error(error.response.data.error);
      } else if (error.errors) {
        Message.error('表单验证失败，请检查输入');
      } else {
        Message.error(isEdit.value ? '编辑失败' : '新增失败');
      }
    }
  } finally {
    if (!isUnmounted.value) {
      loading.value = false;
    }
  }
};

// 取消角色表单
const handleRoleFormCancel = () => {
  if (!initialized.value || isUnmounted.value) return;
  showRoleForm.value = false;
  if (roleFormRef.value) {
    roleFormRef.value.resetFields();
  }
};

// 取消权限表单
const handlePermissionFormCancel = () => {
  if (!initialized.value || isUnmounted.value) return;
  showPermissionForm.value = false;
  permissionForm.roleId = 0;
  permissionForm.permissionIds = [];
};

// 提交权限表单
const handlePermissionFormSubmit = async () => {
  if (!initialized.value || isUnmounted.value) return;
  try {
    permissionLoading.value = true;
    await roleApi.updateRolePermissions(permissionForm.roleId, permissionForm.permissionIds);
    if (!isUnmounted.value) {
      Message.success('权限更新成功');
      showPermissionForm.value = false;
      permissionForm.roleId = 0;
      permissionForm.permissionIds = [];
    }
  } catch (error: any) {
    if (!isUnmounted.value) {
      Message.error(error.response?.data?.error || '权限更新失败');
    }
  } finally {
    if (!isUnmounted.value) {
      permissionLoading.value = false;
    }
  }
};

// 使用页面刷新功能
usePageRefresh(() => {
  if (initialized.value && !isUnmounted.value) {
    fetchRoles();
  }
});

// 页面加载时获取数据
onMounted(async () => {
  if (isUnmounted.value) return;

  const initializeData = async () => {
    try {
      loading.value = true;
      // Sequential loading to prevent race conditions
      await fetchRoles();
      if (!isUnmounted.value) {
        await fetchPermissionTree();
        await nextTick();
        initialized.value = true;
      }
    } catch (error) {
      if (!isUnmounted.value) {
        Message.error('页面初始化失败');
        console.error('Initialization error:', error);
      }
    } finally {
      if (!isUnmounted.value) {
        loading.value = false;
      }
    }
  };

  await initializeData();
});

// 组件卸载前清理
onBeforeUnmount(() => {
  isUnmounted.value = true;
  showRoleForm.value = false;
  showPermissionForm.value = false;
  // Reset all reactive states
  roles.value = [];
  permissionTree.value = [];
  permissionForm.permissionIds = [];
  loading.value = false;
  initialized.value = false;
});
</script>

<style scoped>
.roles {
  padding: 16px;
}
</style>
