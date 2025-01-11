import { defineStore } from 'pinia';
import * as userApi from '@/api/user';

interface UserState {
  token: string | null;
  userInfo: any | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: localStorage.getItem('token'),
    userInfo: null,
  }),

  actions: {
    setToken(token: string) {
      this.token = token;
      localStorage.setItem('token', token);
    },

    clearToken() {
      this.token = null;
      localStorage.removeItem('token');
    },

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
      this.clearToken();
      this.userInfo = null;
    },
  },
});
