import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { UserInfo } from '../types/user';
import * as userApi from '../api/user';

export const useUserStore = defineStore('user', () => {
  const token = ref<string>('');
  const userInfo = ref<UserInfo | null>(null);

  async function login(username: string, password: string) {
    const res = await userApi.login({ username, password });
    token.value = res.data.token;
    localStorage.setItem('token', res.data.token);
    await getUserInfo();
  }

  async function register(username: string, password: string, nickname?: string, email?: string) {
    await userApi.register({ username, password, nickname, email });
  }

  async function getUserInfo() {
    const res = await userApi.getUserInfo();
    userInfo.value = res.data;
  }

  async function updatePassword(oldPassword: string, newPassword: string) {
    await userApi.updatePassword({ old_password: oldPassword, new_password: newPassword });
  }

  async function updateProfile(nickname?: string, email?: string, avatar?: string) {
    await userApi.updateProfile({ nickname, email, avatar });
    await getUserInfo();
  }

  function logout() {
    token.value = '';
    userInfo.value = null;
    localStorage.removeItem('token');
  }

  // 初始化时从 localStorage 获取 token
  const storedToken = localStorage.getItem('token');
  if (storedToken) {
    token.value = storedToken;
  }

  return {
    token,
    userInfo,
    login,
    register,
    getUserInfo,
    updatePassword,
    updateProfile,
    logout,
  };
});
