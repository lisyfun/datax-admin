import axios from 'axios';
import type { InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { Message } from '@arco-design/web-vue';
import router from '@/router';

// 创建 axios 实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 600000,
  headers: {
    'Content-Type': 'application/json',
  },
  // 添加 withCredentials 配置，使请求带上 cookie
  withCredentials: true,
});

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 移除 JWT token 认证相关代码
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  (error) => {
    if (error.response) {
      // 如果是blob类型的响应，需要特殊处理
      if (error.response.config.responseType === 'blob') {
        // 尝试将blob转换为json以读取错误信息
        return new Promise((resolve, reject) => {
          const reader = new FileReader();
          reader.onload = () => {
            try {
              const errorData = JSON.parse(reader.result as string);
              error.response.data = errorData;
              Message.error(errorData.message || '下载失败');
            } catch (e) {
              Message.error('下载失败');
            }
            reject(error);
          };
          reader.onerror = () => {
            Message.error('下载失败');
            reject(error);
          };
          reader.readAsText(error.response.data);
        });
      }

      const { status } = error.response;
      switch (status) {
        case 401:
          Message.error('登录已过期，请重新登录');
          // 不再需要清除token
          router.push('/login');
          break;
        case 403:
          Message.error('没有权限访问');
          break;
        case 404:
          Message.error('请求的资源不存在');
          break;
        case 500:
          Message.error('服务器错误');
          break;
        default:
          Message.error(error.response.data?.message || '请求失败');
      }
    } else {
      Message.error('网络错误，请检查网络连接');
    }
    return Promise.reject(error);
  }
);

export default request;
