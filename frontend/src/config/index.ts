// 后端服务器配置
export const backendConfig = {
  // API 服务器地址
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8080',
  // WebSocket 服务器地址
  wsBaseUrl: `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${import.meta.env.VITE_API_HOST || '127.0.0.1:8080'}`,
};

// 其他配置
export const config = {
  // 应用名称
  appName: 'DataX Admin',
  // 版本号
  version: '1.0.0',
};
