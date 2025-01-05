<template>
  <div class="users">
    <a-card>
      <template #title>用户管理</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchKeyword"
            placeholder="请输入用户名或昵称"
            style="width: 300px"
            @search="handleSearch"
          />
          <a-button type="primary" @click="handleAdd">
            <template #icon><icon-plus /></template>
            新增用户
          </a-button>
        </a-space>
      </template>

      <a-table
        :data="users"
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
          <a-table-column title="用户名" data-index="username" />
          <a-table-column title="昵称" data-index="nickname" />
          <a-table-column title="邮箱" data-index="email" />
          <a-table-column title="状态" align="center">
            <template #cell="{ record }">
              <a-switch
                :model-value="record.status === 1"
                @change="(value: boolean) => handleStatusChange(record, value)"
              />
            </template>
          </a-table-column>
          <a-table-column title="操作" align="center">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="handleEdit(record)">
                  <template #icon><icon-edit /></template>
                  编辑
                </a-button>
                <a-button type="text" size="small" @click="handleAssignRoles(record)">
                  <template #icon><icon-user-group /></template>
                  分配角色
                </a-button>
                <a-popconfirm
                  content="确定要删除该用户吗？"
                  @ok="handleDelete(record)"
                >
                  <a-button type="text" status="danger" size="small">
                    <template #icon><icon-delete /></template>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>

    <!-- 用户表单对话框 -->
    <a-modal
      v-model:visible="showUserForm"
      :title="isEdit ? '编辑用户' : '新增用户'"
      @ok="handleUserFormSubmit"
      @cancel="handleUserFormCancel"
    >
      <a-form ref="userForm" :model="userForm" :rules="userFormRules">
        <a-form-item field="username" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model="userForm.username" placeholder="请输入用户名" :disabled="isEdit" />
        </a-form-item>
        <a-form-item
          v-if="!isEdit"
          field="password"
          label="密码"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <a-input-password v-model="userForm.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item field="nickname" label="昵称">
          <a-input v-model="userForm.nickname" placeholder="请输入昵称" />
        </a-form-item>
        <a-form-item field="email" label="邮箱" :rules="[{ type: 'email', message: '请输入正确的邮箱地址' }]">
          <a-input v-model="userForm.email" placeholder="请输入邮箱" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 分配角色对话框 -->
    <a-modal
      v-model:visible="showRoleForm"
      title="分配角色"
      @ok="handleRoleFormSubmit"
      @cancel="handleRoleFormCancel"
    >
      <a-form :model="roleForm">
        <a-form-item field="roleIds" label="角色">
          <a-select
            v-model="roleForm.roleIds"
            placeholder="请选择角色"
            multiple
            :options="roles"
            :field-names="{
              value: 'id',
              label: 'name',
            }"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { UserInfo } from '@/types/user';
import type { RoleInfo } from '@/types/role';
import * as userApi from '@/api/user';
import * as roleApi from '@/api/role';

// 表格数据
const users = ref<UserInfo[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);
const loading = ref(false);
const searchKeyword = ref('');

// 用户表单
const showUserForm = ref(false);
const isEdit = ref(false);
const userForm = reactive({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
});

// 角色表单
const showRoleForm = ref(false);
const roles = ref<RoleInfo[]>([]);
const roleForm = reactive({
  userId: 0,
  roleIds: [] as number[],
});

// 获取用户列表
const fetchUsers = async () => {
  try {
    loading.value = true;
    const res = await userApi.getUserList({
      page: page.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value,
    });
    users.value = res.data.items;
    total.value = res.data.total;
  } catch (error: any) {
    Message.error(error.response?.data?.error || '获取用户列表失败');
  } finally {
    loading.value = false;
  }
};

// 获取角色列表
const fetchRoles = async () => {
  try {
    const res = await roleApi.getRoleList({
      page: 1,
      page_size: 100,
    });
    roles.value = res.data.items;
  } catch (error: any) {
    Message.error(error.response?.data?.error || '获取角色列表失败');
  }
};

// 获取用户角色
const fetchUserRoles = async (userId: number) => {
  try {
    const res = await roleApi.getUserRoles(userId);
    roleForm.roleIds = res.data.roles;
  } catch (error: any) {
    Message.error(error.response?.data?.error || '获取用户角色失败');
  }
};

// 搜索
const handleSearch = () => {
  page.value = 1;
  fetchUsers();
};

// 分页
const handlePageChange = (current: number) => {
  page.value = current;
  fetchUsers();
};

// 新增用户
const handleAdd = () => {
  isEdit.value = false;
  userForm.id = 0;
  userForm.username = '';
  userForm.password = '';
  userForm.nickname = '';
  userForm.email = '';
  showUserForm.value = true;
};

// 编辑用户
const handleEdit = (record: UserInfo) => {
  isEdit.value = true;
  userForm.id = record.id;
  userForm.username = record.username;
  userForm.nickname = record.nickname || '';
  userForm.email = record.email || '';
  showUserForm.value = true;
};

// 分配角色
const handleAssignRoles = async (record: UserInfo) => {
  roleForm.userId = record.id;
  await fetchUserRoles(record.id);
  showRoleForm.value = true;
};

// 更新状态
const handleStatusChange = async (record: UserInfo, value: boolean) => {
  try {
    await userApi.updateUserStatus(record.id, value ? 1 : 0);
    Message.success('状态更新成功');
    fetchUsers();
  } catch (error: any) {
    Message.error(error.response?.data?.error || '状态更新失败');
  }
};

// 删除用户
const handleDelete = async (record: UserInfo) => {
  try {
    // TODO: 实现删除用户的 API
    Message.success('删除成功');
    fetchUsers();
  } catch (error: any) {
    Message.error(error.response?.data?.error || '删除失败');
  }
};

// 提交用户表单
const handleUserFormSubmit = async () => {
  try {
    // TODO: 实现新增/编辑用户的 API
    Message.success(isEdit.value ? '编辑成功' : '新增成功');
    showUserForm.value = false;
    fetchUsers();
  } catch (error: any) {
    Message.error(error.response?.data?.error || (isEdit.value ? '编辑失败' : '新增失败'));
  }
};

// 取消用户表单
const handleUserFormCancel = () => {
  showUserForm.value = false;
};

// 提交角色表单
const handleRoleFormSubmit = async () => {
  try {
    await roleApi.updateUserRoles(roleForm.userId, { role_ids: roleForm.roleIds });
    Message.success('角色分配成功');
    showRoleForm.value = false;
  } catch (error: any) {
    Message.error(error.response?.data?.error || '角色分配失败');
  }
};

// 取消角色表单
const handleRoleFormCancel = () => {
  showRoleForm.value = false;
};

onMounted(() => {
  fetchUsers();
  fetchRoles();
});
</script>

<style scoped>
.users {
  padding: 16px;
}
</style>
