// 终端信息
export interface TerminalInfo {
  id: number;
  name: string;
  host: string;
  port: number;
  username: string;
  status: 'online' | 'offline';
  lastSeen: string;
  createdAt: string;
}

// 创建终端请求参数
export interface CreateTerminalData {
  name: string;
  host: string;
  port: number;
  username: string;
  password: string;
}

// 更新终端请求参数
export interface UpdateTerminalData {
  name: string;
  host: string;
  port: number;
  username: string;
  password?: string;
}

// 终端列表请求参数
export interface TerminalListParams {
  page: number;
  pageSize: number;
  name?: string;
  host?: string;
  status?: string;
}

// 终端列表响应数据
export interface TerminalListResponse {
  list: TerminalInfo[];
  total: number;
}

// WebSocket 消息类型
export interface TerminalMessage {
  type: 'input' | 'output';
  data: string;
}
