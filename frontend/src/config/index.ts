// 后端服务器配置
export const backendConfig = {
  // API 基础路径
  apiBaseUrl: '/datax/api',
  // WebSocket 基础路径
  wsBaseUrl: `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}`,
};

// 其他配置
export const appConfig = {
  // 应用名称
  appName: 'DATAX ADMIN',
  // 应用基础路径
  basePath: '/datax',
  // 应用版本
  version: '1.0.0',
};
