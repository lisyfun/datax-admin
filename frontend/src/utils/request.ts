import axios from 'axios';
import { Message } from '@arco-design/web-vue';
import { useUserStore } from '../stores/user';

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 10000,
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore();
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response;
      switch (status) {
        case 400:
          Message.error(data.error || '请求参数错误');
          break;
        case 401:
          Message.error('未登录或登录已过期');
          const userStore = useUserStore();
          userStore.logout();
          window.location.href = '/login';
          break;
        case 403:
          Message.error('没有权限');
          break;
        case 404:
          Message.error('请求的资源不存在');
          break;
        case 500:
          Message.error('服务器错误');
          break;
        default:
          Message.error('未知错误');
      }
    } else {
      Message.error('网络错误');
    }
    return Promise.reject(error);
  }
);

export default request;
