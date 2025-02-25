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
            <template #icon><IconPlus /></template>
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
          pageSize,
          showTotal: true,
          showJumper: true,
          showPageSize: true,
        }"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      >
        <template #columns>
          <a-table-column title="ID" data-index="id" />
          <a-table-column title="用户名" data-index="username" />
          <a-table-column title="昵称" data-index="nickname" />
          <a-table-column title="邮箱" data-index="email" />
          <a-table-column title="角色" align="center">
            <template #cell="{ record }">
              <a-space wrap>
                <a-tag
                  v-for="role in record.roles"
                  :key="role.id"
                  color="blue"
                >
                  {{ role.name }}
                </a-tag>
                <a-tag v-if="!record.roles?.length" color="gray">暂无角色</a-tag>
              </a-space>
            </template>
          </a-table-column>
          <a-table-column title="状态" align="center">
            <template #cell="{ record }">
              <a-switch
                :model-value="record.status === 1"
                @change="(value) => handleStatusChange(record, value ? 1 : 0)"
              />
            </template>
          </a-table-column>
          <a-table-column title="操作" align="center" :width="300">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="handleEdit(record)">
                  <template #icon><IconEdit /></template>
                  编辑
                </a-button>
                <a-button type="text" size="small" @click="handleAssignRoles(record)">
                  <template #icon><IconUserGroup /></template>
                  分配角色
                </a-button>
                <a-button type="text" size="small" @click="handleChangePassword(record)">
                  <template #icon><IconLock /></template>
                  修改密码
                </a-button>
                <a-button type="text" size="small" @click="handleResetPassword(record)">
                  <template #icon><IconRefresh /></template>
                  重置密码
                </a-button>
                <a-popconfirm
                  content="确定要删除该用户吗？"
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

    <!-- 用户表单对话框 -->
    <a-modal
      v-model:visible="showUserForm"
      :title="isEdit ? '编辑用户' : '新增用户'"
      @ok="handleUserFormSubmit"
      @cancel="handleUserFormCancel"
    >
      <a-form ref="userFormRef" :model="userForm">
        <a-form-item field="username" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model="userForm.username" placeholder="请输入用户名" :disabled="isEdit">
            <template #prefix>
              <IconUser />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item
          v-if="!isEdit"
          field="password"
          label="密码"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <a-input-password v-model="userForm.password" placeholder="请输入密码">
            <template #prefix>
              <IconLock />
            </template>
          </a-input-password>
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

    <!-- 密码表单对话框 -->
    <a-modal
      v-model:visible="showPasswordForm"
      title="修改密码"
      @ok="handlePasswordFormSubmit"
      @cancel="handlePasswordFormCancel"
    >
      <a-form ref="passwordFormRef" :model="passwordForm">
        <a-form-item field="oldPassword" label="旧密码" :rules="[{ required: true, message: '请输入旧密码' }]">
          <a-input-password v-model="passwordForm.oldPassword" placeholder="请输入旧密码" />
        </a-form-item>
        <a-form-item field="newPassword" label="新密码" :rules="[{ required: true, message: '请输入新密码' }]">
          <a-input-password v-model="passwordForm.newPassword" placeholder="请输入新密码" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 重置密码表单对话框 -->
    <a-modal
      v-model:visible="showResetPasswordForm"
      title="重置密码"
      @ok="handleResetPasswordFormSubmit"
      @cancel="handleResetPasswordFormCancel"
    >
      <a-form ref="resetPasswordFormRef" :model="resetPasswordForm">
        <a-form-item field="password" label="新密码" :rules="[{ required: true, message: '请输入新密码' }]">
          <a-input-password v-model="resetPasswordForm.password" placeholder="请输入新密码" />
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
import { encryptPassword } from '@/utils/crypto';
import { usePageRefresh } from '@/composables/usePageRefresh';
import {
  IconPlus,
  IconEdit,
  IconDelete,
  IconUserGroup,
  IconLock,
  IconRefresh,
  IconUser,
} from '@arco-design/web-vue/es/icon';

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
const userFormRef = ref();
const userForm = reactive({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
});

// 密码表单
const showPasswordForm = ref(false);
const passwordFormRef = ref();
const passwordForm = reactive({
  id: 0,
  oldPassword: '',
  newPassword: '',
});

