import request from '@/utils/request';
import type { TerminalInfo } from '@/types/terminal';

export interface CreateTerminalData {
  name: string;
  host: string;
  port: number;
  username: string;
  password?: string;
}

export interface UpdateTerminalData {
  name: string;
  host: string;
  port: number;
  username: string;
  password?: string;
}

export interface TerminalListParams {
  page: number;
  pageSize: number;
  name?: string;
  host?: string;
  status?: string;
}

export interface TerminalListResponse {
  total: number;
  list: TerminalInfo[];
}

export default {
  // 创建终端
  createTerminal(data: CreateTerminalData) {
    return request.post<{ message: string }>('/terminals', data);
  },

  // 更新终端
  updateTerminal(id: number, data: UpdateTerminalData) {
    return request.put<{ message: string }>(`/terminals/${id}`, data);
  },

  // 删除终端
  deleteTerminal(id: number) {
    return request.delete<{ message: string }>(`/api/v1/terminals/${id}`);
  },

  // 获取终端列表
  getTerminalList(params: TerminalListParams) {
    return request.get<TerminalListResponse>('/terminals', { params });
  },

  // 获取终端详情
  getTerminalById(id: number) {
    return request.get<TerminalInfo>(`/terminals/${id}`);
  },

  // 更新终端状态
  updateTerminalStatus(id: number, status: 'online' | 'offline') {
    return request.put<{ message: string }>(`/terminals/${id}/status`, { status });
  },
};
