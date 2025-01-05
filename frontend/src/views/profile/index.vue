<template>
  <div class="profile">
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card>
          <template #title>个人信息</template>
          <div class="profile-info">
            <div class="avatar-wrapper">
              <a-upload
                action="/api/v1/upload"
                :show-file-list="false"
                accept="image/*"
                @success="handleAvatarSuccess"
              >
                <a-avatar :size="120" class="avatar">
                  <img v-if="userInfo?.avatar" :src="userInfo.avatar" />
                  <icon-user v-else :style="{ fontSize: '60px' }" />
                </a-avatar>
                <div class="avatar-mask">
                  <icon-camera />
                  <span>更换头像</span>
                </div>
              </a-upload>
            </div>
            <a-descriptions :data="userDescriptions" layout="inline-vertical" />
          </div>
        </a-card>
      </a-col>
      <a-col :span="16">
        <a-card>
          <template #title>修改信息</template>
          <a-form :model="form" @submit="handleUpdateProfile">
            <a-form-item field="nickname" label="昵称">
              <a-input v-model="form.nickname" placeholder="请输入昵称" />
            </a-form-item>
            <a-form-item field="email" label="邮箱" :rules="[{ type: 'email', message: '请输入正确的邮箱地址' }]">
              <a-input v-model="form.email" placeholder="请输入邮箱" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="updating">
                保存修改
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>

        <a-card style="margin-top: 16px">
          <template #title>修改密码</template>
          <a-form :model="passwordForm" @submit="handleUpdatePassword">
            <a-form-item
              field="oldPassword"
              label="原密码"
              :rules="[{ required: true, message: '请输入原密码' }]"
            >
              <a-input-password v-model="passwordForm.oldPassword" placeholder="请输入原密码" />
            </a-form-item>
            <a-form-item
              field="newPassword"
              label="新密码"
              :rules="[{ required: true, message: '请输入新密码' }]"
            >
              <a-input-password v-model="passwordForm.newPassword" placeholder="请输入新密码" />
            </a-form-item>
            <a-form-item
              field="confirmPassword"
              label="确认密码"
              :rules="[
                { required: true, message: '请确认新密码' },
                {
                  validator: (value) => value === passwordForm.newPassword,
                  message: '两次输入的密码不一致',
                },
              ]"
            >
              <a-input-password v-model="passwordForm.confirmPassword" placeholder="请确认新密码" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="updatingPassword">
                修改密码
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useUserStore } from '../../stores/user';

const userStore = useUserStore();
const userInfo = computed(() => userStore.userInfo);

const userDescriptions = computed(() => [
  {
    label: '用户名',
    value: userInfo.value?.username || '-',
  },
  {
    label: '昵称',
    value: userInfo.value?.nickname || '-',
  },
  {
    label: '邮箱',
    value: userInfo.value?.email || '-',
  },
]);

const form = reactive({
  nickname: userInfo.value?.nickname || '',
  email: userInfo.value?.email || '',
});

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
});

const updating = ref(false);
const updatingPassword = ref(false);

const handleUpdateProfile = async () => {
  try {
    updating.value = true;
    await userStore.updateProfile(form.nickname, form.email);
    Message.success('个人信息更新成功');
  } catch (error: any) {
    Message.error(error.response?.data?.error || '更新失败');
  } finally {
    updating.value = false;
  }
};

const handleUpdatePassword = async () => {
  try {
    updatingPassword.value = true;
    await userStore.updatePassword(passwordForm.oldPassword, passwordForm.newPassword);
    Message.success('密码修改成功');
    passwordForm.oldPassword = '';
    passwordForm.newPassword = '';
    passwordForm.confirmPassword = '';
  } catch (error: any) {
    Message.error(error.response?.data?.error || '修改失败');
  } finally {
    updatingPassword.value = false;
  }
};

const handleAvatarSuccess = async (response: any) => {
  try {
    await userStore.updateProfile(undefined, undefined, response.data.url);
    Message.success('头像更新成功');
  } catch (error: any) {
    Message.error(error.response?.data?.error || '更新失败');
  }
};
</script>

<style scoped>
.profile {
  padding: 16px;
}

.profile-info {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-wrapper {
  position: relative;
  margin-bottom: 24px;
  cursor: pointer;
}

.avatar {
  border: 2px solid var(--color-primary-light-4);
}

.avatar-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
  background: rgba(0, 0, 0, 0.6);
  border-radius: 50%;
  opacity: 0;
  transition: opacity 0.3s;
}

.avatar-wrapper:hover .avatar-mask {
  opacity: 1;
}

.avatar-mask .icon {
  font-size: 24px;
  margin-bottom: 8px;
}
</style>