// 重置密码表单
const showResetPasswordForm = ref(false);
const resetPasswordFormRef = ref();
const resetPasswordForm = reactive({
  id: 0,
  password: '',
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
  loading.value = true;
  try {
    const res = await userApi.getUserList({
      page: page.value,
      pageSize: pageSize.value,
      username: searchKeyword.value || undefined,
    });
    users.value = res.data.list;
    total.value = res.data.total;
  } catch (error) {
    console.error('获取用户列表失败:', error);
  } finally {
    loading.value = false;
  }
};

// 使用页面刷新功能
usePageRefresh(() => {
  fetchUsers();
});

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
    // 直接使用后端返回的角色ID数组
    roleForm.roleIds = res.data.roles || [];
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

// 处理每页条数变化
const handlePageSizeChange = (size: number) => {
  pageSize.value = size;
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
  userForm.nickname = record.nickname;
  userForm.email = record.email;
  userForm.password = ''; // 编辑时不需要密码
  showUserForm.value = true;
};

// 分配角色
const handleAssignRoles = async (record: UserInfo) => {
  roleForm.userId = record.id;
  // 获取用户当前的角色
  await fetchUserRoles(record.id);
  showRoleForm.value = true;
};

// 更新状态
const handleStatusChange = async (record: UserInfo, value: number): Promise<void> => {
  try {
    // 先更新UI状态
    const targetUser = users.value.find(u => u.id === record.id);
    if (targetUser) {
      targetUser.status = value;
    }

    // 发送请求
    await userApi.updateUserStatus(record.id, value);
    Message.success(value === 1 ? '已启用' : '已禁用');
  } catch (error: any) {
    // 发生错误时恢复原状态
    const targetUser = users.value.find(u => u.id === record.id);
    if (targetUser) {
      targetUser.status = record.status;
    }
    Message.error(error.response?.data?.error || '状态更新失败');
  }
};

// 删除用户
const handleDelete = async (record: UserInfo) => {
  try {
    await userApi.deleteUser(record.id);
    Message.success('删除成功');
    fetchUsers();
  } catch (error: any) {
    Message.error(error.response?.data?.error || '删除用户失败');
  }
};

// 提交用户表单
const handleUserFormSubmit = async () => {
  try {
    if (!userFormRef.value) return;
    await userFormRef.value.validate();

    if (isEdit.value) {
      await userApi.updateProfile({
        nickname: userForm.nickname,
        email: userForm.email
      });
      Message.success('编辑成功');
    } else {
      await userApi.register({
        username: userForm.username,
        password: encryptPassword(userForm.password),
        nickname: userForm.nickname,
        email: userForm.email
      });
      Message.success('新增成功');
    }
    showUserForm.value = false;
    fetchUsers();
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

// 取消用户表单
const handleUserFormCancel = () => {
  showUserForm.value = false;
  if (userFormRef.value) {
    userFormRef.value.resetFields();
  }
};

// 提交角色表单
const handleRoleFormSubmit = async () => {
  try {
    await roleApi.updateUserRoles(roleForm.userId, roleForm.roleIds);
    Message.success('角色分配成功');
    showRoleForm.value = false;
    // 刷新用户列表以显示最新的角色信息
    await fetchUsers();
  } catch (error: any) {
    Message.error(error.response?.data?.error || '角色分配失败');
  }
};

// 取消角色表单
const handleRoleFormCancel = () => {
  showRoleForm.value = false;
  // 清空角色选择
  roleForm.roleIds = [];
};

// 修改密码
const handleChangePassword = (record: UserInfo) => {
  passwordForm.id = record.id;
  passwordForm.oldPassword = '';
  passwordForm.newPassword = '';
  showPasswordForm.value = true;
};

// 重置密码
const handleResetPassword = (record: UserInfo) => {
  resetPasswordForm.id = record.id;
  resetPasswordForm.password = '';
  showResetPasswordForm.value = true;
};

// 提交密码表单
const handlePasswordFormSubmit = async () => {
  try {
    if (!passwordFormRef.value) return;
    await passwordFormRef.value.validate();

    await userApi.updatePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    });
    Message.success('密码修改成功');
    showPasswordForm.value = false;
  } catch (error: any) {
    if (error.response?.data?.error) {
      Message.error(error.response.data.error);
    } else if (error.errors) {
      Message.error('表单验证失败，请检查输入');
    } else {
      Message.error('密码修改失败');
    }
  }
};

// 取消密码表单
const handlePasswordFormCancel = () => {
  showPasswordForm.value = false;
  if (passwordFormRef.value) {
    passwordFormRef.value.resetFields();
  }
};

// 提交重置密码表单
const handleResetPasswordFormSubmit = async () => {
  try {
    if (!resetPasswordFormRef.value) return;
    await resetPasswordFormRef.value.validate();

    await userApi.resetPassword(resetPasswordForm.id, {
      password: (resetPasswordForm.password)
    });
    Message.success('密码重置成功');
    showResetPasswordForm.value = false;
  } catch (error: any) {
    if (error.response?.data?.error) {
      Message.error(error.response.data.error);
    } else if (error.errors) {
      Message.error('表单验证失败，请检查输入');
    } else {
      Message.error('密码重置失败');
    }
  }
};

// 取消重置密码表单
const handleResetPasswordFormCancel = () => {
  showResetPasswordForm.value = false;
  if (resetPasswordFormRef.value) {
    resetPasswordFormRef.value.resetFields();
  }
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
