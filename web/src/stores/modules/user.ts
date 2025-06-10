import { defineStore } from 'pinia';
import * as userApi from '@/api/user';

interface UserState {
  userInfo: any | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    userInfo: null,
  }),

  actions: {
    async getUserInfo() {
      try {
        const { data } = await userApi.getUserInfo();
        this.userInfo = data;
        return data;
      } catch (error) {
        console.error('获取用户信息失败:', error);
        throw error;
      }
    },

    async logout() {
      await userApi.logout(); // 调用登出API
      this.userInfo = null;
    },
  },
});
