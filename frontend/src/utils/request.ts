import axios from 'axios';
import { Message } from '@arco-design/web-vue';
import router from '@/router';

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api',
  timeout: 10000,
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
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
      const { status } = error.response;
      if (status === 401) {
        Message.error('登录已过期，请重新登录');
        localStorage.removeItem('token');
        router.push('/login');
      } else if (status === 403) {
        Message.error('没有权限访问');
      } else if (status === 404) {
        Message.error('请求的资源不存在');
      } else if (status === 500) {
        Message.error('服务器错误');
      } else {
        Message.error(error.response.data?.error || '请求失败');
      }
    } else {
      Message.error('网络错误，请检查网络连接');
    }
    return Promise.reject(error);
  }
);

export default request;
