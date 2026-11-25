import axios from 'axios';
import type { InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { Message } from '@arco-design/web-vue';
import router from '@/router';
import { encryptData, decryptData, isEncryptionEnabled, type EncryptedResponse } from '@/utils/crypto';

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

// 后端错误响应可能是 { error: "..." }，也可能是 { code, message: "..." }，
// 统一提取可展示的文案。
function extractErrorMessage(data: any, fallback: string): string {
  if (data && typeof data === 'object') {
    return data.error || data.message || fallback;
  }
  return fallback;
}

// 如果响应被加密，尝试解密（成功/失败响应都需要处理）
function maybeDecryptResponse(response: AxiosResponse | undefined): void {
  if (!response) return;
  if (response.headers?.['content-encryption'] !== 'true' || !isEncryptionEnabled()) {
    return;
  }
  try {
    const encryptedResponse = response.data as EncryptedResponse;
    if (encryptedResponse?.data && encryptedResponse?.signature && encryptedResponse?.timestamp) {
      response.data = decryptData(encryptedResponse);
    }
  } catch (error) {
    // 解密失败时保留原始数据
  }
}

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 如果启用了加密
    if (isEncryptionEnabled()) {
      // 设置加密标识头，表示客户端支持加密响应
      config.headers = config.headers || {};
      config.headers['Accept-Encryption'] = 'true';

      // 如果有请求数据且是POST/PUT/PATCH方法，则加密请求数据
      if (config.data && (config.method === 'post' || config.method === 'put' || config.method === 'patch')) {
        try {
          // 检查是否是FormData，FormData不需要加密（文件上传等场景）
          if (!(config.data instanceof FormData)) {
            // 加密请求数据
            config.data = encryptData(config.data);
          }
        } catch (error) {
          // 如果加密失败，继续使用原始数据
        }
      }
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    maybeDecryptResponse(response);
    return response;
  },
  (error) => {
    if (error.response) {
      // 错误响应也可能被应用层加密，先尝试解密再取文案
      maybeDecryptResponse(error.response);

      // 如果是blob类型的响应，需要特殊处理
      if (error.response.config.responseType === 'blob') {
        // 尝试将blob转换为json以读取错误信息
        return new Promise((resolve, reject) => {
          const reader = new FileReader();
          reader.onload = () => {
            try {
              const errorData = JSON.parse(reader.result as string);
              error.response.data = errorData;
              Message.error(extractErrorMessage(errorData, '下载失败'));
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
      const message = extractErrorMessage(error.response.data, '');
      switch (status) {
        case 401:
          Message.error('登录已过期，请重新登录');
          // 不再需要清除token
          router.push('/login');
          break;
        case 403:
          Message.error(message || '没有权限访问');
          break;
        case 404:
          Message.error(message || '请求的资源不存在');
          break;
        case 500:
          Message.error(message || '服务器错误');
          break;
        default:
          Message.error(message || '请求失败');
      }
    } else {
      Message.error('网络错误，请检查网络连接');
    }
    return Promise.reject(error);
  }
);

export default request;
